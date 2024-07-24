package resource

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/api"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/shared"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/util"
)

const (
	RuleFile         = "rule.go"
	UsedRuleJsonFile = "used_rules.json"
	EmbededFs        = "embededfs"
)

func findFiles(dir fs.FS, path string) (map[string][]string, error) {
	files := make(map[string][]string, 0)
	err := fs.WalkDir(dir, path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			parent := filepath.Dir(path)
			files[parent] = append(files[parent], path)
		}
		return nil
	})
	return files, err
}

func ReadRuleFile(path string) (string, error) {
	data, err := pkg.ExportRuleFS().ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func isRuleDefined(text string, rule *api.InstFuncRule) bool {
	if rule.Version != "" &&
		!strings.Contains(text, util.StringQuote(rule.GetVersion())) {
		return false
	}
	util.Assert(rule.ImportPath != "", "import path must be set")
	if !strings.Contains(text, util.StringQuote(rule.GetImportPath())) {
		return false
	}
	if rule.ReceiverType != "" &&
		!strings.Contains(text, util.StringQuote(rule.ReceiverType)) {
		return false
	}
	if rule.Function != "" &&
		!strings.Contains(text, util.StringQuote(rule.Function)) {
		return false
	}
	if rule.OnEnter != "" &&
		!strings.Contains(text, util.StringQuote(rule.OnEnter)) {
		return false
	}
	if rule.OnExit != "" &&
		!strings.Contains(text, util.StringQuote(rule.OnExit)) {
		return false
	}
	return true
}

func isHookDefined(text string, rule *api.InstFuncRule) bool {
	util.Assert(rule.OnEnter != "" || rule.OnExit != "", "hook must be set")
	root, err := shared.ParseAstFromSource(text)
	if err != nil {
		log.Printf("failed to parse ast from source: %v", err)
		return false
	}
	if rule.OnEnter != "" {
		if shared.FindFuncDecl(root, rule.OnEnter) == nil {
			return false
		}
	}
	if rule.OnExit != "" {
		if shared.FindFuncDecl(root, rule.OnExit) == nil {
			return false
		}
	}
	return true
}

func FindRuleFile(path string) (string, error) {
	files, err := findFiles(pkg.ExportRuleFS(), ".")
	if err != nil {
		return "", err
	}
	for _, paths := range files {
		for _, p := range paths {
			if strings.HasSuffix(p, path) {
				return p, nil
			}
		}
	}
	return "", nil
}

func FindRuleFiles(rule api.InstRule) ([]string, error) {
	files, err := findFiles(pkg.ExportRuleFS(), ".")
	if err != nil {
		return nil, err
	}
	switch rt := rule.(type) {
	case *api.InstFuncRule:
		if rt.UseRaw {
			util.ShouldNotReachHereT("insane rule type")
		}
		// For function rule, we need to find files where onEnter and onExit
		// are defined
		for _, paths := range files {
			// Now all path files are in same directory, iterate over them to see
			// if there is a rule.go file which indicates an instrumentation rule
			// definition. If so, link it to the rule resource
			for i, path := range paths {
				if strings.HasSuffix(path, RuleFile) {
					text, err := ReadRuleFile(path)
					if err == nil && isRuleDefined(text, rt) {
						// No rule.go please
						paths = append(paths[:i], paths[i+1:]...)
						// Find files where onEnter and onExit are defined
						for _, p := range paths {
							text, err := ReadRuleFile(p)
							if err != nil {
								return nil, err
							}
							if isHookDefined(text, rt) {
								return []string{p}, nil
							}
						}
						return nil, nil
					}
					break
				}
			}
		}
	case *api.InstFileRule:
		// For file rule, we need to find the file with the same name
		// as the rule
		for _, paths := range files {
			for _, p := range paths {
				if strings.HasSuffix(p, rt.FileName) {
					return []string{p}, nil
				}
			}
		}
		return nil, nil
	case *api.InstStructRule:
		util.ShouldNotReachHereT("insane rule type")
	}
	return nil, nil
}

var hash2Rules = make(map[uint64]api.InstRule)

// Rationale of localizeFileRule
// All file dependencies of file rule are located within go/embed file system,
// while func rule dependencies are located within the local file system. Any
// time we want to use a rule, we need to determine whether it is a file rule
// or a func rule, and then locate the file in the corresponding file system.
// That's somewhat inconvenient, to simplify the process, we localize all file
// dependencies of file rules to the local file system. In this way, all kinds
// of rules uses the local file system.
func localizeFileRule(rule *api.InstFileRule) (string, error) {
	target := shared.GetPreprocessLogPath(filepath.Join(EmbededFs, rule.FileName))

	if shared.InPreprocess() {
		exist, err := util.PathExists(target)
		if err != nil {
			return "", fmt.Errorf("failed to check file existence: %w", err)
		}
		if exist {
			return target, nil
		}
		res, err := FindRuleFiles(rule)
		if err != nil {
			return "", fmt.Errorf("failed to find rule file: %w", err)
		}
		if len(res) == 0 {
			return "", fmt.Errorf("rule file not found")
		}
		content, err := ReadRuleFile(res[0])
		if err != nil {
			return "", fmt.Errorf("failed to read rule file: %w", err)
		}
		err = os.MkdirAll(filepath.Dir(target), 0777)
		if err != nil {
			return "", fmt.Errorf("failed to create directory: %w", err)
		}
		_, err = util.WriteStringToFile(target, content)
		if err != nil {
			return "", fmt.Errorf("failed to write file: %w", err)
		}
	}
	return target, nil
}

func InitRules() error {
	for _, rule := range api.Rules {
		// Localize file rule first to get a consistent hash
		if rl, ok := rule.(*api.InstFileRule); ok {
			target, err := localizeFileRule(rl)
			if err != nil {
				return fmt.Errorf("failed to localize file rule: %w", err)
			}
			rl.FileName = target
		}
		h, err := shared.HashStruct(rule)
		if err != nil {
			return fmt.Errorf("failed to hash rule: %w", err)
		}
		if shared.Verbose {
			log.Printf("Rule %v hashed to %d", rule, h)
		}
		hash2Rules[h] = rule
	}
	return nil
}

func FindFuncRuleByHash(hash uint64) *api.InstFuncRule {
	return FindRuleByHash(hash).(*api.InstFuncRule)
}

func FindFileRuleByHash(hash uint64) *api.InstFileRule {
	return FindRuleByHash(hash).(*api.InstFileRule)
}

func FindStructRuleByHash(hash uint64) *api.InstStructRule {
	return FindRuleByHash(hash).(*api.InstStructRule)
}

func FindRuleByHash(hash uint64) api.InstRule {
	util.Assert(len(hash2Rules) > 0, "rule hash not initialized")
	rl, ok := hash2Rules[hash]
	util.Assert(ok, "rule not found")
	return rl
}

func StoreRuleBundles(bundles []*RuleBundle) error {
	shared.GuaranteeInPreprocess()

	ruleLines := make([]string, 0)
	for _, bundle := range bundles {
		bs, err := json.Marshal(*bundle)
		if err != nil {
			return fmt.Errorf("failed to marshal bundle: %w", err)
		}
		ruleLines = append(ruleLines, string(bs))
	}
	ruleFile := shared.GetPreprocessLogPath(UsedRuleJsonFile)
	_, err := util.WriteStringToFile(ruleFile, strings.Join(ruleLines, "\n"))
	if err != nil {
		return fmt.Errorf("failed to write used rules: %w", err)
	}
	return nil
}

func LoadRuleBundles() ([]*RuleBundle, error) {
	shared.GuaranteeInInstrument()

	ruleFile := shared.GetPreprocessLogPath(UsedRuleJsonFile)
	data, err := util.ReadFile(ruleFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read used rules: %w", err)
	}
	lines := strings.Split(data, "\n")
	bundles := make([]*RuleBundle, 0)
	for _, line := range lines {
		if line == "" {
			continue
		}
		bundle := &RuleBundle{}
		err := json.Unmarshal([]byte(line), bundle)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal bundle: %w", err)
		}
		bundles = append(bundles, bundle)
	}
	return bundles, nil
}

package resource

import (
	"fmt"
	"go/token"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/dave/dst"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/api"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/shared"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/util"
)

// RuleBundle is a collection of rules that matched with one compilation action
type RuleBundle struct {
	PackageName      string   // Short package name, e.g. "echo"
	ImportPath       string   // Full import path, e.g. "github.com/labstack/echo/v4"
	FileRules        []uint64 // File rules
	File2FuncRules   map[string]map[string][]uint64
	File2StructRules map[string]map[string][]uint64
}

func NewRuleBundle(importPath string) *RuleBundle {
	return &RuleBundle{
		PackageName:      "",
		ImportPath:       importPath,
		FileRules:        make([]uint64, 0),
		File2FuncRules:   make(map[string]map[string][]uint64),
		File2StructRules: make(map[string]map[string][]uint64),
	}
}

func (rb *RuleBundle) Merge(new *RuleBundle) (*RuleBundle, error) {
	if !new.IsValid() {
		return rb, nil
	}
	util.Assert(rb.ImportPath == new.ImportPath, "inconsistent import path")
	util.Assert(rb.PackageName == new.PackageName, "inconsistent package name")
	fileRules := make(map[uint64]bool)
	for _, h := range rb.FileRules {
		fileRules[h] = true
	}
	for _, h := range new.FileRules {
		if _, exist := fileRules[h]; !exist {
			rb.FileRules = append(rb.FileRules, h)
		}
	}

	for file, rules := range new.File2FuncRules {
		if _, exist := rb.File2FuncRules[file]; !exist {
			rb.File2FuncRules[file] = make(map[string][]uint64)
		}
		for fn, hashes := range rules {
			if _, exist := rb.File2FuncRules[file][fn]; !exist {
				rb.File2FuncRules[file][fn] = make([]uint64, 0)
			}
			rb.File2FuncRules[file][fn] =
				append(rb.File2FuncRules[file][fn], hashes...)
		}
	}
	for file, rules := range new.File2StructRules {
		if _, exist := rb.File2StructRules[file]; !exist {
			rb.File2StructRules[file] = make(map[string][]uint64)
		}
		for st, hashes := range rules {
			if _, exist := rb.File2StructRules[file][st]; !exist {
				rb.File2StructRules[file][st] = make([]uint64, 0)
			}
			rb.File2StructRules[file][st] =
				append(rb.File2StructRules[file][st], hashes...)
		}
	}
	return rb, nil
}

func (rb *RuleBundle) AddFile2FuncRule(file string, rule *api.InstFuncRule) {
	fn := rule.Function + "," + rule.ReceiverType
	util.Assert(fn != "", "sanity check")
	h, err := shared.HashStruct(*rule)
	if err != nil {
		log.Fatalf("Failed to hash struct %v", rule)
	}
	if _, exist := rb.File2FuncRules[file]; !exist {
		rb.File2FuncRules[file] = make(map[string][]uint64)
		rb.File2FuncRules[file][fn] = []uint64{h}
	} else {
		rb.File2FuncRules[file][fn] = append(rb.File2FuncRules[file][fn], h)
	}
}

func (rb *RuleBundle) AddFile2StructRule(file string, rule *api.InstStructRule) {
	st := rule.StructType
	util.Assert(st != "", "sanity check")
	h, err := shared.HashStruct(*rule)
	if err != nil {
		log.Fatalf("Failed to hash struct %v", rule)
	}
	if _, exist := rb.File2StructRules[file]; !exist {
		rb.File2StructRules[file] = make(map[string][]uint64)
		rb.File2StructRules[file][st] = []uint64{h}
	} else {
		rb.File2StructRules[file][st] = append(rb.File2StructRules[file][st], h)
	}
}

func (rb *RuleBundle) AddFileRule(rule *api.InstFileRule) {
	h, err := shared.HashStruct(*rule)
	if err != nil {
		log.Fatalf("Failed to hash struct %v", rule)
	}
	rb.FileRules = append(rb.FileRules, h)
}

// splitVersion splits the version string into three parts, major, minor and patch.
func splitVersion(version string) (int, int, int) {
	util.Assert(strings.Contains(version, "."), "invalid version format")
	var (
		majorVersionStr string
		minorVersionStr string
		patchVersionStr string
		majorVersion    int
		minorVersion    int
		patchVersion    int
		err             error
	)

	dotIdx := strings.Index(version, ".")
	lastDotIdx := strings.LastIndex(version, ".")

	majorVersionStr = version[:dotIdx]
	majorVersion, err = strconv.Atoi(majorVersionStr)
	util.Assert(err == nil,
		"invalid version format, major version is %s", majorVersionStr)

	minorVersionStr = version[dotIdx+1 : lastDotIdx]
	minorVersion, err = strconv.Atoi(minorVersionStr)
	util.Assert(err == nil,
		"invalid version format, minor version is %s", minorVersionStr)

	patchVersionStr = version[lastDotIdx+1:]
	patchVersion, err = strconv.Atoi(patchVersionStr)
	util.Assert(err == nil,
		"invalid version format, patch version is %s", patchVersionStr)

	return majorVersion, minorVersion, patchVersion
}

// splitVersionRange splits the version range into two parts, start and end.
func splitVersionRange(vr string) (string, string) {
	util.Assert(strings.Contains(vr, ","), "invalid version range format")
	util.Assert(strings.Contains(vr, "["), "invalid version range format")
	util.Assert(strings.Contains(vr, ")"), "invalid version range format")

	start := vr[1:strings.Index(vr, ",")]
	end := vr[strings.Index(vr, ",")+1 : len(vr)-1]
	return start, end
}

// verifyVersion splits version into three parts, and then verify their format.
func verifyVersion(version string) bool {
	dotIdx := strings.Index(version, ".")
	lastDotIdx := strings.LastIndex(version, ".")
	majorVersionStr := version[:dotIdx]
	minorVersionStr := version[dotIdx+1 : lastDotIdx]
	patchVersionStr := version[lastDotIdx+1:]
	return verifyPureNumber(patchVersionStr) &&
		verifyPureNumber(minorVersionStr) &&
		verifyPureNumber(majorVersionStr)
}

func verifyPureNumber(subVersion string) bool {
	for _, c := range subVersion {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// findVersionFromPath extracts the version number from file path. For example
// for the path "github.com/gin-gonic/gin@v1.9.1", it returns "v1.9.1". If the
// path does not contain version number, it returns an empty string.
var versionRegexp = regexp.MustCompile(`@v\d+\.\d+\.\d+(-.*?)?/`)

func extractVersion(path string) string {
	version := versionRegexp.FindString(path)
	if version == "" {
		return ""
	}
	// Extract version number from the string
	return version[1 : len(version)-1]
}

// matchVersion checks if the version string matches the version range in the
// rule. The version range is in format [start, end), where start is inclusive
// and end is exclusive. If the rule version string is empty, it always matches.
func matchVersion(version string, ruleVersion string) (bool, error) {
	// Fast path, always match if the rule version is not specified
	if ruleVersion == "" {
		return true, nil
	}
	// Check if both rule version and package version are in sane
	if !strings.Contains(version, "v") {
		return false, fmt.Errorf("invalid version %v %v",
			version, ruleVersion)
	}
	if !strings.Contains(ruleVersion, "[") ||
		!strings.Contains(ruleVersion, ")") ||
		!strings.Contains(ruleVersion, ",") ||
		strings.Contains(ruleVersion, "v") {
		return false, fmt.Errorf("invalid version format in rule %v %v",
			version, ruleVersion)
	}
	// Remove extra whitespace from the rule version string
	ruleVersion = strings.ReplaceAll(ruleVersion, " ", "")

	// Ignore the leading "v" in the version string
	version = version[1:]

	if !verifyVersion(version) {
		return false, fmt.Errorf("matched snapshot version: v%v", version)
	}

	// Extract version number from the string
	majorVersion, minorVersion, patchVersion := splitVersion(version)
	if majorVersion > 999 || minorVersion > 999 || patchVersion > 999 {
		return false, fmt.Errorf("illegal version number")
	}

	// Compare the version with the rule version, the rule version is in the
	// format [start, end), where start is inclusive and end is exclusive
	versionStart, versionEnd := splitVersionRange(ruleVersion)
	majorStart, minorStart, patchStart := splitVersion(versionStart)
	majorEnd, minorEnd, patchEnd := splitVersion(versionEnd)

	U1, U2, U3 := 1000000, 1000, 1
	ruleStart := majorStart*U1 + minorStart*U2 + patchStart*U3
	ruleEnd := majorEnd*U1 + minorEnd*U2 + patchEnd*U3
	v := majorVersion*U1 + minorVersion*U2 + patchVersion*U3
	if v >= ruleStart && v < ruleEnd {
		return true, nil
	}
	return false, nil
}

func findAvailableRules() []api.InstRule {
	// Disable all rules
	if shared.DisableRules == "*" {
		return make([]api.InstRule, 0)
	}

	availables := make([]api.InstRule, len(api.Rules))
	copy(availables, api.Rules)
	if shared.DisableRules == "" {
		return availables
	}

	list := strings.Split(shared.DisableRules, ",")
	rules := make([]api.InstRule, 0)
	for _, v := range availables {
		disabled := false
		for _, disable := range list {
			if v.GetRuleName() != "" && disable == v.GetRuleName() {
				disabled = true
				break
			}
			if disable == v.GetImportPath() {
				disabled = true
				break
			}
		}
		if !disabled {
			rules = append(rules, v)
		}
	}
	if shared.Verbose {
		log.Printf("Available rule: %v", rules)
	}
	return rules
}

func (rb *RuleBundle) IsValid() bool {
	return rb != nil &&
		(len(rb.FileRules) > 0 ||
			len(rb.File2FuncRules) > 0 ||
			len(rb.File2StructRules) > 0)
}

func MatchFuncDecl(decl dst.Decl, function string, receiverType string) bool {
	funcDecl, ok := decl.(*dst.FuncDecl)
	if !ok {
		return false
	}
	if funcDecl.Name.Name != function {
		return false
	}
	if receiverType != "" {
		if !shared.HasReceiver(funcDecl) {
			return false
		}
		switch recvTypeExpr := funcDecl.Recv.List[0].Type.(type) {
		case *dst.StarExpr:
			return "*"+recvTypeExpr.X.(*dst.Ident).Name == receiverType
		case *dst.Ident:
			return recvTypeExpr.Name == receiverType
		default:
			util.Unimplemented()
		}
	} else {
		if shared.HasReceiver(funcDecl) {
			return false
		}
	}

	return true
}

func MatchStructDecl(decl dst.Decl, structType string) bool {
	if genDecl, ok := decl.(*dst.GenDecl); ok {
		if genDecl.Tok == token.TYPE {
			if typeSpec, ok := genDecl.Specs[0].(*dst.TypeSpec); ok {
				if typeSpec.Name.Name == structType {
					return true
				}
			}
		}
	}
	return false
}

type RuleMatcher struct {
	AvailableRules map[string][]api.InstRule
}

func NewRuleMatcher() *RuleMatcher {
	rules := make(map[string][]api.InstRule)
	for _, rule := range findAvailableRules() {
		rules[rule.GetImportPath()] = append(rules[rule.GetImportPath()], rule)
	}
	return &RuleMatcher{AvailableRules: rules}
}

// MatchRuleBundle gives compilation arguments and finds out all interested rules
// for it.
func (rm *RuleMatcher) MatchRuleBundle(importPath string, candidates []string) *RuleBundle {
	util.Assert(importPath != "", "sanity check")
	availables := make([]api.InstRule, len(rm.AvailableRules[importPath]))

	// Okay, we are interested in these candidates, let's read it and match with
	// the instrumentation rule, but first we need to check if the package name
	// are already registered, to avoid futile effort
	copy(availables, rm.AvailableRules[importPath])
	if len(availables) == 0 {
		return nil // fast fail
	}
	bundle := NewRuleBundle(importPath)
	for _, candidate := range candidates {
		// It's not a go file, ignore silently
		if !shared.IsGoFile(candidate) {
			continue
		}

		// Parse the file content
		file := candidate
		fileAst, _ := shared.ParseAstFromFile(file)
		if fileAst == nil {
			// Failed to parse the file, stop here and log only
			// sicne it's a tolerant failure
			log.Printf("Failed to parse file %s from local fs", file)
			continue
		}
		if bundle.PackageName == "" {
			bundle.PackageName = fileAst.Name.Name
		} else {
			util.Assert(bundle.PackageName == fileAst.Name.Name,
				"inconsistent package name")
		}
		// Match the rules with the file
		for i := len(availables) - 1; i >= 0; i-- {
			rule := availables[i]
			util.Assert(rule.GetImportPath() == importPath, "sanity check")
			matched, err := matchVersion(extractVersion(file), rule.GetVersion())
			if err != nil {
				log.Printf("Failed to match version %v", err)
				continue
			}
			if !matched {
				continue
			}
			// Basic check passed, let's match with the rule precisely
			if rl, ok := rule.(*api.InstFileRule); ok {
				// Rule is valid nevertheless, save it
				log.Printf("Matched file rule %s", rule)
				bundle.AddFileRule(rl)
				availables = append(availables[:i], availables[i+1:]...)
			} else {
				valid := false
				for _, decl := range fileAst.Decls {
					if genDecl, ok := decl.(*dst.GenDecl); ok {
						// We are only interested in struct type declaration
						if rl, ok := rule.(*api.InstStructRule); ok {
							if MatchStructDecl(genDecl, rl.StructType) {
								log.Printf("Matched struct rule %s", rule)
								bundle.AddFile2StructRule(file, rl)
								valid = true
							}
						}
					} else if funcDecl, ok := decl.(*dst.FuncDecl); ok {
						// We are only interested in function declaration for func rule
						if rl, ok := rule.(*api.InstFuncRule); ok {
							if MatchFuncDecl(funcDecl, rl.Function, rl.ReceiverType) {
								log.Printf("Matched func rule %s", rule)
								bundle.AddFile2FuncRule(file, rl)
								valid = true
							}
						}

					}
				}
				if valid {
					// Remove the rule from the available rules
					availables = append(availables[:i], availables[i+1:]...)
				}
			}
		}
	}
	return bundle
}

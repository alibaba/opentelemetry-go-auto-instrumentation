// Copyright (c) 2024 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package preprocess

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/config"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/errc"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/resource"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/util"
	"github.com/dave/dst"
)

type ruleMatcher struct {
	availableRules map[string][]resource.InstRule
	moduleVersions []*vendorModule // vendor used only
}

func newRuleMatcher() *ruleMatcher {
	rules := make(map[string][]resource.InstRule)
	for _, rule := range findAvailableRules() {
		rules[rule.GetImportPath()] = append(rules[rule.GetImportPath()], rule)
	}
	if config.GetConf().Verbose {
		util.Log("Available rules: %v", rules)
	}
	return &ruleMatcher{availableRules: rules}
}

type ruleHolder struct {
	resource.InstBaseRule
	resource.InstFileRule
	resource.InstStructRule
	resource.InstFuncRule
}

func loadRuleFile(path string) ([]resource.InstRule, error) {
	content, err := util.ReadFile(path)
	if err != nil {
		currentDir, _ := os.Getwd()
		err = errc.Adhere(err, "pwd", currentDir)
		return nil, err
	}
	return loadRuleRaw(content)
}

func loadRuleRaw(content string) ([]resource.InstRule, error) {
	var h []*ruleHolder
	err := json.Unmarshal([]byte(content), &h)
	if err != nil {
		return nil, errc.New(errc.ErrInvalidJSON, err.Error())
	}
	rules := make([]resource.InstRule, 0)
	for _, rule := range h {
		if rule.StructType != "" {
			r := &rule.InstStructRule
			r.InstBaseRule = rule.InstBaseRule
			rules = append(rules, r)
		} else if rule.Function != "" {
			r := &rule.InstFuncRule
			r.InstBaseRule = rule.InstBaseRule
			rules = append(rules, r)
		} else if rule.FileName != "" {
			r := &rule.InstFileRule
			r.InstBaseRule = rule.InstBaseRule
			rules = append(rules, r)
		} else {
			util.ShouldNotReachHereT("invalid rule type")
		}
	}
	return rules, nil
}

func loadDefaultRules() []resource.InstRule {
	rules, err := loadRuleRaw(pkg.ExportDefaultRuleJson())
	if err != nil {
		util.Log("Failed to load default rules: %v", err)
		return nil
	}
	return rules
}

func findAvailableRules() []resource.InstRule {
	util.GuaranteeInPreprocess()
	// Disable all instrumentation rules and rebuild the whole project to restore
	// all instrumentation actions, this also reverts the modification on Golang
	// runtime package.
	if config.GetConf().Restore {
		return nil
	}

	rules := make([]resource.InstRule, 0)

	// Load default rules unless explicitly disabled
	if !config.GetConf().IsDisableDefault() {
		defaultRules := loadDefaultRules()
		rules = append(rules, defaultRules...)
	}

	// If rule files are provided, load them
	if config.GetConf().RuleJsonFiles != "" {
		// Load multiple rule files
		if strings.Contains(config.GetConf().RuleJsonFiles, ",") {
			ruleFiles := strings.Split(config.GetConf().RuleJsonFiles, ",")
			for _, ruleFile := range ruleFiles {
				r, err := loadRuleFile(ruleFile)
				if err != nil {
					util.Log("Failed to load rules: %v", err)
					continue
				}
				rules = append(rules, r...)
			}
			return rules
		}
		// Load the one rule file
		rs, err := loadRuleFile(config.GetConf().RuleJsonFiles)
		if err != nil {
			util.Log("Failed to load rules: %v", err)
			return nil
		}
		rules = append(rules, rs...)
	}
	return rules
}

// match gives compilation arguments and finds out all interested rules
// for it.
func (rm *ruleMatcher) match(cmdArgs []string) *resource.RuleBundle {
	importPath := findFlagValue(cmdArgs, util.BuildPattern)
	util.Assert(importPath != "", "sanity check")
	util.Log("RunMatch: %v (%v)", importPath, cmdArgs)
	availables := make([]resource.InstRule, len(rm.availableRules[importPath]))

	// Okay, we are interested in these candidates, let's read it and match with
	// the instrumentation rule, but first we need to check if the package name
	// are already registered, to avoid futile effort
	copy(availables, rm.availableRules[importPath])
	if len(availables) == 0 {
		return nil // fast fail
	}
	parsedAst := make(map[string]*dst.File)
	bundle := resource.NewRuleBundle(importPath)

	goVersion := findFlagValue(cmdArgs, util.BuildGoVer)
	util.Assert(goVersion != "", "sanity check")
	util.Assert(strings.HasPrefix(goVersion, "go"), "sanity check")
	goVersion = strings.Replace(goVersion, "go", "v", 1)
	for _, candidate := range cmdArgs {
		// It's not a go file, ignore silently
		if !util.IsGoFile(candidate) {
			continue
		}
		file := candidate

		// If it's a vendor build, we need to extract the version of the module
		// from vendor/modules.txt, otherwise we find the version from source
		// code file path
		version := util.ExtractVersion(file)
		if rm.moduleVersions != nil {
			recorded := findVendorModuleVersion(rm.moduleVersions, importPath)
			if recorded != "" {
				version = recorded
			}
		}

		for i := len(availables) - 1; i >= 0; i-- {
			rule := availables[i]

			// Check if the version is supported
			matched, err := util.MatchVersion(version, rule.GetVersion())
			if err != nil {
				util.Log("Bad match: file %s, rule %s, version %s:\n%v",
					file, rule, version, err)
				continue
			}
			if !matched {
				continue
			}
			// Check if the rule requires a specific Go version(range)
			if rule.GetGoVersion() != "" {
				matched, err = util.MatchVersion(goVersion, rule.GetGoVersion())
				if err != nil {
					util.Log("Bad match: file %s, rule %s, go version %s:\n%v",
						file, rule, goVersion, err)
					continue
				}
				if !matched {
					continue
				}
			}

			// Check if it matches with file rule early as we try to avoid
			// parsing the file content, which is time consuming
			if _, ok := rule.(*resource.InstFileRule); ok {
				ast, err := util.ParseAstFromFileOnlyPackage(file)
				if ast == nil || err != nil {
					util.Log("Failed to parse %s: %v", file, err)
					continue
				}
				util.Log("Match file rule %s", rule)
				bundle.AddFileRule(rule.(*resource.InstFileRule))
				bundle.SetPackageName(ast.Name.Name)
				availables = append(availables[:i], availables[i+1:]...)
				continue
			}

			// Fair enough, parse the file content
			var tree *dst.File
			if _, ok := parsedAst[file]; !ok {
				fileAst, err := util.ParseAstFromFileFast(file)
				if fileAst == nil || err != nil {
					util.Log("failed to parse file %s: %v", file, err)
					continue
				}
				parsedAst[file] = fileAst
				util.Assert(fileAst.Name.Name != "", "empty package name")
				bundle.SetPackageName(fileAst.Name.Name)
				tree = fileAst
			} else {
				tree = parsedAst[file]
			}

			if tree == nil {
				// Failed to parse the file, stop here and log only
				// sicne it's a tolerant failure
				util.Log("Failed to parse file %s", file)
				continue
			}

			// Let's match with the rule precisely
			valid := false
			for _, decl := range tree.Decls {
				if genDecl, ok := decl.(*dst.GenDecl); ok {
					if rl, ok := rule.(*resource.InstStructRule); ok {
						if util.MatchStructDecl(genDecl, rl.StructType) {
							util.Log("Match struct rule %s", rule)
							err = bundle.AddFile2StructRule(file, rl)
							if err != nil {
								util.Log("Failed to add struct rule: %v", err)
								continue
							}
							valid = true
							break
						}
					}
				} else if funcDecl, ok := decl.(*dst.FuncDecl); ok {
					if rl, ok := rule.(*resource.InstFuncRule); ok {
						if util.MatchFuncDecl(funcDecl, rl.Function,
							rl.ReceiverType) {
							util.Log("Match func rule %s", rule)
							err = bundle.AddFile2FuncRule(file, rl)
							if err != nil {
								util.Log("Failed to add func rule: %v", err)
								continue
							}
							valid = true
							break
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
	return bundle
}

func findFlagValue(cmd []string, flag string) string {
	for i, v := range cmd {
		if v == flag {
			return cmd[i+1]
		}
	}
	return ""
}

// vendorModule represents a module in vendor/modules.txt file, it contains
// the module name, version and all submodules of the module, which looks like
//
// # golang.org/x/text v0.21.0
// ## explicit; go 1.18
// golang.org/x/text/secure/bidirule
// golang.org/x/text/transform
// golang.org/x/text/unicode/bidi
// golang.org/x/text/unicode/norm
// # golang.org/x/time v0.5.0
// ## explicit; go 1.18
// golang.org/x/time/rate
//
// The module name is the first line of the module, the version is the second
// part of the first line, and all submodules are listed in the following lines
// starting with the module name.
type vendorModule struct {
	module     string
	version    string
	submodules []string
}

func findVendorModuleVersion(modules []*vendorModule, importPath string) string {
	for _, module := range modules {
		if module.module == importPath {
			return module.version
		}
		for _, submodule := range module.submodules {
			if submodule == importPath {
				return module.version
			}
		}
	}
	return ""
}

func parseVendorModules(projDir string) ([]*vendorModule, error) {
	vendorFile := filepath.Join(projDir, "vendor", "modules.txt")
	if util.PathNotExists(vendorFile) {
		return nil, errc.New(errc.ErrNotExist, "vendor/modules.txt not found")
	}
	// Read the vendor/modules.txt file line by line and parse it in form of
	// #ImportPath Version, all lines following this line are considered as
	// submodules of the same module
	file, err := os.Open(vendorFile)
	if err != nil {
		return nil, errc.New(errc.ErrOpenFile, err.Error())
	}
	defer func(dryRunLog *os.File) {
		err := dryRunLog.Close()
		if err != nil {
			util.Log("Failed to close dry run log file: %v", err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	// 10MB should be enough to accommodate most long line
	buffer := make([]byte, 0, 10*1024*1024)
	scanner.Buffer(buffer, cap(buffer))

	vms := make([]*vendorModule, 0)
	var vm *vendorModule

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "# ") {
			parts := strings.Split(line, " ")
			// In form of #ImportPath Version?
			if len(parts) == 3 {
				util.Assert(parts[0] == "#", "sanity check")
				util.Assert(strings.HasPrefix(parts[2], "v"), "sanity check")
				vm = &vendorModule{
					module:     parts[1],
					version:    parts[2],
					submodules: make([]string, 0),
				}
				vms = append(vms, vm)
			} else {
				// For other lines, we just ignore them
			}
		} else if !strings.HasPrefix(line, "## ") {
			vm.submodules = append(vm.submodules, line)
		} else {
			util.Assert(strings.HasPrefix(line, "## "), "why not otherwise")
		}
	}
	err = scanner.Err()
	if err != nil {
		return nil, errc.New(errc.ErrParseCode,
			"cannot parse vendor/modules.txt")
	}
	return vms, nil
}

func runMatch(matcher *ruleMatcher, cmd string, ch chan *resource.RuleBundle) {
	bundle := matcher.match(util.SplitCmds(cmd))
	ch <- bundle
}

func (dp *DepProcessor) matchRules(compileCmds []string) error {
	defer util.PhaseTimer("Match")()
	matcher := newRuleMatcher()

	// If we are in vendor mode, we need to parse the vendor/modules.txt file
	// to get the version of each module for future matching
	if dp.vendorBuild {
		modules, err := parseVendorModules(dp.getGoModDir())
		if err != nil {
			return err
		}
		if config.GetConf().Verbose {
			util.Log("Vendor modules: %v", modules)
		}
		matcher.moduleVersions = modules
	}

	// Find used instrumentation rule according to compile commands
	ch := make(chan *resource.RuleBundle)
	for _, cmd := range compileCmds {
		go runMatch(matcher, cmd, ch)
	}
	cnt := 0
	for cnt < len(compileCmds) {
		bundle := <-ch
		if bundle.IsValid() {
			dp.bundles = append(dp.bundles, bundle)
		}
		cnt++
	}
	return nil
}

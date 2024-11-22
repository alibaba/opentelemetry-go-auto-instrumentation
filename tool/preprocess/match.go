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
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/config"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/resource"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/shared"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/util"
	"github.com/dave/dst"
)

type ruleMatcher struct {
	availableRules map[string][]resource.InstRule
}

func newRuleMatcher() *ruleMatcher {
	rules := make(map[string][]resource.InstRule)
	for _, rule := range findAvailableRules() {
		rules[rule.GetImportPath()] = append(rules[rule.GetImportPath()], rule)
	}
	if config.GetConf().Verbose {
		log.Printf("Available rules: %v", rules)
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
		return nil, fmt.Errorf("failed to read rule file: %w %v",
			err, currentDir)
	}
	return loadRuleRaw(content)
}

func loadRuleRaw(content string) ([]resource.InstRule, error) {
	var h []*ruleHolder
	err := json.Unmarshal([]byte(content), &h)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal rules: %w", err)
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
		log.Printf("Failed to load default rules: %v", err)
		return nil
	}
	return rules
}

func findAvailableRules() []resource.InstRule {
	shared.GuaranteeInPreprocess()
	// Disable all instrumentation rules and rebuild the whole project to restore
	// all instrumentation actions, this also reverts the modification on Golang
	// runtime package.
	if config.GetConf().Restore {
		return nil
	}

	// If rule file is not set, we will use the default rules
	if config.GetConf().RuleJsonFiles == "" {
		return loadDefaultRules()
	}

	rules := make([]resource.InstRule, 0)

	// Load default rules unless explicitly disabled
	if !config.GetConf().IsDisableDefaultRules() {
		defaultRules := loadDefaultRules()
		rules = append(rules, defaultRules...)
	}

	// Load multiple rule files if provided
	if strings.Contains(config.GetConf().RuleJsonFiles, ",") {
		ruleFiles := strings.Split(config.GetConf().RuleJsonFiles, ",")
		for _, ruleFile := range ruleFiles {
			r, err := loadRuleFile(ruleFile)
			if err != nil {
				log.Printf("Failed to load rules: %v", err)
				continue
			}
			rules = append(rules, r...)
		}
		return rules
	}

	// Load the one rule file if provided
	rs, err := loadRuleFile(config.GetConf().RuleJsonFiles)
	if err != nil {
		log.Printf("Failed to load rules: %v", err)
		return nil
	}
	rules = append(rules, rs...)
	return rules
}

// match gives compilation arguments and finds out all interested rules
// for it.
func (rm *ruleMatcher) match(importPath string,
	candidates []string) *resource.RuleBundle {
	util.Assert(importPath != "", "sanity check")
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
	for _, candidate := range candidates {
		// It's not a go file, ignore silently
		if !shared.IsGoFile(candidate) {
			continue
		}
		file := candidate
		version := shared.ExtractVersion(file)

		for i := len(availables) - 1; i >= 0; i-- {
			rule := availables[i]

			// Check if the version is supported
			matched, err := shared.MatchVersion(version, rule.GetVersion())
			if err != nil {
				log.Printf("Failed to match version %v between %v and %v",
					err, file, rule)
				continue
			}
			if !matched {
				continue
			}
			// Check if it matches with file rule early as we try to avoid
			// parsing the file content, which is time consuming
			if _, ok := rule.(*resource.InstFileRule); ok {
				ast, err := shared.ParseAstFromFileOnlyPackage(file)
				if ast == nil || err != nil {
					log.Printf("Failed to parse %s: %v", file, err)
					continue
				}
				log.Printf("Match file rule %s", rule)
				bundle.AddFileRule(rule.(*resource.InstFileRule))
				bundle.SetPackageName(ast.Name.Name)
				availables = append(availables[:i], availables[i+1:]...)
				continue
			}

			// Fair enough, parse the file content
			var tree *dst.File
			if _, ok := parsedAst[file]; !ok {
				fileAst, err := shared.ParseAstFromFileFast(file)
				if fileAst == nil || err != nil {
					log.Printf("failed to parse file %s: %v", file, err)
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
				log.Printf("Failed to parse file %s", file)
				continue
			}

			// Let's match with the rule precisely
			valid := false
			for _, decl := range tree.Decls {
				if genDecl, ok := decl.(*dst.GenDecl); ok {
					if rl, ok := rule.(*resource.InstStructRule); ok {
						if shared.MatchStructDecl(genDecl, rl.StructType) {
							log.Printf("Match struct rule %s", rule)
							bundle.AddFile2StructRule(file, rl)
							valid = true
							break
						}
					}
				} else if funcDecl, ok := decl.(*dst.FuncDecl); ok {
					if rl, ok := rule.(*resource.InstFuncRule); ok {
						if shared.MatchFuncDecl(funcDecl, rl.Function,
							rl.ReceiverType) {
							log.Printf("Match func rule %s", rule)
							bundle.AddFile2FuncRule(file, rl)
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

func readImportPath(cmd []string) string {
	var pkg string
	for i, v := range cmd {
		if v == "-p" {
			return cmd[i+1]
		}
	}
	return pkg
}

func runMatch(matcher *ruleMatcher, cmd string, ch chan *resource.RuleBundle) {
	cmdArgs := shared.SplitCmds(cmd)
	importPath := readImportPath(cmdArgs)
	util.Assert(importPath != "", "sanity check")
	if config.GetConf().Verbose {
		log.Printf("Matching %v with %v\n", importPath, cmdArgs)
	}
	bundle := matcher.match(importPath, cmdArgs)
	ch <- bundle
}

func (dp *DepProcessor) matchRules(compileCmds []string) error {
	defer util.PhaseTimer("Match")()
	matcher := newRuleMatcher()
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

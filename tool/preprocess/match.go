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
	"log"
	"strings"

	"github.com/dave/dst"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/api"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/resource"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/shared"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/util"
)

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

type ruleMatcher struct {
	availableRules map[string][]api.InstRule
}

func newRuleMatcher() *ruleMatcher {
	rules := make(map[string][]api.InstRule)
	for _, rule := range findAvailableRules() {
		rules[rule.GetImportPath()] = append(rules[rule.GetImportPath()], rule)
	}
	return &ruleMatcher{availableRules: rules}
}

// matchRuleBundle gives compilation arguments and finds out all interested rules
// for it. N.B. This is performance critical, be careful to modify it.
func (rm *ruleMatcher) matchRuleBundle(importPath string,
	candidates []string) *resource.RuleBundle {
	util.Assert(importPath != "", "sanity check")
	availables := make([]api.InstRule, len(rm.availableRules[importPath]))

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
				log.Printf("Failed to match version %v", err)
				continue
			}
			if !matched {
				continue
			}
			// Check if it matches with file rule early as we try to avoid
			// parsing the file content, which is time consuming
			if _, ok := rule.(*api.InstFileRule); ok {
				log.Printf("Match file rule %s", rule)
				bundle.AddFileRule(rule.(*api.InstFileRule))
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
				if bundle.PackageName == "" {
					bundle.PackageName = fileAst.Name.Name
				}
				util.Assert(bundle.PackageName == fileAst.Name.Name,
					"inconsistent package name")
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
					if rl, ok := rule.(*api.InstStructRule); ok {
						if shared.MatchStructDecl(genDecl, rl.StructType) {
							log.Printf("Match struct rule %s", rule)
							bundle.AddFile2StructRule(file, rl)
							valid = true
							break
						}
					}
				} else if funcDecl, ok := decl.(*dst.FuncDecl); ok {
					if rl, ok := rule.(*api.InstFuncRule); ok {
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

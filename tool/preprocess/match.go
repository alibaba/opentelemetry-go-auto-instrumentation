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
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/dave/dst"
	"golang.org/x/mod/semver"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/api"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/resource"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/shared"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/util"
)

// splitVersionRange splits the version range into two parts, start and end.
func splitVersionRange(vr string) (string, string) {
	util.Assert(strings.Contains(vr, ","), "invalid version range format")
	util.Assert(strings.Contains(vr, "["), "invalid version range format")
	util.Assert(strings.Contains(vr, ")"), "invalid version range format")

	start := vr[1:strings.Index(vr, ",")]
	end := vr[strings.Index(vr, ",")+1 : len(vr)-1]
	return "v" + start, "v" + end
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
		return false, fmt.Errorf("invalid version format in rule %v",
			ruleVersion)
	}
	// Remove extra whitespace from the rule version string
	ruleVersion = strings.ReplaceAll(ruleVersion, " ", "")

	// Compare the version with the rule version, the rule version is in the
	// format [start, end), where start is inclusive and end is exclusive
	ruleVersionStart, ruleVersionEnd := splitVersionRange(ruleVersion)
	if semver.Compare(version, ruleVersionStart) >= 0 &&
		semver.Compare(version, ruleVersionEnd) < 0 {
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
// for it.
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
	bundle := resource.NewRuleBundle(importPath)
	for _, candidate := range candidates {
		// It's not a go file, ignore silently
		if !shared.IsGoFile(candidate) {
			continue
		}

		// Parse the file content
		file := candidate
		fileAst, err := shared.ParseAstFromFile(file)
		if fileAst == nil {
			// Failed to parse the file, stop here and log only
			// sicne it's a tolerant failure
			log.Printf("Failed to parse file %s from local fs: %v", file, err)
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
				log.Printf("Match file rule %s", rule)
				bundle.AddFileRule(rl)
				availables = append(availables[:i], availables[i+1:]...)
			} else {
				valid := false
				for _, decl := range fileAst.Decls {
					if genDecl, ok := decl.(*dst.GenDecl); ok {
						if rl, ok := rule.(*api.InstStructRule); ok {
							if shared.MatchStructDecl(genDecl, rl.StructType) {
								log.Printf("Match struct rule %s", rule)
								bundle.AddFile2StructRule(file, rl)
								valid = true
							}
						}
					} else if funcDecl, ok := decl.(*dst.FuncDecl); ok {
						if rl, ok := rule.(*api.InstFuncRule); ok {
							if shared.MatchFuncDecl(funcDecl, rl.Function,
								rl.ReceiverType) {
								log.Printf("Match func rule %s", rule)
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

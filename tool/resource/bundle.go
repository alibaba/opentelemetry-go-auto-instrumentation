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

package resource

import (
	"encoding/json"
	"fmt"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/shared"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/util"
	"github.com/dave/dst"
)

const (
	RuleBundleJsonFile = "rule_bundle.json"
)

// RuleBundle is a collection of rules that matched with one compilation action
type RuleBundle struct {
	PackageName      string
	ImportPath       string
	FileRules        []*InstFileRule
	File2FuncRules   map[string]map[string][]*InstFuncRule
	File2StructRules map[string]map[string][]*InstStructRule
}

func NewRuleBundle(importPath string) *RuleBundle {
	return &RuleBundle{
		PackageName:      "",
		ImportPath:       importPath,
		FileRules:        make([]*InstFileRule, 0),
		File2FuncRules:   make(map[string]map[string][]*InstFuncRule),
		File2StructRules: make(map[string]map[string][]*InstStructRule),
	}
}

func (rb *RuleBundle) String() string {
	bs, _ := json.Marshal(rb)
	return string(bs)
}

func (rb *RuleBundle) IsValid() bool {
	return rb != nil &&
		(len(rb.FileRules) > 0 ||
			len(rb.File2FuncRules) > 0 ||
			len(rb.File2StructRules) > 0)
}

func (rb *RuleBundle) AddFile2FuncRule(file string, rule *InstFuncRule) {
	fn := rule.Function + "," + rule.ReceiverType
	util.Assert(fn != "", "sanity check")
	if _, exist := rb.File2FuncRules[file]; !exist {
		rb.File2FuncRules[file] = make(map[string][]*InstFuncRule)
		rb.File2FuncRules[file][fn] = []*InstFuncRule{rule}
	} else {
		rb.File2FuncRules[file][fn] =
			append(rb.File2FuncRules[file][fn], rule)
	}
}

func (rb *RuleBundle) AddFile2StructRule(file string, rule *InstStructRule) {
	st := rule.StructType
	util.Assert(st != "", "sanity check")
	if _, exist := rb.File2StructRules[file]; !exist {
		rb.File2StructRules[file] = make(map[string][]*InstStructRule)
		rb.File2StructRules[file][st] = []*InstStructRule{rule}
	} else {
		rb.File2StructRules[file][st] =
			append(rb.File2StructRules[file][st], rule)
	}
}

func (rb *RuleBundle) SetPackageName(name string) {
	rb.PackageName = name
}

func (rb *RuleBundle) AddFileRule(rule *InstFileRule) {
	rb.FileRules = append(rb.FileRules, rule)
}

func isHookDefined(root *dst.File, rule *InstFuncRule) bool {
	util.Assert(rule.OnEnter != "" || rule.OnExit != "", "hook must be set")
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

func FindHookFile(rule *InstFuncRule) (string, error) {
	files, err := FindRuleFiles(rule)
	if err != nil {
		return "", fmt.Errorf("failed to find rule files: %w", err)
	}
	for _, file := range files {
		if !shared.IsGoFile(file) {
			continue
		}
		root, err := shared.ParseAstFromFileFast(file)
		if err != nil {
			return "", fmt.Errorf("failed to read hook file: %w", err)
		}
		if isHookDefined(root, rule) {
			return file, nil
		}
	}
	return "", nil
}

func FindRuleFiles(rule InstRule) ([]string, error) {
	files, err := util.ListFilesFlat(rule.GetPath())
	if err != nil {
		return nil, err
	}
	switch rule.(type) {
	case *InstFuncRule, *InstFileRule:
		return files, nil
	case *InstStructRule:
		util.ShouldNotReachHereT("insane rule type")
	}
	return nil, nil
}

func StoreRuleBundles(bundles []*RuleBundle) error {
	shared.GuaranteeInPreprocess()
	ruleFile := shared.GetPreprocessLogPath(RuleBundleJsonFile)
	bs, err := json.Marshal(bundles)
	if err != nil {
		return fmt.Errorf("failed to store used rules: %w", err)
	}
	_, err = util.WriteFile(ruleFile, string(bs))
	if err != nil {
		return fmt.Errorf("failed to write used rules: %w", err)
	}
	return nil
}

func LoadRuleBundles() ([]*RuleBundle, error) {
	shared.GuaranteeInInstrument()

	ruleFile := shared.GetPreprocessLogPath(RuleBundleJsonFile)
	data, err := util.ReadFile(ruleFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read used rules: %w", err)
	}
	var bundles []*RuleBundle
	err = json.Unmarshal([]byte(data), &bundles)
	if err != nil {
		return nil, fmt.Errorf("failed to load used rules: %w", err)
	}
	return bundles, nil
}

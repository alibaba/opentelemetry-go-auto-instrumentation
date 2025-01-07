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

package instrument

import (
	"fmt"
	"path/filepath"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/resource"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/shared"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/util"
	"github.com/dave/dst"
)

func addStructField(rule *resource.InstStructRule, decl dst.Decl) {
	util.Assert(rule.FieldName != "" && rule.FieldType != "",
		"rule must have field and type")
	util.Log("Apply struct rule %v", rule)
	shared.AddStructField(decl, rule.FieldName, rule.FieldType)
}

func (rp *RuleProcessor) applyStructRules(bundle *resource.RuleBundle) error {
	for file, struct2Rules := range bundle.File2StructRules {
		util.Assert(filepath.IsAbs(file), "file path must be absolute")
		// Apply struct rules to the file
		astRoot, err := rp.loadAst(file)
		if err != nil {
			return fmt.Errorf("failed to load ast from file: %w", err)
		}
		for _, decl := range astRoot.Decls {
			for structName, rules := range struct2Rules {
				if shared.MatchStructDecl(decl, structName) {
					for _, rule := range rules {
						addStructField(rule, decl)
					}
				}
			}
		}
		// Once all struct rules are applied, we restore AST to file and use it
		// in future compilation
		newFile, err := rp.restoreAst(file, astRoot)
		if err != nil {
			return fmt.Errorf("failed to restore ast: %w", err)
		}
		rp.saveDebugFile(newFile)
	}
	return nil
}

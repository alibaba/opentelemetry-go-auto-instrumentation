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
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/util"
	"github.com/dave/dst"
)

func (rp *RuleProcessor) addStructField(rule *resource.InstStructRule, decl dst.Decl) {
	util.Assert(rule.FieldName != "" && rule.FieldType != "",
		"rule must have field and type")
	// Find last field in the struct and record its position
	stType := decl.(*dst.GenDecl).Specs[0].(*dst.TypeSpec).Type.(*dst.StructType)
	fields := stType.Fields.List
	// Find the last field containing the line number in the struct
	for i := len(fields) - 1; i >= 0; i-- {
		field := fields[i]
		fieldPos := rp.parser.FindPosition(field)
		if fieldPos.IsValid() {
			tag := fmt.Sprintf("//line %s", fieldPos.String())
			field.Decs.Before = dst.EmptyLine
			field.Decs.Start.Append(tag)
			util.Log("== Last field %v", tag)
			break
		}
	}
	util.Log("Apply struct rule %v", rule)
	util.AddStructField(decl, rule.FieldName, rule.FieldType)
}

func (rp *RuleProcessor) applyStructRules(bundle *resource.RuleBundle) error {
	for file, struct2Rules := range bundle.File2StructRules {
		util.Assert(filepath.IsAbs(file), "file path must be absolute")
		// Apply struct rules to the file
		astRoot, err := rp.loadAst(file)
		if err != nil {
			return err
		}
		for _, decl := range astRoot.Decls {
			for structName, rules := range struct2Rules {
				if util.MatchStructDecl(decl, structName) {
					for _, rule := range rules {
						rp.addStructField(rule, decl)
					}
				}
			}
		}
		// Once all struct rules are applied, we restore AST to file and use it
		// in future compilation
		newFile, err := rp.restoreAst(file, astRoot)
		if err != nil {
			return err
		}
		// Line directive must be placed at the beginning of the line, otherwise
		// it will be ignored by the compiler
		err = rp.enableLineDirective(newFile)
		if err != nil {
			return err
		}
		rp.saveDebugFile(newFile)
	}
	return nil
}

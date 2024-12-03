// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package instrument

import (
	"fmt"
	"log"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/resource"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/shared"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/util"
	"github.com/dave/dst"
)

func addStructField(rule *resource.InstStructRule, decl dst.Decl) {
	util.Assert(rule.FieldName != "" && rule.FieldType != "",
		"rule must have field and type")
	log.Printf("Apply struct rule %v", rule)
	shared.AddStructField(decl, rule.FieldName, rule.FieldType)
}

func (rp *RuleProcessor) applyStructRules(bundle *resource.RuleBundle) error {
	for file, struct2Rules := range bundle.File2StructRules {
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

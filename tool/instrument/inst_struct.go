package instrument

import (
	"fmt"
	"log"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/resource"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/shared"
)

func (rp *RuleProcessor) applyStructRules(bundle *resource.RuleBundle) error {
	for file, struct2Rules := range bundle.File2StructRules {
		// Apply struct rules to the file
		astRoot, err := rp.loadAst(file)
		if err != nil {
			return fmt.Errorf("failed to load ast from file: %w", err)
		}
		for _, decl := range astRoot.Decls {
			for structName, rules := range struct2Rules {
				if resource.MatchStructDecl(decl, structName) {
					for _, ruleHash := range rules {
						rule := resource.FindStructRuleByHash(ruleHash)
						if rule.FieldName == "" || rule.FieldType == "" {
							return fmt.Errorf("rule must have field and type")
						}
						log.Printf("Apply struct rule %v", rule)
						shared.AddStructField(decl, rule.FieldName, rule.FieldType)
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
		shared.SaveDebugFile("struct_", newFile)
	}
	return nil
}

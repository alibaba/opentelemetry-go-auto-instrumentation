package instrument

import (
	"fmt"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/shared"
	"log"
	"path/filepath"
	"strings"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/resource"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/util"
)

func (rp *RuleProcessor) applyFileRules(bundle *resource.RuleBundle) (err error) {
	for _, ruleHash := range bundle.FileRules {
		rule := resource.FindFileRuleByHash(ruleHash)
		if rule.FileName == "" {
			return fmt.Errorf("file rule must have a file name")
		}
		// Decorate the source code to remove //go:build exclude
		// and rename package name
		source, err := util.ReadFile(rule.FileName)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", rule.FileName, err)
		}
		source = shared.RemoveGoBuildComment(source)
		source = shared.RenamePackage(source, bundle.PackageName)

		// Get last section of file path as file name
		fileName := filepath.Base(rule.FileName)
		target := filepath.Join(rp.workDir, fmt.Sprintf("otel_inst_file_%s", fileName))
		_, err = util.WriteStringToFile(target, source)
		if err != nil {
			return fmt.Errorf("failed to write extra file %s: %w", target, err)
		}
		// Relocate the file dependency of the rule, any rules targeting the
		// file dependency specified by the rule should be updated to target the
		// new file
		rp.setRelocated(rule.FileName, target)

		// Append or replace the file to the compile arguments
		mode := "REPLACE"
		if rule.Replace {
			err = rp.replaceCompileArg(target, func(arg string) bool {
				return strings.HasSuffix(arg, fileName)
			})
			if err != nil {
				return fmt.Errorf("failed to replace %v %w", fileName, err)
			}
		} else {
			mode = "APPEND"
			rp.addCompileArg(target)
		}
		log.Printf("Apply file rule %v by %s mode", fileName, mode)
		shared.SaveDebugFile("file_", target)
	}
	return nil
}

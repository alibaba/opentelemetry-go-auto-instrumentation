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
	"strings"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/resource"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/shared"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/util"
)

func (rp *RuleProcessor) applyFileRules(bundle *resource.RuleBundle) (err error) {
	for _, rule := range bundle.FileRules {
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

		// Get last section of file path as file name
		fileName := filepath.Base(rule.FileName)
		target := filepath.Join(rp.workDir,
			fmt.Sprintf("otel_inst_file_%s", fileName))
		_, err = util.WriteFile(target, source)
		if err != nil {
			return fmt.Errorf("failed to write extra file %s: %w", target, err)
		}
		// Relocate the file dependency of the rule, any rules targeting the
		// file dependency specified by the rule should be updated to target the
		// new file
		rp.setRelocated(rule.FileName, target)

		// Append or replace the file to the compile arguments
		if rule.Replace {
			err = rp.replaceCompileArg(target, func(arg string) bool {
				return strings.HasSuffix(arg, fileName)
			})
			if err != nil {
				return fmt.Errorf("failed to replace %v %w", fileName, err)
			}
		} else {
			rp.addCompileArg(target)
		}
		util.Log("Apply file rule %v", rule)
		rp.saveDebugFile(target)
	}
	return nil
}

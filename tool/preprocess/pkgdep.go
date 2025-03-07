// Copyright (c) 2024 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package preprocess

import (
	"strings"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/resource"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/util"
)

const (
	PkgDep     = "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg"
	OtelPkgDep = "otel_pkgdep"
)

func replaceImport(importPath string, code string) string {
	code = strings.ReplaceAll(code, PkgDep, importPath)
	return code
}

func (dp *DepProcessor) replaceOtelImports() error {
	moduleName, err := dp.getImportPathOf(OtelPkgDep)
	if err != nil {
		return err
	}
	// Replace imports in generated files
	generated := []string{dp.generatedOf(OtelRules), dp.generatedOf(OtelPkgDep)}
	generated = append(generated, dp.sources...)
	for _, dep := range generated {
		files, err := util.ListFiles(dep)
		if err != nil {
			return err
		}
		for _, file := range files {
			// Skip non-go files as no imports within them
			if !util.IsGoFile(file) {
				continue
			}
			// Read file content and replace content then write back
			content, err := util.ReadFile(file)
			if err != nil {
				return err
			}
			content = replaceImport(moduleName, content)
			_, err = util.WriteFile(file, content)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (dp *DepProcessor) copyPkgDep() error {
	return resource.CopyPkgTo(dp.generatedOf(OtelPkgDep))
}

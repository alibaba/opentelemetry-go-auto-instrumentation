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
	"fmt"
	"log"
	"strings"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/config"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/resource"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/shared"
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
		return fmt.Errorf("failed to get import path of otel_pkg: %w", err)
	}

	for _, dep := range []string{OtelRules, OtelPkgDep} {
		files, err := util.ListFiles(dep)
		if err != nil {
			return fmt.Errorf("failed to list files: %w", err)
		}
		for _, file := range files {
			// Skip non-go files as no imports within them
			if !shared.IsGoFile(file) {
				continue
			}
			// Read file content and replace content then write back
			content, err := util.ReadFile(file)
			if err != nil {
				return fmt.Errorf("failed to read file content: %w", err)
			}
			if config.GetConf().Verbose {
				log.Printf("Replace import path of %s to %s", file, moduleName)
			}

			content = replaceImport(moduleName, content)
			_, err = util.WriteFile(file, content)
			if err != nil {
				return fmt.Errorf("failed to write file content: %w", err)
			}
		}
	}
	return nil
}

func (dp *DepProcessor) copyPkgDep() error {
	err := resource.CopyPkgTo(OtelPkgDep)
	if err != nil {
		return fmt.Errorf("failed to copy pkg deps: %w", err)
	}
	return nil
}

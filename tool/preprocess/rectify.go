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
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/config"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/errc"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/util"
)

const (
	pkgPrefix = "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg"
)

type moduleInfo struct {
	Path  string `json:"Path"`
	Error string `json:"Error"`
	Dir   string `json:"Dir"`
}

func (dp *DepProcessor) findModCacheDir() (string, error) {
	if config.BuildPath != "" &&
		util.PathExists(config.BuildPath) &&
		config.GetConf().PkgModule == "" {
		// In development mode, there is a high probability that we may have
		// added new rules. In such cases, we prefer to use the pkg directory
		// with the newly added rules rather than the remote pkg module. Our
		// approach is to embed the source code directory into the tool during
		// the packaging process. The tool will then check whether this directory
		// exists. If it does exist, the tool will use it. Otherwise, we will
		// fall back to using the remote pkg module.
		return config.BuildPath, nil
	}
	pkgVersion := config.UsedPkg
	if config.GetConf().PkgModule != "" {
		// If the user has specified a custom pkg module by using otel set -pkgmodule
		// we should use it instead of the default one, it has the highest priority
		pkgVersion = config.GetConf().PkgModule
	}
	modulePath := pkgPrefix + "@" + pkgVersion
	output, err := runCmdCombinedOutput(dp.getGoModDir(),
		"go", "mod", "download", "-json", modulePath)
	if err != nil {
		return "", err
	}
	var moduleInfo moduleInfo
	if err := json.Unmarshal([]byte(output), &moduleInfo); err != nil {
		return "", errc.New(errc.ErrPreprocess, "failed to unmarshal module info")
	}
	if moduleInfo.Error != "" {
		return "", errc.New(errc.ErrPreprocess,
			fmt.Sprintf("error downloading module: %s", moduleInfo.Error))
	}
	return moduleInfo.Dir, nil
}

// rectifyRule rectifies the file rules path to the local module cache path.
func (dp *DepProcessor) rectifyRule() error {
	util.GuaranteeInPreprocess()
	defer util.PhaseTimer("Fetch")()
	rectified := map[string]bool{}
	for _, bundle := range dp.bundles {
		for _, funcRules := range bundle.File2FuncRules {
			for _, rs := range funcRules {
				for _, rule := range rs {
					if rule.UseRaw {
						continue
					}
					if rectified[rule.GetPath()] {
						continue
					}
					p := strings.TrimPrefix(rule.Path, pkgPrefix)
					p = filepath.Join(dp.pkgLocalCache, p)
					rule.SetPath(p)
					rectified[p] = true
				}
			}
		}
		for _, fileRule := range bundle.FileRules {
			if rectified[fileRule.GetPath()] {
				continue
			}
			p := strings.TrimPrefix(fileRule.Path, pkgPrefix)
			p = filepath.Join(dp.pkgLocalCache, p)
			fileRule.SetPath(p)
			fileRule.FileName = filepath.Join(p, fileRule.FileName)
			rectified[p] = true
		}
	}
	return nil
}

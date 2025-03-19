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

func (dp *DepProcessor) findModCacheDir(modulePath string) (string, error) {
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
	modCacheDir, err := dp.findModCacheDir("github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg@2da1f02")
	if err != nil {
		return err
	}
	if modCacheDir == "" {
		return errc.New(errc.ErrPreprocess, "cannot find rule cache dir")
	}
	util.Log("Local module cache: %s", modCacheDir)
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
					p = filepath.Join(modCacheDir, p)
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
			p = filepath.Join(modCacheDir, p)
			fileRule.SetPath(p)
			fileRule.FileName = filepath.Join(p, fileRule.FileName)
			rectified[p] = true
		}
	}
	return nil
}

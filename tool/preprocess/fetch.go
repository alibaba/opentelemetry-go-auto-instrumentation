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
	"path/filepath"
	"strings"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/errc"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/util"
)

const (
	pkgPrefix = "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg"
)

// fetchRules fetches the rules via the network
func (dp *DepProcessor) fetchRules(compileCmds []string) error {
	util.GuaranteeInPreprocess()
	defer util.PhaseTimer("Fetch")()
	ruleCacheDir := ""
	for _, compileCmd := range compileCmds {
		if !util.IsCompileCommand(compileCmd) {
			continue
		}
		args := util.SplitCmds(compileCmd)
		for i, arg := range args {
			if arg == "-p" &&
				strings.Contains(args[i+1], pkgPrefix) {
				for j := i + 1; j < len(args); j++ {
					if util.IsGoFile(args[j]) {
						dir := filepath.Dir(args[j])
						ruleCacheDir = dir
						break
					}
				}
				p := ruleCacheDir
				for {
					if strings.HasSuffix(p, "/pkg") {
						break
					}
					if p == "/" {
						break
					}
					p = filepath.Dir(p)
				}
				ruleCacheDir = p
				break
			}
		}
	}
	if ruleCacheDir == "" {
		return errc.New(errc.ErrPreprocess, "cannot find rule cache dir")
	}
	util.Log("rule cache dir: %s", ruleCacheDir)
	for _, bundle := range dp.bundles {
		for _, funcRules := range bundle.File2FuncRules {
			for _, rs := range funcRules {
				for _, rule := range rs {
					if rule.UseRaw {
						continue
					}
					p := strings.ReplaceAll(rule.Path, pkgPrefix, ruleCacheDir)
					rule.SetPath(p)
				}
			}
		}
		for _, fileRule := range bundle.FileRules {
			p := strings.ReplaceAll(fileRule.Path, pkgPrefix, ruleCacheDir)
			fileName := filepath.Join(p, fileRule.FileName)
			fileRule.SetPath(p)
			fileRule.FileName = fileName
		}
	}
	return nil
}

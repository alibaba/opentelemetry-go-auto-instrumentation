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
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/shared"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/util"
)

func runModDownload(path string) (string, error) {
	return util.RunCmdOutput("go", "mod", "download", "-json", path)
}

type moduleHolder struct {
	Error string // error loading module
	Dir   string // absolute path to cached source root directory
}

func fetchNetwork(path string) (string, error) {
	text, err := runModDownload(path)
	if err != nil {
		return "", fmt.Errorf("failed to download rule: %w %v %v",
			err, text, path)
	}
	var mod *moduleHolder
	err = json.Unmarshal([]byte(text), &mod)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal rule: %w %v",
			err, path)
	}
	if mod.Error != "" {
		return "", fmt.Errorf("failed to download rule: %s %v",
			mod.Error, path)
	}
	exist, err := util.PathExists(mod.Dir)
	if err != nil {
		return "", fmt.Errorf("failed to check path: %w %v",
			err, mod.Dir)
	}
	if !exist {
		return "", fmt.Errorf("failed to find module path: %s",
			mod.Dir)
	}
	return mod.Dir, nil
}

func isStdRulePath(path string) bool {
	return strings.HasPrefix(path, StdRulesPath)
}

func (dp *DepProcessor) fetchEmbed(path string) (string, error) {
	// Mangle the rule path to the local file system, e.g.
	// github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/foo@v1
	// => rules/foo
	path = strings.Replace(path, StdRulesPrefix, "", 1)
	if strings.Contains(path, "@") {
		path = strings.Split(path, "@")[0]
	}
	// Walk through the rule cache and copy all rule files we matched to the
	// local file system
	walkFn := func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		data, err := dp.ruleCache.ReadFile(p)
		if err != nil {
			return err
		}
		target := shared.GetPreprocessLogPath(filepath.Join(OtelRuleCache, p))
		err = os.MkdirAll(filepath.Dir(target), 0777)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
		_, err = util.WriteFile(target, string(data))
		if err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}
		if shared.Verbose {
			log.Printf("Copy embed file %v to %v", p, target)
		}
		return nil
	}

	err := fs.WalkDir(dp.ruleCache, path, walkFn)
	if err != nil {
		return "", fmt.Errorf("failed to walk dir: %w", err)
	}
	// Now all rule files are copied to the local file system, we can return
	// the path to corresponding local file system
	dir := shared.GetPreprocessLogPath(filepath.Join(OtelRuleCache, path))
	return dir, nil
}

func (dp *DepProcessor) fetchFrom(path string) (string, error) {
	// Path to local file system, use local directory directly
	exist, err := util.PathExists(path)
	if err != nil {
		return "", fmt.Errorf("failed to check path: %w", err)
	}
	if exist {
		log.Printf("Fetch %s from local file system", path)
		return path, nil
	}
	// Path to network
	if shared.IsModPath(path) {
		// If the path points to the network but is an officially provided
		// module, then we can retrieve it from the embed cache instead of
		// downloading it from the network.
		if isStdRulePath(path) {
			dir, err := dp.fetchEmbed(path)
			if err != nil {
				return "", fmt.Errorf("failed to fetch embed: %w", err)
			}
			log.Printf("Fetch %s from embed cache", path)
			return dir, nil
		}

		// Download the module to the local file system
		dir, err := fetchNetwork(path)
		if err != nil {
			return "", fmt.Errorf("failed to download module: %w", err)
		}
		// Get path to the local module cache
		log.Printf("Fetch %s from network %s", path, dir)
		return dir, nil
	}

	// Best effort to find the path but not found, give up
	return "", fmt.Errorf("can not fetch path %s", path)
}

// fetchRules fetches the rules via the network
func (dp *DepProcessor) fetchRules() error {
	shared.GuaranteeInPreprocess()
	defer util.PhaseTimer("Fetch")()
	// Different rules may share the same path, we dont want to fetch the same
	// path multiple times, so we use a map to record the resolved paths
	resolved := map[string]string{}
	for _, bundle := range dp.bundles {
		// For func rules, we fetch from either local fs or network directly
		for _, funcRules := range bundle.File2FuncRules {
			for _, rs := range funcRules {
				for _, rule := range rs {
					if rule.UseRaw {
						continue
					}
					if path, ok := resolved[rule.GetPath()]; ok {
						rule.SetPath(path)
						continue
					}
					util.Assert(rule.GetPath() != "", "sanity check")
					path, err := dp.fetchFrom(rule.GetPath())
					if err != nil {
						return fmt.Errorf("failed to fetch func rules: %w", err)
					}
					resolved[rule.GetPath()] = path
					rule.SetPath(path)
				}
			}
		}
		// For file rules, we fetch from either local fs or network
		// and concatenate the path with the file name as the final local path
		for _, fileRule := range bundle.FileRules {
			util.Assert(fileRule.GetPath() != "", "sanity check")
			var path string

			// If the path is already resolved, use it directly
			if p, ok := resolved[fileRule.GetPath()]; ok {
				path = p
			} else {
				p, err := dp.fetchFrom(fileRule.GetPath())
				if err != nil {
					return fmt.Errorf("failed to fetch file rules: %w", err)
				}
				path = p
				resolved[fileRule.GetPath()] = path
			}

			// Further check if the joined file exists
			file := filepath.Join(path, fileRule.FileName)
			exist, err := util.PathExists(file)
			if err != nil {
				return fmt.Errorf("failed to check file path: %w %v",
					err, fileRule.FileName)
			}
			if !exist {
				return fmt.Errorf("failed to find file %v", fileRule.FileName)
			}
			fileRule.FileName = file
			fileRule.SetPath(path)
		}
	}
	return nil
}

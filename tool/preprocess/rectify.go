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
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/data"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/errc"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/resource"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/util"
)

const (
	pkgPrefix = "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg"
)

var otelDeps = map[string]string{
	"go.opentelemetry.io/otel":                                          "v1.35.0",
	"go.opentelemetry.io/otel/sdk":                                      "v1.35.0",
	"go.opentelemetry.io/otel/trace":                                    "v1.35.0",
	"go.opentelemetry.io/otel/metric":                                   "v1.35.0",
	"go.opentelemetry.io/otel/sdk/metric":                               "v1.35.0",
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace":                 "v1.35.0",
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc":   "v1.35.0",
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp":   "v1.35.0",
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc": "v1.35.0",
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp": "v1.35.0",
	"go.opentelemetry.io/otel/exporters/prometheus":                     "v0.57.0",
	"go.opentelemetry.io/contrib/instrumentation/runtime":               "v0.60.0",
	"google.golang.org/protobuf":                                        "v1.35.2",
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric":            "v1.35.0",
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace":             "v1.35.0",
	"go.opentelemetry.io/otel/exporters/zipkin":                         "v1.35.0",
}

func extractGZip(data []byte, targetDir string) error {
	err := os.MkdirAll(targetDir, 0755)
	if err != nil {
		return errc.New(errc.ErrMkdirAll, err.Error())
	}

	gzReader, err := gzip.NewReader(strings.NewReader(string(data)))
	if err != nil {
		return errc.New(errc.ErrReadDir, fmt.Sprintf("gzip error: %v", err))
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errc.New(errc.ErrReadDir, fmt.Sprintf("gzip error: %v", err))
		}

		// Skip AppleDouble files (._filename) and other hidden files
		if strings.HasPrefix(filepath.Base(header.Name), "._") ||
			strings.HasPrefix(filepath.Base(header.Name), ".") {
			continue
		}

		targetPath := filepath.Join(targetDir, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			err = os.MkdirAll(targetPath, os.FileMode(header.Mode))
			if err != nil {
				return errc.New(errc.ErrMkdirAll, err.Error())
			}

		case tar.TypeReg:
			file, err := os.OpenFile(targetPath, os.O_CREATE|os.O_RDWR,
				os.FileMode(header.Mode))
			if err != nil {
				return errc.New(errc.ErrOpenFile, err.Error())
			}
			defer file.Close()

			_, err = io.Copy(file, tarReader)
			if err != nil {
				return errc.New(errc.ErrCopyFile, err.Error())
			}

		default:
			return errc.New(errc.ErrPreprocess,
				fmt.Sprintf("unsupported file type: %c in %s",
					header.Typeflag, header.Name))
		}
	}

	return nil
}

// Fetch the zipped pkg module from the embedded data section and extract it to
// a temporary directory, then return the path to the pkg directory.
func findModCacheDir() (string, error) {
	bs, err := data.UseEmbededPkg()
	if err != nil {
		return "", errc.New(errc.ErrPreprocess,
			fmt.Sprintf("error reading embedded pkg: %v", err))
	}
	tempPkg := util.GetTempBuildDirWith("alibaba-pkg")
	if util.PathExists(tempPkg) {
		_ = os.RemoveAll(tempPkg)
	}
	err = extractGZip(bs, tempPkg)
	if err != nil {
		return "", err
	}
	return filepath.Join(tempPkg, "pkg"), nil
}

// rectifyRule rectifies the file rules path to the local module cache path.
func (dp *DepProcessor) rectifyRule(bundles []*resource.RuleBundle) error {
	util.GuaranteeInPreprocess()
	defer util.PhaseTimer("Fetch")()
	rectified := map[string]bool{}
	for _, bundle := range bundles {
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

func (dp *DepProcessor) rectifyMod() error {
	// Backup go.mod and go.sum files
	gomodDir := dp.getGoModDir()
	files := []string{}
	files = append(files, filepath.Join(gomodDir, util.GoModFile))
	files = append(files, filepath.Join(gomodDir, util.GoSumFile))
	files = append(files, filepath.Join(gomodDir, util.GoWorkSumFile))
	for _, file := range files {
		if util.PathExists(file) {
			err := dp.backupFile(file)
			if err != nil {
				return err
			}
		}
	}
	// Since we haven't published the alibaba-otel pkg module, we need to add
	// a replace directive to tell the go tool to use the local module cache
	// instead of the remote module. This is a workaround for the case that
	// the remote module is not available(published).
	replaceMap := map[string][2]string{
		pkgPrefix: {dp.pkgLocalCache, ""},
	}
	// OTel dependencies may publish new versions that are not compatible
	// with the otel tool. In such cases, we need to add a replace directive
	// to use certain versions of the OTel dependencies for given otel tool,
	// otherwise, the otel tool may fail to run.
	for path, version := range otelDeps {
		replaceMap[path] = [2]string{path, version}
	}
	err := addModReplace(dp.getGoModPath(), replaceMap)
	if err != nil {
		return err
	}
	return nil
}

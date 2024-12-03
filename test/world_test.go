// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/shared"
)

const WorldAppName = "world"

func TestCompileTheWorld(t *testing.T) {
	UseApp(WorldAppName)

	RunGoBuild(t, "go", "build")
	RunApp(t, WorldAppName)
	text := ReadPreprocessLog(t, shared.DebugLogFile)

	regex := `\"ImportPath\":\"([^"]+)\"`
	r := regexp.MustCompile(regex)

	importPaths := make(map[string]struct{})
	matches := r.FindAllStringSubmatch(text, -1)

	for _, match := range matches {
		importPath := match[1]
		importPaths[importPath] = struct{}{}
	}
	if len(importPaths) != 23 {
		t.Log("Matched import paths:")
		for path := range importPaths {
			fmt.Println(path)
		}
		t.Fatalf("Rule matching is not complete")
	}
}

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"embed"
	_ "embed"
)

//go:embed core inst-api inst-api-semconv
var embedDir embed.FS

func ExportPkgDirList() []string {
	return []string{"core", "inst-api", "inst-api-semconv"}
}
func ExportPkgFS() embed.FS {
	return embedDir
}

//go:embed otel_setup.go
var otelSetupSDKTemplate string

func ExportOtelSetupSDKTemplate() string {
	return otelSetupSDKTemplate
}

//go:embed api/api.go
var apiSnippet string

func ExportAPISnippet() string { return apiSnippet }

//go:embed data/default.json
var defaultRuleJson string

func ExportDefaultRuleJson() string { return defaultRuleJson }

// This is simply a cache for the embedded rules directory. Technically, now any
// rules can by fetched from either network or local file system. But we want
// to keep the rules in tool binary, in this way, we dont need to fetch them from
// network, which is a very slow process.
//
//go:embed rules
var ruleCache embed.FS

func ExportRuleCache() embed.FS {
	return ruleCache
}

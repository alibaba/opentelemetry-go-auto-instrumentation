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

package pkg

import (
	"embed"
	_ "embed"
)

//go:embed core inst-api inst-api-semconv test
var embedDir embed.FS

func ExportPkgDirList() []string {
	return []string{"core", "inst-api", "inst-api-semconv", "test"}
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

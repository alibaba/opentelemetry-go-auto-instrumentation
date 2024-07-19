package pkg

import (
	"embed"
	_ "embed"
)

//go:embed *
var embedDir embed.FS

func ExportPkgDirList() []string {
	return []string{}
}
func ExportPkgFS() embed.FS {
	return embedDir
}

//go:embed otel_setup.go
var otelSetupSDKTemplate string

func ExportOtelSetupSDKTemplate() string {
	return otelSetupSDKTemplate
}

//go:embed rules
var ruleFS embed.FS

func ExportRuleFS() embed.FS {
	return ruleFS
}

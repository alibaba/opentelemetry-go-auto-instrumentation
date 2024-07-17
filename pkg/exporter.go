package pkg

import (
	"embed"
	_ "embed"
)

//go:embed rules
var ruleFS embed.FS

func ExportRuleFS() embed.FS {
	return ruleFS
}

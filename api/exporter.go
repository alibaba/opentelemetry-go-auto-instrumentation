package api

import (
	_ "embed"
)

//go:embed api.go
var apiTemplate string

func ExportAPITemplate() string {
	return apiTemplate
}

var Rules = make([]InstRule, 0)

func (rule *InstFuncRule) Register() {
	Rules = append(Rules, rule)
}

func (rule *InstFileRule) Register() {
	Rules = append(Rules, rule)
}

func (rule *InstStructRule) Register() {
	Rules = append(Rules, rule)
}

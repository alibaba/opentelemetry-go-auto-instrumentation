package api

import (
	"fmt"
)

// -----------------------------------------------------------------------------
// Execution Order
//
// The ExecOrder type is used to define the order of the rule execution. The pair
// of onEnter and onExit callbacks are executed from outermost to innermost. For
// example, given the following rules:
//  1. Rule1: {ExecOrderOutermost, OnEnter1, OnExit1}
//  2. Rule2: {ExecOrderInner, OnEnter2, OnExit2}
// The order of execution will be as follows
// OnEnter1 -> OnEnter2 -> Original Function -> OnExit2 -> OnExit1

type ExecOrder int

const (
	ExecOrderOutermost ExecOrder = 0
	ExecOrderOuter     ExecOrder = 1
	ExecOrderInner     ExecOrder = 2
	ExecOrderInnermost ExecOrder = 3
	ExecOrderDefault             = ExecOrderOutermost
)

// -----------------------------------------------------------------------------
// Instrumentation Rule
//
// Instrumentation rules are used to define the behavior of the instrumentation
// for a specific function call. The rules are defined in the init() function
// of rule.go in each package directory. The rules are then used by the instrument
// package to generate the instrumentation code. Multiple rules can be defined
// for a single function call, and the rules are executed in the order of their
// priority. The Order is defined by the ExecOrder type. The rules are executed
// in the order of their priority, from high to low.
// There are several types of rules for different purposes:
// - InstFuncRule: Instrumentation rule for a specific function call
// - InstStructRule: Instrumentation rule for a specific struct type
// - InstFileRule: Instrumentation rule for a specific file

type InstRule interface {
	GetVersion() string    // GetVersion returns the version of the rule
	GetImportPath() string // GetImportPath returns the import path of the rule
	GetRuleName() string   // GetRuleName returns the rule name
	String() string        // String returns string representation of the rule
}

type InstBaseRule struct {
	RuleName   string // Optional rule name
	Version    string // Version of the rule, e.g. "[1.9.1,1.9.2)" or "" for general match
	ImportPath string // Import path of the rule, e.g. "github.com/gin-gonic/gin"
}

func (rule *InstBaseRule) GetVersion() string {
	return rule.Version
}

func (rule *InstBaseRule) GetImportPath() string {
	return rule.ImportPath
}

func (rule *InstBaseRule) GetRuleName() string {
	return rule.RuleName
}

// InstFuncRule finds specific function call and instrument by adding new code
type InstFuncRule struct {
	InstBaseRule
	Function     string    // Function name, e.g. "New"
	ReceiverType string    // Receiver type name, e.g. "*gin.Engine"
	Order        ExecOrder // Order of the rule
	UseRaw       bool      // UseRaw indicates whether to insert raw code string
	OnEnter      string    // OnEnter callback, called before original function
	OnExit       string    // OnExit callback, called after original function
	FileDeps     []string  // File dependencies, add custom file into project
	PackageDeps  []string  // Overwrite package dependencies
}

// InstStructRule finds specific struct type and instrument by adding new field
type InstStructRule struct {
	InstBaseRule
	StructType string // Struct type name, e.g. "Engine"
	FieldName  string // New field name, e.g. "Logger"
	FieldType  string // New field type, e.g. "zap.Logger"
}

// InstFileRule adds user file into compilation unit and do further compilation
type InstFileRule struct {
	InstBaseRule
	FileName string // File name, e.g. "engine.go"
	Replace  bool   // Replace indicates whether to replace the original file
}

func NewRule(importPath, funcName, recvTypeName string,
	onEnter, onExit string) *InstFuncRule {
	rule := &InstFuncRule{
		InstBaseRule: InstBaseRule{
			RuleName:   "",
			Version:    "",
			ImportPath: importPath,
		},
		Function:     funcName,
		ReceiverType: recvTypeName,
		Order:        ExecOrderOutermost,
		UseRaw:       false,
		OnEnter:      onEnter,
		OnExit:       onExit,
		FileDeps:     make([]string, 0),
		PackageDeps:  make([]string, 0),
	}
	return rule
}

func (rule *InstFuncRule) WithVersion(version string) *InstFuncRule {
	rule.Version = version
	return rule
}

func (rule *InstFuncRule) WithExecOrder(priority ExecOrder) *InstFuncRule {
	rule.Order = priority
	return rule
}

func (rule *InstFuncRule) WithUseRaw(useRaw bool) *InstFuncRule {
	rule.UseRaw = useRaw
	return rule
}

func (rule *InstFuncRule) WithFileDeps(deps ...string) *InstFuncRule {
	rule.FileDeps = append(rule.FileDeps, deps...)
	return rule
}

func (rule *InstFuncRule) WithPackageDep(dep, version string) *InstFuncRule {
	rule.PackageDeps = append(rule.PackageDeps, dep+"@"+version)
	return rule
}

func NewFileRule(importPath, fileName string) *InstFileRule {
	rule := &InstFileRule{
		InstBaseRule: InstBaseRule{
			RuleName:   "",
			Version:    "",
			ImportPath: importPath,
		},
		FileName: fileName,
		Replace:  false,
	}
	return rule
}

func NewStructRule(importPath, structName, fieldName, fieldType string) *InstStructRule {
	rule := &InstStructRule{
		InstBaseRule: InstBaseRule{
			RuleName:   "",
			Version:    "",
			ImportPath: importPath,
		},
		StructType: structName,
		FieldName:  fieldName,
		FieldType:  fieldType,
	}
	return rule
}

func (rule *InstFileRule) WithReplace(replace bool) *InstFileRule {
	rule.Replace = replace
	return rule
}

func (rule *InstFileRule) WithVersion(version string) *InstFileRule {
	rule.Version = version
	return rule
}

func (rule *InstFuncRule) String() string {
	if rule.ReceiverType == "" {
		return fmt.Sprintf("%s@%s@%s@%s {%s %s}",
			rule.RuleName,
			rule.ImportPath, rule.Version,
			rule.Function,
			rule.OnEnter, rule.OnExit)
	}
	return fmt.Sprintf("%s@%s@%s@(%s).%s {%s %s}",
		rule.RuleName,
		rule.ImportPath, rule.Version,
		rule.ReceiverType, rule.Function,
		rule.OnEnter, rule.OnExit)
}

func (rule *InstStructRule) String() string {
	return fmt.Sprintf("%s@%s@%s {%s}",
		rule.RuleName,
		rule.ImportPath, rule.Version,
		rule.StructType)
}

func (rule *InstFileRule) String() string {
	return fmt.Sprintf("%s@%s@%s {%s}",
		rule.RuleName,
		rule.ImportPath, rule.Version,
		rule.FileName)
}

func (rule *InstFuncRule) WithRuleName(name string) *InstFuncRule {
	rule.RuleName = name
	return rule
}

func (rule *InstStructRule) WithRuleName(name string) *InstStructRule {
	rule.RuleName = name
	return rule
}

func (rule *InstFileRule) WithRuleName(name string) *InstFileRule {
	rule.RuleName = name
	return rule
}

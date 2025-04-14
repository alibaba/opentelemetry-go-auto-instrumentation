package test

import (
	"testing"
)

const mcp_dependency_name = "github.com/mark3labs/mcp-go/mcp"
const mcp_module_name = "mcp"

func init() {
	TestCases = append(TestCases,
		NewGeneralTestCase("mcp-0.20.0-sse-tool-test", mcp_module_name, "0.20.0", "0.20.0", "1.22.0", "", TestMcpTool),
		NewGeneralTestCase("mcp-0.20.0-sse-prompt-test", mcp_module_name, "0.20.0", "0.20.0", "1.22.0", "", TestMcpPrompt),
		NewGeneralTestCase("mcp-0.20.0-sse-resource-test", mcp_module_name, "0.20.0", "0.20.0", "1.22.0", "", TestMcpResource),
	)

}

func TestMcpTool(t *testing.T, env ...string) {
	UseApp("mcp/v0.20.0")
	RunGoBuild(t, "go", "build", "test_sse_tool.go", "ext.go")
	RunApp(t, "test_sse_tool", env...)
}
func TestMcpPrompt(t *testing.T, env ...string) {
	UseApp("mcp/v0.20.0")
	RunGoBuild(t, "go", "build", "test_sse_prompt.go", "ext.go")
	RunApp(t, "test_sse_prompt", env...)
}
func TestMcpResource(t *testing.T, env ...string) {
	UseApp("mcp/v0.20.0")
	RunGoBuild(t, "go", "build", "test_sse_resource.go", "ext.go")
	RunApp(t, "test_sse_resource", env...)
}

// 由于标准输入输出通信通信无法在此处test中实现，会挂住测试进程，所以此方法只留作后续stdio可用时使用，目前不使用
// Since standard input/output communication cannot be implemented in the test here and will cause the test process to hang, this method is reserved for future use when stdio becomes available, and is currently not used.
/*func TestStdioTool(t *testing.T, env ...string) {
	UseApp("mcp/v0.20.0")
	RunGoBuild(t, "go", "build", "test_stdio_server.go", "ext.go")
	fmt.Println("a1")
	RunGoBuild(t, "go", "build", "test_stdio_tool.go")
	fmt.Println("a2")
	RunApp(t, "test_stdio_tool", env...)
}*/

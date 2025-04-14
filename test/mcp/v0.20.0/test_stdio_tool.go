package main

import (
	"context"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"time"
)

// 由于标准输入输出通信通信无法在此处test中实现，会挂住测试进程，所以此方法只留作后续stdio可用时使用，目前不使用
// Since standard input/output communication cannot be implemented in the test here and will cause the test process to hang, this method is reserved for future use when stdio becomes available, and is currently not used.

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	c, err := client.NewStdioMCPClient(
		"./test_stdio_server",
		[]string{"IN_OTEL_TEST=true"}, // Empty ENV
	)
	if err != nil {
		panic(err)
	}
	defer c.Close()
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "example-client",
		Version: "1.0.0",
	}
	_, err = c.Initialize(ctx, initRequest)
	if err != nil {
		panic(err)
	}
	listTmpRequest := mcp.CallToolRequest{}
	listTmpRequest.Params.Name = "hello_world"
	listTmpRequest.Params.Arguments = map[string]interface{}{
		"name": "abc",
	}
	_, err = c.CallTool(ctx, listTmpRequest)
	if err != nil {
		panic(err)
	}
	toolsRequest := mcp.ListToolsRequest{}
	_, err = c.ListTools(ctx, toolsRequest)
	if err != nil {
		panic(err)
	}
	/*verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		verifier.VerifyLLMCommonAttributes(stubs[0][2], "execute_other:initialize", "mcp")
		verifier.VerifyLLMCommonAttributes(stubs[3][2], "execute_tool", "mcp")
		verifier.VerifyLLMCommonAttributes(stubs[4][2], "execute_other:tools/list", "mcp")
	}, 0)*/
}

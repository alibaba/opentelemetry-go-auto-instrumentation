package main

import (
	"context"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"time"
)

func main() {
	mcpServer := server.NewMCPServer("test", "1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithToolCapabilities(true),
	)
	tool := mcp.NewTool("hello_world",
		mcp.WithDescription("Say hello to someone"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the person to greet"),
		),
	)
	// Add tool handler
	mcpServer.AddTool(tool, helloHandler)
	testServer := server.NewTestServer(mcpServer)
	defer testServer.Close()
	// Connect to SSE endpoint
	c, err := client.NewSSEMCPClient(testServer.URL + "/sse")
	if err != nil {
		panic(err)
	}
	defer c.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Start the mcpClient
	if err := c.Start(ctx); err != nil {
		panic(err)
	}

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

	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		verifier.VerifyLLMCommonAttributes(stubs[1][2], "execute_other:initialize", "mcp")
		verifier.VerifyLLMCommonAttributes(stubs[3][2], "execute_tool", "mcp")
		verifier.VerifyLLMCommonAttributes(stubs[4][2], "execute_other:tools/list", "mcp")
	}, 5)
}

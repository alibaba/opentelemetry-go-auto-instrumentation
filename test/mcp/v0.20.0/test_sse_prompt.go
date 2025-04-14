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
	mcpServer.AddPrompt(mcp.NewPrompt(SIMPLE,
		mcp.WithPromptDescription("A simple prompt"),
	), handleSimplePrompt)

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

	GetPrompt := mcp.GetPromptRequest{}
	GetPrompt.Params.Name = SIMPLE
	_, err = c.GetPrompt(ctx, GetPrompt)
	if err != nil {
		panic(err)
	}

	promptsRequest := mcp.ListPromptsRequest{}
	_, err = c.ListPrompts(ctx, promptsRequest)
	if err != nil {
		panic(err)
	}

	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		verifier.VerifyLLMCommonAttributes(stubs[1][2], "execute_other:initialize", "mcp")
		verifier.VerifyLLMCommonAttributes(stubs[3][2], "execute_other:prompts/get", "mcp")
		verifier.VerifyLLMCommonAttributes(stubs[4][2], "execute_other:prompts/list", "mcp")
	}, 5)
}

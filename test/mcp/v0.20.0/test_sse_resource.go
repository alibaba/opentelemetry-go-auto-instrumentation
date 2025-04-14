package main

import (
	"context"
	"fmt"
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
	mcpServer.AddResource(mcp.NewResource("test://static/resource",
		"Static Resource",
		mcp.WithMIMEType("text/plain"),
	), handleReadResource)

	mcpServer.AddResourceTemplate(
		mcp.NewResourceTemplate(
			"test://dynamic/resource/{id}",
			"Dynamic Resource",
		),
		handleResourceTemplate,
	)

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

	fmt.Println("ReadResource...")
	readResourceRequest := mcp.ReadResourceRequest{}
	readResourceRequest.Params.URI = "test://static/resource"
	_, err = c.ReadResource(ctx, readResourceRequest)
	if err != nil {
		panic(err)
	}

	listResourcesRequest := mcp.ListResourcesRequest{}
	_, err = c.ListResources(ctx, listResourcesRequest)
	if err != nil {
		panic(err)
	}

	listResourceTemplatesRequest := mcp.ListResourceTemplatesRequest{}
	_, err = c.ListResourceTemplates(ctx, listResourceTemplatesRequest)
	if err != nil {
		panic(err)
	}

	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		verifier.VerifyLLMCommonAttributes(stubs[1][2], "execute_other:initialize", "mcp")
		verifier.VerifyLLMCommonAttributes(stubs[3][2], "execute_other:resources/read", "mcp")
		verifier.VerifyLLMCommonAttributes(stubs[4][2], "execute_other:resources/list", "mcp")
		verifier.VerifyLLMCommonAttributes(stubs[5][2], "execute_other:resources/templates/list", "mcp")

	}, 7)
}

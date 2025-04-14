package main

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// 由于标准输入输出通信通信无法在此处test中实现，会挂住测试进程，所以此方法只留作后续stdio可用时使用，目前不使用
// Since standard input/output communication cannot be implemented in the test here and will cause the test process to hang, this method is reserved for future use when stdio becomes available, and is currently not used.

func main() {
	mcpServer := server.NewMCPServer("test", "1.0.0",
		//server.WithResourceCapabilities(true, true),
		//server.WithPromptCapabilities(true),
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

	mcpServer.AddPrompt(mcp.NewPrompt(SIMPLE,
		mcp.WithPromptDescription("A simple prompt"),
	), handleSimplePrompt)
	if err := server.ServeStdio(mcpServer); err != nil {
		panic(err)
	}

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
	if err := server.ServeStdio(mcpServer); err != nil {
		panic(err)
	}
}

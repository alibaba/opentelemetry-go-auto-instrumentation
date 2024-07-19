package main

import (
	"log"
	"os"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/shared"
)

func main() {
	shared.ParseOptions()
	if shared.PrintVersion {
		shared.PrintTheVersion()
		os.Exit(0)
	}
	err := shared.InitOptions()
	if err != nil {
		log.Printf("failed to init options: %v", err)
		os.Exit(1)

	}
	err = tool.Run()
	if err != nil {
		log.Printf("failed to run the tool: %v", err)
		os.Exit(1)
	}
}

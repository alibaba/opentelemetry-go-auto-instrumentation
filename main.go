package main

// Entry

import (
	"flag"
	"fmt"
	"os"
	"otel-auto-instrumentation/internal"
	"otel-auto-instrumentation/internal/shared"
)

func main() {
	flag.BoolVar(&shared.InToolexec, shared.NameOfInToolexec, false, shared.UsageOfIntoolexec)
	flag.Parse()

	err := internal.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

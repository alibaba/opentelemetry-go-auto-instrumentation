package test

import (
	"testing"
)

const ollama_dependency_name = "github.com/ollama/ollama"
const ollama_module_name = "ollama"

func init() {
	TestCases = append(TestCases,
		NewGeneralTestCase("ollama-0.3.14-chat-test", ollama_module_name, "0.3.14", "0.3.14", "1.22.0", "", TestOllamaChat),
		NewGeneralTestCase("ollama-0.3.14-generate-test", ollama_module_name, "0.3.14", "0.3.14", "1.22.0", "", TestOllamaGenerate),
	)
}

func TestOllamaChat(t *testing.T, env ...string) {
	UseApp("ollama/v0.3.14")
	RunGoBuild(t, "go", "build", "test_chat.go")
	RunApp(t, "test_chat", env...)
}

func TestOllamaGenerate(t *testing.T, env ...string) {
	UseApp("ollama/v0.3.14")
	RunGoBuild(t, "go", "build", "test_generate.go")
	RunApp(t, "test_generate", env...)
}
package test

import "testing"

const eino_dependency_name = "github.com/cloudwego/eino"
const eino_module_name = "eino"

func init() {
	TestCases = append(TestCases,
		NewGeneralTestCase("eino-react-agent-test", eino_module_name, "v0.3.51", "", "1.18", "", TestReactAgentEino),
		NewGeneralTestCase("eino-openai-invoke-test", eino_module_name, "v0.3.51", "", "1.18", "", TestOpenAIInvokeEino),
		NewGeneralTestCase("eino-openai-stream-test", eino_module_name, "v0.3.51", "", "1.18", "", TestOpenAIStreamEino),
		NewGeneralTestCase("eino-ollama-invoke-test", eino_module_name, "v0.3.51", "", "1.18", "", TestOllamaInvokeEino),
		NewGeneralTestCase("eino-ollama-stream-test", eino_module_name, "v0.3.51", "", "1.18", "", TestOllamaStreamEino),
		NewGeneralTestCase("eino-document-test", eino_module_name, "v0.3.51", "", "1.18", "", TestDocumentEino),
		NewLatestDepthTestCase("eino-latest-depth-test", eino_dependency_name, eino_module_name, "v0.3.51", "", "1.21", "", TestOpenAIInvokeEino),
		NewMuzzleTestCase("eino-muzzle-test-react-agent", eino_dependency_name, eino_module_name, "v0.3.51", "", "1.21", "1.24", []string{"go", "build", "test_react_agent.go", "eino_common.go"}),
		NewMuzzleTestCase("eino-muzzle-test-openai-invoke", eino_dependency_name, eino_module_name, "v0.3.51", "", "1.21", "1.24", []string{"go", "build", "test_openai_invoke_chatmodel.go", "eino_common.go"}),
		NewMuzzleTestCase("eino-muzzle-test-openai-stream", eino_dependency_name, eino_module_name, "v0.3.51", "", "1.21", "1.24", []string{"go", "build", "test_openai_stream_chatmodel.go", "eino_common.go"}),
		NewMuzzleTestCase("eino-muzzle-test-ollama-invoke", eino_dependency_name, eino_module_name, "v0.3.51", "", "1.21", "", []string{"go", "build", "test_ollama_invoke_chatmodel.go", "eino_common.go"}),
		NewMuzzleTestCase("eino-muzzle-test-ollama-stream", eino_dependency_name, eino_module_name, "v0.3.51", "", "1.21", "", []string{"go", "build", "test_ollama_stream_chatmodel.go", "eino_common.go"}),
		NewMuzzleTestCase("eino-muzzle-test-document-test", eino_dependency_name, eino_module_name, "v0.3.51", "", "1.21", "", []string{"go", "build", "test_document_graph.go", "eino_common.go"}),
	)
}

func TestReactAgentEino(t *testing.T, env ...string) {
	UseApp("eino/v0.3.51")
	RunGoBuild(t, "go", "build", "test_react_agent.go", "eino_common.go")
	RunApp(t, "test_react_agent", env...)
}

func TestDocumentEino(t *testing.T, env ...string) {
	UseApp("eino/v0.3.51")
	RunGoBuild(t, "go", "build", "test_document_graph.go", "eino_common.go")
	RunApp(t, "test_document_graph", env...)
}

func TestOpenAIInvokeEino(t *testing.T, env ...string) {
	UseApp("eino/v0.3.51")
	RunGoBuild(t, "go", "build", "test_openai_invoke_chatmodel.go", "eino_common.go")
	RunApp(t, "test_openai_invoke_chatmodel", env...)
}

func TestOpenAIStreamEino(t *testing.T, env ...string) {
	UseApp("eino/v0.3.51")
	RunGoBuild(t, "go", "build", "test_openai_stream_chatmodel.go", "eino_common.go")
	RunApp(t, "test_openai_stream_chatmodel", env...)
}

func TestOllamaInvokeEino(t *testing.T, env ...string) {
	UseApp("eino/v0.3.51")
	RunGoBuild(t, "go", "build", "test_ollama_invoke_chatmodel.go", "eino_common.go")
	RunApp(t, "test_ollama_invoke_chatmodel", env...)
}

func TestOllamaStreamEino(t *testing.T, env ...string) {
	UseApp("eino/v0.3.51")
	RunGoBuild(t, "go", "build", "test_ollama_stream_chatmodel.go", "eino_common.go")
	RunApp(t, "test_ollama_stream_chatmodel", env...)
}

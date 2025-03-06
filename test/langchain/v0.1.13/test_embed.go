package main

import (
	"context"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier"
	"github.com/tmc/langchaingo/embeddings"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func main() {
	/*llmEm, err := ollama.New([]ollama.Option{
		ollama.WithModel("snowflake-arctic-embed:22m"),
		ollama.WithServerURL("http://127.0.0.1:" + os.Getenv("OLLAMA_EMBD_PORT")),
	}...)
	if err != nil {
		panic(err)
	}*/

	emb, err := embeddings.NewEmbedder(embedFake{},
		embeddings.WithBatchSize(15), //批量上传限制
	)
	if err != nil {
		panic(err)
	}
	txt1 := `君不见黄河之水天上来，奔流到海不复回。`
	txt2 := `君不见高堂明镜悲白发，朝如青丝暮成雪。`
	_, err = emb.EmbedQuery(context.Background(), txt1)
	if err != nil {
		panic(err)
	}
	_, err = emb.EmbedDocuments(context.Background(), []string{txt1, txt2})
	if err != nil {
		panic(err)
	}
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		verifier.VerifyLLMAttributes(stubs[0][0], "singleEmbed")
		verifier.VerifyLLMAttributes(stubs[1][0], "batchedEmbed")
	}, 3)
}

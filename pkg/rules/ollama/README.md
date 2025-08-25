# Ollama Instrumentation

This package provides OpenTelemetry instrumentation for [Ollama](https://github.com/ollama/ollama), a framework for running large language models locally.

## Features

### Core Instrumentation
- **Chat API**: Full instrumentation for chat completions
- **Generate API**: Full instrumentation for text generation
- **OpenTelemetry Spans**: Automatic span creation with appropriate attributes
- **Token Metrics**: Capture input and output token counts

### Streaming Support (NEW)
- **Real-time Streaming**: Full support for Ollama's callback-based streaming
- **TTFT Measurement**: Time To First Token tracking for streaming responses
- **Chunk Tracking**: Monitor streaming progress with chunk counts
- **Token Rate**: Calculate tokens per second throughput
- **Span Events**: Record streaming milestones (first token, progress, completion)

## Instrumented Methods

The following Ollama API methods are instrumented:

- `Client.Generate()` - Text generation with optional streaming
- `Client.Chat()` - Chat completions with optional streaming

## Attributes Captured

### Standard OpenTelemetry GenAI Attributes
- `gen_ai.system`: "ollama"
- `gen_ai.request.model`: Model name (e.g., "llama3:8b")
- `gen_ai.operation.name`: "chat" or "generate"
- `gen_ai.usage.input_tokens`: Input token count
- `gen_ai.usage.output_tokens`: Output token count

### Streaming-Specific Attributes
- `gen_ai.response.streaming`: Boolean indicating if streaming was used
- `gen_ai.response.ttft_ms`: Time to first token in milliseconds
- `gen_ai.response.chunk_count`: Total number of streaming chunks
- `gen_ai.response.tokens_per_second`: Token generation throughput
- `gen_ai.response.stream_duration_ms`: Total streaming duration

## Streaming Behavior

Ollama streams by default when the `Stream` field is:
- `nil` (not specified) - **defaults to streaming**
- `true` - explicitly enables streaming
- `false` - explicitly disables streaming

## Span Events

For streaming requests, the following span events are recorded:

1. **First token received**: Records TTFT when first content chunk arrives
2. **Streaming progress**: Periodic updates every 10 chunks or 500ms
3. **Streaming completed**: Final metrics including total chunks and token rate

## Usage Examples

### Non-Streaming Request
```go
streamFlag := false
req := &api.GenerateRequest{
    Model:  "llama3:8b",
    Prompt: "Hello, world!",
    Stream: &streamFlag, // Explicitly disable streaming
}
```

### Streaming Request (Default)
```go
req := &api.GenerateRequest{
    Model:  "llama3:8b",
    Prompt: "Write a poem",
    // Stream not set - defaults to streaming
}
```

### Explicit Streaming Request
```go
streamFlag := true
req := &api.GenerateRequest{
    Model:  "llama3:8b",
    Prompt: "Write a poem",
    Stream: &streamFlag, // Explicitly enable streaming
}
```

## Testing

Run the instrumented tests with:

```bash
# Non-streaming tests
cd test/ollama/v0.3.14
../../../otel go build -o test_generate test_generate.go
../../../otel go build -o test_chat test_chat.go

# Streaming tests
../../../otel go build -o test_generate_stream test_generate_stream.go
../../../otel go build -o test_chat_stream test_chat_stream.go

# Backward compatibility test
../../../otel go build -o test_backward_compat test_backward_compat.go

# Run with OpenTelemetry export
OTEL_EXPORTER_OTLP_ENDPOINT="http://localhost:4318" ./test_generate_stream
```

## Limitations

- **Input Token Counts**: Due to Ollama API design, input token counts are only available in the response, not the request
- **Streaming Token Counts**: Token counts are cumulative and only accurate in the final chunk
- **Model Parameters**: Advanced parameters (temperature, top_p, etc.) are not currently captured

## Version Support

- Minimum Ollama version: v0.3.14
- Maximum Ollama version: Latest (no upper limit)
package eino

type (
	promptRequestKey    struct{}
	llmRequestKey       struct{}
	embeddingRequestKey struct{}
	indexerRequestKey   struct{}
	retrieverRequestKey struct{}
	loaderRequestKey    struct{}
	toolRequestKey      struct{}
)

type einoRequest struct {
	operationName string
	input         map[string]string
}

type einoResponse struct {
	operationName string
	output        map[string]string
}

type einoLLMRequest struct {
	operationName    string
	modelName        string
	encodingFormats  []string
	frequencyPenalty float64
	presencePenalty  float64
	maxTokens        int64
	usageInputTokens int64
	stopSequences    []string
	temperature      float64
	topK             float64
	topP             float64
	serverAddress    string
	seed             int64
}
type einoLLMResponse struct {
	responseFinishReasons []string
	responseModel         string
	usageOutputTokens     int64
	responseID            string
}

type ChatModelConfig struct {
	BaseURL          string
	PresencePenalty  float64
	Seed             int64
	FrequencyPenalty float64
}

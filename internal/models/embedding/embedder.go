package embedding

import (
	"context"
	"fmt"
	"strings"

	"github.com/Tencent/WeKnora/internal/models/utils/ollama"
	"github.com/Tencent/WeKnora/internal/runtime"
	"github.com/Tencent/WeKnora/internal/types"
)

// Embedder defines the interface for text vectorization
type Embedder interface {
	// Embed converts text to vector
	Embed(ctx context.Context, text string) ([]float32, error)

	// BatchEmbed converts multiple texts to vectors in batch
	BatchEmbed(ctx context.Context, texts []string) ([][]float32, error)

	// GetModelName returns the model name
	GetModelName() string

	// GetDimensions returns the vector dimensions
	GetDimensions() int

	// GetModelID returns the model ID
	GetModelID() string

	EmbedderPooler
}

type EmbedderPooler interface {
	BatchEmbedWithPool(ctx context.Context, model Embedder, texts []string) ([][]float32, error)
}

// EmbedderType represents the embedder type
type EmbedderType string

// Config represents the embedder configuration
type Config struct {
	Source               types.ModelSource `json:"source"`
	BaseURL              string            `json:"base_url"`
	ModelName            string            `json:"model_name"`
	APIKey               string            `json:"api_key"`
	TruncatePromptTokens int               `json:"truncate_prompt_tokens"`
	Dimensions           int               `json:"dimensions"`
	ModelID              string            `json:"model_id"`
}

// NewEmbedder creates an embedder based on the configuration
func NewEmbedder(config Config) (Embedder, error) {
	var embedder Embedder
	var err error
	switch strings.ToLower(string(config.Source)) {
	case string(types.ModelSourceLocal):
		runtime.GetContainer().Invoke(func(pooler EmbedderPooler, ollamaService *ollama.OllamaService) {
			embedder, err = NewOllamaEmbedder(config.BaseURL,
				config.ModelName, config.TruncatePromptTokens, config.Dimensions, config.ModelID, pooler, ollamaService)
		})
		return embedder, err
	case string(types.ModelSourceRemote):
		runtime.GetContainer().Invoke(func(pooler EmbedderPooler) {
			embedder, err = NewOpenAIEmbedder(config.APIKey,
				config.BaseURL,
				config.ModelName,
				config.TruncatePromptTokens,
				config.Dimensions,
				config.ModelID,
				pooler)
		})
		return embedder, err
	default:
		return nil, fmt.Errorf("unsupported embedder source: %s", config.Source)
	}
}

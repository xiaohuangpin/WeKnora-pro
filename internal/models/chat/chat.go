package chat

import (
	"context"
	"fmt"
	"strings"

	"github.com/Tencent/WeKnora/internal/models/utils/ollama"
	"github.com/Tencent/WeKnora/internal/runtime"
	"github.com/Tencent/WeKnora/internal/types"
)

// ChatOptions 聊天选项
type ChatOptions struct {
	Temperature         float64 `json:"temperature"`           // 温度参数
	TopP                float64 `json:"top_p"`                 // Top P 参数
	Seed                int     `json:"seed"`                  // 随机种子
	MaxTokens           int     `json:"max_tokens"`            // 最大 token 数
	MaxCompletionTokens int     `json:"max_completion_tokens"` // 最大完成 token 数
	FrequencyPenalty    float64 `json:"frequency_penalty"`     // 频率惩罚
	PresencePenalty     float64 `json:"presence_penalty"`      // 存在惩罚
	Thinking            *bool   `json:"thinking"`              // 是否启用思考
}

// Message 表示聊天消息
type Message struct {
	Role    string `json:"role"`    // 角色：system, user, assistant
	Content string `json:"content"` // 消息内容
}

// Chat 定义了聊天接口
type Chat interface {
	// Chat 进行非流式聊天
	Chat(ctx context.Context, messages []Message, opts *ChatOptions) (*types.ChatResponse, error)

	// ChatStream 进行流式聊天
	ChatStream(ctx context.Context, messages []Message, opts *ChatOptions) (<-chan types.StreamResponse, error)

	// GetModelName 获取模型名称
	GetModelName() string

	// GetModelID 获取模型ID
	GetModelID() string
}

type ChatConfig struct {
	Source    types.ModelSource
	BaseURL   string
	ModelName string
	APIKey    string
	ModelID   string
}

// NewChat 创建聊天实例
func NewChat(config *ChatConfig) (Chat, error) {
	var chat Chat
	var err error
	switch strings.ToLower(string(config.Source)) {
	case string(types.ModelSourceLocal):
		runtime.GetContainer().Invoke(func(ollamaService *ollama.OllamaService) {
			chat, err = NewOllamaChat(config, ollamaService)
		})
		if err != nil {
			return nil, err
		}
		return chat, nil
	case string(types.ModelSourceRemote):
		return NewRemoteAPIChat(config)
	default:
		return nil, fmt.Errorf("unsupported chat model source: %s", config.Source)
	}
}

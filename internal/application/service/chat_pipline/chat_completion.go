package chatpipline

import (
	"context"

	"github.com/Tencent/WeKnora/internal/logger"
	"github.com/Tencent/WeKnora/internal/types"
	"github.com/Tencent/WeKnora/internal/types/interfaces"
)

// PluginChatCompletion implements chat completion functionality
// as a plugin that can be registered to EventManager
type PluginChatCompletion struct {
	modelService interfaces.ModelService // Interface for model operations
}

// NewPluginChatCompletion creates a new PluginChatCompletion instance
// and registers it with the EventManager
func NewPluginChatCompletion(eventManager *EventManager, modelService interfaces.ModelService) *PluginChatCompletion {
	res := &PluginChatCompletion{
		modelService: modelService,
	}
	eventManager.Register(res)
	return res
}

// ActivationEvents returns the event types this plugin handles
func (p *PluginChatCompletion) ActivationEvents() []types.EventType {
	return []types.EventType{types.CHAT_COMPLETION}
}

// OnEvent handles chat completion events
// It prepares the chat model, messages, and calls the model to generate responses
func (p *PluginChatCompletion) OnEvent(
	ctx context.Context, eventType types.EventType, chatManage *types.ChatManage, next func() *PluginError,
) *PluginError {
	logger.Info(ctx, "Starting chat completion")

	// Prepare chat model and options
	chatModel, opt, err := prepareChatModel(ctx, p.modelService, chatManage)
	if err != nil {
		return ErrGetChatModel.WithError(err)
	}

	// Prepare messages including conversation history
	logger.Info(ctx, "Preparing chat messages with history")
	chatMessages := prepareMessagesWithHistory(chatManage)

	// Call the chat model to generate response
	logger.Info(ctx, "Calling chat model")
	chatResponse, err := chatModel.Chat(ctx, chatMessages, opt)
	if err != nil {
		logger.Errorf(ctx, "Failed to call chat model: %v", err)
		return ErrModelCall.WithError(err)
	}

	logger.Info(ctx, "Chat completion successful")
	chatManage.ChatResponse = chatResponse
	return next()
}

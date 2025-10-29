package chatpipline

import (
	"context"
	"strings"

	"github.com/Tencent/WeKnora/internal/logger"
	"github.com/Tencent/WeKnora/internal/types"
)

// PluginStreamFilter implements stream filtering functionality for chat pipeline
type PluginStreamFilter struct{}

// NewPluginStreamFilter creates a new stream filter plugin instance
func NewPluginStreamFilter(eventManager *EventManager) *PluginStreamFilter {
	res := &PluginStreamFilter{}
	eventManager.Register(res)
	return res
}

// ActivationEvents returns the event types this plugin handles
func (p *PluginStreamFilter) ActivationEvents() []types.EventType {
	return []types.EventType{types.STREAM_FILTER}
}

// OnEvent handles stream filtering events in the chat pipeline
func (p *PluginStreamFilter) OnEvent(ctx context.Context,
	eventType types.EventType, chatManage *types.ChatManage, next func() *PluginError,
) *PluginError {
	logger.Info(ctx, "Starting stream filter")
	logger.Info(ctx, "Creating new stream channel")

	// Create new stream channel and replace the original one
	oldStream := chatManage.ResponseChan
	newStream := make(chan types.StreamResponse)
	chatManage.ResponseChan = newStream

	// Initialize response builder and check if no-match prefix filtering is needed
	responseBuilder := &strings.Builder{}
	matchNoMatchBuilderPrefix := chatManage.SummaryConfig.NoMatchPrefix != ""

	if matchNoMatchBuilderPrefix {
		logger.Infof(ctx, "Using no match prefix filter: %s", chatManage.SummaryConfig.NoMatchPrefix)
	}

	// Start goroutine to filter the stream
	go func() {
		logger.Info(ctx, "Starting stream filter goroutine")
		for resp := range oldStream {
			// Accumulate answer content
			if resp.ResponseType == types.ResponseTypeAnswer {
				responseBuilder.WriteString(resp.Content)
			}

			// Skip filtering if no prefix matching is required
			if !matchNoMatchBuilderPrefix {
				newStream <- resp
				continue
			}

			// Check if content matches the no-match prefix
			if !strings.HasPrefix(chatManage.SummaryConfig.NoMatchPrefix, responseBuilder.String()) {
				resp.Content = responseBuilder.String()
				newStream <- resp
				logger.Info(
					ctx, "Content does not match no-match prefix, passing through, content: ",
					responseBuilder.String(),
				)
				matchNoMatchBuilderPrefix = false
			}
		}

		// Handle NO_MATCH case when stream ends
		if matchNoMatchBuilderPrefix {
			logger.Info(ctx, "Content matches no-match prefix, using fallback response")
			newStream <- NewFallback(ctx, chatManage.FallbackResponse)
		}
		logger.Info(ctx, "Stream filter completed, closing new stream")
		close(newStream)
	}()

	logger.Info(ctx, "Stream filter initialized")
	return next()
}

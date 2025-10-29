package chatpipline

import (
	"context"

	"github.com/Tencent/WeKnora/internal/logger"
	"github.com/Tencent/WeKnora/internal/types"
)

// PluginFilterTopK is a plugin that filters search results to keep only the top K items
type PluginFilterTopK struct{}

// NewPluginFilterTopK creates a new instance of PluginFilterTopK and registers it with the event manager
func NewPluginFilterTopK(eventManager *EventManager) *PluginFilterTopK {
	res := &PluginFilterTopK{}
	eventManager.Register(res)
	return res
}

// ActivationEvents returns the event types that this plugin responds to
func (p *PluginFilterTopK) ActivationEvents() []types.EventType {
	return []types.EventType{types.FILTER_TOP_K}
}

// OnEvent handles the FILTER_TOP_K event by filtering results to keep only the top K items
// It can filter MergeResult, RerankResult, or SearchResult depending on which is available
func (p *PluginFilterTopK) OnEvent(ctx context.Context,
	eventType types.EventType, chatManage *types.ChatManage, next func() *PluginError,
) *PluginError {
	logger.Info(ctx, "Starting filter top-K process")
	logger.Infof(ctx, "Filter configuration: top-K = %d", chatManage.RerankTopK)

	filterTopK := func(searchResult []*types.SearchResult, topK int) []*types.SearchResult {
		if topK > 0 && len(searchResult) > topK {
			logger.Infof(ctx, "Filtering results: before=%d, after=%d", len(searchResult), topK)
			searchResult = searchResult[:topK]
		}
		return searchResult
	}

	if len(chatManage.MergeResult) > 0 {
		chatManage.MergeResult = filterTopK(chatManage.MergeResult, chatManage.RerankTopK)
	} else if len(chatManage.RerankResult) > 0 {
		chatManage.RerankResult = filterTopK(chatManage.RerankResult, chatManage.RerankTopK)
	} else if len(chatManage.SearchResult) > 0 {
		chatManage.SearchResult = filterTopK(chatManage.SearchResult, chatManage.RerankTopK)
	} else {
		logger.Info(ctx, "No results to filter")
	}

	logger.Info(ctx, "Filter top-K process completed")
	return next()
}

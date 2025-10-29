package chatpipline

import (
	"context"
	"encoding/json"
	"sort"

	"github.com/Tencent/WeKnora/internal/logger"
	"github.com/Tencent/WeKnora/internal/types"
)

// PluginMerge handles merging of search result chunks
type PluginMerge struct{}

// NewPluginMerge creates and registers a new PluginMerge instance
func NewPluginMerge(eventManager *EventManager) *PluginMerge {
	res := &PluginMerge{}
	eventManager.Register(res)
	return res
}

// ActivationEvents returns the event types this plugin handles
func (p *PluginMerge) ActivationEvents() []types.EventType {
	return []types.EventType{types.CHUNK_MERGE}
}

// OnEvent processes the CHUNK_MERGE event to merge search result chunks
func (p *PluginMerge) OnEvent(ctx context.Context,
	eventType types.EventType, chatManage *types.ChatManage, next func() *PluginError,
) *PluginError {
	logger.Info(ctx, "Starting chunk merge process")

	// Use rerank results if available, fallback to search results
	searchResult := chatManage.RerankResult
	if len(searchResult) == 0 {
		logger.Info(ctx, "No rerank results available, using search results")
		searchResult = chatManage.SearchResult
	}

	logger.Infof(ctx, "Processing %d chunks for merging", len(searchResult))

	if len(searchResult) == 0 {
		logger.Info(ctx, "No chunks available for merging")
		return next()
	}

	// Group chunks by their knowledge source ID
	knowledgeGroup := make(map[string][]*types.SearchResult)
	for _, chunk := range searchResult {
		knowledgeGroup[chunk.KnowledgeID] = append(knowledgeGroup[chunk.KnowledgeID], chunk)
	}

	logger.Infof(ctx, "Grouped chunks by knowledge ID, %d knowledge sources", len(knowledgeGroup))

	mergedChunks := []*types.SearchResult{}
	// Process each knowledge source separately
	for knowledgeID, chunks := range knowledgeGroup {
		logger.Infof(ctx, "Processing knowledge ID: %s with %d chunks", knowledgeID, len(chunks))

		// Sort chunks by their start position in the original document
		sort.Slice(chunks, func(i, j int) bool {
			if chunks[i].StartAt == chunks[j].StartAt {
				return chunks[i].EndAt < chunks[j].EndAt
			}
			return chunks[i].StartAt < chunks[j].StartAt
		})

		// Merge overlapping or adjacent chunks
		knowledgeMergedChunks := []*types.SearchResult{}
		if chunks[0].ChunkType == types.ChunkTypeSummary {
			knowledgeMergedChunks = append(knowledgeMergedChunks, chunks[0])
			// skip the first chunk if it is summary chunk
			// This is to avoid merging the summary chunk with the first content chunk
			chunks = chunks[1:]
		}
		if len(chunks) > 0 {
			knowledgeMergedChunks = append(knowledgeMergedChunks, chunks[0])
		}
		for i := 1; i < len(chunks); i++ {
			lastChunk := knowledgeMergedChunks[len(knowledgeMergedChunks)-1]
			// If the current chunk starts after the last chunk ends, add it to the merged chunks
			if chunks[i].StartAt > lastChunk.EndAt {
				knowledgeMergedChunks = append(knowledgeMergedChunks, chunks[i])
				continue
			}
			// Merge overlapping chunks
			if chunks[i].EndAt > lastChunk.EndAt {
				lastChunk.Content = lastChunk.Content +
					string([]rune(chunks[i].Content)[lastChunk.EndAt-chunks[i].StartAt:])
				lastChunk.EndAt = chunks[i].EndAt
				lastChunk.SubChunkID = append(lastChunk.SubChunkID, chunks[i].ID)

				// 合并 ImageInfo
				if err := mergeImageInfo(ctx, lastChunk, chunks[i]); err != nil {
					logger.Warnf(ctx, "Failed to merge ImageInfo: %v", err)
				}
			}
			if chunks[i].Score > lastChunk.Score {
				lastChunk.Score = chunks[i].Score
			}
		}

		logger.Infof(ctx, "Merged %d chunks into %d chunks for knowledge ID: %s",
			len(chunks), len(knowledgeMergedChunks), knowledgeID)

		mergedChunks = append(mergedChunks, knowledgeMergedChunks...)
	}

	// Sort merged chunks by their score (highest first)
	sort.Slice(mergedChunks, func(i, j int) bool {
		return mergedChunks[i].Score > mergedChunks[j].Score
	})

	logger.Infof(ctx, "Final merged result: %d chunks, sorted by score", len(mergedChunks))

	chatManage.MergeResult = mergedChunks
	return next()
}

// mergeImageInfo 合并两个chunk的ImageInfo
func mergeImageInfo(ctx context.Context, target *types.SearchResult, source *types.SearchResult) error {
	// 如果source没有ImageInfo，不需要合并
	if source.ImageInfo == "" {
		return nil
	}

	var sourceImageInfos []types.ImageInfo
	if err := json.Unmarshal([]byte(source.ImageInfo), &sourceImageInfos); err != nil {
		logger.Warnf(ctx, "Failed to unmarshal source ImageInfo: %v", err)
		return err
	}

	// 如果source的ImageInfo为空，不需要合并
	if len(sourceImageInfos) == 0 {
		return nil
	}

	// 处理target的ImageInfo
	var targetImageInfos []types.ImageInfo
	if target.ImageInfo != "" {
		if err := json.Unmarshal([]byte(target.ImageInfo), &targetImageInfos); err != nil {
			logger.Warnf(ctx, "Failed to unmarshal target ImageInfo: %v", err)
			// 如果目标解析失败，直接使用源数据
			target.ImageInfo = source.ImageInfo
			return nil
		}
	}

	// 合并ImageInfo
	targetImageInfos = append(targetImageInfos, sourceImageInfos...)

	// 去重
	uniqueMap := make(map[string]bool)
	uniqueImageInfos := make([]types.ImageInfo, 0, len(targetImageInfos))

	for _, imgInfo := range targetImageInfos {
		// 使用URL作为唯一标识
		if imgInfo.URL != "" && !uniqueMap[imgInfo.URL] {
			uniqueMap[imgInfo.URL] = true
			uniqueImageInfos = append(uniqueImageInfos, imgInfo)
		}
	}

	// 序列化合并后的ImageInfo
	mergedImageInfoJSON, err := json.Marshal(uniqueImageInfos)
	if err != nil {
		logger.Warnf(ctx, "Failed to marshal merged ImageInfo: %v", err)
		return err
	}

	// 更新目标chunk的ImageInfo
	target.ImageInfo = string(mergedImageInfoJSON)
	logger.Infof(ctx, "Successfully merged ImageInfo, total count: %d", len(uniqueImageInfos))
	return nil
}

package elasticsearch

import (
	"maps"
	"slices"

	"github.com/Tencent/WeKnora/internal/types"
)

// VectorEmbedding defines the Elasticsearch document structure for vector embeddings
type VectorEmbedding struct {
	Content         string    `json:"content" gorm:"column:content;not null"`            // Text content of the chunk
	SourceID        string    `json:"source_id" gorm:"column:source_id;not null"`        // ID of the source document
	SourceType      int       `json:"source_type" gorm:"column:source_type;not null"`    // Type of the source document
	ChunkID         string    `json:"chunk_id" gorm:"column:chunk_id"`                   // Unique ID of the text chunk
	KnowledgeID     string    `json:"knowledge_id" gorm:"column:knowledge_id"`           // ID of the knowledge item
	KnowledgeBaseID string    `json:"knowledge_base_id" gorm:"column:knowledge_base_id"` // ID of the knowledge base
	Embedding       []float32 `json:"embedding" gorm:"column:embedding;not null"`        // Vector embedding of the content
}

// VectorEmbeddingWithScore extends VectorEmbedding with similarity score
type VectorEmbeddingWithScore struct {
	VectorEmbedding
	Score float64 // Similarity score from vector search
}

// ToDBVectorEmbedding converts IndexInfo to Elasticsearch document format
func ToDBVectorEmbedding(embedding *types.IndexInfo, additionalParams map[string]interface{}) *VectorEmbedding {
	vector := &VectorEmbedding{
		Content:         embedding.Content,
		SourceID:        embedding.SourceID,
		SourceType:      int(embedding.SourceType),
		ChunkID:         embedding.ChunkID,
		KnowledgeID:     embedding.KnowledgeID,
		KnowledgeBaseID: embedding.KnowledgeBaseID,
	}
	// Add embedding data if available in additionalParams
	if additionalParams != nil && slices.Contains(slices.Collect(maps.Keys(additionalParams)), "embedding") {
		if embeddingMap, ok := additionalParams["embedding"].(map[string][]float32); ok {
			vector.Embedding = embeddingMap[embedding.SourceID]
		}
	}
	return vector
}

// FromDBVectorEmbeddingWithScore converts Elasticsearch document to IndexWithScore domain model
func FromDBVectorEmbeddingWithScore(id string,
	embedding *VectorEmbeddingWithScore,
	matchType types.MatchType,
) *types.IndexWithScore {
	return &types.IndexWithScore{
		ID:              id,
		SourceID:        embedding.SourceID,
		SourceType:      types.SourceType(embedding.SourceType),
		ChunkID:         embedding.ChunkID,
		KnowledgeID:     embedding.KnowledgeID,
		KnowledgeBaseID: embedding.KnowledgeBaseID,
		Content:         embedding.Content,
		Score:           embedding.Score,
		MatchType:       matchType,
	}
}

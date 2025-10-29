package types

// RetrieverEngineType represents the type of retriever engine
type RetrieverEngineType string

// RetrieverEngineType constants
const (
	PostgresRetrieverEngineType      RetrieverEngineType = "postgres"
	ElasticsearchRetrieverEngineType RetrieverEngineType = "elasticsearch"
	InfinityRetrieverEngineType      RetrieverEngineType = "infinity"
	ElasticFaissRetrieverEngineType  RetrieverEngineType = "elasticfaiss"
)

// RetrieverType represents the type of retriever
type RetrieverType string

// RetrieverType constants
const (
	KeywordsRetrieverType  RetrieverType = "keywords"  // Keywords retriever
	VectorRetrieverType    RetrieverType = "vector"    // Vector retriever
	WebSearchRetrieverType RetrieverType = "websearch" // Web search retriever
)

// RetrieveParams represents the parameters for retrieval
type RetrieveParams struct {
	// Query text
	Query string
	// Query embedding (used for vector retrieval)
	Embedding []float32
	// Knowledge base IDs
	KnowledgeBaseIDs []string
	// Excluded knowledge IDs
	ExcludeKnowledgeIDs []string
	// Excluded chunk IDs
	ExcludeChunkIDs []string
	// Number of results to return
	TopK int
	// Similarity threshold
	Threshold float64
	// Additional parameters, different retrievers may require different parameters
	AdditionalParams map[string]interface{}
	// Retriever type
	RetrieverType RetrieverType // Retriever type
}

// RetrieverEngineParams represents the parameters for retriever engine
type RetrieverEngineParams struct {
	// Retriever engine type
	RetrieverEngineType RetrieverEngineType `yaml:"retriever_engine_type" json:"retriever_engine_type"`
	// Retriever type
	RetrieverType RetrieverType `yaml:"retriever_type" json:"retriever_type"`
}

// IndexWithScore represents the index with score
type IndexWithScore struct {
	// ID
	ID string
	// Content
	Content string
	// Source ID
	SourceID string
	// Source type
	SourceType SourceType
	// Chunk ID
	ChunkID string
	// Knowledge ID
	KnowledgeID string
	// Knowledge base ID
	KnowledgeBaseID string
	// Score
	Score float64
	// Match type
	MatchType MatchType
}

// RetrieveResult represents the result of retrieval
type RetrieveResult struct {
	Results             []*IndexWithScore   // Retrieval results
	RetrieverEngineType RetrieverEngineType // Retrieval source type
	RetrieverType       RetrieverType       // Retrieval type
	Error               error               // Retrieval error
}

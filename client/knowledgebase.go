// Package client provides the implementation for interacting with the WeKnora API
// The KnowledgeBase related interfaces are used to manage knowledge bases
// Knowledge bases are collections of knowledge entries that can be used for question-answering
// They can also be searched and queried using hybrid search
package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// KnowledgeBase represents a knowledge base
type KnowledgeBase struct {
	ID                    string                `json:"id"`
	Name                  string                `json:"name"` // Name must be unique within the same tenant
	Description           string                `json:"description"`
	TenantID              uint                  `json:"tenant_id"` // Changed to uint type
	ChunkingConfig        ChunkingConfig        `json:"chunking_config"`
	ImageProcessingConfig ImageProcessingConfig `json:"image_processing_config"`
	EmbeddingModelID      string                `json:"embedding_model_id"`
	SummaryModelID        string                `json:"summary_model_id"` // Summary model ID
	CreatedAt             time.Time             `json:"created_at"`
	UpdatedAt             time.Time             `json:"updated_at"`
}

// KnowledgeBaseConfig represents knowledge base configuration
type KnowledgeBaseConfig struct {
	ChunkingConfig        ChunkingConfig        `json:"chunking_config"`
	ImageProcessingConfig ImageProcessingConfig `json:"image_processing_config"`
}

// ChunkingConfig represents document chunking configuration
type ChunkingConfig struct {
	ChunkSize        int      `json:"chunk_size"`        // Chunk size
	ChunkOverlap     int      `json:"chunk_overlap"`     // Overlap size
	Separators       []string `json:"separators"`        // Separators
	EnableMultimodal bool     `json:"enable_multimodal"` // Whether to enable multimodal processing
}

// ImageProcessingConfig represents image processing configuration
type ImageProcessingConfig struct {
	ModelID string `json:"model_id"` // Multimodal model ID
}

// KnowledgeBaseResponse knowledge base response
type KnowledgeBaseResponse struct {
	Success bool          `json:"success"`
	Data    KnowledgeBase `json:"data"`
}

// KnowledgeBaseListResponse knowledge base list response
type KnowledgeBaseListResponse struct {
	Success bool            `json:"success"`
	Data    []KnowledgeBase `json:"data"`
}

// SearchResult represents search result
type SearchResult struct {
	ID                string            `json:"id"`
	Content           string            `json:"content"`
	KnowledgeID       string            `json:"knowledge_id"`
	ChunkIndex        int               `json:"chunk_index"`
	KnowledgeTitle    string            `json:"knowledge_title"`
	StartAt           int               `json:"start_at"`
	EndAt             int               `json:"end_at"`
	Seq               int               `json:"seq"`
	Score             float64           `json:"score"`
	ChunkType         string            `json:"chunk_type"`
	ImageInfo         string            `json:"image_info"`
	Metadata          map[string]string `json:"metadata"`
	KnowledgeFilename string            `json:"knowledge_filename"`
	KnowledgeSource   string            `json:"knowledge_source"`
}

// HybridSearchResponse hybrid search response
type HybridSearchResponse struct {
	Success bool            `json:"success"`
	Data    []*SearchResult `json:"data"`
}

type CopyKnowledgeBaseRequest struct {
	SourceID string `json:"source_id"`
	TargetID string `json:"target_id"`
}

// CreateKnowledgeBase creates a knowledge base
func (c *Client) CreateKnowledgeBase(ctx context.Context, knowledgeBase *KnowledgeBase) (*KnowledgeBase, error) {
	resp, err := c.doRequest(ctx, http.MethodPost, "/api/v1/knowledge-bases", knowledgeBase, nil)
	if err != nil {
		return nil, err
	}

	var response KnowledgeBaseResponse
	if err := parseResponse(resp, &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

// GetKnowledgeBase gets a knowledge base
func (c *Client) GetKnowledgeBase(ctx context.Context, knowledgeBaseID string) (*KnowledgeBase, error) {
	path := fmt.Sprintf("/api/v1/knowledge-bases/%s", knowledgeBaseID)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}

	var response KnowledgeBaseResponse
	if err := parseResponse(resp, &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

// ListKnowledgeBases lists knowledge bases
func (c *Client) ListKnowledgeBases(ctx context.Context) ([]KnowledgeBase, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "/api/v1/knowledge-bases", nil, nil)
	if err != nil {
		return nil, err
	}

	var response KnowledgeBaseListResponse
	if err := parseResponse(resp, &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

// UpdateKnowledgeBaseRequest update knowledge base request
type UpdateKnowledgeBaseRequest struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Config      *KnowledgeBaseConfig `json:"config"`
}

// UpdateKnowledgeBase updates a knowledge base
func (c *Client) UpdateKnowledgeBase(ctx context.Context,
	knowledgeBaseID string,
	request *UpdateKnowledgeBaseRequest,
) (*KnowledgeBase, error) {
	path := fmt.Sprintf("/api/v1/knowledge-bases/%s", knowledgeBaseID)
	resp, err := c.doRequest(ctx, http.MethodPut, path, request, nil)
	if err != nil {
		return nil, err
	}

	var response KnowledgeBaseResponse
	if err := parseResponse(resp, &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

// DeleteKnowledgeBase deletes a knowledge base
func (c *Client) DeleteKnowledgeBase(ctx context.Context, knowledgeBaseID string) error {
	path := fmt.Sprintf("/api/v1/knowledge-bases/%s", knowledgeBaseID)
	resp, err := c.doRequest(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return err
	}

	var response struct {
		Success bool   `json:"success"`
		Message string `json:"message,omitempty"`
	}

	return parseResponse(resp, &response)
}

// HybridSearch performs hybrid search
func (c *Client) HybridSearch(ctx context.Context, knowledgeBaseID string, query string) ([]*SearchResult, error) {
	path := fmt.Sprintf("/api/v1/knowledge-bases/%s/hybrid-search", knowledgeBaseID)

	queryParams := url.Values{}
	queryParams.Add("query", query)

	resp, err := c.doRequest(ctx, http.MethodGet, path, nil, queryParams)
	if err != nil {
		return nil, err
	}

	var response HybridSearchResponse
	if err := parseResponse(resp, &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

func (c *Client) CopyKnowledgeBase(ctx context.Context, request *CopyKnowledgeBaseRequest) error {
	path := "/api/v1/knowledge-bases/copy"

	resp, err := c.doRequest(ctx, http.MethodPost, path, request, nil)
	if err != nil {
		return err
	}

	var response struct {
		Success bool   `json:"success"`
		Message string `json:"message,omitempty"`
	}

	return parseResponse(resp, &response)
}

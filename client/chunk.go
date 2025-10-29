// Package client provides the implementation for interacting with the WeKnora API
// This package encapsulates CRUD operations for server resources and provides a friendly interface for callers
// The Chunk related interfaces are used to manage document chunks in the knowledge base
package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// Chunk represents the information about a document chunk
// Chunks are the basic units of storage and indexing in the knowledge base
type Chunk struct {
	ID              string    `json:"id"`               // Unique identifier of the chunk
	KnowledgeID     string    `json:"knowledge_id"`     // Identifier of the parent knowledge
	TenantID        uint      `json:"tenant_id"`        // Tenant ID
	Content         string    `json:"content"`          // Text content of the chunk
	Embedding       []float32 `json:"embedding"`        // Vector embedding representation
	ChunkIndex      int       `json:"chunk_index"`      // Index position of chunk in the document
	TotalChunks     int       `json:"total_chunks"`     // Total number of chunks in the document
	IsEnabled       bool      `json:"is_enabled"`       // Whether this chunk is enabled
	StartAt         int       `json:"start_at"`         // Starting position in original text
	EndAt           int       `json:"end_at"`           // Ending position in original text
	VectorStoreID   string    `json:"vector_store_id"`  // Vector storage ID
	KeywordStoreID  string    `json:"keyword_store_id"` // Keyword storage ID
	EmbeddingStatus int       `json:"embedding_status"` // Embedding status: 0-unprocessed, 1-processing, 2-completed
	ChunkType       string    `json:"chunk_type"`
	ImageInfo       string    `json:"image_info"`
	CreatedAt       string    `json:"created_at"` // Creation time
	UpdatedAt       string    `json:"updated_at"` // Last update time
}

// ChunkResponse represents the response for a single chunk
// API response structure containing a single chunk information
type ChunkResponse struct {
	Success bool  `json:"success"` // Whether operation was successful
	Data    Chunk `json:"data"`    // Chunk data
}

// ChunkListResponse represents the response for a list of chunks
// API response structure for returning a list of chunks
type ChunkListResponse struct {
	Success  bool    `json:"success"`   // Whether operation was successful
	Data     []Chunk `json:"data"`      // List of chunks
	Total    int64   `json:"total"`     // Total count
	Page     int     `json:"page"`      // Current page
	PageSize int     `json:"page_size"` // Items per page
}

// UpdateChunkRequest represents the request structure for updating a chunk
// Used for requesting chunk information updates
type UpdateChunkRequest struct {
	Content    string    `json:"content"`     // Chunk content
	Embedding  []float32 `json:"embedding"`   // Vector embedding
	ChunkIndex int       `json:"chunk_index"` // Chunk index
	IsEnabled  bool      `json:"is_enabled"`  // Whether enabled
	StartAt    int       `json:"start_at"`    // Start position
	EndAt      int       `json:"end_at"`      // End position
}

// ListKnowledgeChunks lists all chunks under a knowledge document
// Queries all chunks by knowledge ID with pagination support
// Parameters:
//   - ctx: Context
//   - knowledgeID: Knowledge ID
//   - page: Page number, starts from 1
//   - pageSize: Number of items per page
//
// Returns:
//   - []Chunk: List of chunks
//   - int64: Total count
//   - error: Error information
func (c *Client) ListKnowledgeChunks(ctx context.Context,
	knowledgeID string, page int, pageSize int,
) ([]Chunk, int64, error) {
	path := fmt.Sprintf("/api/v1/chunks/%s", knowledgeID)

	queryParams := url.Values{}
	queryParams.Add("page", strconv.Itoa(page))
	queryParams.Add("page_size", strconv.Itoa(pageSize))

	resp, err := c.doRequest(ctx, http.MethodGet, path, nil, queryParams)
	if err != nil {
		return nil, 0, err
	}

	var response ChunkListResponse
	if err := parseResponse(resp, &response); err != nil {
		return nil, 0, err
	}

	return response.Data, response.Total, nil
}

// UpdateChunk updates a chunk's information
// Updates information for a specific chunk under a knowledge document
// Parameters:
//   - ctx: Context
//   - knowledgeID: Knowledge ID
//   - chunkID: Chunk ID
//   - request: Update request
//
// Returns:
//   - *Chunk: Updated chunk
//   - error: Error information
func (c *Client) UpdateChunk(ctx context.Context,
	knowledgeID string, chunkID string, request *UpdateChunkRequest,
) (*Chunk, error) {
	path := fmt.Sprintf("/api/v1/chunks/%s/%s", knowledgeID, chunkID)
	resp, err := c.doRequest(ctx, http.MethodPut, path, request, nil)
	if err != nil {
		return nil, err
	}

	var response ChunkResponse
	if err := parseResponse(resp, &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

// DeleteChunk deletes a specific chunk
// Deletes a specific chunk under a knowledge document
// Parameters:
//   - ctx: Context
//   - knowledgeID: Knowledge ID
//   - chunkID: Chunk ID
//
// Returns:
//   - error: Error information
func (c *Client) DeleteChunk(ctx context.Context, knowledgeID string, chunkID string) error {
	path := fmt.Sprintf("/api/v1/chunks/%s/%s", knowledgeID, chunkID)
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

// DeleteChunksByKnowledgeID deletes all chunks under a knowledge document
// Batch deletes all chunks under the specified knowledge document
// Parameters:
//   - ctx: Context
//   - knowledgeID: Knowledge ID
//
// Returns:
//   - error: Error information
func (c *Client) DeleteChunksByKnowledgeID(ctx context.Context, knowledgeID string) error {
	path := fmt.Sprintf("/api/v1/chunks/%s", knowledgeID)
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

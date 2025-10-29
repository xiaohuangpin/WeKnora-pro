// Package service provides business logic implementations for WeKnora application
// This package contains service layer implementations that coordinate between
// repositories and handlers, applying business rules and transaction management
package service

import (
	"context"

	"github.com/Tencent/WeKnora/internal/logger"
	"github.com/Tencent/WeKnora/internal/types"
	"github.com/Tencent/WeKnora/internal/types/interfaces"
)

// chunkService implements the ChunkService interface
// It provides operations for managing document chunks in the knowledge base
// Chunks are segments of documents that have been processed and prepared for indexing
type chunkService struct {
	chunkRepository interfaces.ChunkRepository // Repository for chunk data persistence
	kbRepository    interfaces.KnowledgeBaseRepository
	modelService    interfaces.ModelService
}

// NewChunkService creates a new chunk service
// It initializes a service with the provided chunk repository
// Parameters:
//   - chunkRepository: Repository for chunk operations
//
// Returns:
//   - interfaces.ChunkService: Initialized chunk service implementation
func NewChunkService(
	chunkRepository interfaces.ChunkRepository,
	kbRepository interfaces.KnowledgeBaseRepository,
	modelService interfaces.ModelService,
) interfaces.ChunkService {
	return &chunkService{
		chunkRepository: chunkRepository,
		kbRepository:    kbRepository,
		modelService:    modelService,
	}
}

// CreateChunks creates multiple chunks
// This method persists a batch of document chunks to the repository
// Parameters:
//   - ctx: Context with authentication and request information
//   - chunks: Slice of document chunks to create
//
// Returns:
//   - error: Any error encountered during chunk creation
func (s *chunkService) CreateChunks(ctx context.Context, chunks []*types.Chunk) error {
	logger.Info(ctx, "Start creating chunks")
	logger.Infof(ctx, "Creating %d chunks", len(chunks))

	err := s.chunkRepository.CreateChunks(ctx, chunks)
	if err != nil {
		logger.ErrorWithFields(ctx, err, map[string]interface{}{
			"chunk_count": len(chunks),
		})
		return err
	}

	logger.Info(ctx, "Chunks created successfully")
	return nil
}

// GetChunkByID retrieves a chunk by its ID
// This method fetches a specific chunk using its ID and validates tenant access
// Parameters:
//   - ctx: Context with authentication and request information
//   - knowledgeID: ID of the knowledge document containing the chunk
//   - id: ID of the chunk to retrieve
//
// Returns:
//   - *types.Chunk: Retrieved chunk if found
//   - error: Any error encountered during retrieval
func (s *chunkService) GetChunkByID(ctx context.Context, knowledgeID string, id string) (*types.Chunk, error) {
	logger.Info(ctx, "Start getting chunk by ID")
	logger.Infof(ctx, "Getting chunk, ID: %s, knowledge ID: %s", id, knowledgeID)

	tenantID := ctx.Value(types.TenantIDContextKey).(uint)
	logger.Infof(ctx, "Tenant ID: %d", tenantID)

	chunk, err := s.chunkRepository.GetChunkByID(ctx, tenantID, id)
	if err != nil {
		logger.ErrorWithFields(ctx, err, map[string]interface{}{
			"chunk_id":     id,
			"knowledge_id": knowledgeID,
			"tenant_id":    tenantID,
		})
		return nil, err
	}

	logger.Info(ctx, "Chunk retrieved successfully")
	return chunk, nil
}

// ListChunksByKnowledgeID lists all chunks for a knowledge ID
// This method retrieves all chunks belonging to a specific knowledge document
// Parameters:
//   - ctx: Context with authentication and request information
//   - knowledgeID: ID of the knowledge document
//
// Returns:
//   - []*types.Chunk: List of chunks belonging to the knowledge document
//   - error: Any error encountered during retrieval
func (s *chunkService) ListChunksByKnowledgeID(ctx context.Context, knowledgeID string) ([]*types.Chunk, error) {
	logger.Info(ctx, "Start listing chunks by knowledge ID")
	logger.Infof(ctx, "Knowledge ID: %s", knowledgeID)

	tenantID := ctx.Value(types.TenantIDContextKey).(uint)
	logger.Infof(ctx, "Tenant ID: %d", tenantID)

	chunks, err := s.chunkRepository.ListChunksByKnowledgeID(ctx, tenantID, knowledgeID)
	if err != nil {
		logger.ErrorWithFields(ctx, err, map[string]interface{}{
			"knowledge_id": knowledgeID,
			"tenant_id":    tenantID,
		})
		return nil, err
	}

	logger.Infof(ctx, "Retrieved %d chunks successfully", len(chunks))
	return chunks, nil
}

// ListPagedChunksByKnowledgeID lists chunks for a knowledge ID with pagination
// This method retrieves chunks with pagination support for better performance with large datasets
// Parameters:
//   - ctx: Context with authentication and request information
//   - knowledgeID: ID of the knowledge document
//   - page: Pagination parameters including page number and page size
//
// Returns:
//   - *types.PageResult: Paginated result containing chunks and pagination metadata
//   - error: Any error encountered during retrieval
func (s *chunkService) ListPagedChunksByKnowledgeID(ctx context.Context,
	knowledgeID string, page *types.Pagination,
) (*types.PageResult, error) {
	logger.Info(ctx, "Start listing paged chunks by knowledge ID")
	logger.Infof(ctx, "Knowledge ID: %s, page: %d, page size: %d", knowledgeID, page.Page, page.PageSize)

	tenantID := ctx.Value(types.TenantIDContextKey).(uint)
	logger.Infof(ctx, "Tenant ID: %d", tenantID)
	chunkType := []types.ChunkType{types.ChunkTypeText}
	chunks, total, err := s.chunkRepository.ListPagedChunksByKnowledgeID(ctx, tenantID, knowledgeID, page, chunkType)
	if err != nil {
		logger.ErrorWithFields(ctx, err, map[string]interface{}{
			"knowledge_id": knowledgeID,
			"tenant_id":    tenantID,
			"page":         page.Page,
			"page_size":    page.PageSize,
		})
		return nil, err
	}

	logger.Infof(ctx, "Retrieved %d chunks out of %d total chunks", len(chunks), total)
	return types.NewPageResult(total, page, chunks), nil
}

// updateChunk updates a chunk
// This method updates an existing chunk in the repository
// Parameters:
//   - ctx: Context with authentication and request information
//   - chunk: Chunk with updated fields
//
// Returns:
//   - error: Any error encountered during update
//
// This method handles the actual update logic for a chunk, including updating the vector database representation
func (s *chunkService) UpdateChunk(ctx context.Context, chunk *types.Chunk) error {
	logger.Infof(ctx, "Updating chunk, ID: %s, knowledge ID: %s", chunk.ID, chunk.KnowledgeID)

	// Update the chunk in the repository
	err := s.chunkRepository.UpdateChunk(ctx, chunk)
	if err != nil {
		logger.ErrorWithFields(ctx, err, map[string]interface{}{
			"chunk_id":     chunk.ID,
			"knowledge_id": chunk.KnowledgeID,
		})
		return err
	}

	logger.Info(ctx, "Chunk updated successfully")
	return nil
}

// DeleteChunk deletes a chunk by ID
// This method removes a specific chunk from the repository
// Parameters:
//   - ctx: Context with authentication and request information
//   - id: ID of the chunk to delete
//
// Returns:
//   - error: Any error encountered during deletion
func (s *chunkService) DeleteChunk(ctx context.Context, id string) error {
	logger.Info(ctx, "Start deleting chunk")
	logger.Infof(ctx, "Deleting chunk, ID: %s", id)

	tenantID := ctx.Value(types.TenantIDContextKey).(uint)
	logger.Infof(ctx, "Tenant ID: %d", tenantID)

	err := s.chunkRepository.DeleteChunk(ctx, tenantID, id)
	if err != nil {
		logger.ErrorWithFields(ctx, err, map[string]interface{}{
			"chunk_id":  id,
			"tenant_id": tenantID,
		})
		return err
	}

	logger.Info(ctx, "Chunk deleted successfully")
	return nil
}

// DeleteChunksByKnowledgeID deletes all chunks for a knowledge ID
// This method removes all chunks belonging to a specific knowledge document
// Parameters:
//   - ctx: Context with authentication and request information
//   - knowledgeID: ID of the knowledge document
//
// Returns:
//   - error: Any error encountered during bulk deletion
func (s *chunkService) DeleteChunksByKnowledgeID(ctx context.Context, knowledgeID string) error {
	logger.Info(ctx, "Start deleting all chunks by knowledge ID")
	logger.Infof(ctx, "Knowledge ID: %s", knowledgeID)

	tenantID := ctx.Value(types.TenantIDContextKey).(uint)
	logger.Infof(ctx, "Tenant ID: %d", tenantID)

	err := s.chunkRepository.DeleteChunksByKnowledgeID(ctx, tenantID, knowledgeID)
	if err != nil {
		logger.ErrorWithFields(ctx, err, map[string]interface{}{
			"knowledge_id": knowledgeID,
			"tenant_id":    tenantID,
		})
		return err
	}

	logger.Info(ctx, "All chunks under knowledge deleted successfully")
	return nil
}

func (s *chunkService) DeleteByKnowledgeList(ctx context.Context, ids []string) error {
	logger.Info(ctx, "Start deleting all chunks by knowledge IDs")
	logger.Infof(ctx, "Knowledge IDs: %v", ids)

	tenantID := ctx.Value(types.TenantIDContextKey).(uint)
	logger.Infof(ctx, "Tenant ID: %d", tenantID)

	err := s.chunkRepository.DeleteByKnowledgeList(ctx, tenantID, ids)
	if err != nil {
		logger.ErrorWithFields(ctx, err, map[string]interface{}{
			"knowledge_id": ids,
			"tenant_id":    tenantID,
		})
		return err
	}

	logger.Info(ctx, "All chunks under knowledge deleted successfully")
	return nil
}

func (s *chunkService) ListChunkByParentID(ctx context.Context, tenantID uint, parentID string) ([]*types.Chunk, error) {
	logger.Info(ctx, "Start listing chunk by parent ID")
	logger.Infof(ctx, "Parent ID: %s", parentID)

	chunks, err := s.chunkRepository.ListChunkByParentID(ctx, tenantID, parentID)
	if err != nil {
		logger.ErrorWithFields(ctx, err, map[string]interface{}{
			"parent_id": parentID,
			"tenant_id": tenantID,
		})
		return nil, err
	}

	logger.Info(ctx, "Chunk listed successfully")
	return chunks, nil
}

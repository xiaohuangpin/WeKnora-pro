package interfaces

import (
	"context"

	"github.com/Tencent/WeKnora/internal/types"
)

// SessionService defines the session service interface
type SessionService interface {
	// CreateSession creates a session
	CreateSession(ctx context.Context, session *types.Session) (*types.Session, error)
	// GetSession gets a session
	GetSession(ctx context.Context, id string) (*types.Session, error)
	// GetSessionsByTenant gets all sessions of a tenant
	GetSessionsByTenant(ctx context.Context) ([]*types.Session, error)
	// GetPagedSessionsByTenant gets paged sessions of a tenant
	GetPagedSessionsByTenant(ctx context.Context, page *types.Pagination) (*types.PageResult, error)
	// UpdateSession updates a session
	UpdateSession(ctx context.Context, session *types.Session) error
	// DeleteSession deletes a session
	DeleteSession(ctx context.Context, id string) error
	// GenerateTitle generates a title for the current conversation
	GenerateTitle(ctx context.Context, sessionID string, messages []types.Message) (string, error)
	// KnowledgeQA performs knowledge-based question answering
	KnowledgeQA(ctx context.Context,
		sessionID, query string,
	) ([]*types.SearchResult, <-chan types.StreamResponse, error)
	// KnowledgeQAByEvent performs knowledge-based question answering by event
	KnowledgeQAByEvent(ctx context.Context, chatManage *types.ChatManage, eventList []types.EventType) error
	// SearchKnowledge performs knowledge-based search, without summarization
	SearchKnowledge(ctx context.Context, knowledgeBaseID, query string) ([]*types.SearchResult, error)
}

// SessionRepository defines the session repository interface
type SessionRepository interface {
	// Create creates a session
	Create(ctx context.Context, session *types.Session) (*types.Session, error)
	// Get gets a session
	Get(ctx context.Context, tenantID uint, id string) (*types.Session, error)
	// GetByTenantID gets all sessions of a tenant
	GetByTenantID(ctx context.Context, tenantID uint) ([]*types.Session, error)
	// GetPagedByTenantID gets paged sessions of a tenant
	GetPagedByTenantID(ctx context.Context, tenantID uint, page *types.Pagination) ([]*types.Session, int64, error)
	// Update updates a session
	Update(ctx context.Context, session *types.Session) error
	// Delete deletes a session
	Delete(ctx context.Context, tenantID uint, id string) error
}

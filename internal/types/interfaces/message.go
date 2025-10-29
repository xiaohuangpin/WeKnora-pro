package interfaces

import (
	"context"
	"time"

	"github.com/Tencent/WeKnora/internal/types"
)

// MessageService defines the message service interface
type MessageService interface {
	// CreateMessage creates a message
	CreateMessage(ctx context.Context, message *types.Message) (*types.Message, error)

	// GetMessage gets a message
	GetMessage(ctx context.Context, sessionID string, id string) (*types.Message, error)

	// GetMessagesBySession gets all messages of a session
	GetMessagesBySession(ctx context.Context, sessionID string, page int, pageSize int) ([]*types.Message, error)

	// GetRecentMessagesBySession gets recent messages of a session
	GetRecentMessagesBySession(ctx context.Context, sessionID string, limit int) ([]*types.Message, error)

	// GetMessagesBySessionBeforeTime gets messages before a specific time of a session
	GetMessagesBySessionBeforeTime(
		ctx context.Context, sessionID string, beforeTime time.Time, limit int,
	) ([]*types.Message, error)

	// UpdateMessage updates a message
	UpdateMessage(ctx context.Context, message *types.Message) error

	// DeleteMessage deletes a message
	DeleteMessage(ctx context.Context, sessionID string, id string) error
}

// MessageRepository defines the message repository interface
type MessageRepository interface {
	MessageService
	// GetFirstMessageOfUser gets the first message of a user
	GetFirstMessageOfUser(ctx context.Context, sessionID string) (*types.Message, error)
}

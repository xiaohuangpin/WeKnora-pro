package interfaces

import (
	"context"
	"time"

	"github.com/Tencent/WeKnora/internal/types"
)

// StreamInfo stream information
type StreamInfo struct {
	SessionID           string           // session ID
	RequestID           string           // request ID
	Query               string           // query content
	Content             string           // current content
	KnowledgeReferences types.References // knowledge references
	LastUpdated         time.Time        // last updated time
	IsCompleted         bool             // whether completed
}

// StreamManager stream manager interface
type StreamManager interface {
	// RegisterStream registers a new stream
	RegisterStream(ctx context.Context, sessionID, requestID, query string) error

	// UpdateStream updates the stream content
	UpdateStream(ctx context.Context, sessionID, requestID string, content string, references types.References) error

	// CompleteStream completes the stream
	CompleteStream(ctx context.Context, sessionID, requestID string) error

	// GetStream gets a specific stream
	GetStream(ctx context.Context, sessionID, requestID string) (*StreamInfo, error)
}

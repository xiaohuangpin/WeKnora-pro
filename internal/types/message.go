// Package types defines data structures and types used throughout the system
package types

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// History represents a conversation history entry
// Contains query-answer pairs and associated knowledge references
// Used for tracking conversation context and history
type History struct {
	Query               string     // User query text
	Answer              string     // System response text
	CreateAt            time.Time  // When this history entry was created
	KnowledgeReferences References // Knowledge references used in the answer
}

// Message represents a conversation message
// Each message belongs to a conversation session and can be from either user or system
// Messages can contain references to knowledge chunks used to generate responses
type Message struct {
	// Unique identifier for the message
	ID string `json:"id" gorm:"type:varchar(36);primaryKey"`
	// ID of the session this message belongs to
	SessionID string `json:"session_id"`
	// Request identifier for tracking API requests
	RequestID string `json:"request_id"`
	// Message text content
	Content string `json:"content"`
	// Message role: "user", "assistant", "system"
	Role string `json:"role"`
	// References to knowledge chunks used in the response
	KnowledgeReferences References `json:"knowledge_references" gorm:"type:json,column:knowledge_references"`
	// Whether message generation is complete
	IsCompleted bool `json:"is_completed"`
	// Message creation timestamp
	CreatedAt time.Time `json:"created_at"`
	// Last update timestamp
	UpdatedAt time.Time `json:"updated_at"`
	// Soft delete timestamp
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// BeforeCreate is a GORM hook that runs before creating a new message record
// Automatically generates a UUID for new messages and initializes knowledge references
// Parameters:
//   - tx: GORM database transaction
//
// Returns:
//   - error: Any error encountered during the hook execution
func (m *Message) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New().String()
	if m.KnowledgeReferences == nil {
		m.KnowledgeReferences = make(References, 0)
	}
	return nil
}

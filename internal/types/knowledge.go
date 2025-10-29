package types

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Knowledge represents a knowledge entity in the system.
// It contains metadata about the knowledge source, its processing status,
// and references to the physical file if applicable.
type Knowledge struct {
	// Unique identifier of the knowledge
	ID string `json:"id" gorm:"type:varchar(36);primaryKey"`
	// Tenant ID
	TenantID uint `json:"tenant_id"`
	// ID of the knowledge base
	KnowledgeBaseID string `json:"knowledge_base_id"`
	// Type of the knowledge
	Type string `json:"type"`
	// Title of the knowledge
	Title string `json:"title"`
	// Description of the knowledge
	Description string `json:"description"`
	// Source of the knowledge
	Source string `json:"source"`
	// Parse status of the knowledge
	ParseStatus string `json:"parse_status"`
	// Enable status of the knowledge
	EnableStatus string `json:"enable_status"`
	// ID of the embedding model
	EmbeddingModelID string `json:"embedding_model_id"`
	// File name of the knowledge
	FileName string `json:"file_name"`
	// File type of the knowledge
	FileType string `json:"file_type"`
	// File size of the knowledge
	FileSize int64 `json:"file_size"`
	// File hash of the knowledge
	FileHash string `json:"file_hash"`
	// File path of the knowledge
	FilePath string `json:"file_path"`
	// Storage size of the knowledge
	StorageSize int64 `json:"storage_size"`
	// Metadata of the knowledge
	Metadata JSON `json:"metadata" gorm:"type:json"`
	// Creation time of the knowledge
	CreatedAt time.Time `json:"created_at"`
	// Last updated time of the knowledge
	UpdatedAt time.Time `json:"updated_at"`
	// Processed time of the knowledge
	ProcessedAt *time.Time `json:"processed_at"`
	// Error message of the knowledge
	ErrorMessage string `json:"error_message"`
	// Deletion time of the knowledge
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// GetMetadata returns the metadata as a map[string]string.
func (k *Knowledge) GetMetadata() map[string]string {
	metadata := make(map[string]string)
	metadataMap, err := k.Metadata.Map()
	if err != nil {
		return nil
	}
	for k, v := range metadataMap {
		metadata[k] = fmt.Sprintf("%v", v)
	}
	return metadata
}

// BeforeCreate hook generates a UUID for new Knowledge entities before they are created.
func (k *Knowledge) BeforeCreate(tx *gorm.DB) (err error) {
	k.ID = uuid.New().String()
	return nil
}

// KnowledgeCheckParams defines parameters used to check if knowledge already exists.
type KnowledgeCheckParams struct {
	// File parameters
	FileName string
	FileSize int64
	FileHash string
	// URL parameters
	URL string
	// Text passage parameters
	Passages []string
	// Knowledge type
	Type string
}

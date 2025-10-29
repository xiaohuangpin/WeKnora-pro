package interfaces

import (
	"context"
	"io"
	"mime/multipart"
)

// FileService is the interface for file services.
// FileService provides methods to save, retrieve, and delete files.
type FileService interface {
	// SaveFile saves a file.
	SaveFile(ctx context.Context, file *multipart.FileHeader, tenantID uint, knowledgeID string) (string, error)
	// GetFile retrieves a file.
	GetFile(ctx context.Context, filePath string) (io.ReadCloser, error)
	// DeleteFile deletes a file.
	DeleteFile(ctx context.Context, filePath string) error
}

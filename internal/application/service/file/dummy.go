package file

import (
	"context"
	"errors"
	"io"
	"mime/multipart"

	"github.com/Tencent/WeKnora/internal/types/interfaces"
	"github.com/google/uuid"
)

// DummyFileService is a no-op implementation of the FileService interface
// used for testing or when file storage is not required
type DummyFileService struct{}

// NewDummyFileService creates a new instance of DummyFileService
func NewDummyFileService() interfaces.FileService {
	return &DummyFileService{}
}

// SaveFile pretends to save a file but just returns a random UUID
// This is useful for testing without actual file operations
func (s *DummyFileService) SaveFile(ctx context.Context,
	file *multipart.FileHeader, tenantID uint, knowledgeID string,
) (string, error) {
	return uuid.New().String(), nil
}

// GetFile always returns an error as dummy service doesn't store files
func (s *DummyFileService) GetFile(ctx context.Context, filePath string) (io.ReadCloser, error) {
	return nil, errors.New("not implemented")
}

// DeleteFile is a no-op operation that always succeeds
func (s *DummyFileService) DeleteFile(ctx context.Context, filePath string) error {
	return nil
}

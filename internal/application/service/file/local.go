package file

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/Tencent/WeKnora/internal/logger"
	"github.com/Tencent/WeKnora/internal/types/interfaces"
)

// localFileService implements the FileService interface for local file system storage
type localFileService struct {
	baseDir string // Base directory for file storage
}

// NewLocalFileService creates a new local file service instance
func NewLocalFileService(baseDir string) interfaces.FileService {
	return &localFileService{
		baseDir: baseDir,
	}
}

// SaveFile stores an uploaded file to the local file system
// The file is stored in a directory structure: baseDir/tenantID/knowledgeID/filename
// Returns the full file path or an error if saving fails
func (s *localFileService) SaveFile(ctx context.Context,
	file *multipart.FileHeader, tenantID uint, knowledgeID string,
) (string, error) {
	logger.Info(ctx, "Starting to save file locally")
	logger.Infof(ctx, "File information: name=%s, size=%d, tenant ID=%d, knowledge ID=%s",
		file.Filename, file.Size, tenantID, knowledgeID)

	// Create storage directory with tenant and knowledge ID
	dir := filepath.Join(s.baseDir, fmt.Sprintf("%d", tenantID), knowledgeID)
	logger.Infof(ctx, "Creating directory: %s", dir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		logger.Errorf(ctx, "Failed to create directory: %v", err)
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Generate unique filename using timestamp
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(dir, filename)
	logger.Infof(ctx, "Generated file path: %s", filePath)

	// Open source file for reading
	logger.Info(ctx, "Opening source file")
	src, err := file.Open()
	if err != nil {
		logger.Errorf(ctx, "Failed to open source file: %v", err)
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Create destination file for writing
	logger.Info(ctx, "Creating destination file")
	dst, err := os.Create(filePath)
	if err != nil {
		logger.Errorf(ctx, "Failed to create destination file: %v", err)
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// Copy content from source to destination
	logger.Info(ctx, "Copying file content")
	if _, err := io.Copy(dst, src); err != nil {
		logger.Errorf(ctx, "Failed to copy file content: %v", err)
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	logger.Infof(ctx, "File saved successfully: %s", filePath)
	return filePath, nil
}

// GetFile retrieves a file from the local file system by its path
// Returns a ReadCloser for reading the file content
func (s *localFileService) GetFile(ctx context.Context, filePath string) (io.ReadCloser, error) {
	logger.Infof(ctx, "Getting file: %s", filePath)

	// Open the file for reading
	file, err := os.Open(filePath)
	if err != nil {
		logger.Errorf(ctx, "Failed to open file: %v", err)
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	logger.Info(ctx, "File opened successfully")
	return file, nil
}

// DeleteFile removes a file from the local file system
// Returns an error if deletion fails
func (s *localFileService) DeleteFile(ctx context.Context, filePath string) error {
	logger.Infof(ctx, "Deleting file: %s", filePath)

	// Remove the file
	err := os.Remove(filePath)
	if err != nil {
		logger.Errorf(ctx, "Failed to delete file: %v", err)
		return fmt.Errorf("failed to delete file: %w", err)
	}

	logger.Info(ctx, "File deleted successfully")
	return nil
}

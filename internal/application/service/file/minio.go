package file

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"

	"github.com/Tencent/WeKnora/internal/types/interfaces"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// minioFileService MinIO file service implementation
type minioFileService struct {
	client     *minio.Client
	bucketName string
}

// NewMinioFileService creates a MinIO file service
func NewMinioFileService(endpoint,
	accessKeyID, secretAccessKey, bucketName string, useSSL bool) (interfaces.FileService, error) {
	// Initialize MinIO client
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MinIO client: %w", err)
	}

	// Check if bucket exists, create if not
	exists, err := client.BucketExists(context.Background(), bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket: %w", err)
	}

	if !exists {
		err = client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return &minioFileService{
		client:     client,
		bucketName: bucketName,
	}, nil
}

// SaveFile saves a file to MinIO
func (s *minioFileService) SaveFile(ctx context.Context,
	file *multipart.FileHeader, tenantID uint, knowledgeID string,
) (string, error) {
	// Generate object name
	ext := filepath.Ext(file.Filename)
	objectName := fmt.Sprintf("%d/%s/%s%s", tenantID, knowledgeID, uuid.New().String(), ext)

	// Open file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Upload file to MinIO
	_, err = s.client.PutObject(ctx, s.bucketName, objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to MinIO: %w", err)
	}

	// Return the complete path to the object
	return fmt.Sprintf("minio://%s/%s", s.bucketName, objectName), nil
}

// GetFile gets a file from MinIO
func (s *minioFileService) GetFile(ctx context.Context, filePath string) (io.ReadCloser, error) {
	// Parse MinIO path
	// Format: minio://bucketName/objectName
	if len(filePath) < 9 || filePath[:8] != "minio://" {
		return nil, fmt.Errorf("invalid MinIO file path: %s", filePath)
	}

	// Extract object name
	objectName := filePath[9+len(s.bucketName):]
	if objectName[0] == '/' {
		objectName = objectName[1:]
	}

	// Get object
	obj, err := s.client.GetObject(ctx, s.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get file from MinIO: %w", err)
	}

	return obj, nil
}

// DeleteFile deletes a file
func (s *minioFileService) DeleteFile(ctx context.Context, filePath string) error {
	// Parse MinIO path
	// Format: minio://bucketName/objectName
	if len(filePath) < 9 || filePath[:8] != "minio://" {
		return fmt.Errorf("invalid MinIO file path: %s", filePath)
	}

	// Extract object name
	objectName := filePath[9+len(s.bucketName):]
	if objectName[0] == '/' {
		objectName = objectName[1:]
	}

	// Delete object
	err := s.client.RemoveObject(ctx, s.bucketName, objectName, minio.RemoveObjectOptions{
		GovernanceBypass: true,
	})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

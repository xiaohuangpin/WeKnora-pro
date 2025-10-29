package service

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"strings"
	"time"

	"github.com/Tencent/WeKnora/internal/logger"
	"github.com/Tencent/WeKnora/internal/types"
	"github.com/Tencent/WeKnora/internal/types/interfaces"
)

var apiKeySecret = func() []byte {
	return []byte(os.Getenv("TENANT_AES_KEY"))
}

// ListTenantsParams defines parameters for listing tenants with filtering and pagination
type ListTenantsParams struct {
	Page     int    // Page number for pagination
	PageSize int    // Number of items per page
	Status   string // Filter by tenant status
	Name     string // Filter by tenant name
}

// tenantService implements the TenantService interface
type tenantService struct {
	repo interfaces.TenantRepository // Repository for tenant data operations
}

// NewTenantService creates a new tenant service instance
func NewTenantService(repo interfaces.TenantRepository) interfaces.TenantService {
	return &tenantService{repo: repo}
}

// CreateTenant creates a new tenant
func (s *tenantService) CreateTenant(ctx context.Context, tenant *types.Tenant) (*types.Tenant, error) {
	logger.Info(ctx, "Start creating tenant")

	if tenant.Name == "" {
		logger.Error(ctx, "Tenant name cannot be empty")
		return nil, errors.New("tenant name cannot be empty")
	}

	logger.Infof(ctx, "Creating tenant, name: %s", tenant.Name)

	// Create tenant with initial values
	tenant.APIKey = s.generateApiKey(0)
	tenant.Status = "active"
	tenant.CreatedAt = time.Now()
	tenant.UpdatedAt = time.Now()

	logger.Info(ctx, "Saving tenant information to database")
	if err := s.repo.CreateTenant(ctx, tenant); err != nil {
		logger.ErrorWithFields(ctx, err, map[string]interface{}{
			"tenant_name": tenant.Name,
		})
		return nil, err
	}

	logger.Infof(ctx, "Tenant created successfully, ID: %d, generating official API Key", tenant.ID)
	tenant.APIKey = s.generateApiKey(tenant.ID)
	if err := s.repo.UpdateTenant(ctx, tenant); err != nil {
		logger.ErrorWithFields(ctx, err, map[string]interface{}{
			"tenant_id":   tenant.ID,
			"tenant_name": tenant.Name,
		})
		return nil, err
	}

	logger.Infof(ctx, "Tenant creation and update completed, ID: %d, name: %s", tenant.ID, tenant.Name)
	return tenant, nil
}

// GetTenantByID retrieves a tenant by their ID
func (s *tenantService) GetTenantByID(ctx context.Context, id uint) (*types.Tenant, error) {
	logger.Info(ctx, "Start retrieving tenant")

	if id == 0 {
		logger.Error(ctx, "Tenant ID cannot be 0")
		return nil, errors.New("tenant ID cannot be 0")
	}

	logger.Infof(ctx, "Retrieving tenant, ID: %d", id)

	tenant, err := s.repo.GetTenantByID(ctx, id)
	if err != nil {
		logger.ErrorWithFields(ctx, err, map[string]interface{}{
			"tenant_id": id,
		})
		return nil, err
	}

	logger.Infof(ctx, "Tenant retrieved successfully, ID: %d, name: %s", tenant.ID, tenant.Name)
	return tenant, nil
}

// ListTenants retrieves a list of all tenants
func (s *tenantService) ListTenants(ctx context.Context) ([]*types.Tenant, error) {
	logger.Info(ctx, "Start retrieving tenant list")

	tenants, err := s.repo.ListTenants(ctx)
	if err != nil {
		logger.ErrorWithFields(ctx, err, nil)
		return nil, err
	}

	logger.Infof(ctx, "Tenant list retrieved successfully, total: %d", len(tenants))
	return tenants, nil
}

// UpdateTenant updates an existing tenant's information
func (s *tenantService) UpdateTenant(ctx context.Context, tenant *types.Tenant) (*types.Tenant, error) {
	if tenant.ID == 0 {
		logger.Error(ctx, "Tenant ID cannot be 0")
		return nil, errors.New("tenant ID cannot be 0")
	}

	logger.Infof(ctx, "Updating tenant, ID: %d, name: %s", tenant.ID, tenant.Name)

	// Generate new API key if empty
	if tenant.APIKey == "" {
		logger.Info(ctx, "API Key is empty, generating new API Key")
		tenant.APIKey = s.generateApiKey(tenant.ID)
	}

	tenant.UpdatedAt = time.Now()
	logger.Info(ctx, "Saving tenant information to database")

	if err := s.repo.UpdateTenant(ctx, tenant); err != nil {
		logger.ErrorWithFields(ctx, err, map[string]interface{}{
			"tenant_id": tenant.ID,
		})
		return nil, err
	}

	logger.Infof(ctx, "Tenant updated successfully, ID: %d", tenant.ID)
	return tenant, nil
}

// DeleteTenant removes a tenant by their ID
func (s *tenantService) DeleteTenant(ctx context.Context, id uint) error {
	logger.Info(ctx, "Start deleting tenant")

	if id == 0 {
		logger.Error(ctx, "Tenant ID cannot be 0")
		return errors.New("tenant ID cannot be 0")
	}

	logger.Infof(ctx, "Deleting tenant, ID: %d", id)

	// Get tenant information for logging
	tenant, err := s.repo.GetTenantByID(ctx, id)
	if err != nil {
		if err.Error() == "record not found" {
			logger.Warnf(ctx, "Tenant to be deleted does not exist, ID: %d", id)
		} else {
			logger.ErrorWithFields(ctx, err, map[string]interface{}{
				"tenant_id": id,
			})
			return err
		}
	} else {
		logger.Infof(ctx, "Deleting tenant, ID: %d, name: %s", id, tenant.Name)
	}

	err = s.repo.DeleteTenant(ctx, id)
	if err != nil {
		logger.ErrorWithFields(ctx, err, map[string]interface{}{
			"tenant_id": id,
		})
		return err
	}

	logger.Infof(ctx, "Tenant deleted successfully, ID: %d", id)
	return nil
}

// UpdateAPIKey updates the API key for a specific tenant
func (s *tenantService) UpdateAPIKey(ctx context.Context, id uint) (string, error) {
	logger.Info(ctx, "Start updating tenant API Key")

	if id == 0 {
		logger.Error(ctx, "Tenant ID cannot be 0")
		return "", errors.New("tenant ID cannot be 0")
	}

	logger.Infof(ctx, "Retrieving tenant information, ID: %d", id)

	tenant, err := s.repo.GetTenantByID(ctx, id)
	if err != nil {
		logger.ErrorWithFields(ctx, err, map[string]interface{}{
			"tenant_id": id,
		})
		return "", err
	}

	logger.Infof(ctx, "Generating new API Key for tenant, ID: %d", id)
	tenant.APIKey = s.generateApiKey(tenant.ID)

	if err := s.repo.UpdateTenant(ctx, tenant); err != nil {
		logger.ErrorWithFields(ctx, err, map[string]interface{}{
			"tenant_id": id,
		})
		return "", err
	}

	logger.Infof(ctx, "Tenant API Key updated successfully, ID: %d", id)
	return tenant.APIKey, nil
}

// generateApiKey generates a secure API key for tenant authentication
func (r *tenantService) generateApiKey(tenantID uint) string {
	// 1. Convert tenant_id to bytes
	idBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(idBytes, uint64(tenantID))

	// 2. Encrypt tenant_id using AES-GCM
	block, err := aes.NewCipher(apiKeySecret())
	if err != nil {
		panic("Failed to create AES cipher: " + err.Error())
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic("Failed to create GCM cipher: " + err.Error())
	}

	ciphertext := aesgcm.Seal(nil, nonce, idBytes, nil)

	// 3. Combine nonce and ciphertext, then encode with base64
	combined := append(nonce, ciphertext...)
	encoded := base64.RawURLEncoding.EncodeToString(combined)

	// Create final API Key in format: sk-{encrypted_part}
	return "sk-" + encoded
}

// ExtractTenantIDFromAPIKey extracts the tenant ID from an API key
func (r *tenantService) ExtractTenantIDFromAPIKey(apiKey string) (uint, error) {
	// 1. Validate format and extract encrypted part
	parts := strings.SplitN(apiKey, "-", 2)
	if len(parts) != 2 || parts[0] != "sk" {
		return 0, errors.New("invalid API key format")
	}

	// 2. Decode the base64 part
	encryptedData, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return 0, errors.New("invalid API key encoding")
	}

	// 3. Separate nonce and ciphertext
	if len(encryptedData) < 12 {
		return 0, errors.New("invalid API key length")
	}
	nonce, ciphertext := encryptedData[:12], encryptedData[12:]

	// 4. Decrypt
	block, err := aes.NewCipher(apiKeySecret())
	if err != nil {
		return 0, errors.New("decryption error")
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return 0, errors.New("decryption error")
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return 0, errors.New("API key is invalid or has been tampered with")
	}

	// 5. Convert back to tenant_id
	tenantID := binary.LittleEndian.Uint64(plaintext)

	return uint(tenantID), nil
}

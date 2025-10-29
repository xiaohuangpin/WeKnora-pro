package interfaces

import (
	"context"

	"github.com/Tencent/WeKnora/internal/types"
)

// TenantService defines the tenant service interface
type TenantService interface {
	// CreateTenant creates a tenant
	CreateTenant(ctx context.Context, tenant *types.Tenant) (*types.Tenant, error)
	// GetTenantByID gets a tenant by ID
	GetTenantByID(ctx context.Context, id uint) (*types.Tenant, error)
	// ListTenants lists all tenants
	ListTenants(ctx context.Context) ([]*types.Tenant, error)
	// UpdateTenant updates a tenant
	UpdateTenant(ctx context.Context, tenant *types.Tenant) (*types.Tenant, error)
	// DeleteTenant deletes a tenant
	DeleteTenant(ctx context.Context, id uint) error
	// UpdateAPIKey updates the API key
	UpdateAPIKey(ctx context.Context, id uint) (string, error)
	// ExtractTenantIDFromAPIKey extracts the tenant ID from the API key
	ExtractTenantIDFromAPIKey(apiKey string) (uint, error)
}

// TenantRepository defines the tenant repository interface
type TenantRepository interface {
	// CreateTenant creates a tenant
	CreateTenant(ctx context.Context, tenant *types.Tenant) error
	// GetTenantByID gets a tenant by ID
	GetTenantByID(ctx context.Context, id uint) (*types.Tenant, error)
	// ListTenants lists all tenants
	ListTenants(ctx context.Context) ([]*types.Tenant, error)
	// UpdateTenant updates a tenant
	UpdateTenant(ctx context.Context, tenant *types.Tenant) error
	// DeleteTenant deletes a tenant
	DeleteTenant(ctx context.Context, id uint) error
	// AdjustStorageUsed adjusts the storage used for a tenant
	AdjustStorageUsed(ctx context.Context, tenantID uint, delta int64) error
}

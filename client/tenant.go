// Package client provides the implementation for interacting with the WeKnora API
// The Tenant related interfaces are used to manage tenants in the system
// Tenants can be created, retrieved, updated, deleted, and queried
// They can also be used to manage retriever engines for different tasks
package client

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// RetrieverEngines defines a collection of retriever engine parameters
type RetrieverEngines struct {
	Engines []RetrieverEngineParams `json:"engines"`
}

// RetrieverEngineParams contains configuration for retriever engines
type RetrieverEngineParams struct {
	RetrieverType       string `json:"retriever_type"`        // Type of retriever (e.g., keywords, vector)
	RetrieverEngineType string `json:"retriever_engine_type"` // Type of engine implementing the retriever
}

// Tenant represents tenant information in the system
type Tenant struct {
	ID uint `yaml:"id" json:"id" gorm:"primaryKey"`
	// Tenant name
	Name string `yaml:"name" json:"name"`
	// Tenant description
	Description string `yaml:"description" json:"description"`
	// API key for authentication
	APIKey string `yaml:"api_key" json:"api_key"`
	// Tenant status (active, inactive)
	Status string `yaml:"status" json:"status" gorm:"default:'active'"`
	// Configured retrieval engines
	RetrieverEngines RetrieverEngines `yaml:"retriever_engines" json:"retriever_engines" gorm:"type:json"`
	// Business/department information
	Business string `yaml:"business" json:"business"`
	// Creation timestamp
	CreatedAt time.Time `yaml:"created_at" json:"created_at"`
	// Last update timestamp
	UpdatedAt time.Time `yaml:"updated_at" json:"updated_at"`
}

// TenantResponse represents the API response structure for tenant operations
type TenantResponse struct {
	Success bool   `json:"success"` // Whether the operation was successful
	Data    Tenant `json:"data"`    // Tenant data
}

// TenantListResponse represents the API response structure for listing tenants
type TenantListResponse struct {
	Success bool `json:"success"` // Whether the operation was successful
	Data    struct {
		Items []Tenant `json:"items"` // List of tenant items
	} `json:"data"`
}

// CreateTenant creates a new tenant
func (c *Client) CreateTenant(ctx context.Context, tenant *Tenant) (*Tenant, error) {
	resp, err := c.doRequest(ctx, http.MethodPost, "/api/v1/tenants", tenant, nil)
	if err != nil {
		return nil, err
	}

	var response TenantResponse
	if err := parseResponse(resp, &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

// GetTenant retrieves a tenant by ID
func (c *Client) GetTenant(ctx context.Context, tenantID uint) (*Tenant, error) {
	path := fmt.Sprintf("/api/v1/tenants/%d", tenantID)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}

	var response TenantResponse
	if err := parseResponse(resp, &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

// UpdateTenant updates an existing tenant
func (c *Client) UpdateTenant(ctx context.Context, tenant *Tenant) (*Tenant, error) {
	path := fmt.Sprintf("/api/v1/tenants/%d", tenant.ID)
	resp, err := c.doRequest(ctx, http.MethodPut, path, tenant, nil)
	if err != nil {
		return nil, err
	}

	var response TenantResponse
	if err := parseResponse(resp, &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

// DeleteTenant removes a tenant by ID
func (c *Client) DeleteTenant(ctx context.Context, tenantID uint) error {
	path := fmt.Sprintf("/api/v1/tenants/%d", tenantID)
	resp, err := c.doRequest(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return err
	}

	var response struct {
		Success bool   `json:"success"`
		Message string `json:"message,omitempty"`
	}

	return parseResponse(resp, &response)
}

// ListTenants retrieves all tenants
func (c *Client) ListTenants(ctx context.Context) ([]Tenant, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "/api/v1/tenants", nil, nil)
	if err != nil {
		return nil, err
	}

	var response TenantListResponse
	if err := parseResponse(resp, &response); err != nil {
		return nil, err
	}

	return response.Data.Items, nil
}

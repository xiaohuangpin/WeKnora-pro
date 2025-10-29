package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Tencent/WeKnora/internal/errors"
	"github.com/Tencent/WeKnora/internal/logger"
	"github.com/Tencent/WeKnora/internal/types"
	"github.com/Tencent/WeKnora/internal/types/interfaces"
)

// TenantHandler implements HTTP request handlers for tenant management
// Provides functionality for creating, retrieving, updating, and deleting tenants
// through the REST API endpoints
type TenantHandler struct {
	service interfaces.TenantService
}

// NewTenantHandler creates a new tenant handler instance with the provided service
// Parameters:
//   - service: An implementation of the TenantService interface for business logic
//
// Returns a pointer to the newly created TenantHandler
func NewTenantHandler(service interfaces.TenantService) *TenantHandler {
	return &TenantHandler{
		service: service,
	}
}

// CreateTenant handles the HTTP request for creating a new tenant
// It deserializes the request body into a tenant object, validates it,
// calls the service to create the tenant, and returns the result
// Parameters:
//   - c: Gin context for the HTTP request
func (h *TenantHandler) CreateTenant(c *gin.Context) {
	ctx := c.Request.Context()

	logger.Info(ctx, "Start creating tenant")

	var tenantData types.Tenant
	if err := c.ShouldBindJSON(&tenantData); err != nil {
		logger.Error(ctx, "Failed to parse request parameters", err)
		appErr := errors.NewValidationError("Invalid request parameters").WithDetails(err.Error())
		c.Error(appErr)
		return
	}

	logger.Infof(ctx, "Creating tenant, name: %s", tenantData.Name)

	createdTenant, err := h.service.CreateTenant(ctx, &tenantData)
	if err != nil {
		// Check if this is an application-specific error
		if appErr, ok := errors.IsAppError(err); ok {
			logger.Error(ctx, "Failed to create tenant: application error", appErr)
			c.Error(appErr)
		} else {
			logger.ErrorWithFields(ctx, err, nil)
			c.Error(errors.NewInternalServerError("Failed to create tenant").WithDetails(err.Error()))
		}
		return
	}

	logger.Infof(ctx, "Tenant created successfully, ID: %d, name: %s", createdTenant.ID, createdTenant.Name)
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    createdTenant,
	})
}

// GetTenant handles the HTTP request for retrieving a tenant by ID
// It extracts and validates the tenant ID from the URL parameter,
// retrieves the tenant from the service, and returns it in the response
// Parameters:
//   - c: Gin context for the HTTP request
func (h *TenantHandler) GetTenant(c *gin.Context) {
	ctx := c.Request.Context()

	logger.Info(ctx, "Start retrieving tenant")

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		logger.Errorf(ctx, "Invalid tenant ID: %s", c.Param("id"))
		c.Error(errors.NewBadRequestError("Invalid tenant ID"))
		return
	}

	logger.Infof(ctx, "Retrieving tenant, ID: %d", id)

	tenant, err := h.service.GetTenantByID(ctx, uint(id))
	if err != nil {
		// Check if this is an application-specific error
		if appErr, ok := errors.IsAppError(err); ok {
			logger.Error(ctx, "Failed to retrieve tenant: application error", appErr)
			c.Error(appErr)
		} else {
			logger.ErrorWithFields(ctx, err, nil)
			c.Error(errors.NewInternalServerError("Failed to retrieve tenant").WithDetails(err.Error()))
		}
		return
	}

	logger.Infof(ctx, "Retrieved tenant successfully, ID: %d, Name: %s", tenant.ID, tenant.Name)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    tenant,
	})
}

// UpdateTenant handles the HTTP request for updating an existing tenant
// It extracts the tenant ID from the URL parameter, deserializes the request body,
// validates the data, updates the tenant through the service, and returns the result
// Parameters:
//   - c: Gin context for the HTTP request
func (h *TenantHandler) UpdateTenant(c *gin.Context) {
	ctx := c.Request.Context()

	logger.Info(ctx, "Start updating tenant")

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		logger.Errorf(ctx, "Invalid tenant ID: %s", c.Param("id"))
		c.Error(errors.NewBadRequestError("Invalid tenant ID"))
		return
	}

	var tenantData types.Tenant
	if err := c.ShouldBindJSON(&tenantData); err != nil {
		logger.Error(ctx, "Failed to parse request parameters", err)
		c.Error(errors.NewValidationError("Invalid request data").WithDetails(err.Error()))
		return
	}

	logger.Infof(ctx, "Updating tenant, ID: %d, Name: %s", id, tenantData.Name)

	tenantData.ID = uint(id)
	updatedTenant, err := h.service.UpdateTenant(ctx, &tenantData)
	if err != nil {
		// Check if this is an application-specific error
		if appErr, ok := errors.IsAppError(err); ok {
			logger.Error(ctx, "Failed to update tenant: application error", appErr)
			c.Error(appErr)
		} else {
			logger.ErrorWithFields(ctx, err, nil)
			c.Error(errors.NewInternalServerError("Failed to update tenant").WithDetails(err.Error()))
		}
		return
	}

	logger.Infof(ctx, "Tenant updated successfully, ID: %d, Name: %s", updatedTenant.ID, updatedTenant.Name)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    updatedTenant,
	})
}

// DeleteTenant handles the HTTP request for deleting a tenant
// It extracts and validates the tenant ID from the URL parameter,
// calls the service to delete the tenant, and returns the result
// Parameters:
//   - c: Gin context for the HTTP request
func (h *TenantHandler) DeleteTenant(c *gin.Context) {
	ctx := c.Request.Context()

	logger.Info(ctx, "Start deleting tenant")

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		logger.Errorf(ctx, "Invalid tenant ID: %s", c.Param("id"))
		c.Error(errors.NewBadRequestError("Invalid tenant ID"))
		return
	}

	logger.Infof(ctx, "Deleting tenant, ID: %d", id)

	if err := h.service.DeleteTenant(ctx, uint(id)); err != nil {
		// Check if this is an application-specific error
		if appErr, ok := errors.IsAppError(err); ok {
			logger.Error(ctx, "Failed to delete tenant: application error", appErr)
			c.Error(appErr)
		} else {
			logger.ErrorWithFields(ctx, err, nil)
			c.Error(errors.NewInternalServerError("Failed to delete tenant").WithDetails(err.Error()))
		}
		return
	}

	logger.Infof(ctx, "Tenant deleted successfully, ID: %d", id)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Tenant deleted successfully",
	})
}

// ListTenants handles the HTTP request for retrieving a list of all tenants
// It calls the service to fetch the tenant list and returns it in the response
// Parameters:
//   - c: Gin context for the HTTP request
func (h *TenantHandler) ListTenants(c *gin.Context) {
	ctx := c.Request.Context()

	logger.Info(ctx, "Start retrieving tenant list")

	tenants, err := h.service.ListTenants(ctx)
	if err != nil {
		// Check if this is an application-specific error
		if appErr, ok := errors.IsAppError(err); ok {
			logger.Error(ctx, "Failed to retrieve tenant list: application error", appErr)
			c.Error(appErr)
		} else {
			logger.ErrorWithFields(ctx, err, nil)
			c.Error(errors.NewInternalServerError("Failed to retrieve tenant list").WithDetails(err.Error()))
		}
		return
	}

	logger.Infof(ctx, "Retrieved tenant list successfully, Total: %d tenants", len(tenants))
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"items": tenants,
		},
	})
}

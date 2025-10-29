package handler

import (
	"net/http"

	"github.com/Tencent/WeKnora/internal/application/service"
	"github.com/Tencent/WeKnora/internal/errors"
	"github.com/Tencent/WeKnora/internal/logger"
	"github.com/Tencent/WeKnora/internal/types"
	"github.com/Tencent/WeKnora/internal/types/interfaces"
	"github.com/gin-gonic/gin"
)

// ModelHandler handles HTTP requests for model-related operations
// It implements the necessary methods to create, retrieve, update, and delete models
type ModelHandler struct {
	service interfaces.ModelService
}

// NewModelHandler creates a new instance of ModelHandler
// It requires a model service implementation that handles business logic
// Parameters:
//   - service: An implementation of the ModelService interface
//
// Returns a pointer to the newly created ModelHandler
func NewModelHandler(service interfaces.ModelService) *ModelHandler {
	return &ModelHandler{service: service}
}

// CreateModelRequest defines the structure for model creation requests
// Contains all fields required to create a new model in the system
type CreateModelRequest struct {
	Name        string                `json:"name" binding:"required"`
	Type        types.ModelType       `json:"type" binding:"required"`
	Source      types.ModelSource     `json:"source" binding:"required"`
	Description string                `json:"description"`
	Parameters  types.ModelParameters `json:"parameters" binding:"required"`
	IsDefault   bool                  `json:"is_default"`
}

// CreateModel handles the HTTP request to create a new model
// It validates the request, processes it using the model service,
// and returns the created model to the client
// Parameters:
//   - c: Gin context for the HTTP request
func (h *ModelHandler) CreateModel(c *gin.Context) {
	ctx := c.Request.Context()

	logger.Info(ctx, "Start creating model")

	var req CreateModelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error(ctx, "Failed to parse request parameters", err)
		c.Error(errors.NewBadRequestError(err.Error()))
		return
	}

	tenantID := c.GetUint(types.TenantIDContextKey.String())
	if tenantID == 0 {
		logger.Error(ctx, "Tenant ID is empty")
		c.Error(errors.NewBadRequestError("Tenant ID cannot be empty"))
		return
	}

	logger.Infof(ctx, "Creating model, Tenant ID: %d, Model name: %s, Model type: %s",
		tenantID, req.Name, req.Type)

	model := &types.Model{
		TenantID:    tenantID,
		Name:        req.Name,
		Type:        req.Type,
		Source:      req.Source,
		Description: req.Description,
		Parameters:  req.Parameters,
		IsDefault:   req.IsDefault,
	}

	if err := h.service.CreateModel(ctx, model); err != nil {
		logger.ErrorWithFields(ctx, err, nil)
		c.Error(errors.NewInternalServerError(err.Error()))
		return
	}

	logger.Infof(ctx, "Model created successfully, ID: %s, Name: %s", model.ID, model.Name)
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    model,
	})
}

// GetModel handles the HTTP request to retrieve a model by its ID
// It fetches the model from the service and returns it to the client,
// or returns appropriate error messages if the model cannot be found
// Parameters:
//   - c: Gin context for the HTTP request
func (h *ModelHandler) GetModel(c *gin.Context) {
	ctx := c.Request.Context()

	logger.Info(ctx, "Start retrieving model")

	id := c.Param("id")
	if id == "" {
		logger.Error(ctx, "Model ID is empty")
		c.Error(errors.NewBadRequestError("Model ID cannot be empty"))
		return
	}

	logger.Infof(ctx, "Retrieving model, ID: %s", id)
	model, err := h.service.GetModelByID(ctx, id)
	if err != nil {
		if err == service.ErrModelNotFound {
			logger.Warnf(ctx, "Model not found, ID: %s", id)
			c.Error(errors.NewNotFoundError("Model not found"))
			return
		}
		logger.ErrorWithFields(ctx, err, nil)
		c.Error(errors.NewInternalServerError(err.Error()))
		return
	}

	logger.Infof(ctx, "Retrieved model successfully, ID: %s, Name: %s", model.ID, model.Name)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    model,
	})
}

// ListModels handles the HTTP request to retrieve all models for a tenant
// It validates the tenant ID, fetches models from the service, and returns them to the client
// Parameters:
//   - c: Gin context for the HTTP request
func (h *ModelHandler) ListModels(c *gin.Context) {
	ctx := c.Request.Context()

	logger.Info(ctx, "Start retrieving model list")

	tenantID := c.GetUint(types.TenantIDContextKey.String())
	if tenantID == 0 {
		logger.Error(ctx, "Tenant ID is empty")
		c.Error(errors.NewBadRequestError("Tenant ID cannot be empty"))
		return
	}

	logger.Infof(ctx, "Retrieving model list, Tenant ID: %d", tenantID)
	models, err := h.service.ListModels(ctx)
	if err != nil {
		logger.ErrorWithFields(ctx, err, nil)
		c.Error(errors.NewInternalServerError(err.Error()))
		return
	}

	logger.Infof(ctx, "Retrieved model list successfully, Tenant ID: %d, Total: %d models", tenantID, len(models))
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    models,
	})
}

// UpdateModelRequest defines the structure for model update requests
// Contains fields that can be updated for an existing model
type UpdateModelRequest struct {
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Parameters  types.ModelParameters `json:"parameters"`
	IsDefault   bool                  `json:"is_default"`
}

// UpdateModel handles the HTTP request to update an existing model
// It validates the request, retrieves the current model, applies changes,
// and updates the model in the service
// Parameters:
//   - c: Gin context for the HTTP request
func (h *ModelHandler) UpdateModel(c *gin.Context) {
	ctx := c.Request.Context()

	logger.Info(ctx, "Start updating model")

	id := c.Param("id")
	if id == "" {
		logger.Error(ctx, "Model ID is empty")
		c.Error(errors.NewBadRequestError("Model ID cannot be empty"))
		return
	}

	var req UpdateModelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error(ctx, "Failed to parse request parameters", err)
		c.Error(errors.NewBadRequestError(err.Error()))
		return
	}

	logger.Infof(ctx, "Retrieving model information, ID: %s", id)
	model, err := h.service.GetModelByID(ctx, id)
	if err != nil {
		if err == service.ErrModelNotFound {
			logger.Warnf(ctx, "Model not found, ID: %s", id)
			c.Error(errors.NewNotFoundError("Model not found"))
			return
		}
		logger.ErrorWithFields(ctx, err, nil)
		c.Error(errors.NewInternalServerError(err.Error()))
		return
	}

	// Update model fields if they are provided in the request
	if req.Name != "" {
		model.Name = req.Name
	}
	if req.Description != "" {
		model.Description = req.Description
	}
	if req.Parameters != (types.ModelParameters{}) {
		model.Parameters = req.Parameters
	}
	model.IsDefault = req.IsDefault

	logger.Infof(ctx, "Updating model, ID: %s, Name: %s", id, model.Name)
	if err := h.service.UpdateModel(ctx, model); err != nil {
		logger.ErrorWithFields(ctx, err, nil)
		c.Error(errors.NewInternalServerError(err.Error()))
		return
	}

	logger.Infof(ctx, "Model updated successfully, ID: %s", id)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    model,
	})
}

// DeleteModel handles the HTTP request to delete a model by its ID
// It validates the model ID, attempts to delete the model through the service,
// and returns appropriate status and messages
// Parameters:
//   - c: Gin context for the HTTP request
func (h *ModelHandler) DeleteModel(c *gin.Context) {
	ctx := c.Request.Context()

	logger.Info(ctx, "Start deleting model")

	id := c.Param("id")
	if id == "" {
		logger.Error(ctx, "Model ID is empty")
		c.Error(errors.NewBadRequestError("Model ID cannot be empty"))
		return
	}

	logger.Infof(ctx, "Deleting model, ID: %s", id)
	if err := h.service.DeleteModel(ctx, id); err != nil {
		if err == service.ErrModelNotFound {
			logger.Warnf(ctx, "Model not found, ID: %s", id)
			c.Error(errors.NewNotFoundError("Model not found"))
			return
		}
		logger.ErrorWithFields(ctx, err, nil)
		c.Error(errors.NewInternalServerError(err.Error()))
		return
	}

	logger.Infof(ctx, "Model deleted successfully, ID: %s", id)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Model deleted",
	})
}

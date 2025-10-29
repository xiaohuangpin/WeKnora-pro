package repository

import (
	"context"
	"errors"

	"github.com/Tencent/WeKnora/internal/types"
	"github.com/Tencent/WeKnora/internal/types/interfaces"
	"gorm.io/gorm"
)

// modelRepository implements the model repository interface
type modelRepository struct {
	db *gorm.DB
}

// NewModelRepository creates a new model repository
func NewModelRepository(db *gorm.DB) interfaces.ModelRepository {
	return &modelRepository{db: db}
}

// Create creates a new model
func (r *modelRepository) Create(ctx context.Context, m *types.Model) error {
	return r.db.WithContext(ctx).Create(m).Error
}

// GetByID retrieves a model by ID
func (r *modelRepository) GetByID(ctx context.Context, tenantID uint, id string) (*types.Model, error) {
	var m types.Model
	if err := r.db.WithContext(ctx).Where("id = ?", id).Where(
		"tenant_id = ? or tenant_id = 0", tenantID,
	).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

// List lists models with optional filtering
func (r *modelRepository) List(
	ctx context.Context, tenantID uint, modelType types.ModelType, source types.ModelSource,
) ([]*types.Model, error) {
	var models []*types.Model
	query := r.db.WithContext(ctx).Where(
		"tenant_id = ? or tenant_id = 0", tenantID,
	)

	if modelType != "" {
		query = query.Where("type = ?", modelType)
	}

	if source != "" {
		query = query.Where("source = ?", source)
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	return models, nil
}

// Update updates a model
func (r *modelRepository) Update(ctx context.Context, m *types.Model) error {
	return r.db.WithContext(ctx).Debug().Model(&types.Model{}).Where(
		"id = ? AND tenant_id = ?", m.ID, m.TenantID,
	).Updates(m).Error
}

// Delete deletes a model
func (r *modelRepository) Delete(ctx context.Context, tenantID uint, id string) error {
	return r.db.WithContext(ctx).Where(
		"id = ? AND tenant_id = ?", id, tenantID,
	).Delete(&types.Model{}).Error
}

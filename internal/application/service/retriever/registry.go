package retriever

import (
	"fmt"
	"sync"

	"github.com/Tencent/WeKnora/internal/types"
	"github.com/Tencent/WeKnora/internal/types/interfaces"
)

// RetrieveEngineRegistry implements the retrieval engine registry
type RetrieveEngineRegistry struct {
	repositories map[types.RetrieverEngineType]interfaces.RetrieveEngineService
	mu           sync.RWMutex
}

// NewRetrieveEngineRegistry creates a new retrieval engine registry
func NewRetrieveEngineRegistry() interfaces.RetrieveEngineRegistry {
	return &RetrieveEngineRegistry{
		repositories: make(map[types.RetrieverEngineType]interfaces.RetrieveEngineService),
	}
}

// Register registers a retrieval engine service
func (r *RetrieveEngineRegistry) Register(repo interfaces.RetrieveEngineService) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.repositories[repo.EngineType()]; exists {
		return fmt.Errorf("Repository type %s already registered", repo.EngineType())
	}

	r.repositories[repo.EngineType()] = repo
	return nil
}

// GetRetrieveEngineService retrieves a retrieval engine service by type
func (r *RetrieveEngineRegistry) GetRetrieveEngineService(repoType types.RetrieverEngineType) (
	interfaces.RetrieveEngineService, error,
) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	repo, exists := r.repositories[repoType]
	if !exists {
		return nil, fmt.Errorf("Repository of type %s not found", repoType)
	}

	return repo, nil
}

// GetAllRetrieveEngineServices retrieves all registered retrieval engine services
func (r *RetrieveEngineRegistry) GetAllRetrieveEngineServices() []interfaces.RetrieveEngineService {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Create a copy to avoid modifying the original map
	result := make([]interfaces.RetrieveEngineService, 0, len(r.repositories))
	for _, v := range r.repositories {
		result = append(result, v)
	}

	return result
}

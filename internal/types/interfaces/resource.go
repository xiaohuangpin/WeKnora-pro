package interfaces

import (
	"context"

	"github.com/Tencent/WeKnora/internal/types"
)

// ResourceCleaner defines the resource cleaner interface
type ResourceCleaner interface {
	// Register registers a resource cleanup function
	Register(cleanup types.CleanupFunc)

	// RegisterWithName registers a resource cleanup function with a name
	RegisterWithName(name string, cleanup types.CleanupFunc)

	// Cleanup executes all resource cleanup functions
	Cleanup(ctx context.Context) []error
}

package reload

//go:generate mockgen -destination mock_reloader.go -package reload dmm-aggregator-backend/pkg/reload Reloader

import (
	"context"
)

// Reloader knows how to reload a resource.
type Reloader interface {
	Reload(ctx context.Context, id string) error
}

// ReloaderFunc is a helper to create reloaders based on functions.
type ReloaderFunc func(ctx context.Context, id string) error

// Reload satisfies Reloader interface.
func (r ReloaderFunc) Reload(ctx context.Context, id string) error { return r(ctx, id) }

package assetdomain

import (
	"context"

	"github.com/google/uuid"
)

// Repository is the interface that wraps the basic asset repository methods.
type Repository interface {
	// Save saves all the asset uncommited events to the event store
	Save(ctx context.Context, asset *Asset) error

	// GetByID returns the asset by the given ID
	GetByID(ctx context.Context, id uuid.UUID) (*Asset, error)

	// GetAll returns all the assets
	GetAll(ctx context.Context) ([]*Asset, error)

	// Exists checks if an asset with the given ID exists
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
}

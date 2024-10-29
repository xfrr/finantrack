package assetsmongo

import (
	"context"

	"github.com/google/uuid"
	"github.com/xfrr/finantrack/internal/shared/xmongo"

	assetDomain "github.com/xfrr/finantrack/internal/contexts/assets/domain"
)

var _ assetDomain.Repository = (*Repository)(nil)

// Repository represents the MongoDB repository for assets.
type Repository struct {
	eventStore xmongo.EventStore
}

// NewRepository creates a new MongoDB repository for assets.
func NewRepository(eventStore xmongo.EventStore) *Repository {
	return &Repository{
		eventStore: eventStore,
	}
}

// Save saves the asset changes into the event store.
func (r *Repository) Save(ctx context.Context, asset *assetDomain.Asset) error {
	changes := asset.AggregateChanges()
	if len(changes) == 0 {
		return nil
	}

	return r.eventStore.Save(ctx, changes...)
}

// GetAll retrieves all assets from the event store.
func (r *Repository) GetAll(ctx context.Context) ([]*assetDomain.Asset, error) {
	// TODO: Implement
	return nil, nil
}

// GetByID retrieves an asset by its ID from the event store.
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*assetDomain.Asset, error) {
	events, err := r.eventStore.Get(ctx, xmongo.WithAggregateIDCriteria(id.String())())
	if err != nil {
		return nil, err
	}

	asset, err := assetDomain.HydrateAsset(id, events)
	if err != nil {
		return nil, err
	}

	if asset.IsDeleted() {
		return nil, assetDomain.ErrAssetNotFound
	}

	return asset, nil
}

// Exists checks if an asset with the given ID exists in the event store.
func (r *Repository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	return r.eventStore.ExistsByAggregateID(ctx, id)
}

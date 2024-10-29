package assetimmudb

import (
	"context"
	"database/sql"
	"strings"

	"github.com/google/uuid"
	assetdomain "github.com/xfrr/finantrack/internal/contexts/assets/domain"
)

// Asset represents the asset entity.
type Asset struct {
	ID uuid.UUID `json:"id"`
}

// Repository implements the Repository interface using ImmuDB.
type Repository struct {
	db *sql.DB
}

// NewImmuRepository creates a new ImmuRepository with the given ImmuDB client.
func NewImmuRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Save saves the given asset to ImmuDB.
func (r *Repository) Save(ctx context.Context, asset *assetdomain.Asset) error {
	// TODO: Implement Save method
	return nil
}

// GetByID retrieves an asset by its ID from ImmuDB.
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*assetdomain.Asset, error) {
	// TODO: Implement GetByID method
	return nil, nil
}

// GetAll retrieves all assets from ImmuDB.
func (r *Repository) GetAll(ctx context.Context) ([]*assetdomain.Asset, error) {
	// TODO: Implement GetAll method
	return nil, nil
}

// Exists checks whether an asset exists in ImmuDB by its ID.
func (r *Repository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	// TODO: Implement Exists method
	return false, nil
}

// isNotFoundError checks if the error corresponds to a missing key in ImmuDB.
func isNotFoundError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "key not found")
}

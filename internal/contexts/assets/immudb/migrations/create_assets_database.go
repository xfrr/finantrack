package assetimmudbmigrations

import (
	"database/sql"

	"github.com/xfrr/finantrack/internal/shared/ximmudb"
)

var _ ximmudb.Migration = (*CreateAssetsDatabase)(nil)

// CreateAssetsDatabase represents a immudb SQL migration.
type CreateAssetsDatabase struct {
}

// NewCreateAssetsDatabase creates a new migration.
func NewCreateAssetsDatabase() ximmudb.Migration {
	return &CreateAssetsDatabase{}
}

// Up applies the migration.
func (m *CreateAssetsDatabase) Up(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE DATABASE IF NOT EXISTS assets;
	`)
	if err != nil {
		return err
	}

	return nil
}

// Down reverts the migration.
// Note: migration cannot be reverted in immudb.
func (m *CreateAssetsDatabase) Down() error {
	return nil
}

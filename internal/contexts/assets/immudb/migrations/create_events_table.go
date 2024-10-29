package assetimmudbmigrations

import (
	"database/sql"

	"github.com/xfrr/finantrack/internal/shared/ximmudb"
)

var _ ximmudb.Migration = (*CreateAssetEventsTable)(nil)

// CreateAssetEventsTable represents a immudb SQL migration.
type CreateAssetEventsTable struct {
}

// NewCreateAssetEventsTable creates a new migration.
func NewCreateAssetEventsTable() ximmudb.Migration {
	return &CreateAssetEventsTable{}
}

// Up applies the migration.
func (m *CreateAssetEventsTable) Up(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS events (
			id VARCHAR[100],
			aggregate_id VARCHAR[100],
			aggregate_name VARCHAR[100],
			created_at TIMESTAMP,
			payload VARCHAR[4096],
			PRIMARY KEY id
		);
	`)
	if err != nil {
		return err
	}

	return nil
}

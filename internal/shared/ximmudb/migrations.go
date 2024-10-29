package ximmudb

import "database/sql"

// Migration represents a immudb SQL migration.
type Migration interface {
	Up(db *sql.DB) error
}

// Migrate applies the migrations in order.
func Migrate(db *sql.DB, migrations []Migration) error {
	for _, m := range migrations {
		if err := m.Up(db); err != nil {
			return err
		}
	}
	return nil
}

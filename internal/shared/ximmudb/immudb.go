package ximmudb

import (
	"context"
	"database/sql"
	"time"

	immudb "github.com/codenotary/immudb/pkg/client"
	"github.com/codenotary/immudb/pkg/stdlib"
)

type Config struct {
	Host         string
	Port         int
	User         string
	Pass         string
	DB           string
	Timeout      time.Duration
	MaxOpenConns int
	MaxLife      time.Duration
	MaxIdleCons  int
}

func NewImmuDBClient(ctx context.Context, cfg Config) (*sql.DB, error) {
	// Stablish database connection configuration
	iopts := immudb.DefaultOptions().
		WithAddress(cfg.Host).
		WithPort(cfg.Port).
		WithUsername(cfg.User).
		WithPassword(cfg.Pass).
		WithDatabase(cfg.DB)

	// Create a new standard database/sql client and stablish pool connection
	dbClient := stdlib.OpenDB(iopts)
	dbClient.SetConnMaxLifetime(time.Duration(cfg.MaxLife) * time.Minute)
	dbClient.SetMaxOpenConns(cfg.MaxOpenConns)
	dbClient.SetMaxIdleConns(cfg.MaxIdleCons)

	pingErr := dbClient.Ping()
	if pingErr != nil {
		return nil, pingErr
	}

	return dbClient, nil
}

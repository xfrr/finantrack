package assets

import (
	"context"
	"strconv"
	"time"

	"github.com/xfrr/finantrack/internal/shared/xevent"
	"github.com/xfrr/finantrack/internal/shared/ximmudb"
	"github.com/xfrr/finantrack/services"

	assetdomain "github.com/xfrr/finantrack/internal/contexts/assets/domain"
	assetimmudb "github.com/xfrr/finantrack/internal/contexts/assets/immudb"
	assetimmudbmigrations "github.com/xfrr/finantrack/internal/contexts/assets/immudb/migrations"
)

const (
	DefaultTimeout      = 5 * time.Second
	DefaultMaxOpenConns = 10
	DefaultMaxLife      = 10 * time.Minute
	DefaultMaxIdleCons  = 10
)

type immudbRepositoryFactory struct {
	dbHost         string
	dbPort         string
	dbUser         string
	dbPass         string
	dbName         string
	eventsRegistry xevent.Registry
}

func (f immudbRepositoryFactory) NewAssetEventRepository() services.RepositoryFactoryFunc[assetdomain.Repository] {
	return func(ctx context.Context) (assetdomain.Repository, func() error, error) {
		port, err := strconv.Atoi(f.dbPort)
		if err != nil {
			return nil, nil, err
		}

		db, err := ximmudb.NewImmuDBClient(ctx, ximmudb.Config{
			Host:         f.dbHost,
			Port:         port,
			User:         f.dbUser,
			Pass:         f.dbPass,
			DB:           "defaultdb",
			Timeout:      DefaultTimeout,
			MaxOpenConns: DefaultMaxOpenConns,
			MaxLife:      DefaultMaxLife,
			MaxIdleCons:  DefaultMaxIdleCons,
		})
		if err != nil {
			return nil, nil, err
		}

		err = ximmudb.Migrate(db, []ximmudb.Migration{
			assetimmudbmigrations.NewCreateAssetsDatabase(),
			assetimmudbmigrations.NewCreateAssetEventsTable(),
		})
		if err != nil {
			return nil, nil, err
		}

		repo, err := assetimmudb.NewImmuRepository(db)
		if err != nil {
			return nil, nil, err
		}

		return repo, func() error {
			return db.Close()
		}, nil
	}
}

func newImmuDBRepositoryFactory(
	cfg services.Config,
	eventsRegistry xevent.Registry,
) immudbRepositoryFactory {
	return immudbRepositoryFactory{
		dbHost:         cfg.DatabaseHost,
		dbPort:         cfg.DatabasePort,
		dbUser:         cfg.DatabaseUser,
		dbPass:         cfg.DatabasePass,
		dbName:         cfg.DatabaseName,
		eventsRegistry: eventsRegistry,
	}
}

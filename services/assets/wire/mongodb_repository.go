package assets

import (
	"context"
	"fmt"

	"github.com/xfrr/finantrack/internal/shared/xevent"
	"github.com/xfrr/finantrack/internal/shared/xmongo"
	"github.com/xfrr/finantrack/services"

	assetdomain "github.com/xfrr/finantrack/internal/contexts/assets/domain"
	assetsmongo "github.com/xfrr/finantrack/internal/contexts/assets/mongodb"
)

type mongoRepositoryFactory struct {
	dbHost         string
	dbPort         string
	dbUser         string
	dbPass         string
	dbName         string
	eventsRegistry xevent.Registry
}

func (f mongoRepositoryFactory) NewAssetEventRepository() services.RepositoryFactoryFunc[assetdomain.Repository] {
	return func(ctx context.Context) (assetdomain.Repository, func() error, error) {
		var repo *assetsmongo.Repository

		ctx, cancel := context.WithTimeout(ctx, xmongo.MongoConnectDefaultTimeout)
		defer cancel()

		mongoClient, err := xmongo.NewClient(ctx, f.buildURI(), f.dbName)
		if err != nil {
			return repo, nil, err
		}

		eventStore, err := xmongo.NewMongoEventStore(ctx, mongoClient, f.eventsRegistry)
		if err != nil {
			return repo, nil, err
		}

		closer := func() error {
			err = mongoClient.Close(ctx)
			if err != nil {
				return err
			}
			return err
		}

		return assetsmongo.NewRepository(eventStore), closer, nil
	}
}

func (f mongoRepositoryFactory) buildURI() string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s",
		f.dbUser,
		f.dbPass,
		f.dbHost,
		f.dbPort,
	)
}

func newMongoRepositoryFactory(
	cfg services.Config,
	eventsRegistry xevent.Registry,
) mongoRepositoryFactory {
	return mongoRepositoryFactory{
		dbHost:         cfg.DatabaseHost,
		dbPort:         cfg.DatabasePort,
		dbUser:         cfg.DatabaseUser,
		dbPass:         cfg.DatabasePass,
		dbName:         cfg.DatabaseName,
		eventsRegistry: eventsRegistry,
	}
}

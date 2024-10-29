package assets

import (
	"context"

	"github.com/xfrr/finantrack/internal/shared/xevent"
	"github.com/xfrr/finantrack/internal/shared/xmongo"
	"github.com/xfrr/finantrack/services"

	assetdomain "github.com/xfrr/finantrack/internal/contexts/assets/domain"
	assetsmongo "github.com/xfrr/finantrack/internal/contexts/assets/mongodb"
)

type mongoRepositoryFactory struct {
	dbURI          string
	dbName         string
	eventsRegistry xevent.Registry
}

func (f mongoRepositoryFactory) NewAssetEventRepository() services.RepositoryFactoryFunc[assetdomain.Repository] {
	return func(ctx context.Context) (assetdomain.Repository, func() error, error) {
		var repo *assetsmongo.Repository

		ctx, cancel := context.WithTimeout(ctx, xmongo.MongoConnectDefaultTimeout)
		defer cancel()

		mongoClient, err := xmongo.NewClient(ctx, f.dbURI, f.dbName)
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

func newMongoRepositoryFactory(
	dbURI, dbName string,
	eventsRegistry xevent.Registry,
) mongoRepositoryFactory {
	return mongoRepositoryFactory{
		dbURI:          dbURI,
		dbName:         dbName,
		eventsRegistry: eventsRegistry,
	}
}

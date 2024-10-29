package assets

import (
	"github.com/xfrr/finantrack/internal/shared/xevent"
	"github.com/xfrr/finantrack/services"

	assetdomain "github.com/xfrr/finantrack/internal/contexts/assets/domain"
)

func newRepositoryFactory(
	dbURI string,
	dbName string,
	eventsRegistry xevent.Registry,
) (services.RepositoryFactory[assetdomain.Repository], error) {
	repoFactory := services.NewRepositoryFactory[assetdomain.Repository]()

	// Register the MongoDB repository
	err := repoFactory.RegisterRepository(
		services.MongoDatabaseEngine,
		newMongoRepositoryFactory(
			dbURI,
			dbName,
			eventsRegistry,
		).NewAssetEventRepository(),
	)
	if err != nil {
		return nil, err
	}

	// Register the MongoDB repository
	err = repoFactory.RegisterRepository(
		services.ImmuDBDatabaseEngine,
		newImmuDBRepositoryFactory(
			dbURI,
			dbName,
			eventsRegistry,
		).NewAssetEventRepository(),
	)

	return repoFactory, nil
}

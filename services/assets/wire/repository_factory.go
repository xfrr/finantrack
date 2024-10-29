package assets

import (
	"github.com/xfrr/finantrack/internal/shared/xevent"
	"github.com/xfrr/finantrack/services"

	assetdomain "github.com/xfrr/finantrack/internal/contexts/assets/domain"
)

func newRepositoryFactory(
	cfg services.Config,
	eventsRegistry xevent.Registry,
) (services.RepositoryFactory[assetdomain.Repository], error) {
	repoFactory := services.NewRepositoryFactory[assetdomain.Repository]()

	// Register the MongoDB repository
	err := repoFactory.RegisterRepository(
		services.MongoDatabaseEngine,
		newMongoRepositoryFactory(
			cfg,
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
			cfg,
			eventsRegistry,
		).NewAssetEventRepository(),
	)
	if err != nil {
		return nil, err
	}

	return repoFactory, nil
}

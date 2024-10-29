package assets

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/xfrr/finantrack/internal/shared/xevent"
	"github.com/xfrr/finantrack/services"

	assetdomain "github.com/xfrr/finantrack/internal/contexts/assets/domain"
	assetimmudb "github.com/xfrr/finantrack/internal/contexts/assets/immudb"
)

type immudbRepositoryFactory struct {
	dbURI          string
	dbName         string
	eventsRegistry xevent.Registry
}

func (f immudbRepositoryFactory) NewAssetEventRepository() services.RepositoryFactoryFunc[assetdomain.Repository] {
	return func(ctx context.Context) (assetdomain.Repository, func() error, error) {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		db, err := sql.Open("immudb", f.dbURI)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		return assetimmudb.NewImmuRepository(db), func() error {
			return nil
		}, nil
	}
}

func newImmuDBRepositoryFactory(
	dbURI string,
	dbName string,
	eventsRegistry xevent.Registry,
) immudbRepositoryFactory {
	return immudbRepositoryFactory{
		dbURI:          dbURI,
		dbName:         dbName,
		eventsRegistry: eventsRegistry,
	}
}

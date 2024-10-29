package assets

import (
	"github.com/xfrr/finantrack/internal/shared/xevent"

	assetevents "github.com/xfrr/finantrack/internal/contexts/assets/domain/events"
)

func newAssetEventsRegistry() xevent.Registry {
	eventsRegistry := xevent.NewPayloadRegistry()
	xevent.Register(eventsRegistry, assetevents.AssetCreatedEventType, func() interface{} {
		return &assetevents.AssetCreatedEvent{}
	})
	xevent.Register(eventsRegistry, assetevents.AssetDeletedEventType, func() interface{} {
		return &assetevents.AssetDeletedEvent{}
	})
	return eventsRegistry
}

package assetdomain

import (
	"errors"

	"github.com/google/uuid"
	"github.com/xfrr/go-cqrsify/aggregate"

	assetevents "github.com/xfrr/finantrack/internal/contexts/assets/domain/events"
)

// AggregateType represents the asset aggregate type.
const AggregateType = "asset"

var (
	// ErrAssetNameIsRequired represents the error when the asset name is required.
	ErrAssetNameIsRequired = errors.New("asset name is required")

	// ErrAssetTypeIsRequired represents the error when the asset type is required.
	ErrAssetTypeIsRequired = errors.New("asset type is required")

	// ErrAssetMoneyIsInvalid represents the error when the asset money is invalid.
	ErrAssetMoneyIsInvalid = errors.New("asset money is invalid")

	// ErrAssetNotFound represents the error when the asset is not found.
	ErrAssetNotFound = errors.New("asset not found")

	// ErrAssetAlreadyExists represents the error when the asset already exists.
	ErrAssetAlreadyExists = errors.New("asset already exists with given identifier")
)

// Asset represents any resource owned or controlled by a business
// or an economic entity that can be converted into cash.
type Asset struct {
	*aggregate.Base[uuid.UUID]

	name      string
	assetType AssetType
	money     Money
	deleted   bool
}

// NewAsset creates a new Asset instance with the given data.
func NewAsset(id uuid.UUID, name string, assetType AssetType, money Money) (*Asset, error) {
	asset := &Asset{
		Base:      aggregate.New(id, AggregateType),
		name:      name,
		assetType: assetType,
		money:     money,
	}

	// Create the asset created event
	aggregate.NextChange(
		asset,
		uuid.New(),
		assetevents.AssetCreatedEventType,
		&assetevents.AssetCreatedEvent{
			AssetID:            id.String(),
			AssetName:          name,
			AssetType:          assetType.String(),
			AssetMoneyAmount:   money.Amount,
			AssetMoneyCurrency: money.Currency.String(),
		},
	)

	// Register the event handlers
	asset.When(assetevents.AssetCreatedEventType, asset.assetCreatedEventHandler)
	asset.When(assetevents.AssetDeletedEventType, asset.assetDeletedEventHandler)

	err := asset.Validate()
	if err != nil {
		return nil, err
	}

	return asset, nil
}

// ID returns the asset ID.
func (a *Asset) ID() uuid.UUID {
	return a.Base.AggregateID()
}

// Name returns the asset name.
func (a *Asset) Name() string {
	return a.name
}

// Type returns the asset type.
func (a *Asset) Type() AssetType {
	return a.assetType
}

// Money returns the asset money.
func (a *Asset) Money() Money {
	return a.money
}

// IsDeleted checks if the asset is deleted.
func (a *Asset) IsDeleted() bool {
	return a.deleted
}

// MarkAsDeleted deletes the asset by given ID.
func (a *Asset) MarkAsDeleted() {
	if a.IsDeleted() {
		return
	}

	aggregate.NextChange(
		a,
		uuid.New(),
		assetevents.AssetDeletedEventType,
		&assetevents.AssetDeletedEvent{
			AssetID: a.ID().String(),
		},
	)
}

// Validate validates the asset.
func (a *Asset) Validate() error {
	if _, err := uuid.Parse(a.ID().String()); err != nil {
		return err
	}

	if a.name == "" {
		return ErrAssetNameIsRequired
	}

	err := a.assetType.Validate()
	if err != nil {
		return err
	}

	err = a.money.Validate()
	if err != nil {
		return err
	}

	return nil
}

func HydrateAsset(id uuid.UUID, events []aggregate.Change) (*Asset, error) {
	asset := &Asset{
		Base: aggregate.New(id, AggregateType),
	}

	// Register the event handlers
	asset.When(assetevents.AssetCreatedEventType, asset.assetCreatedEventHandler)
	asset.When(assetevents.AssetDeletedEventType, asset.assetDeletedEventHandler)

	err := aggregate.Hydrate(asset, events)
	if err != nil {
		return nil, err
	}

	err = asset.Validate()
	if err != nil {
		return nil, err
	}

	return asset, nil
}

// assetCreatedEventHandler is the event handler for the asset created event.
func (a *Asset) assetCreatedEventHandler(event aggregate.Change) {
	evt, ok := event.Payload().(*assetevents.AssetCreatedEvent)
	if !ok {
		return
	}

	a.name = evt.AssetName
	a.assetType = AssetType(evt.AssetType)
	a.money = Money{
		Amount:   evt.AssetMoneyAmount,
		Currency: Currency(evt.AssetMoneyCurrency),
	}
}

// assetDeletedEventHandler is the event handler for the asset deleted event.
func (a *Asset) assetDeletedEventHandler(_ aggregate.Change) {
	a.deleted = true
}

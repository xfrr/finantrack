package assetevents

const AssetDeletedEventType = "asset.deleted"

type AssetDeletedEvent struct {
	AssetID string
}

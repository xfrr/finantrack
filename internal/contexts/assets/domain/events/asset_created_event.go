package assetevents

const AssetCreatedEventType = "asset.created"

type AssetCreatedEvent struct {
	AssetID            string
	AssetType          string
	AssetName          string
	AssetMoneyAmount   float64
	AssetMoneyCurrency string
}

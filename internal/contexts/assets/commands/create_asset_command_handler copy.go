package assetscommands

import (
	"context"

	"github.com/google/uuid"
	assets "github.com/xfrr/finantrack/internal/contexts/assets/domain"
)

type CreateAssetCommand struct {
	AssetID            string
	AssetName          string
	AssetType          string
	AssetMoneyAmount   float64
	AssetMoneyCurrency string
}

func (c CreateAssetCommand) CommandName() string {
	return "CreateAssetCommand"
}

type CreateAssetCommandHandler struct {
	assets assets.Repository
}

func NewCreateAssetCommandHandler(assets assets.Repository) *CreateAssetCommandHandler {
	return &CreateAssetCommandHandler{
		assets: assets,
	}
}

func (h *CreateAssetCommandHandler) Handle(ctx context.Context, cmd CreateAssetCommand) (interface{}, error) {
	var (
		err        error
		asset      *assets.Asset
		assetMoney assets.Money

		assetName = cmd.AssetName
		assetType = assets.AssetType(cmd.AssetType)
	)

	assetID, err := uuid.Parse(cmd.AssetID)
	if err != nil {
		return nil, err
	}

	// Check if the asset already exists
	var ok bool
	if ok, err = h.assets.Exists(ctx, assetID); err != nil {
		return nil, err
	} else if ok {
		return nil, assets.ErrAssetAlreadyExists
	}

	// Creates a new money value object from the given amount and currency
	assetMoney, err = assets.NewMoney(cmd.AssetMoneyAmount, cmd.AssetMoneyCurrency)
	if err != nil {
		return nil, err
	}

	// Creates a new asset entity from the given data
	asset, err = assets.NewAsset(
		assetID,
		assetName,
		assetType,
		assetMoney,
	)
	if err != nil {
		return nil, err
	}

	// TODO: Publish event

	// Save the asset
	return nil, h.assets.Save(ctx, asset)
}

package assetscommands

import (
	"context"

	"github.com/google/uuid"
	assets "github.com/xfrr/finantrack/internal/contexts/assets/domain"
)

type DeleteAssetCommand struct {
	AssetID string
}

func (c DeleteAssetCommand) CommandName() string {
	return "DeleteAssetCommand"
}

type DeleteAssetCommandHandler struct {
	assets assets.Repository
}

func NewDeleteAssetCommandHandler(assets assets.Repository) *DeleteAssetCommandHandler {
	return &DeleteAssetCommandHandler{
		assets: assets,
	}
}

func (h *DeleteAssetCommandHandler) Handle(ctx context.Context, cmd DeleteAssetCommand) (interface{}, error) {
	var (
		err   error
		asset *assets.Asset
	)

	// Parse the asset ID
	assetID, err := uuid.Parse(cmd.AssetID)
	if err != nil {
		return nil, err
	}

	// Get the asset by ID
	asset, err = h.assets.GetByID(ctx, assetID)
	if err != nil {
		return nil, err
	}

	if asset == nil || asset.IsDeleted() {
		return nil, assets.ErrAssetNotFound
	}

	// Mark the asset as deleted
	asset.MarkAsDeleted()

	// Save the asset
	return nil, h.assets.Save(ctx, asset)
}

package assetshttp

import (
	"github.com/gin-gonic/gin"
	"github.com/xfrr/go-cqrsify/cqrs"
)

const ModifyAssetPath = "/assets/:id"

type ModifyAssetHandler struct {
	bus cqrs.Bus
}

func (h *ModifyAssetHandler) Method() string {
	return "PUT"
}

func (h *ModifyAssetHandler) Path() string {
	return ModifyAssetPath
}

func NewModifyAssetHandler(cmdbus cqrs.Bus) *ModifyAssetHandler {
	return &ModifyAssetHandler{
		bus: cmdbus,
	}
}

// @Summary		Modify an asset
// @Description	Modify an asset
// @Tags			assets
// @Accept			json
// @Produce		json
// @Success		200	{object}	string
// @Router			/assets/{id} [put]
// @Param			id		path	string				true	"Asset ID"	default(00000000-0000-0000-0000-000000000000)
// @Param			body	body	ModifyAssetRequest	true	"Asset data"
func (h *ModifyAssetHandler) Handle(c *gin.Context) {
	// TODO: dispatch command to modify asset
}

// ModifyAssetRequest represents the request to modify an asset
type ModifyAssetRequest struct {
	AssetName          string  `json:"asset_name" binding:"required"`
	AssetType          string  `json:"asset_type" binding:"required"`
	AssetMoneyAmount   float64 `json:"asset_money_amount" binding:"required"`
	AssetMoneyCurrency string  `json:"asset_money_currency" binding:"required"`
}

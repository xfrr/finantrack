package assetshttp

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xfrr/go-cqrsify/cqrs"

	assetscommands "github.com/xfrr/finantrack/internal/contexts/assets/commands"
)

const CreateAssetPath = "/assets/:id"

type CreateAssetHandler struct {
	bus cqrs.Bus
}

func (h *CreateAssetHandler) Method() string {
	return "POST"
}

func (h *CreateAssetHandler) Path() string {
	return CreateAssetPath
}

// @Summary		Create a new asset
// @Description	Create a new asset
// @Tags			assets
// @Accept			json
// @Produce		json
// @Success		200	{object}	string
// @Router			/assets/{id} [post]
// @Param			id		path	string				true	"Asset ID"	default(00000000-0000-0000-0000-000000000000)
// @Param			body	body	CreateAssetRequest	true	"Asset data"
func (h *CreateAssetHandler) Handle(c *gin.Context) {
	var req CreateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// dispatch command to create asset
	_, err := cqrs.Dispatch(c.Request.Context(), h.bus, assetscommands.CreateAssetCommand{
		AssetID:            c.Param("id"),
		AssetName:          req.AssetName,
		AssetType:          req.AssetType,
		AssetMoneyAmount:   req.AssetMoneyAmount,
		AssetMoneyCurrency: req.AssetMoneyCurrency,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Asset created"})
}

func NewCreateAssetHandler(cmdbus cqrs.Bus) *CreateAssetHandler {
	return &CreateAssetHandler{
		bus: cmdbus,
	}
}

type CreateAssetRequest struct {
	AssetName          string  `json:"assetName" example:"My Asset"`
	AssetType          string  `json:"assetType" example:"cash"`
	AssetMoneyAmount   float64 `json:"assetMoneyAmount" example:"1000.00"`
	AssetMoneyCurrency string  `json:"assetMoneyCurrency" example:"USD"`
}

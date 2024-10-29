package assetshttp

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xfrr/go-cqrsify/cqrs"

	assetscommands "github.com/xfrr/finantrack/internal/contexts/assets/commands"
)

const DeleteAssetPath = "/assets/:id"

type DeleteAssetHandler struct {
	bus cqrs.Bus
}

func (h *DeleteAssetHandler) Method() string {
	return "DELETE"
}

func (h *DeleteAssetHandler) Path() string {
	return DeleteAssetPath
}

func NewDeleteAssetHandler(cmdbus cqrs.Bus) *DeleteAssetHandler {
	return &DeleteAssetHandler{
		bus: cmdbus,
	}
}

// @Summary		Delete an asset
// @Description	Delete an asset
// @Tags			assets
// @Accept			json
// @Produce		json
// @Success		200	{object}	string
// @Router			/assets/{id} [delete]
// @Param			id	path	string	true	"Asset ID"	default(00000000-0000-0000-0000-000000000000)
func (h *DeleteAssetHandler) Handle(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing asset ID"})
		return
	}

	// dispatch command to create asset
	_, err := cqrs.Dispatch(c.Request.Context(), h.bus, assetscommands.DeleteAssetCommand{
		AssetID: id,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Asset deleted"})
}

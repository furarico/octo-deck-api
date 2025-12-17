package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 指定したコミュニティ取得
// (GET /communities/{id})
func (h *Handler) GetCommunity(c *gin.Context, id string) {
	community, err := h.communityService.GetCommunityByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, convertCommunityToAPI(*community))
}

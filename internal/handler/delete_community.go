package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// コミュニティを削除
// (DELETE /communities/{id})
func (h *Handler) DeleteCommunity(c *gin.Context, id string) {
	if err := h.communityService.DeleteCommunity(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

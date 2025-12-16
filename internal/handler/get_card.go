package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// (GET /cards/{githubId})
func (h *Handler) GetCard(c *gin.Context, githubId string) {
	card, err := h.cardService.GetCardByID(githubId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, convertCardToAPI(*card))
}

package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// (GET /cards/me)
func (h *Handler) GetMyCard(c *gin.Context) {
	card, err := h.cardService.GetMyCard()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, convertCardToAPI(*card))
}

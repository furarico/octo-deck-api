package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// (GET /cards)
func (h *Handler) GetCards(c *gin.Context) {
	cards, err := h.cardService.GetAllCards()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cards": cards,
	})
}

package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// (GET /cards/{id})
func (h *Handler) GetCard(c *gin.Context, id string) {
	card, err := h.cardService.GetCardByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, convertCardToAPI(*card))
}

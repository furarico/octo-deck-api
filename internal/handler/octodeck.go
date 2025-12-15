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

// (GET /cards/me)
func (h *Handler) GetMyCard(c *gin.Context) {
	card, err := h.cardService.GetMyCard()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"card": card,
	})
}

// (GET /cards/{id})
func (h *Handler) GetCard(c *gin.Context, id string) {
	card, err := h.cardService.GetCardByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"card": card,
	})
}

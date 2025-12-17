package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 指定したコミュニティに自分のカードを追加
// (POST /communities/{id}/cards)
func (h *Handler) AddCardToCommunity(c *gin.Context, id string) {
	githubID := c.GetString("github_id")
	if githubID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "github_id is missing from context",
		})
		return
	}

	// github_idから自分のカードを取得
	card, err := h.cardService.GetMyCard(githubID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "card not found",
		})
		return
	}

	// コミュニティにカードを追加
	cardID := uuid.UUID(card.ID).String()
	if err := h.communityService.AddCardToCommunity(id, cardID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

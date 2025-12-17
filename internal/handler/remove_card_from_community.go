package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 指定したコミュニティの自分のカードを削除
// (DELETE /communities/{id}/cards)
func (h *Handler) RemoveCardFromCommunity(c *gin.Context, id string) {
	ctx := c.Request.Context()
	githubClient := getGitHubClient(c)
	githubID := c.GetString("github_id")
	if githubID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "github_id is missing from context",
		})
		return
	}

	// github_idから自分のカードを取得
	card, err := h.cardService.GetMyCard(ctx, githubID, githubClient)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "card not found",
		})
		return
	}

	// コミュニティからカードを削除
	cardID := uuid.UUID(card.ID).String()
	if err := h.communityService.RemoveCardFromCommunity(id, cardID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

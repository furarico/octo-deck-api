package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// (GET /cards/me)
func (h *Handler) GetMyCard(c *gin.Context) {
	ctx := c.Request.Context()
	githubClient := getGitHubClient(c)

	githubID := c.GetString("github_id")
	if githubID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "github_id is missing from context",
		})
		return
	}

	card, err := h.cardService.GetMyCard(ctx, githubID, githubClient)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"card": convertCardToAPI(*card),
	})
}

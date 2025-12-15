package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// (GET /cards/me)
func (h *Handler) GetMyCard(c *gin.Context) {
	// TODO: 認証情報からGitHubIDを取得する
	githubID := c.GetString("github_id")

	card, err := h.cardService.GetMyCard(githubID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, convertCardToAPI(*card))
}

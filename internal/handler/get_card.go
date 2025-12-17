package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// (GET /cards/{githubId})
func (h *Handler) GetCard(c *gin.Context, githubId string) {
	ctx := c.Request.Context()
	githubClient := getGitHubClient(c)

	card, err := h.cardService.GetCardByID(ctx, githubId, githubClient)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, convertCardToAPI(*card))
}

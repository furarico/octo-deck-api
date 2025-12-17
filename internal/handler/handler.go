package handler

import (
	"github.com/furarico/octo-deck-api/internal/github"
	"github.com/furarico/octo-deck-api/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	cardService *service.CardService
}

func NewHandler(cardService *service.CardService) *Handler {
	return &Handler{
		cardService: cardService,
	}
}

func getGitHubClient(c *gin.Context) *github.Client {
	return c.MustGet("github_client").(*github.Client)
}

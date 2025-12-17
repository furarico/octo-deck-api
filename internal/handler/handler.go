package handler

import (
	"context"

	"github.com/furarico/octo-deck-api/internal/domain"
	"github.com/furarico/octo-deck-api/internal/github"
	"github.com/gin-gonic/gin"
)

// CardServiceInterface はハンドラーが必要とするサービスのインターフェース
type CardServiceInterface interface {
	GetAllCards(ctx context.Context, githubID string, githubClient *github.Client) ([]domain.Card, error)
	GetCardByGitHubID(ctx context.Context, githubID string, githubClient *github.Client) (*domain.Card, error)
	GetMyCard(ctx context.Context, githubID string, githubClient *github.Client) (*domain.Card, error)
}

type Handler struct {
	cardService CardServiceInterface
}

func NewHandler(cardService CardServiceInterface) *Handler {
	return &Handler{
		cardService: cardService,
	}
}

func getGitHubClient(c *gin.Context) *github.Client {
	return c.MustGet("github_client").(*github.Client)
}

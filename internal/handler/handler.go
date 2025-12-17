package handler

import (
	"context"

	"github.com/furarico/octo-deck-api/internal/domain"
	"github.com/furarico/octo-deck-api/internal/github"
	"github.com/furarico/octo-deck-api/internal/service"
	"github.com/gin-gonic/gin"
)

// CardServiceInterface はハンドラーが必要とするサービスのインターフェース
type CardServiceInterface interface {
	GetAllCards(ctx context.Context, githubID string, githubClient *github.Client) ([]domain.Card, error)
	GetCardByGitHubID(ctx context.Context, githubID string, githubClient *github.Client) (*domain.Card, error)
	GetMyCard(ctx context.Context, githubID string, githubClient *github.Client) (*domain.Card, error)
}

type Handler struct {
	cardService      CardServiceInterface
	communityService *service.CommunityService
}

func NewHandler(cardService CardServiceInterface, communityService *service.CommunityService) *Handler {
	return &Handler{
		cardService:      cardService,
		communityService: communityService,
	}
}

func NewCardHandler(cardService CardServiceInterface) *Handler {
	return &Handler{cardService: cardService}
}

func NewCommunityHandler(communityService *service.CommunityService) *Handler {
	return &Handler{communityService: communityService}
}

func getGitHubClient(c *gin.Context) *github.Client {
	return c.MustGet("github_client").(*github.Client)
}

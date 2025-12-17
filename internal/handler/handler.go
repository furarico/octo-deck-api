package handler

import (
	"context"
	"errors"

	"github.com/furarico/octo-deck-api/internal/domain"
	"github.com/furarico/octo-deck-api/internal/github"
	"github.com/furarico/octo-deck-api/internal/service"
	"github.com/gin-gonic/gin"
)

// Context keys for storing values
type contextKey string

const (
	GitHubClientKey contextKey = "github_client"
	GitHubIDKey     contextKey = "github_id"
	GitHubLoginKey  contextKey = "github_login"
	GitHubNameKey   contextKey = "github_name"
	GitHubAvatarKey contextKey = "github_avatar_url"
)

// CardServiceInterface はハンドラーが必要とするサービスのインターフェース
type CardServiceInterface interface {
	GetAllCards(ctx context.Context, githubID string, githubClient *github.Client) ([]domain.Card, error)
	GetCardByGitHubID(ctx context.Context, githubID string, githubClient *github.Client) (*domain.Card, error)
	GetMyCard(ctx context.Context, githubID string, githubClient *github.Client) (*domain.Card, error)
	GetOrCreateMyCard(ctx context.Context, githubID string, githubClient *github.Client) (*domain.Card, error)
	AddCardToDeck(ctx context.Context, collectorGithubID string, targetGithubID string, githubClient *github.Client) (*domain.Card, error)
}

// StatsServiceInterface はハンドラーが必要とする統計サービスのインターフェース
type StatsServiceInterface interface {
	GetUserStats(ctx context.Context, githubID string, githubClient *github.Client) (*github.ContributionStats, error)
}

type Handler struct {
	cardService      CardServiceInterface
	communityService *service.CommunityService
	statsService     StatsServiceInterface
}

func NewHandler(cardService CardServiceInterface, communityService *service.CommunityService, statsService StatsServiceInterface) *Handler {
	return &Handler{
		cardService:      cardService,
		communityService: communityService,
		statsService:     statsService,
	}
}

func NewCardHandler(cardService CardServiceInterface) *Handler {
	return &Handler{cardService: cardService}
}

func NewCommunityHandler(communityService *service.CommunityService) *Handler {
	return &Handler{communityService: communityService}
}

func NewStatsHandler(statsService StatsServiceInterface) *Handler {
	return &Handler{statsService: statsService}
}

// getRequestContext extracts the underlying request context from gin.Context or returns ctx as-is
func getRequestContext(ctx context.Context) context.Context {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		return ginCtx.Request.Context()
	}
	return ctx
}

// getGitHubClient retrieves the GitHub client from context
func getGitHubClient(ctx context.Context) (*github.Client, error) {
	reqCtx := getRequestContext(ctx)
	client, ok := reqCtx.Value(GitHubClientKey).(*github.Client)
	if !ok {
		return nil, errors.New("github_client not found in context")
	}
	return client, nil
}

// getGitHubID retrieves the GitHub ID from context
func getGitHubID(ctx context.Context) (string, error) {
	reqCtx := getRequestContext(ctx)
	id, ok := reqCtx.Value(GitHubIDKey).(string)
	if !ok || id == "" {
		return "", errors.New("github_id not found in context")
	}
	return id, nil
}

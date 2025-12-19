package handler

import (
	"context"
	"errors"

	"github.com/furarico/octo-deck-api/internal/domain"
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
	GetAllCards(ctx context.Context, githubID string, githubClient service.GitHubClient) ([]domain.Card, error)
	GetCardByGitHubID(ctx context.Context, githubID string, githubClient service.GitHubClient) (*domain.Card, error)
	GetMyCard(ctx context.Context, githubID string, githubClient service.GitHubClient) (*domain.Card, error)
	GetOrCreateMyCard(ctx context.Context, githubID string, githubClient service.GitHubClient) (*domain.Card, error)
	AddCardToDeck(ctx context.Context, collectorGithubID string, targetGithubID string, githubClient service.GitHubClient) (*domain.Card, error)
	RemoveCardFromDeck(ctx context.Context, collectorGithubID string, targetGithubID string, githubClient service.GitHubClient) (*domain.Card, error)
}

// StatsServiceInterface はハンドラーが必要とする統計サービスのインターフェース
type StatsServiceInterface interface {
	GetUserStats(ctx context.Context, githubID string, githubClient service.GitHubClient) (*domain.Stats, error)
}

// CommunityServiceInterface はハンドラーが必要とするコミュニティサービスのインターフェース
type CommunityServiceInterface interface {
	GetAllCommunities(githubID string) ([]domain.Community, error)
	GetCommunityByID(id string) (*domain.Community, error)
	GetCommunityWithHighlightedCard(ctx context.Context, id string, githubClient service.GitHubClient) (*domain.Community, *domain.HighlightedCard, error)
	GetCommunityCards(ctx context.Context, id string, githubClient service.GitHubClient) ([]domain.Card, error)
	CreateCommunity(name string) (*domain.Community, error)
	DeleteCommunity(id string) error
	AddCardToCommunity(communityID string, cardID string) error
	RemoveCardFromCommunity(communityID string, cardID string) error
}

type Handler struct {
	cardService      CardServiceInterface
	communityService CommunityServiceInterface
	statsService     StatsServiceInterface
}

func NewHandler(cardService CardServiceInterface, communityService CommunityServiceInterface, statsService StatsServiceInterface) *Handler {
	return &Handler{
		cardService:      cardService,
		communityService: communityService,
		statsService:     statsService,
	}
}

func NewCardHandler(cardService CardServiceInterface) *Handler {
	return &Handler{cardService: cardService}
}

func NewCommunityHandler(communityService CommunityServiceInterface) *Handler {
	return &Handler{communityService: communityService}
}

func NewStatsHandler(statsService StatsServiceInterface) *Handler {
	return &Handler{statsService: statsService}
}

// gin.Contextからcontext.Contextを取得するためのヘルパー関数
func getRequestContext(ctx context.Context) context.Context {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		return ginCtx.Request.Context()
	}
	return ctx
}

// context.ContextからGitHub Clientを取得するためのヘルパー関数
func getGitHubClient(ctx context.Context) (service.GitHubClient, error) {
	reqCtx := getRequestContext(ctx)
	client, ok := reqCtx.Value(GitHubClientKey).(service.GitHubClient)
	if !ok {
		return nil, errors.New("github_client not found in context")
	}
	return client, nil
}

// context.ContextからGitHub IDを取得するためのヘルパー関数
func getGitHubID(ctx context.Context) (string, error) {
	reqCtx := getRequestContext(ctx)
	id, ok := reqCtx.Value(GitHubIDKey).(string)
	if !ok || id == "" {
		return "", errors.New("github_id not found in context")
	}
	return id, nil
}

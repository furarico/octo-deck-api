package service

import (
	"context"

	"github.com/furarico/octo-deck-api/internal/domain"
	"github.com/furarico/octo-deck-api/internal/github"
)

// MockCardService はテスト用のモックサービス
type MockCardService struct {
	GetAllCardsFunc       func(ctx context.Context, githubID string, githubClient *github.Client) ([]domain.Card, error)
	GetCardByGitHubIDFunc func(ctx context.Context, githubID string, githubClient *github.Client) (*domain.Card, error)
	GetMyCardFunc         func(ctx context.Context, githubID string, githubClient *github.Client) (*domain.Card, error)
	GetOrCreateMyCardFunc func(ctx context.Context, githubID string, githubClient *github.Client) (*domain.Card, error)
	AddCardToDeckFunc     func(ctx context.Context, collectorGithubID string, targetGithubID string, githubClient *github.Client) (*domain.Card, error)
}

func NewMockCardService() *MockCardService {
	return &MockCardService{}
}

func (m *MockCardService) GetAllCards(ctx context.Context, githubID string, githubClient *github.Client) ([]domain.Card, error) {
	if m.GetAllCardsFunc != nil {
		return m.GetAllCardsFunc(ctx, githubID, githubClient)
	}
	return []domain.Card{}, nil
}

func (m *MockCardService) GetCardByGitHubID(ctx context.Context, githubID string, githubClient *github.Client) (*domain.Card, error) {
	if m.GetCardByGitHubIDFunc != nil {
		return m.GetCardByGitHubIDFunc(ctx, githubID, githubClient)
	}
	return &domain.Card{}, nil
}

func (m *MockCardService) GetMyCard(ctx context.Context, githubID string, githubClient *github.Client) (*domain.Card, error) {
	if m.GetMyCardFunc != nil {
		return m.GetMyCardFunc(ctx, githubID, githubClient)
	}
	return nil, nil
}

func (m *MockCardService) GetOrCreateMyCard(ctx context.Context, githubID string, githubClient *github.Client) (*domain.Card, error) {
	if m.GetOrCreateMyCardFunc != nil {
		return m.GetOrCreateMyCardFunc(ctx, githubID, githubClient)
	}
	return nil, nil
}

func (m *MockCardService) AddCardToDeck(ctx context.Context, collectorGithubID string, targetGithubID string, githubClient *github.Client) (*domain.Card, error) {
	if m.AddCardToDeckFunc != nil {
		return m.AddCardToDeckFunc(ctx, collectorGithubID, targetGithubID, githubClient)
	}
	return nil, nil
}

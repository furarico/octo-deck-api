package service

import (
	"context"

	"github.com/furarico/octo-deck-api/internal/domain"
)

// MockCardService はテスト用のモックサービス
type MockCardService struct {
	GetAllCardsFunc        func(githubID string) ([]domain.Card, error)
	GetCardByGitHubIDFunc  func(ctx context.Context, githubID string, githubClient GitHubClient) (*domain.Card, error)
	GetMyCardFunc          func(ctx context.Context, githubID string, githubClient GitHubClient) (*domain.Card, error)
	GetOrCreateMyCardFunc  func(ctx context.Context, githubID string, nodeID string, githubClient GitHubClient) (*domain.Card, error)
	AddCardToDeckFunc      func(ctx context.Context, collectorGithubID string, targetGithubID string, githubClient GitHubClient) (*domain.Card, error)
	RemoveCardFromDeckFunc func(ctx context.Context, collectorGithubID string, targetGithubID string, githubClient GitHubClient) (*domain.Card, error)
}

func NewMockCardService() *MockCardService {
	return &MockCardService{}
}

func (m *MockCardService) GetAllCards(githubID string) ([]domain.Card, error) {
	if m.GetAllCardsFunc != nil {
		return m.GetAllCardsFunc(githubID)
	}
	return []domain.Card{}, nil
}

func (m *MockCardService) GetCardByGitHubID(ctx context.Context, githubID string, githubClient GitHubClient) (*domain.Card, error) {
	if m.GetCardByGitHubIDFunc != nil {
		return m.GetCardByGitHubIDFunc(ctx, githubID, githubClient)
	}
	return &domain.Card{}, nil
}

func (m *MockCardService) GetMyCard(ctx context.Context, githubID string, githubClient GitHubClient) (*domain.Card, error) {
	if m.GetMyCardFunc != nil {
		return m.GetMyCardFunc(ctx, githubID, githubClient)
	}
	return nil, nil
}

func (m *MockCardService) GetOrCreateMyCard(ctx context.Context, githubID string, nodeID string, githubClient GitHubClient) (*domain.Card, error) {
	if m.GetOrCreateMyCardFunc != nil {
		return m.GetOrCreateMyCardFunc(ctx, githubID, nodeID, githubClient)
	}
	return nil, nil
}

func (m *MockCardService) AddCardToDeck(ctx context.Context, collectorGithubID string, targetGithubID string, githubClient GitHubClient) (*domain.Card, error) {
	if m.AddCardToDeckFunc != nil {
		return m.AddCardToDeckFunc(ctx, collectorGithubID, targetGithubID, githubClient)
	}
	return nil, nil
}

func (m *MockCardService) RemoveCardFromDeck(ctx context.Context, collectorGithubID string, targetGithubID string, githubClient GitHubClient) (*domain.Card, error) {
	if m.RemoveCardFromDeckFunc != nil {
		return m.RemoveCardFromDeckFunc(ctx, collectorGithubID, targetGithubID, githubClient)
	}
	return nil, nil
}

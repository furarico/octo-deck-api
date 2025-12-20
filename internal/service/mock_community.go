package service

import (
	"context"
	"time"

	"github.com/furarico/octo-deck-api/internal/domain"
)

// MockCommunityService はテスト用のモックサービス
type MockCommunityService struct {
	GetAllCommunitiesFunc               func(ctx context.Context, githubID string) ([]domain.Community, error)
	GetCommunityByIDFunc                func(ctx context.Context, id string) (*domain.Community, error)
	GetCommunityWithHighlightedCardFunc func(ctx context.Context, id string) (*domain.Community, *domain.HighlightedCard, error)
	RefreshHighlightedCardFunc          func(ctx context.Context, id string, githubClient GitHubClient) (*domain.Community, *domain.HighlightedCard, error)
	GetCommunityCardsFunc               func(ctx context.Context, id string) ([]domain.Card, error)
	CreateCommunityWithPeriodFunc       func(ctx context.Context, name string, startDateTime, endDateTime time.Time) (*domain.Community, error)
	DeleteCommunityFunc                 func(ctx context.Context, id string) error
	AddCardToCommunityFunc              func(ctx context.Context, communityID string, cardID string) error
	RemoveCardFromCommunityFunc         func(ctx context.Context, communityID string, cardID string) error
}

func NewMockCommunityService() *MockCommunityService {
	return &MockCommunityService{}
}

func (m *MockCommunityService) GetAllCommunities(ctx context.Context, githubID string) ([]domain.Community, error) {
	if m.GetAllCommunitiesFunc != nil {
		return m.GetAllCommunitiesFunc(ctx, githubID)
	}
	return []domain.Community{}, nil
}

func (m *MockCommunityService) GetCommunityByID(ctx context.Context, id string) (*domain.Community, error) {
	if m.GetCommunityByIDFunc != nil {
		return m.GetCommunityByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockCommunityService) GetCommunityWithHighlightedCard(ctx context.Context, id string) (*domain.Community, *domain.HighlightedCard, error) {
	if m.GetCommunityWithHighlightedCardFunc != nil {
		return m.GetCommunityWithHighlightedCardFunc(ctx, id)
	}
	return nil, nil, nil
}

func (m *MockCommunityService) RefreshHighlightedCard(ctx context.Context, id string, githubClient GitHubClient) (*domain.Community, *domain.HighlightedCard, error) {
	if m.RefreshHighlightedCardFunc != nil {
		return m.RefreshHighlightedCardFunc(ctx, id, githubClient)
	}
	return nil, nil, nil
}

func (m *MockCommunityService) GetCommunityCards(ctx context.Context, id string) ([]domain.Card, error) {
	if m.GetCommunityCardsFunc != nil {
		return m.GetCommunityCardsFunc(ctx, id)
	}
	return []domain.Card{}, nil
}

func (m *MockCommunityService) CreateCommunityWithPeriod(ctx context.Context, name string, startDateTime, endDateTime time.Time) (*domain.Community, error) {
	if m.CreateCommunityWithPeriodFunc != nil {
		return m.CreateCommunityWithPeriodFunc(ctx, name, startDateTime, endDateTime)
	}
	return nil, nil
}

func (m *MockCommunityService) DeleteCommunity(ctx context.Context, id string) error {
	if m.DeleteCommunityFunc != nil {
		return m.DeleteCommunityFunc(ctx, id)
	}
	return nil
}

func (m *MockCommunityService) AddCardToCommunity(ctx context.Context, communityID string, cardID string) error {
	if m.AddCardToCommunityFunc != nil {
		return m.AddCardToCommunityFunc(ctx, communityID, cardID)
	}
	return nil
}

func (m *MockCommunityService) RemoveCardFromCommunity(ctx context.Context, communityID string, cardID string) error {
	if m.RemoveCardFromCommunityFunc != nil {
		return m.RemoveCardFromCommunityFunc(ctx, communityID, cardID)
	}
	return nil
}

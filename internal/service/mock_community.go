package service

import (
	"context"

	"github.com/furarico/octo-deck-api/internal/domain"
	"github.com/furarico/octo-deck-api/internal/github"
)

// MockCommunityService はテスト用のモックサービス
type MockCommunityService struct {
	GetAllCommunitiesFunc              func(githubID string) ([]domain.Community, error)
	GetCommunityByIDFunc               func(id string) (*domain.Community, error)
	GetCommunityWithHighlightedCardFunc func(ctx context.Context, id string, githubClient *github.Client) (*domain.Community, *domain.HighlightedCard, error)
	GetCommunityCardsFunc              func(id string) ([]domain.Card, error)
	CreateCommunityFunc                func(name string) (*domain.Community, error)
	DeleteCommunityFunc                func(id string) error
	AddCardToCommunityFunc             func(communityID string, cardID string) error
	RemoveCardFromCommunityFunc        func(communityID string, cardID string) error
}

func NewMockCommunityService() *MockCommunityService {
	return &MockCommunityService{}
}

func (m *MockCommunityService) GetAllCommunities(githubID string) ([]domain.Community, error) {
	if m.GetAllCommunitiesFunc != nil {
		return m.GetAllCommunitiesFunc(githubID)
	}
	return []domain.Community{}, nil
}

func (m *MockCommunityService) GetCommunityByID(id string) (*domain.Community, error) {
	if m.GetCommunityByIDFunc != nil {
		return m.GetCommunityByIDFunc(id)
	}
	return nil, nil
}

func (m *MockCommunityService) GetCommunityWithHighlightedCard(ctx context.Context, id string, githubClient *github.Client) (*domain.Community, *domain.HighlightedCard, error) {
	if m.GetCommunityWithHighlightedCardFunc != nil {
		return m.GetCommunityWithHighlightedCardFunc(ctx, id, githubClient)
	}
	return nil, nil, nil
}

func (m *MockCommunityService) GetCommunityCards(id string) ([]domain.Card, error) {
	if m.GetCommunityCardsFunc != nil {
		return m.GetCommunityCardsFunc(id)
	}
	return []domain.Card{}, nil
}

func (m *MockCommunityService) CreateCommunity(name string) (*domain.Community, error) {
	if m.CreateCommunityFunc != nil {
		return m.CreateCommunityFunc(name)
	}
	return nil, nil
}

func (m *MockCommunityService) DeleteCommunity(id string) error {
	if m.DeleteCommunityFunc != nil {
		return m.DeleteCommunityFunc(id)
	}
	return nil
}

func (m *MockCommunityService) AddCardToCommunity(communityID string, cardID string) error {
	if m.AddCardToCommunityFunc != nil {
		return m.AddCardToCommunityFunc(communityID, cardID)
	}
	return nil
}

func (m *MockCommunityService) RemoveCardFromCommunity(communityID string, cardID string) error {
	if m.RemoveCardFromCommunityFunc != nil {
		return m.RemoveCardFromCommunityFunc(communityID, cardID)
	}
	return nil
}

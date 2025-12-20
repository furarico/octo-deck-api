package repository

import (
	"context"

	"github.com/furarico/octo-deck-api/internal/domain"
)

type MockCommunityRepository struct {
	FindAllFunc                     func(ctx context.Context, githubID string) ([]domain.Community, error)
	FindByIDFunc                    func(ctx context.Context, id string) (*domain.Community, error)
	FindByIDWithHighlightedCardFunc func(ctx context.Context, id string) (*domain.Community, error)
	FindCardsFunc                   func(ctx context.Context, id string) ([]domain.Card, error)
	CreateFunc                      func(ctx context.Context, community *domain.Community) error
	DeleteFunc                      func(ctx context.Context, id string) error
	AddCardFunc                     func(ctx context.Context, communityID string, cardID string) error
	RemoveCardFunc                  func(ctx context.Context, communityID string, cardID string) error
	UpdateHighlightedCardFunc       func(ctx context.Context, communityID string, highlightedCard *domain.HighlightedCard) error
}

func NewMockCommunityRepository() *MockCommunityRepository {
	return &MockCommunityRepository{}
}

// FindAll はすべてのコミュニティを取得する
func (r *MockCommunityRepository) FindAll(ctx context.Context, githubID string) ([]domain.Community, error) {
	if r.FindAllFunc != nil {
		return r.FindAllFunc(ctx, githubID)
	}
	return []domain.Community{}, nil
}

// FindByID は指定されたIDのコミュニティを取得する
func (r *MockCommunityRepository) FindByID(ctx context.Context, id string) (*domain.Community, error) {
	if r.FindByIDFunc != nil {
		return r.FindByIDFunc(ctx, id)
	}
	return nil, nil
}

// FindCards は指定したコミュニティIDのカード一覧を取得する
func (r *MockCommunityRepository) FindCards(ctx context.Context, id string) ([]domain.Card, error) {
	if r.FindCardsFunc != nil {
		return r.FindCardsFunc(ctx, id)
	}
	return []domain.Card{}, nil
}

// Create はコミュニティを作成する
func (r *MockCommunityRepository) Create(ctx context.Context, community *domain.Community) error {
	if r.CreateFunc != nil {
		return r.CreateFunc(ctx, community)
	}
	return nil
}

// Delete はコミュニティを削除する
func (r *MockCommunityRepository) Delete(ctx context.Context, id string) error {
	if r.DeleteFunc != nil {
		return r.DeleteFunc(ctx, id)
	}
	return nil
}

// AddCard はコミュニティにカードを追加する
func (r *MockCommunityRepository) AddCard(ctx context.Context, communityID string, cardID string) error {
	if r.AddCardFunc != nil {
		return r.AddCardFunc(ctx, communityID, cardID)
	}
	return nil
}

// RemoveCard はコミュニティからカードを削除する
func (r *MockCommunityRepository) RemoveCard(ctx context.Context, communityID string, cardID string) error {
	if r.RemoveCardFunc != nil {
		return r.RemoveCardFunc(ctx, communityID, cardID)
	}
	return nil
}

// FindByIDWithHighlightedCard は指定されたIDのコミュニティをHighlightedCard付きで取得する
func (r *MockCommunityRepository) FindByIDWithHighlightedCard(ctx context.Context, id string) (*domain.Community, error) {
	if r.FindByIDWithHighlightedCardFunc != nil {
		return r.FindByIDWithHighlightedCardFunc(ctx, id)
	}
	return nil, nil
}

// UpdateHighlightedCard はコミュニティのHighlightedCardを更新する
func (r *MockCommunityRepository) UpdateHighlightedCard(ctx context.Context, communityID string, highlightedCard *domain.HighlightedCard) error {
	if r.UpdateHighlightedCardFunc != nil {
		return r.UpdateHighlightedCardFunc(ctx, communityID, highlightedCard)
	}
	return nil
}

package repository

import (
	"context"

	"github.com/furarico/octo-deck-api/internal/domain"
)

type MockCardRepository struct {
	FindAllFunc                  func(ctx context.Context, githubID string) ([]domain.Card, error)
	FindByGitHubIDFunc           func(ctx context.Context, githubID string) (*domain.Card, error)
	FindMyCardFunc               func(ctx context.Context, githubID string) (*domain.Card, error)
	FindAllCardsInDBFunc         func(ctx context.Context) ([]domain.Card, error)
	CreateFunc                   func(ctx context.Context, card *domain.Card) error
	UpdateFunc                   func(ctx context.Context, card *domain.Card) error
	AddToCollectedCardsFunc      func(ctx context.Context, collectorGithubID string, cardID domain.CardID) error
	RemoveFromCollectedCardsFunc func(ctx context.Context, collectorGithubID string, cardID domain.CardID) error
}

func NewMockCardRepository() *MockCardRepository {
	return &MockCardRepository{}
}

// FindAll は自分が集めたカードを全て取得する
func (r *MockCardRepository) FindAll(ctx context.Context, githubID string) ([]domain.Card, error) {
	if r.FindAllFunc != nil {
		return r.FindAllFunc(ctx, githubID)
	}

	return []domain.Card{}, nil
}

// FindByGitHubID はGitHub IDでカードを取得する
func (r *MockCardRepository) FindByGitHubID(ctx context.Context, githubID string) (*domain.Card, error) {
	if r.FindByGitHubIDFunc != nil {
		return r.FindByGitHubIDFunc(ctx, githubID)
	}
	return &domain.Card{}, nil
}

// FindMyCard は自分のカードを取得する
func (r *MockCardRepository) FindMyCard(ctx context.Context, githubID string) (*domain.Card, error) {
	if r.FindMyCardFunc != nil {
		return r.FindMyCardFunc(ctx, githubID)
	}
	return nil, nil
}

// FindAllCardsInDB はデータベース内の全カードを取得する
func (r *MockCardRepository) FindAllCardsInDB(ctx context.Context) ([]domain.Card, error) {
	if r.FindAllCardsInDBFunc != nil {
		return r.FindAllCardsInDBFunc(ctx)
	}
	return []domain.Card{}, nil
}

// Create は新しいカードを作成する
func (r *MockCardRepository) Create(ctx context.Context, card *domain.Card) error {
	if r.CreateFunc != nil {
		return r.CreateFunc(ctx, card)
	}
	return nil
}

// AddToCollectedCards はカードをデッキに追加する
func (r *MockCardRepository) AddToCollectedCards(ctx context.Context, collectorGithubID string, cardID domain.CardID) error {
	if r.AddToCollectedCardsFunc != nil {
		return r.AddToCollectedCardsFunc(ctx, collectorGithubID, cardID)
	}
	return nil
}

// RemoveFromCollectedCards はカードをデッキから削除する
func (r *MockCardRepository) RemoveFromCollectedCards(ctx context.Context, collectorGithubID string, cardID domain.CardID) error {
	if r.RemoveFromCollectedCardsFunc != nil {
		return r.RemoveFromCollectedCardsFunc(ctx, collectorGithubID, cardID)
	}
	return nil
}

// Update はカード情報を更新する
func (r *MockCardRepository) Update(ctx context.Context, card *domain.Card) error {
	if r.UpdateFunc != nil {
		return r.UpdateFunc(ctx, card)
	}
	return nil
}

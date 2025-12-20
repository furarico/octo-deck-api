package repository

import "github.com/furarico/octo-deck-api/internal/domain"

type MockCardRepository struct {
	FindAllFunc                  func(githubID string) ([]domain.Card, error)
	FindByGitHubIDFunc           func(githubID string) (*domain.Card, error)
	FindMyCardFunc               func(githubID string) (*domain.Card, error)
	FindAllCardsInDBFunc         func() ([]domain.Card, error)
	CreateFunc                   func(card *domain.Card) error
	UpdateFunc                   func(card *domain.Card) error
	AddToCollectedCardsFunc      func(collectorGithubID string, cardID domain.CardID) error
	RemoveFromCollectedCardsFunc func(collectorGithubID string, cardID domain.CardID) error
}

func NewMockCardRepository() *MockCardRepository {
	return &MockCardRepository{}
}

// FindAll は自分が集めたカードを全て取得する
func (r *MockCardRepository) FindAll(githubID string) ([]domain.Card, error) {
	if r.FindAllFunc != nil {
		return r.FindAllFunc(githubID)
	}

	return []domain.Card{}, nil
}

// FindByGitHubID はGitHub IDでカードを取得する
func (r *MockCardRepository) FindByGitHubID(githubID string) (*domain.Card, error) {
	if r.FindByGitHubIDFunc != nil {
		return r.FindByGitHubIDFunc(githubID)
	}
	return &domain.Card{}, nil
}

// FindMyCard は自分のカードを取得する
func (r *MockCardRepository) FindMyCard(githubID string) (*domain.Card, error) {
	if r.FindMyCardFunc != nil {
		return r.FindMyCardFunc(githubID)
	}
	return nil, nil
}

// FindAllCardsInDB はデータベース内の全カードを取得する
func (r *MockCardRepository) FindAllCardsInDB() ([]domain.Card, error) {
	if r.FindAllCardsInDBFunc != nil {
		return r.FindAllCardsInDBFunc()
	}
	return []domain.Card{}, nil
}

// Create は新しいカードを作成する
func (r *MockCardRepository) Create(card *domain.Card) error {
	if r.CreateFunc != nil {
		return r.CreateFunc(card)
	}
	return nil
}

// AddToCollectedCards はカードをデッキに追加する
func (r *MockCardRepository) AddToCollectedCards(collectorGithubID string, cardID domain.CardID) error {
	if r.AddToCollectedCardsFunc != nil {
		return r.AddToCollectedCardsFunc(collectorGithubID, cardID)
	}
	return nil
}

// RemoveFromCollectedCards はカードをデッキから削除する
func (r *MockCardRepository) RemoveFromCollectedCards(collectorGithubID string, cardID domain.CardID) error {
	if r.RemoveFromCollectedCardsFunc != nil {
		return r.RemoveFromCollectedCardsFunc(collectorGithubID, cardID)
	}
	return nil
}

// Update はカード情報を更新する
func (r *MockCardRepository) Update(card *domain.Card) error {
	if r.UpdateFunc != nil {
		return r.UpdateFunc(card)
	}
	return nil
}

package repository

import "github.com/furarico/octo-deck-api/internal/domain"

type MockCardRepository struct {
	FindAllFunc        func(githubID string) ([]domain.Card, error)
	FindByGitHubIDFunc func(githubID string) (*domain.Card, error)
	FindMyCardFunc     func(githubID string) (*domain.Card, error)
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

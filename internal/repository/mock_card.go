package repository

import "github.com/furarico/octo-deck-api/internal/domain"

type MockCardRepository struct {
	FindAllFunc    func(githubID string) ([]domain.CardWithOwner, error)
	FindByIDFunc   func(id string) (*domain.CardWithOwner, error)
	FindMyCardFunc func(githubID string) (*domain.CardWithOwner, error)
}

func NewMockCardRepository() *MockCardRepository {
	return &MockCardRepository{}
}

// FindAll は自分が集めたカードを全て取得する
func (r *MockCardRepository) FindAll(githubID string) ([]domain.CardWithOwner, error) {
	if r.FindAllFunc != nil {
		return r.FindAllFunc(githubID)
	}

	return []domain.CardWithOwner{}, nil
}

// FindByID は指定されたIDのカードを取得する
func (r *MockCardRepository) FindByID(id string) (*domain.CardWithOwner, error) {
	if r.FindByIDFunc != nil {
		return r.FindByIDFunc(id)
	}
	return &domain.CardWithOwner{}, nil
}

// FindMyCard は自分のカードを取得する
func (r *MockCardRepository) FindMyCard(githubID string) (*domain.CardWithOwner, error) {
	if r.FindMyCardFunc != nil {
		return r.FindMyCardFunc(githubID)
	}
	return nil, nil
}

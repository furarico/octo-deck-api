package repository

import "github.com/furarico/octo-deck-api/internal/domain"

type MockCardRepository struct {
	FindAllFunc    func() ([]domain.CardWithOwner, error)
	FindByIDFunc   func(id string) (*domain.CardWithOwner, error)
	FindMyCardFunc func() (*domain.CardWithOwner, error)
}

func NewMockCardRepository() *MockCardRepository {
	return &MockCardRepository{}
}

// FindAll は全てのカードを取得する
func (r *MockCardRepository) FindAll() ([]domain.CardWithOwner, error) {

	return []domain.CardWithOwner{}, nil
}

// FindByID は指定されたIDのカードを取得する
func (r *MockCardRepository) FindByID(id string) (*domain.CardWithOwner, error) {

	return &domain.CardWithOwner{}, nil
}

// FindMyCard は自分のカードを取得する
func (r *MockCardRepository) FindMyCard() (*domain.CardWithOwner, error) {

	return &domain.CardWithOwner{}, nil
}

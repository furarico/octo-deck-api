package repository

import "github.com/furarico/octo-deck-api/internal/domain"

// MockCardRepository は service.CardRepository インターフェースを実装する
type MockCardRepository struct{}

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

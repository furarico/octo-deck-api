package repository

import (
	"github.com/furarico/octo-deck-api/internal/domain"
	"gorm.io/gorm"
)

type cardRepository struct {
	db *gorm.DB
}

func NewCardRepository(db *gorm.DB) *cardRepository {
	return &cardRepository{db: db}
}

func (r *cardRepository) FindAll() ([]domain.Card, error) {
	var cards []domain.Card
	if err := r.db.Find(&cards).Error; err != nil {
		return nil, err
	}
	return cards, nil
}

func (r *cardRepository) FindByID(id string) (*domain.Card, error) {
	var card domain.Card
	if err := r.db.First(&card, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &card, nil
}

// TODO: ユーザーIDを取得してカードを取得する
func (r *cardRepository) FindMyCard() (*domain.Card, error) {
	var card domain.Card
	card = domain.Card{
		ID:      domain.NewCardID(),
		OwnerID: domain.NewUserID(),
	}
	return &card, nil
}

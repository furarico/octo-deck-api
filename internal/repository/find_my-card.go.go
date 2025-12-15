package repository

import (
	"github.com/furarico/octo-deck-api/internal/domain"
)

// TODO: ユーザーIDを取得してカードを取得する
func (r *cardRepository) FindMyCard() (*domain.Card, error) {
	var card domain.Card
	card = domain.Card{
		ID:      domain.NewCardID(),
		OwnerID: domain.NewUserID(),
	}
	return &card, nil
}

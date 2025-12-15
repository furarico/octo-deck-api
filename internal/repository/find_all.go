package repository

import (
	"github.com/furarico/octo-deck-api/internal/domain"
)

func (r *cardRepository) FindAll() ([]domain.CardWithOwner, error) {

	return []domain.CardWithOwner{}, nil
}

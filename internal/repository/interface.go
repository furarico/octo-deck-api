package repository

import api "github.com/furarico/octo-deck-api/generated"

type CardRepository interface {
	FindAll() ([]api.Card, error)
	FindByID(id string) (*api.Card, error)
	FindMyCard() (*api.Card, error)
}

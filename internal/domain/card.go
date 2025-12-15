package domain

import (
	"github.com/google/uuid"
)

type CardID uuid.UUID

func NewCardID() CardID {
	return CardID(uuid.New())
}

type Card struct {
	ID      CardID
	OwnerID UserID
}

func NewCard(ownerID UserID) *Card {
	return &Card{
		ID:      NewCardID(),
		OwnerID: ownerID,
	}
}

package domain

import (
	"github.com/google/uuid"
)

type CardID uuid.UUID

func (c CardID) String() string {
	return uuid.UUID(c).String()
}

func NewCardID() CardID {
	return CardID(uuid.New())
}

// Card はカードのエンティティ（DBに保存される単位）
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

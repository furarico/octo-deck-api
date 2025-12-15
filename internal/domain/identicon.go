package domain

import (
	"github.com/google/uuid"
)

type IdenticonID uuid.UUID

type Color string

type Blocks [5][5]bool

type Identicon struct {
	ID     IdenticonID
	UserID UserID
	Color  Color
	Blocks Blocks
}

func NewIdenticon(userID UserID, color Color, blocks Blocks) *Identicon {
	return &Identicon{
		ID:     IdenticonID(uuid.New()),
		UserID: userID,
		Color:  color,
		Blocks: blocks,
	}
}

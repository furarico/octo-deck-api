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

type Color string

type Blocks [5][5]bool

type Card struct {
	ID       CardID
	GithubID string
	UserName string
	FullName string
	IconUrl  string
	Color    Color
	Blocks   Blocks
}

func NewCard(githubID string, color Color, blocks Blocks) *Card {
	return &Card{
		ID:       NewCardID(),
		GithubID: githubID,
		Color:    color,
		Blocks:   blocks,
	}
}

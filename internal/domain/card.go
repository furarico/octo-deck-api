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
	ID               CardID
	GithubID         string
	NodeID           string
	UserName         string
	FullName         string
	IconUrl          string
	Color            Color
	Blocks           Blocks
	MostUsedLanguage Language
}

func NewCard(githubID string, nodeID string, color Color, blocks Blocks, mostUsedLanguage Language, userName string, fullName string, iconUrl string) *Card {
	return &Card{
		ID:               NewCardID(),
		GithubID:         githubID,
		NodeID:           nodeID,
		Color:            color,
		Blocks:           blocks,
		MostUsedLanguage: mostUsedLanguage,
		UserName:         userName,
		FullName:         fullName,
		IconUrl:          iconUrl,
	}
}

package database

import (
	"encoding/json"
	"time"

	"github.com/furarico/octo-deck-api/internal/domain"
)

type Card struct {
	GithubID   string          `gorm:"not null;primaryKey"`
	CreatedAt  time.Time       `gorm:"autoCreateTime"`
	Color      string          `gorm:"not null"`
	BlocksData json.RawMessage `gorm:"type:jsonb;not null"`
}

func (c *Card) ToDomain() *domain.Card {
	var blocks domain.Blocks
	_ = json.Unmarshal(c.BlocksData, &blocks)

	return &domain.Card{
		GithubID: c.GithubID,
		Color:    domain.Color(c.Color),
		Blocks:   blocks,
	}
}

func CardFromDomain(card *domain.Card) *Card {
	blocksData, _ := json.Marshal(card.Blocks)

	return &Card{
		GithubID:   card.GithubID,
		Color:      string(card.Color),
		BlocksData: blocksData,
	}
}

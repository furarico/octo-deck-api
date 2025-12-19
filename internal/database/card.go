package database

import (
	"encoding/json"
	"time"

	"github.com/furarico/octo-deck-api/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Card struct {
	ID         uuid.UUID       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	GithubID   string          `gorm:"not null"`
	NodeID     string          `gorm:"not null"`
	CreatedAt  time.Time       `gorm:"autoCreateTime"`
	Color      string          `gorm:"not null"`
	BlocksData json.RawMessage `gorm:"type:jsonb;not null"`
}

func (c *Card) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

func (c *Card) ToDomain() *domain.Card {
	var blocks domain.Blocks
	_ = json.Unmarshal(c.BlocksData, &blocks)

	return &domain.Card{
		ID:       domain.CardID(c.ID),
		GithubID: c.GithubID,
		NodeID:   c.NodeID,
		Color:    domain.Color(c.Color),
		Blocks:   blocks,
	}
}

func CardFromDomain(card *domain.Card) *Card {
	blocksData, _ := json.Marshal(card.Blocks)

	return &Card{
		ID:         uuid.UUID(card.ID),
		GithubID:   card.GithubID,
		NodeID:     card.NodeID,
		Color:      string(card.Color),
		BlocksData: blocksData,
	}
}

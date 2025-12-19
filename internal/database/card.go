package database

import (
	"encoding/json"
	"time"

	"github.com/furarico/octo-deck-api/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Card struct {
	ID                    uuid.UUID       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	GithubID              string          `gorm:"not null"`
	NodeID                string          `gorm:"not null"`
	CreatedAt             time.Time       `gorm:"autoCreateTime"`
	Color                 string          `gorm:"not null"`
	BlocksData            json.RawMessage `gorm:"type:jsonb;not null"`
	UserName              string          `gorm:"default:''"`
	FullName              string          `gorm:"default:''"`
	IconUrl               string          `gorm:"default:''"`
	MostUsedLanguageName  string          `gorm:"default:''"`
	MostUsedLanguageColor string          `gorm:"default:''"`
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
		UserName: c.UserName,
		FullName: c.FullName,
		IconUrl:  c.IconUrl,
		MostUsedLanguage: domain.Language{
			LanguageName: c.MostUsedLanguageName,
			Color:        c.MostUsedLanguageColor,
		},
	}
}

func CardFromDomain(card *domain.Card) *Card {
	blocksData, _ := json.Marshal(card.Blocks)

	return &Card{
		ID:                    uuid.UUID(card.ID),
		GithubID:              card.GithubID,
		NodeID:                card.NodeID,
		Color:                 string(card.Color),
		BlocksData:            blocksData,
		UserName:              card.UserName,
		FullName:              card.FullName,
		IconUrl:               card.IconUrl,
		MostUsedLanguageName:  card.MostUsedLanguage.LanguageName,
		MostUsedLanguageColor: card.MostUsedLanguage.Color,
	}
}

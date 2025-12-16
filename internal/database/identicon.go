package database

import (
	"encoding/json"
	"time"

	"github.com/furarico/octo-deck-api/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Identicon struct {
	ID         uuid.UUID       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID     uuid.UUID       `gorm:"type:uuid;not null"`
	Color      string          `gorm:"not null"`
	BlocksData json.RawMessage `gorm:"type:jsonb;not null"`
	CreatedAt  time.Time       `gorm:"autoCreateTime"`
}

func (i *Identicon) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}

func (i *Identicon) ToDomain() *domain.Identicon {
	var blocks domain.Blocks
	_ = json.Unmarshal(i.BlocksData, &blocks)

	return &domain.Identicon{
		ID:     domain.IdenticonID(i.ID),
		UserID: domain.UserID(i.UserID),
		Color:  domain.Color(i.Color),
		Blocks: blocks,
	}
}

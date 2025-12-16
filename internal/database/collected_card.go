package database

import (
	"time"

	"github.com/furarico/octo-deck-api/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CollectedCard struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID      uuid.UUID `gorm:"type:uuid;not null"`
	CardID      uuid.UUID `gorm:"type:uuid;not null"`
	CollectedAt time.Time `gorm:"autoCreateTime"`

	Card Card `gorm:"foreignKey:CardID"`
}

func (cc *CollectedCard) BeforeCreate(tx *gorm.DB) error {
	if cc.ID == uuid.Nil {
		cc.ID = uuid.New()
	}
	return nil
}

func (cc *CollectedCard) ToDomain() *domain.CollectedCard {
	return &domain.CollectedCard{
		ID:          domain.CollectedCardID(cc.ID),
		CollectorID: domain.UserID(cc.UserID),
		CardID:      domain.CardID(cc.CardID),
	}
}

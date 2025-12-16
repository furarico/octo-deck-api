package database

import (
	"time"

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

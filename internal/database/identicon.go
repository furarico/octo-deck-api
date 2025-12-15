package database

import (
	"encoding/json"
	"time"

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

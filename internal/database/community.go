package database

import (
	"github.com/furarico/octo-deck-api/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Community struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name string    `gorm:"not null"`
}

func (c *Community) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

func (c *Community) ToDomain() *domain.Community {
	return &domain.Community{
		ID:   domain.CommunityID(c.ID),
		Name: c.Name,
	}
}

package database

import (
	"time"

	"github.com/furarico/octo-deck-api/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CollectedCard struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CollectorGithubID string    `gorm:"not null"`
	GithubID          string    `gorm:"not null"`
	CollectedAt       time.Time `gorm:"autoCreateTime"`

	Card Card `gorm:"foreignKey:GithubID"`
}

func (cc *CollectedCard) BeforeCreate(tx *gorm.DB) error {
	if cc.ID == uuid.Nil {
		cc.ID = uuid.New()
	}
	return nil
}

func (cc *CollectedCard) ToDomain() *domain.CollectedCard {
	return &domain.CollectedCard{
		ID:                domain.CollectedCardID(cc.ID),
		CollectorGithubID: cc.CollectorGithubID,
		GithubID:          cc.Card.GithubID,
	}
}

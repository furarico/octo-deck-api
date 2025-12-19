package database

import (
	"time"

	"github.com/furarico/octo-deck-api/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Community struct {
	ID                      uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name                    string     `gorm:"not null"`
	StartedAt               time.Time  `gorm:"not null"`
	EndedAt                 time.Time  `gorm:"not null"`
	CreatedAt               time.Time  `gorm:"autoCreateTime"`
	BestContributorCardID   *uuid.UUID `gorm:"type:uuid"`
	BestCommitterCardID     *uuid.UUID `gorm:"type:uuid"`
	BestIssuerCardID        *uuid.UUID `gorm:"type:uuid"`
	BestPullRequesterCardID *uuid.UUID `gorm:"type:uuid"`
	BestReviewerCardID      *uuid.UUID `gorm:"type:uuid"`
	// リレーション
	BestContributorCard   *Card `gorm:"foreignKey:BestContributorCardID"`
	BestCommitterCard     *Card `gorm:"foreignKey:BestCommitterCardID"`
	BestIssuerCard        *Card `gorm:"foreignKey:BestIssuerCardID"`
	BestPullRequesterCard *Card `gorm:"foreignKey:BestPullRequesterCardID"`
	BestReviewerCard      *Card `gorm:"foreignKey:BestReviewerCardID"`
}

func (c *Community) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

func (c *Community) ToDomain() *domain.Community {
	community := &domain.Community{
		ID:        domain.CommunityID(c.ID),
		Name:      c.Name,
		StartedAt: c.StartedAt,
		EndedAt:   c.EndedAt,
	}

	// HighlightedCardを構築
	if c.BestContributorCard != nil {
		community.HighlightedCard.BestContributor = *c.BestContributorCard.ToDomain()
	}
	if c.BestCommitterCard != nil {
		community.HighlightedCard.BestCommitter = *c.BestCommitterCard.ToDomain()
	}
	if c.BestIssuerCard != nil {
		community.HighlightedCard.BestIssuer = *c.BestIssuerCard.ToDomain()
	}
	if c.BestPullRequesterCard != nil {
		community.HighlightedCard.BestPullRequester = *c.BestPullRequesterCard.ToDomain()
	}
	if c.BestReviewerCard != nil {
		community.HighlightedCard.BestReviewer = *c.BestReviewerCard.ToDomain()
	}

	return community
}

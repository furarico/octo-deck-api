package database

import (
	"time"

	"github.com/furarico/octo-deck-api/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommunityUser struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CommunityID uuid.UUID `gorm:"type:uuid;not null"`
	UserID      uuid.UUID `gorm:"type:uuid;not null"`
	JoinedAt    time.Time `gorm:"autoCreateTime"`
}

func (cu *CommunityUser) BeforeCreate(tx *gorm.DB) error {
	if cu.ID == uuid.Nil {
		cu.ID = uuid.New()
	}
	return nil
}

func (cu *CommunityUser) ToDomain() *domain.CommunityUser {
	return &domain.CommunityUser{
		ID:          domain.CommunityUserID(cu.ID),
		CommunityID: domain.CommunityID(cu.CommunityID),
		UserID:      domain.UserID(cu.UserID),
		JoinedAt:    cu.JoinedAt,
	}
}

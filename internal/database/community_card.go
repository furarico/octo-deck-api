package database

import (
	"time"

	"github.com/furarico/octo-deck-api/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommunityCard struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CommunityID       uuid.UUID `gorm:"type:uuid;not null"`
	CardID            uuid.UUID `gorm:"type:uuid;not null"`
	JoinedAt          time.Time `gorm:"autoCreateTime"`
	TotalContribution int       `gorm:"default:0"`

	Card      Card      `gorm:"foreignKey:CardID"`
	Community Community `gorm:"foreignKey:CommunityID"`
}

func (cc *CommunityCard) BeforeCreate(tx *gorm.DB) error {
	if cc.ID == uuid.Nil {
		cc.ID = uuid.New()
	}
	return nil
}

func (cc *CommunityCard) ToDomain() *domain.CommunityCard {
	return &domain.CommunityCard{
		ID:                domain.CommunityCardID(cc.ID),
		CommunityID:       domain.CommunityID(cc.CommunityID),
		CardID:            domain.CardID(cc.CardID),
		JoinedAt:          cc.JoinedAt,
		TotalContribution: cc.TotalContribution,
	}
}

package domain

import (
	"time"

	"github.com/google/uuid"
)

type CommunityID uuid.UUID

func NewCommunityID() CommunityID {
	return CommunityID(uuid.New())
}

type Community struct {
	ID             CommunityID
	Name           string
	StartedAt      time.Time
	EndedAt        time.Time
	BestContribute BestContribute
}

func NewCommunity(name string, startedAt time.Time, endedAt time.Time, bestContribute BestContribute) *Community {
	return &Community{
		ID:             NewCommunityID(),
		Name:           name,
		StartedAt:      startedAt,
		EndedAt:        endedAt,
		BestContribute: bestContribute,
	}
}

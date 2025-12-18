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
	ID        CommunityID
	Name      string
	StartedAt time.Time
	EndedAt   time.Time
}

func NewCommunity(name string) *Community {
	return &Community{
		ID:   NewCommunityID(),
		Name: name,
	}
}

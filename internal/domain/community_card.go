package domain

import (
	"time"

	"github.com/google/uuid"
)

type CommunityCardID uuid.UUID

func NewCommunityCardID() CommunityCardID {
	return CommunityCardID(uuid.New())
}

type CommunityCard struct {
	ID          CommunityCardID
	CommunityID CommunityID
	CardID      CardID
	JoinedAt    time.Time
}

func NewCommunityCard(communityID CommunityID, cardID CardID) *CommunityCard {
	return &CommunityCard{
		ID:          NewCommunityCardID(),
		CommunityID: communityID,
		CardID:      cardID,
		JoinedAt:    time.Now(),
	}
}

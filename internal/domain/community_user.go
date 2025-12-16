package domain

import (
	"time"

	"github.com/google/uuid"
)

type CommunityUserID uuid.UUID

func NewCommunityUserID() CommunityUserID {
	return CommunityUserID(uuid.New())
}

type CommunityUser struct {
	ID          CommunityUserID
	CommunityID CommunityID
	UserID      UserID
	JoinedAt    time.Time
}

func NewCommunityUser(communityID CommunityID, userID UserID) *CommunityUser {
	return &CommunityUser{
		ID:          NewCommunityUserID(),
		CommunityID: communityID,
		UserID:      userID,
		JoinedAt:    time.Now(),
	}
}

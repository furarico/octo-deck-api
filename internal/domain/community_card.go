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
	GithubID    string
	JoinedAt    time.Time
}

func NewCommunityCard(communityID CommunityID, githubID string) *CommunityCard {
	return &CommunityCard{
		ID:          NewCommunityCardID(),
		CommunityID: communityID,
		GithubID:    githubID,
		JoinedAt:    time.Now(),
	}
}

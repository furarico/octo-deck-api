package domain

import (
	"github.com/google/uuid"
)

type CollectedCardID uuid.UUID

func NewCollectedCardID() CollectedCardID {
	return CollectedCardID(uuid.New())
}

type CollectedCard struct {
	ID                CollectedCardID
	CollectorGithubID string
	GithubID          string
}

func NewCollectedCard(collectorGithubID string, githubID string) *CollectedCard {
	return &CollectedCard{
		ID:                NewCollectedCardID(),
		CollectorGithubID: collectorGithubID,
		GithubID:          githubID,
	}
}

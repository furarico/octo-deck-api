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
	CardID            CardID
}

func NewCollectedCard(collectorGithubID string, cardID CardID) *CollectedCard {
	return &CollectedCard{
		ID:                NewCollectedCardID(),
		CollectorGithubID: collectorGithubID,
		CardID:            cardID,
	}
}

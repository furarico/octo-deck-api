package domain

import (
	"github.com/google/uuid"
)

type CollectedCardID uuid.UUID

func NewCollectedCardID() CollectedCardID {
	return CollectedCardID(uuid.New())
}

type CollectedCard struct {
	ID          CollectedCardID
	CollectorID UserID
	CardID      CardID
}

func NewCollectedCard(collectorID UserID, cardID CardID) *CollectedCard {
	return &CollectedCard{
		ID:          NewCollectedCardID(),
		CollectorID: collectorID,
		CardID:      cardID,
	}
}

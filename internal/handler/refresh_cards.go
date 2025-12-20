package handler

import (
	"context"
	"fmt"

	api "github.com/furarico/octo-deck-api/generated"
)

// 全カードを更新
// (PUT /cards/refresh)
func (h *Handler) RefreshAllCards(ctx context.Context, request api.RefreshAllCardsRequestObject) (api.RefreshAllCardsResponseObject, error) {
	githubClient, err := getGitHubClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get github client: %w", err)
	}

	cards, err := h.cardService.RefreshAllCards(ctx, githubClient)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh all cards: %w", err)
	}

	apiCards := make([]api.Card, 0, len(cards))
	for _, card := range cards {
		apiCards = append(apiCards, convertCardToAPI(card))
	}

	return api.RefreshAllCards200JSONResponse{
		Card: apiCards,
	}, nil
}

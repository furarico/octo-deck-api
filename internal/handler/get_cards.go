package handler

import (
	"context"
	"fmt"

	api "github.com/furarico/octo-deck-api/generated"
)

// カード一覧取得
// (GET /cards)
func (h *Handler) GetCards(ctx context.Context, request api.GetCardsRequestObject) (api.GetCardsResponseObject, error) {
	githubClient, err := getGitHubClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: %w", err)
	}

	githubID, err := getGitHubID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: %w", err)
	}

	cards, err := h.cardService.GetAllCards(ctx, githubID, githubClient)
	if err != nil {
		return nil, fmt.Errorf("failed to get cards: %w", err)
	}

	cardsAPI := make([]api.Card, len(cards))
	for i, card := range cards {
		cardsAPI[i] = convertCardToAPI(card)
	}

	return api.GetCards200JSONResponse{Cards: cardsAPI}, nil
}

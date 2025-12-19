package handler

import (
	"context"
	"fmt"

	api "github.com/furarico/octo-deck-api/generated"
)

// 指定したコミュニティのカード一覧取得
// (GET /communities/{id}/cards)
func (h *Handler) GetCommunityCards(ctx context.Context, request api.GetCommunityCardsRequestObject) (api.GetCommunityCardsResponseObject, error) {
	githubClient, err := getGitHubClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get github client: %w", err)
	}

	// リクエストコンテキストを取得
	reqCtx := getRequestContext(ctx)

	cards, err := h.communityService.GetCommunityCards(reqCtx, request.Id, githubClient)
	if err != nil {
		return nil, fmt.Errorf("failed to get community cards: %w", err)
	}

	cardsAPI := make([]api.Card, len(cards))
	for i, card := range cards {
		cardsAPI[i] = convertCardToAPI(card)
	}

	return api.GetCommunityCards200JSONResponse{Cards: cardsAPI}, nil
}

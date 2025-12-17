package handler

import (
	"context"
	"fmt"

	api "github.com/furarico/octo-deck-api/generated"
)

// 自分のカード取得
// (GET /cards/me)
func (h *Handler) GetMyCard(ctx context.Context, request api.GetMyCardRequestObject) (api.GetMyCardResponseObject, error) {
	githubClient, err := getGitHubClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: %w", err)
	}

	githubID, err := getGitHubID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: %w", err)
	}

	card, err := h.cardService.GetMyCard(ctx, githubID, githubClient)
	if err != nil {
		return nil, fmt.Errorf("card not found: %w", err)
	}

	return api.GetMyCard200JSONResponse{Card: convertCardToAPI(*card)}, nil
}

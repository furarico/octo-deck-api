package handler

import (
	"context"
	"fmt"

	api "github.com/furarico/octo-deck-api/generated"
)

// 指定したカード取得
// (GET /cards/{githubId})
func (h *Handler) GetCard(ctx context.Context, request api.GetCardRequestObject) (api.GetCardResponseObject, error) {
	githubClient, err := getGitHubClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: %w", err)
	}

	card, err := h.cardService.GetCardByGitHubID(ctx, request.GithubId, githubClient)
	if err != nil {
		return nil, fmt.Errorf("card not found: %w", err)
	}

	return api.GetCard200JSONResponse{Card: convertCardToAPI(*card)}, nil
}

package handler

import (
	"context"
	"fmt"

	api "github.com/furarico/octo-deck-api/generated"
)

// 自分のカード取得（存在しない場合は作成）
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

	card, err := h.cardService.GetOrCreateMyCard(ctx, githubID, githubClient)
	if err != nil {
		return nil, fmt.Errorf("failed to get or create card: %w", err)
	}

	return api.GetMyCard200JSONResponse{Card: convertCardToAPI(*card)}, nil
}

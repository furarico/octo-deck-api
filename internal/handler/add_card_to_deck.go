package handler

import (
	"context"
	"fmt"

	api "github.com/furarico/octo-deck-api/generated"
)

// カードをデッキに追加
// (POST /cards)
func (h *Handler) AddCardToDeck(ctx context.Context, request api.AddCardToDeckRequestObject) (api.AddCardToDeckResponseObject, error) {
	githubClient, err := getGitHubClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: %w", err)
	}

	collectorGithubID, err := getGitHubID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: %w", err)
	}

	// リクエストボディから追加対象のGitHub IDを取得
	if request.Body == nil {
		return nil, fmt.Errorf("request body is required")
	}
	targetGithubID := string(*request.Body)

	card, err := h.cardService.AddCardToDeck(ctx, collectorGithubID, targetGithubID, githubClient)
	if err != nil {
		return nil, fmt.Errorf("failed to add card to deck: %w", err)
	}

	return api.AddCardToDeck200JSONResponse{Card: convertCardToAPI(*card)}, nil
}

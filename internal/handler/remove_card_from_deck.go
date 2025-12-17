package handler

import (
	"context"
	"fmt"

	api "github.com/furarico/octo-deck-api/generated"
)

// カードをデッキから削除
// (DELETE /cards/{githubId})
func (h *Handler) RemoveCardFromDeck(ctx context.Context, request api.RemoveCardFromDeckRequestObject) (api.RemoveCardFromDeckResponseObject, error) {
	githubClient, err := getGitHubClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: %w", err)
	}

	collectorGithubID, err := getGitHubID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: %w", err)
	}

	// パスパラメータから削除対象のGitHub IDを取得
	targetGithubID := request.GithubId

	card, err := h.cardService.RemoveCardFromDeck(ctx, collectorGithubID, targetGithubID, githubClient)
	if err != nil {
		return nil, fmt.Errorf("failed to remove card from deck: %w", err)
	}

	return api.RemoveCardFromDeck200JSONResponse{Card: convertCardToAPI(*card)}, nil
}

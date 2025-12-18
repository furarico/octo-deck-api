package handler

import (
	"context"
	"fmt"

	api "github.com/furarico/octo-deck-api/generated"
)

// 指定したコミュニティに自分のカードを追加
// (POST /communities/{id}/cards)
func (h *Handler) AddCardToCommunity(ctx context.Context, request api.AddCardToCommunityRequestObject) (api.AddCardToCommunityResponseObject, error) {
	githubClient, err := getGitHubClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: %w", err)
	}

	githubID, err := getGitHubID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: %w", err)
	}

	// github_idから自分のカードを取得
	card, err := h.cardService.GetMyCard(ctx, githubID, githubClient)
	if err != nil {
		return nil, fmt.Errorf("card not found: %w", err)
	}

	// コミュニティにカードを追加
	if err := h.communityService.AddCardToCommunity(request.Id, card.GithubID); err != nil {
		return nil, fmt.Errorf("failed to add card to community: %w", err)
	}

	return api.AddCardToCommunity200JSONResponse{Card: convertCardToAPI(*card)}, nil
}

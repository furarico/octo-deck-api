package handler

import (
	"context"
	"fmt"

	api "github.com/furarico/octo-deck-api/generated"
)

// 指定したコミュニティ取得
// (GET /communities/{id})
func (h *Handler) GetCommunity(ctx context.Context, request api.GetCommunityRequestObject) (api.GetCommunityResponseObject, error) {
	githubClient, err := getGitHubClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get github client: %w", err)
	}

	community, highlightedCard, err := h.communityService.GetCommunityWithHighlightedCard(ctx, request.Id, githubClient)
	if err != nil {
		return nil, fmt.Errorf("community not found: %w", err)
	}

	return api.GetCommunity200JSONResponse{
		Community:       convertCommunityToAPI(*community),
		HighlightedCard: convertHighlightedCardToAPI(*highlightedCard),
	}, nil
}

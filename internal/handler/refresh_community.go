package handler

import (
	"context"
	"fmt"

	api "github.com/furarico/octo-deck-api/generated"
)

// コミュニティのHighlightedCardを更新
// (PUT /communities/{id}/refresh)
func (h *Handler) RefreshCommunity(ctx context.Context, request api.RefreshCommunityRequestObject) (api.RefreshCommunityResponseObject, error) {
	githubClient, err := getGitHubClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get github client: %w", err)
	}

	community, highlightedCard, err := h.communityService.RefreshHighlightedCard(ctx, request.Id, githubClient)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh community: %w", err)
	}

	return api.RefreshCommunity200JSONResponse{
		Community:       convertCommunityToAPI(*community),
		HighlightedCard: convertHighlightedCardToAPI(*highlightedCard),
	}, nil
}

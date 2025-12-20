package handler

import (
	"context"
	"fmt"

	api "github.com/furarico/octo-deck-api/generated"
)

// コミュニティ一覧取得
// (GET /communities)
func (h *Handler) GetCommunities(ctx context.Context, request api.GetCommunitiesRequestObject) (api.GetCommunitiesResponseObject, error) {
	githubID, err := getGitHubID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: %w", err)
	}

	communities, err := h.communityService.GetAllCommunities(ctx, githubID)
	if err != nil {
		return nil, fmt.Errorf("failed to get communities: %w", err)
	}

	communitiesAPI := make([]api.Community, len(communities))
	for i, community := range communities {
		communitiesAPI[i] = convertCommunityToAPI(community)
	}

	return api.GetCommunities200JSONResponse{Communities: communitiesAPI}, nil
}

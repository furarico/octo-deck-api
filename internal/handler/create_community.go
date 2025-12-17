package handler

import (
	"context"
	"fmt"

	api "github.com/furarico/octo-deck-api/generated"
)

// コミュニティを作成
// (POST /communities)
func (h *Handler) CreateCommunity(ctx context.Context, request api.CreateCommunityRequestObject) (api.CreateCommunityResponseObject, error) {
	if request.Body == nil || *request.Body == "" {
		return nil, fmt.Errorf("community name is required")
	}

	name := string(*request.Body)

	community, err := h.communityService.CreateCommunity(name)
	if err != nil {
		return nil, fmt.Errorf("failed to create community: %w", err)
	}

	return api.CreateCommunity200JSONResponse{Community: convertCommunityToAPI(*community)}, nil
}

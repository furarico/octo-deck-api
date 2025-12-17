package handler

import (
	"context"
	"fmt"

	api "github.com/furarico/octo-deck-api/generated"
)

// 指定したコミュニティ取得
// (GET /communities/{id})
func (h *Handler) GetCommunity(ctx context.Context, request api.GetCommunityRequestObject) (api.GetCommunityResponseObject, error) {
	community, err := h.communityService.GetCommunityByID(request.Id)
	if err != nil {
		return nil, fmt.Errorf("community not found: %w", err)
	}

	return api.GetCommunity200JSONResponse{Community: convertCommunityToAPI(*community)}, nil
}

package handler

import (
	"context"
	"fmt"

	api "github.com/furarico/octo-deck-api/generated"
)

// コミュニティを削除
// (DELETE /communities/{id})
func (h *Handler) DeleteCommunity(ctx context.Context, request api.DeleteCommunityRequestObject) (api.DeleteCommunityResponseObject, error) {
	// 削除前にコミュニティ情報を取得
	community, err := h.communityService.GetCommunityByID(request.Id)
	if err != nil {
		return nil, fmt.Errorf("community not found: %w", err)
	}

	if err := h.communityService.DeleteCommunity(request.Id); err != nil {
		return nil, fmt.Errorf("failed to delete community: %w", err)
	}

	return api.DeleteCommunity200JSONResponse{Community: convertCommunityToAPI(*community)}, nil
}

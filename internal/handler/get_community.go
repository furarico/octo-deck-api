package handler

import (
	"context"
	"fmt"

	api "github.com/furarico/octo-deck-api/generated"
)

// 指定したコミュニティ取得
// (GET /communities/{id})
func (h *Handler) GetCommunity(_ context.Context, request api.GetCommunityRequestObject) (api.GetCommunityResponseObject, error) {
	community, highlightedCard, err := h.communityService.GetCommunityWithHighlightedCard(request.Id)
	if err != nil {
		return nil, fmt.Errorf("community not found: %w", err)
	}

	return api.GetCommunity200JSONResponse{
		Community:       convertCommunityToAPI(*community),
		HighlightedCard: convertHighlightedCardToAPI(*highlightedCard),
	}, nil
}

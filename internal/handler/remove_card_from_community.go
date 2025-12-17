package handler

import (
	"context"
	"fmt"

	api "github.com/furarico/octo-deck-api/generated"
	"github.com/google/uuid"
)

// 指定したコミュニティの自分のカードを削除
// (DELETE /communities/{id}/cards)
func (h *Handler) RemoveCardFromCommunity(ctx context.Context, request api.RemoveCardFromCommunityRequestObject) (api.RemoveCardFromCommunityResponseObject, error) {
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

	// コミュニティからカードを削除
	cardID := uuid.UUID(card.ID).String()
	if err := h.communityService.RemoveCardFromCommunity(request.Id, cardID); err != nil {
		return nil, fmt.Errorf("failed to remove card from community: %w", err)
	}

	return api.RemoveCardFromCommunity200JSONResponse{Card: convertCardToAPI(*card)}, nil
}

package handler

import (
	"context"
	"fmt"

	api "github.com/furarico/octo-deck-api/generated"
	"github.com/google/uuid"
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
	cardID := uuid.UUID(card.ID).String()
	if err := h.communityService.AddCardToCommunity(ctx, request.Id, cardID); err != nil {
		return nil, fmt.Errorf("failed to add card to community: %w", err)
	}

	return api.AddCardToCommunity200JSONResponse{Card: convertCardToAPI(*card)}, nil
}

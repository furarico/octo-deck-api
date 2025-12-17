package handler

import (
	"context"
	"fmt"

	api "github.com/furarico/octo-deck-api/generated"
)

// カードをデッキから削除
// (DELETE /cards/{githubId})
func (h *Handler) RemoveCardFromDeck(ctx context.Context, request api.RemoveCardFromDeckRequestObject) (api.RemoveCardFromDeckResponseObject, error) {
	// TODO: 実装
	return nil, fmt.Errorf("not implemented")
}

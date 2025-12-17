package handler

import (
	"context"
	"fmt"

	api "github.com/furarico/octo-deck-api/generated"
)

// カードをデッキに追加
// (POST /cards)
func (h *Handler) AddCardToDeck(ctx context.Context, request api.AddCardToDeckRequestObject) (api.AddCardToDeckResponseObject, error) {
	// TODO: 実装
	return nil, fmt.Errorf("not implemented")
}

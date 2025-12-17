package handler

import (
	"context"
	"fmt"

	api "github.com/furarico/octo-deck-api/generated"
)

func (h *Handler) GetUserStats(ctx context.Context, request api.GetUserStatsRequestObject) (api.GetUserStatsResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}

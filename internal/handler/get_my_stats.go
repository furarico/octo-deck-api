package handler

import (
	"context"
	"fmt"

	api "github.com/furarico/octo-deck-api/generated"
)

func (h *Handler) GetMyStats(ctx context.Context, request api.GetMyStatsRequestObject) (api.GetMyStatsResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}

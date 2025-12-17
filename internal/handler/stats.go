package handler

import (
	"context"
	"fmt"

	api "github.com/furarico/octo-deck-api/generated"
)

// 自分の統計情報取得
// (GET /stats/me)
func (h *Handler) GetMyStats(ctx context.Context, request api.GetMyStatsRequestObject) (api.GetMyStatsResponseObject, error) {
	// TODO: 実装
	return nil, fmt.Errorf("not implemented")
}

// ユーザーの統計情報取得
// (GET /stats/{githubId})
func (h *Handler) GetUserStats(ctx context.Context, request api.GetUserStatsRequestObject) (api.GetUserStatsResponseObject, error) {
	// TODO: 実装
	return nil, fmt.Errorf("not implemented")
}

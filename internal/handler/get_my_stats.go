package handler

import (
	"context"

	api "github.com/furarico/octo-deck-api/generated"
)

// 自分の統計情報取得
func (h *Handler) GetMyStats(ctx context.Context, request api.GetMyStatsRequestObject) (api.GetMyStatsResponseObject, error) {
	// GitHub IDを取得
	githubID, err := getGitHubID(ctx)
	if err != nil {
		return nil, err
	}

	// GitHub Clientを取得
	githubClient, err := getGitHubClient(ctx)
	if err != nil {
		return nil, err
	}

	// 統計情報を取得
	stats, err := h.statsService.GetUserStats(ctx, githubID, githubClient)
	if err != nil {
		return nil, err
	}

	// レスポンスに変換
	userStats, err := convertUserStatsToAPI(stats)
	if err != nil {
		return nil, err
	}

	return api.GetMyStats200JSONResponse{
		Stats: userStats,
	}, nil
}

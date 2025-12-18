package handler

import (
	"context"

	api "github.com/furarico/octo-deck-api/generated"
)

// ユーザーの統計情報取得
func (h *Handler) GetUserStats(ctx context.Context, request api.GetUserStatsRequestObject) (api.GetUserStatsResponseObject, error) {
	// パスパラメータからGitHub IDを取得
	githubID := request.GithubId

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

	return api.GetUserStats200JSONResponse{
		Stats: userStats,
	}, nil
}

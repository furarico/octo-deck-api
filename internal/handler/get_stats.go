package handler

import (
	"context"
	"time"

	api "github.com/furarico/octo-deck-api/generated"
	"github.com/oapi-codegen/runtime/types"
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
	contributions := make([]api.Contribution, len(stats.Contributions))
	for i, c := range stats.Contributions {
		// DateをYYYY-MM-DD形式でパース
		date, err := time.Parse("2006-01-02", c.Date)
		if err != nil {
			return nil, err
		}

		contributions[i] = api.Contribution{
			Date:  types.Date{Time: date},
			Count: int32(c.Count),
		}
	}

	return api.GetUserStats200JSONResponse{
		Stats: api.UserStats{
			Contributions: contributions,
		},
	}, nil
}

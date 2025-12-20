package service

import (
	"context"
	"time"

	"github.com/furarico/octo-deck-api/internal/github"
)

// GitHubClient はServiceが必要とするGitHub APIクライアントのインターフェース
type GitHubClient interface {
	GetAuthenticatedUser(ctx context.Context) (*github.UserInfo, error)
	GetUserByID(ctx context.Context, id int64) (*github.UserInfo, error)
	GetUsersByIDs(ctx context.Context, ids []int64) (map[int64]*github.UserInfo, error)
	GetUserStats(ctx context.Context, githubID int64) (*github.UserStats, error)
	GetMostUsedLanguage(ctx context.Context, login string) (string, string, error)
	GetMostUsedLanguages(ctx context.Context, logins []string) (map[string]github.LanguageInfo, error)
	// GetUsersFullInfoByNodeIDs はNodeIDを使ってユーザーの全情報（基本情報、貢献データ、言語情報）を一括取得する
	GetUsersFullInfoByNodeIDs(ctx context.Context, nodeIDs []string, from, to time.Time) ([]github.UserFullInfo, error)
}

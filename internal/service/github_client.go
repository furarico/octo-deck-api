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
	GetUserStats(ctx context.Context, githubID int64) (*github.UserStats, error)
	GetMostUsedLanguage(ctx context.Context, login string) (string, string, error)
	GetContributionStats(ctx context.Context, githubID int64) (*github.ContributionStats, error)
	GetUsersContributions(ctx context.Context, usernames []string, from, to time.Time) ([]github.UserContributionStats, error)
	// バッチ取得メソッド（N+1問題解消用）
	GetUsersByIDs(ctx context.Context, ids []int64) (map[int64]*github.UserInfo, error)
	GetUsersByLogins(ctx context.Context, logins []string) (map[string]*github.UserInfo, error)
	GetUsersLanguages(ctx context.Context, logins []string) (map[string]*github.LanguageInfo, error)
}

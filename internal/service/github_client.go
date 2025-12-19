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
}

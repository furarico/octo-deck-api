package service

import (
	"context"
	"time"

	"github.com/furarico/octo-deck-api/internal/github"
)

// MockGitHubClient はテスト用のモッククライアント
type MockGitHubClient struct {
	GetAuthenticatedUserFunc  func(ctx context.Context) (*github.UserInfo, error)
	GetUserByIDFunc           func(ctx context.Context, id int64) (*github.UserInfo, error)
	GetUserStatsFunc          func(ctx context.Context, githubID int64) (*github.UserStats, error)
	GetMostUsedLanguageFunc   func(ctx context.Context, login string) (string, string, error)
	GetContributionStatsFunc  func(ctx context.Context, githubID int64) (*github.ContributionStats, error)
	GetUsersContributionsFunc func(ctx context.Context, usernames []string, from, to time.Time) ([]github.UserContributionStats, error)
}

// MockGitHubClientがGitHubClientインターフェースを実装していることを確認
var _ GitHubClient = (*MockGitHubClient)(nil)

func NewMockGitHubClient() *MockGitHubClient {
	return &MockGitHubClient{}
}

func (m *MockGitHubClient) GetAuthenticatedUser(ctx context.Context) (*github.UserInfo, error) {
	if m.GetAuthenticatedUserFunc != nil {
		return m.GetAuthenticatedUserFunc(ctx)
	}
	return &github.UserInfo{}, nil
}

func (m *MockGitHubClient) GetUserByID(ctx context.Context, id int64) (*github.UserInfo, error) {
	if m.GetUserByIDFunc != nil {
		return m.GetUserByIDFunc(ctx, id)
	}
	return &github.UserInfo{}, nil
}

func (m *MockGitHubClient) GetUserStats(ctx context.Context, githubID int64) (*github.UserStats, error) {
	if m.GetUserStatsFunc != nil {
		return m.GetUserStatsFunc(ctx, githubID)
	}
	return &github.UserStats{}, nil
}

func (m *MockGitHubClient) GetMostUsedLanguage(ctx context.Context, login string) (string, string, error) {
	if m.GetMostUsedLanguageFunc != nil {
		return m.GetMostUsedLanguageFunc(ctx, login)
	}
	return "Go", "#00ADD8", nil
}

func (m *MockGitHubClient) GetContributionStats(ctx context.Context, githubID int64) (*github.ContributionStats, error) {
	if m.GetContributionStatsFunc != nil {
		return m.GetContributionStatsFunc(ctx, githubID)
	}
	return &github.ContributionStats{}, nil
}

func (m *MockGitHubClient) GetUsersContributions(ctx context.Context, usernames []string, from, to time.Time) ([]github.UserContributionStats, error) {
	if m.GetUsersContributionsFunc != nil {
		return m.GetUsersContributionsFunc(ctx, usernames, from, to)
	}
	return []github.UserContributionStats{}, nil
}

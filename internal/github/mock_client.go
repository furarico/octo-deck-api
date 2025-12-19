package github

import (
	"context"
	"time"
)

// MockClient はテスト用のモッククライアント
// service.GitHubClientインターフェースを実装します
type MockClient struct {
	GetAuthenticatedUserFunc  func(ctx context.Context) (*UserInfo, error)
	GetUserByIDFunc           func(ctx context.Context, id int64) (*UserInfo, error)
	GetUserStatsFunc          func(ctx context.Context, githubID int64) (*UserStats, error)
	GetMostUsedLanguageFunc   func(ctx context.Context, login string) (string, string, error)
	GetContributionStatsFunc  func(ctx context.Context, githubID int64) (*ContributionStats, error)
	GetUsersContributionsFunc func(ctx context.Context, usernames []string, from, to time.Time) ([]UserContributionStats, error)
	// バッチ取得メソッド（N+1問題解消用）
	GetUsersByIDsFunc     func(ctx context.Context, ids []int64) (map[int64]*UserInfo, error)
	GetUsersByLoginsFunc  func(ctx context.Context, logins []string) (map[string]*UserInfo, error)
	GetUsersLanguagesFunc func(ctx context.Context, logins []string) (map[string]*LanguageInfo, error)
}

func NewMockClient() *MockClient {
	return &MockClient{}
}

func (m *MockClient) GetAuthenticatedUser(ctx context.Context) (*UserInfo, error) {
	if m.GetAuthenticatedUserFunc != nil {
		return m.GetAuthenticatedUserFunc(ctx)
	}
	return &UserInfo{}, nil
}

func (m *MockClient) GetUserByID(ctx context.Context, id int64) (*UserInfo, error) {
	if m.GetUserByIDFunc != nil {
		return m.GetUserByIDFunc(ctx, id)
	}
	return &UserInfo{}, nil
}

func (m *MockClient) GetUserStats(ctx context.Context, githubID int64) (*UserStats, error) {
	if m.GetUserStatsFunc != nil {
		return m.GetUserStatsFunc(ctx, githubID)
	}
	return &UserStats{}, nil
}

func (m *MockClient) GetMostUsedLanguage(ctx context.Context, login string) (string, string, error) {
	if m.GetMostUsedLanguageFunc != nil {
		return m.GetMostUsedLanguageFunc(ctx, login)
	}
	return "Go", "#00ADD8", nil
}

func (m *MockClient) GetContributionStats(ctx context.Context, githubID int64) (*ContributionStats, error) {
	if m.GetContributionStatsFunc != nil {
		return m.GetContributionStatsFunc(ctx, githubID)
	}
	return &ContributionStats{}, nil
}

func (m *MockClient) GetUsersContributions(ctx context.Context, usernames []string, from, to time.Time) ([]UserContributionStats, error) {
	if m.GetUsersContributionsFunc != nil {
		return m.GetUsersContributionsFunc(ctx, usernames, from, to)
	}
	return []UserContributionStats{}, nil
}

func (m *MockClient) GetUsersByIDs(ctx context.Context, ids []int64) (map[int64]*UserInfo, error) {
	if m.GetUsersByIDsFunc != nil {
		return m.GetUsersByIDsFunc(ctx, ids)
	}
	result := make(map[int64]*UserInfo)
	for _, id := range ids {
		result[id] = &UserInfo{
			ID:        id,
			Login:     "testuser",
			Name:      "Test User",
			AvatarURL: "https://example.com/avatar.png",
		}
	}
	return result, nil
}

func (m *MockClient) GetUsersByLogins(ctx context.Context, logins []string) (map[string]*UserInfo, error) {
	if m.GetUsersByLoginsFunc != nil {
		return m.GetUsersByLoginsFunc(ctx, logins)
	}
	result := make(map[string]*UserInfo)
	for _, login := range logins {
		result[login] = &UserInfo{
			ID:        12345,
			Login:     login,
			Name:      "Test User",
			AvatarURL: "https://example.com/avatar.png",
		}
	}
	return result, nil
}

func (m *MockClient) GetUsersLanguages(ctx context.Context, logins []string) (map[string]*LanguageInfo, error) {
	if m.GetUsersLanguagesFunc != nil {
		return m.GetUsersLanguagesFunc(ctx, logins)
	}
	result := make(map[string]*LanguageInfo)
	for _, login := range logins {
		result[login] = &LanguageInfo{
			Name:  "Go",
			Color: "#00ADD8",
		}
	}
	return result, nil
}

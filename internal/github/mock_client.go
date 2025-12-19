package github

import (
	"context"
	"time"
)

// MockClient はテスト用のモッククライアント
// service.GitHubClientインターフェースを実装します
type MockClient struct {
	GetAuthenticatedUserFunc   func(ctx context.Context) (*UserInfo, error)
	GetUserByIDFunc            func(ctx context.Context, id int64) (*UserInfo, error)
	GetUsersByIDsFunc          func(ctx context.Context, ids []int64) (map[int64]*UserInfo, error)
	GetUserStatsFunc           func(ctx context.Context, githubID int64) (*UserStats, error)
	GetMostUsedLanguageFunc    func(ctx context.Context, login string) (string, string, error)
	GetMostUsedLanguagesFunc   func(ctx context.Context, logins []string) (map[string]LanguageInfo, error)
	GetContributionStatsFunc   func(ctx context.Context, githubID int64) (*ContributionStats, error)
	GetUsersContributionsFunc  func(ctx context.Context, usernames []string, from, to time.Time) ([]UserContributionStats, error)
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

func (m *MockClient) GetUsersByIDs(ctx context.Context, ids []int64) (map[int64]*UserInfo, error) {
	if m.GetUsersByIDsFunc != nil {
		return m.GetUsersByIDsFunc(ctx, ids)
	}
	// デフォルトでは各IDに対してGetUserByIDを呼び出す
	result := make(map[int64]*UserInfo)
	for _, id := range ids {
		info, err := m.GetUserByID(ctx, id)
		if err != nil {
			continue
		}
		result[id] = info
	}
	return result, nil
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

func (m *MockClient) GetMostUsedLanguages(ctx context.Context, logins []string) (map[string]LanguageInfo, error) {
	if m.GetMostUsedLanguagesFunc != nil {
		return m.GetMostUsedLanguagesFunc(ctx, logins)
	}
	// デフォルトでは各ログイン名に対してGetMostUsedLanguageを呼び出す
	result := make(map[string]LanguageInfo)
	for _, login := range logins {
		name, color, err := m.GetMostUsedLanguage(ctx, login)
		if err != nil {
			result[login] = LanguageInfo{Name: "Unknown", Color: "#586069"}
			continue
		}
		result[login] = LanguageInfo{Name: name, Color: color}
	}
	return result, nil
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

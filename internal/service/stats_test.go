package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/furarico/octo-deck-api/internal/github"
)

// テスト用のヘルパー関数: 正常なUserStatsを返す
func createTestUserStats() *github.UserStats {
	return &github.UserStats{
		Contributions: []github.Contribution{
			{Date: "2024-01-01", Count: 5},
			{Date: "2024-01-02", Count: 3},
		},
		TotalContribution: 100,
		MostUsedLanguage:  "Go",
		MostUsedLanguageColor: "#00ADD8",
		ContributionDetail: github.ContributionDetail{
			CommitCount:      50,
			IssueCount:       20,
			PullRequestCount: 20,
			ReviewCount:      10,
		},
	}
}

// GetUserStats は指定されたGitHub IDのユーザーの統計情報を取得する
func TestGetUserStats(t *testing.T) {
	tests := []struct {
		name        string
		githubID    string
		setupGitHub func() *MockGitHubClient
		wantErr     bool
		wantErrMsg  string
	}{
		{
			name:     "正常に統計情報を取得できる",
			githubID: "12345",
			setupGitHub: func() *MockGitHubClient {
				return &MockGitHubClient{
					GetUserStatsFunc: func(ctx context.Context, githubID int64) (*github.UserStats, error) {
						return createTestUserStats(), nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:     "無効なGitHubIDの場合",
			githubID: "invalid_id",
			setupGitHub: func() *MockGitHubClient {
				return &MockGitHubClient{}
			},
			wantErr:    true,
			wantErrMsg: "invalid github id",
		},
		{
			name:     "GitHubClientエラーが発生した場合",
			githubID: "12345",
			setupGitHub: func() *MockGitHubClient {
				return &MockGitHubClient{
					GetUserStatsFunc: func(ctx context.Context, githubID int64) (*github.UserStats, error) {
						return nil, fmt.Errorf("github api error")
					},
				}
			},
			wantErr:    true,
			wantErrMsg: "failed to get user stats",
		},
		{
			name:     "ToDomainStatsで日付パースエラーが発生した場合",
			githubID: "12345",
			setupGitHub: func() *MockGitHubClient {
				return &MockGitHubClient{
					GetUserStatsFunc: func(ctx context.Context, githubID int64) (*github.UserStats, error) {
						// 無効な日付形式を含むUserStatsを返す
						return &github.UserStats{
							Contributions: []github.Contribution{
								{Date: "invalid-date", Count: 5},
							},
							TotalContribution: 100,
							MostUsedLanguage:  "Go",
							MostUsedLanguageColor: "#00ADD8",
							ContributionDetail: github.ContributionDetail{
								CommitCount:      50,
								IssueCount:       20,
								PullRequestCount: 20,
								ReviewCount:      10,
							},
						}, nil
					},
				}
			},
			wantErr:    true,
			wantErrMsg: "failed to convert stats to domain",
		},
		{
			name:     "空のContributionsの場合でも正常に処理できる",
			githubID: "12345",
			setupGitHub: func() *MockGitHubClient {
				return &MockGitHubClient{
					GetUserStatsFunc: func(ctx context.Context, githubID int64) (*github.UserStats, error) {
						return &github.UserStats{
							Contributions:     []github.Contribution{},
							TotalContribution: 0,
							MostUsedLanguage:  "",
							MostUsedLanguageColor: "",
							ContributionDetail: github.ContributionDetail{},
						}, nil
					},
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			githubClient := tt.setupGitHub()
			service := NewStatsService()
			stats, err := service.GetUserStats(ctx, tt.githubID, githubClient)

			if tt.wantErr {
				if err == nil {
					t.Errorf("エラーが期待されましたが、エラーが発生しませんでした")
					return
				}
				if tt.wantErrMsg != "" && !contains(err.Error(), tt.wantErrMsg) {
					t.Errorf("エラーメッセージが期待と異なります: 期待=%s, 実際=%s", tt.wantErrMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラーが発生しました: %v", err)
					return
				}
				if stats == nil {
					t.Errorf("統計情報がnilです")
					return
				}
			}
		})
	}
}

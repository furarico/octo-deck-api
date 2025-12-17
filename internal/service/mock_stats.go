package service

import (
	"context"

	"github.com/furarico/octo-deck-api/internal/github"
)

// MockStatsService はテスト用のモック統計サービス
type MockStatsService struct {
	GetUserStatsFunc func(ctx context.Context, githubID string, githubClient *github.Client) (*github.ContributionStats, error)
}

func NewMockStatsService() *MockStatsService {
	return &MockStatsService{}
}

func (m *MockStatsService) GetUserStats(ctx context.Context, githubID string, githubClient *github.Client) (*github.ContributionStats, error) {
	if m.GetUserStatsFunc != nil {
		return m.GetUserStatsFunc(ctx, githubID, githubClient)
	}
	return &github.ContributionStats{}, nil
}

package service

import (
	"context"

	"github.com/furarico/octo-deck-api/internal/domain"
)

// MockStatsService はテスト用のモック統計サービス
type MockStatsService struct {
	GetUserStatsFunc func(ctx context.Context, githubID string, githubClient GitHubClient) (*domain.Stats, error)
}

func NewMockStatsService() *MockStatsService {
	return &MockStatsService{}
}

func (m *MockStatsService) GetUserStats(ctx context.Context, githubID string, githubClient GitHubClient) (*domain.Stats, error) {
	if m.GetUserStatsFunc != nil {
		return m.GetUserStatsFunc(ctx, githubID, githubClient)
	}
	return &domain.Stats{}, nil
}

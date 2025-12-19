package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/furarico/octo-deck-api/internal/domain"
)

type StatsService struct{}

func NewStatsService() *StatsService {
	return &StatsService{}
}

// GetUserStats は指定されたGitHub IDのユーザーの統計情報を取得する
func (s *StatsService) GetUserStats(ctx context.Context, githubID string, githubClient GitHubClient) (*domain.Stats, error) {
	id, err := strconv.ParseInt(githubID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid github id: %w", err)
	}

	githubStats, err := githubClient.GetUserStats(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}

	domainStats, err := githubStats.ToDomainStats()
	if err != nil {
		return nil, fmt.Errorf("failed to convert stats to domain: %w", err)
	}

	return domainStats, nil
}

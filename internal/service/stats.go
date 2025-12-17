package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/furarico/octo-deck-api/internal/github"
)

type StatsService struct{}

func NewStatsService() *StatsService {
	return &StatsService{}
}

// GetUserStats は指定されたGitHub IDのユーザーの統計情報を取得する
func (s *StatsService) GetUserStats(ctx context.Context, githubID string, githubClient *github.Client) (*github.ContributionStats, error) {
	id, err := strconv.ParseInt(githubID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid github id: %w", err)
	}

	stats, err := githubClient.GetContributionStats(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get contribution stats: %w", err)
	}

	return stats, nil
}

package handler

import (
	api "github.com/furarico/octo-deck-api/generated"
	"github.com/furarico/octo-deck-api/internal/domain"
	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"
)

// APIのCard型に変換する
func convertCardToAPI(card domain.Card) api.Card {
	return api.Card{
		GithubId: card.GithubID,
		UserName: card.UserName,
		FullName: card.FullName,
		IconUrl:  card.IconUrl,
		Identicon: api.Identicon{
			Blocks: convertBlocks(card.Blocks),
			Color:  string(card.Color),
		},
	}
}

// ドメインのBlocks型をAPIのBlocks型に変換する
func convertBlocks(blocks domain.Blocks) [][]bool {
	blocksArray := make([][]bool, 5)
	for i := range blocksArray {
		blocksArray[i] = make([]bool, 5)
		for j := range blocksArray[i] {
			blocksArray[i][j] = blocks[i][j]
		}
	}
	return blocksArray
}

// APIのCommunity型に変換する
func convertCommunityToAPI(community domain.Community) api.Community {
	return api.Community{
		Id:   uuid.UUID(community.ID).String(),
		Name: community.Name,
	}
}

// UserStatsをAPIのUserStats型に変換する
func convertUserStatsToAPI(stats *domain.Stats) (api.UserStats, error) {
	contributions := make([]api.Contribution, len(stats.Contributions))
	for i, c := range stats.Contributions {
		contributions[i] = api.Contribution{
			Date:  types.Date{Time: c.Date},
			Count: int32(c.Count),
		}
	}

	return api.UserStats{
		Contributions:     contributions,
		TotalContribution: int32(stats.TotalContribution),
		ContributionDetail: api.ContributionDetail{
			ReviewCount:      int32(stats.ContributionDetail.ReviewCount),
			CommitCount:      int32(stats.ContributionDetail.CommitCount),
			IssueCount:       int32(stats.ContributionDetail.IssueCount),
			PullRequestCount: int32(stats.ContributionDetail.PullRequestCount),
		},
		MostUsedLanguage: api.Language{
			Name:  stats.MostUsedLanguage.LanguageName,
			Color: stats.MostUsedLanguage.Color,
		},
	}, nil
}

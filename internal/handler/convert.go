package handler

import (
	"time"

	api "github.com/furarico/octo-deck-api/generated"
	"github.com/furarico/octo-deck-api/internal/domain"
	"github.com/furarico/octo-deck-api/internal/github"
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

// ContributionStatsをAPIのUserStats型に変換する
func convertContributionStatsToAPI(stats *github.ContributionStats) (api.UserStats, error) {
	contributions := make([]api.Contribution, len(stats.Contributions))
	for i, c := range stats.Contributions {
		// DateをYYYY-MM-DD形式でパース
		date, err := time.Parse("2006-01-02", c.Date)
		if err != nil {
			return api.UserStats{}, err
		}

		contributions[i] = api.Contribution{
			Date:  types.Date{Time: date},
			Count: int32(c.Count),
		}
	}

	return api.UserStats{
		Contributions: contributions,
	}, nil
}

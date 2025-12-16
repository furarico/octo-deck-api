package handler

import (
	api "github.com/furarico/octo-deck-api/generated"
	"github.com/furarico/octo-deck-api/internal/domain"
	"github.com/google/uuid"
)

// APIのCard型に変換する
func convertCardToAPI(card domain.Card) api.Card {
	return api.Card{
		GithubId: card.GithubID,
		UserName: "", // TODO: GitHubから取得したユーザー名を設定
		FullName: "", // TODO: GitHubから取得したフルネームを設定
		IconUrl:  "", // TODO: GitHubから取得したアイコンURLを設定
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

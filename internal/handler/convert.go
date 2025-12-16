package handler

import (
	api "github.com/furarico/octo-deck-api/generated"
	"github.com/furarico/octo-deck-api/internal/domain"
	"github.com/google/uuid"
)

// APIのCard型に変換する
func convertCardToAPI(card domain.Card) api.Card {
	return api.Card{
		Id: uuid.UUID(card.ID).String(),
		Identicon: api.Identicon{
			Blocks: convertBlocks(card.Blocks),
			Color:  string(card.Color),
		},
		UserName: card.GithubID,
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

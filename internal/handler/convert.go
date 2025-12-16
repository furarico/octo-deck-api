package handler

import (
	api "github.com/furarico/octo-deck-api/generated"
	"github.com/furarico/octo-deck-api/internal/domain"
	"github.com/google/uuid"
)

// APIのCard型に変換する
func convertCardToAPI(card domain.CardWithOwner) api.Card {
	// Check for nil pointers to avoid panic
	if card.Owner == nil || card.Card == nil {
		return api.Card{}
	}
	return api.Card{
		FullName: card.Owner.FullName,
		IconUrl:  card.Owner.IconURL,
		Id:       uuid.UUID(card.Card.ID).String(),
		Identicon: api.Identicon{
			Blocks: convertBlocks(card.Owner.Identicon.Blocks),
			Color:  string(card.Owner.Identicon.Color),
		},
		UserName: card.Owner.UserName,
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

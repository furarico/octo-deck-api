package repository

import (
	"encoding/json"

	"github.com/furarico/octo-deck-api/internal/database"
	"github.com/furarico/octo-deck-api/internal/domain"
)

func (r *cardRepository) FindByID(id string) (*domain.CardWithOwner, error) {
	var dbCard database.Card
	r.db.First(&dbCard, "id = ?", id)

	var dbUser database.User
	r.db.First(&dbUser, "id = ?", dbCard.UserID)

	var dbIdenticon database.Identicon
	r.db.First(&dbIdenticon, "user_id = ?", dbUser.ID)

	return &domain.CardWithOwner{
		Card: &domain.Card{
			ID:      domain.CardID(dbCard.ID),
			OwnerID: domain.UserID(dbCard.UserID),
		},
		Owner: &domain.User{
			ID:       domain.UserID(dbUser.ID),
			UserName: dbUser.UserName,
			FullName: dbUser.FullName,
			IconURL:  dbUser.IconURL,
			Identicon: domain.Identicon{
				Color:  domain.Color(dbIdenticon.Color),
				Blocks: parseBlocks(dbIdenticon.BlocksData),
			},
		},
	}, nil
}

func parseBlocks(data json.RawMessage) domain.Blocks {
	var blocks domain.Blocks
	json.Unmarshal(data, &blocks)
	return blocks
}

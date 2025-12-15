package repository

import (
	"encoding/json"

	"github.com/furarico/octo-deck-api/internal/database"
	"github.com/furarico/octo-deck-api/internal/domain"
	"gorm.io/gorm"
)

type cardRepository struct {
	db *gorm.DB
}

func NewCardRepository(db *gorm.DB) *cardRepository {
	return &cardRepository{db: db}
}

// FindAll はGitHubIDから自分が集めたカードを全て取得する
func (r *cardRepository) FindAll(githubID string) ([]domain.CardWithOwner, error) {
	// ① GitHubIDからユーザーを特定
	var collector database.User
	if err := r.db.First(&collector, "github_id = ?", githubID).Error; err != nil {
		return nil, err
	}

	// ② collected_cardsテーブルからそのユーザーが集めたカードを取得
	var collectedCards []database.CollectedCard
	if err := r.db.Find(&collectedCards, "user_id = ?", collector.ID).Error; err != nil {
		return nil, err
	}

	// ③ 各カードの情報を取得
	var result []domain.CardWithOwner
	for _, cc := range collectedCards {
		// カードを取得
		var dbCard database.Card
		if err := r.db.First(&dbCard, "id = ?", cc.CardID).Error; err != nil {
			return nil, err
		}

		// カードの所有者を取得
		var dbUser database.User
		if err := r.db.First(&dbUser, "id = ?", dbCard.UserID).Error; err != nil {
			return nil, err
		}

		// Identiconを取得
		var dbIdenticon database.Identicon
		if err := r.db.First(&dbIdenticon, "user_id = ?", dbUser.ID).Error; err != nil {
			return nil, err
		}

		cardWithOwner := domain.CardWithOwner{
			Card: &domain.Card{
				ID:      domain.CardID(dbCard.ID),
				OwnerID: domain.UserID(dbCard.UserID),
			},
			Owner: &domain.User{
				ID:       domain.UserID(dbUser.ID),
				UserName: dbUser.UserName,
				FullName: dbUser.FullName,
				GitHubID: dbUser.GithubID,
				IconURL:  dbUser.IconURL,
				Identicon: domain.Identicon{
					Color:  domain.Color(dbIdenticon.Color),
					Blocks: parseBlocks(dbIdenticon.BlocksData),
				},
			},
		}
		result = append(result, cardWithOwner)
	}

	return result, nil
}

// FindByID は指定されたIDのカードを取得する
func (r *cardRepository) FindByID(id string) (*domain.CardWithOwner, error) {
	var dbCard database.Card
	if err := r.db.First(&dbCard, "id = ?", id).Error; err != nil {
		return nil, err
	}

	var dbUser database.User
	if err := r.db.First(&dbUser, "id = ?", dbCard.UserID).Error; err != nil {
		return nil, err
	}

	var dbIdenticon database.Identicon
	if err := r.db.First(&dbIdenticon, "user_id = ?", dbUser.ID).Error; err != nil {
		return nil, err
	}

	return &domain.CardWithOwner{
		Card: &domain.Card{
			ID:      domain.CardID(dbCard.ID),
			OwnerID: domain.UserID(dbCard.UserID),
		},
		Owner: &domain.User{
			ID:       domain.UserID(dbUser.ID),
			UserName: dbUser.UserName,
			FullName: dbUser.FullName,
			GitHubID: dbUser.GithubID,
			IconURL:  dbUser.IconURL,
			Identicon: domain.Identicon{
				Color:  domain.Color(dbIdenticon.Color),
				Blocks: parseBlocks(dbIdenticon.BlocksData),
			},
		},
	}, nil
}

// FindMyCard はGitHubIDから自分のカードを取得する
func (r *cardRepository) FindMyCard(githubID string) (*domain.CardWithOwner, error) {
	var dbUser database.User
	if err := r.db.First(&dbUser, "github_id = ?", githubID).Error; err != nil {
		return nil, err
	}

	var dbCard database.Card
	if err := r.db.First(&dbCard, "user_id = ?", dbUser.ID).Error; err != nil {
		return nil, err
	}

	var dbIdenticon database.Identicon
	if err := r.db.First(&dbIdenticon, "user_id = ?", dbUser.ID).Error; err != nil {
		return nil, err
	}

	return &domain.CardWithOwner{
		Card: &domain.Card{
			ID:      domain.CardID(dbCard.ID),
			OwnerID: domain.UserID(dbCard.UserID),
		},
		Owner: &domain.User{
			ID:       domain.UserID(dbUser.ID),
			UserName: dbUser.UserName,
			FullName: dbUser.FullName,
			GitHubID: dbUser.GithubID,
			IconURL:  dbUser.IconURL,
			Identicon: domain.Identicon{
				Color:  domain.Color(dbIdenticon.Color),
				Blocks: parseBlocks(dbIdenticon.BlocksData),
			},
		},
	}, nil
}

// parseBlocks はDBのJSON形式をドメインのBlocks型に変換する
func parseBlocks(data json.RawMessage) domain.Blocks {
	var blocks domain.Blocks
	json.Unmarshal(data, &blocks)
	return blocks
}

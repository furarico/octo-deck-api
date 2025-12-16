package repository

import (
	"encoding/json"
	"fmt"

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

	// ② Preloadで関連データを一括取得
	var collectedCards []database.CollectedCard
	if err := r.db.
		Preload("Card").
		Preload("Card.User").
		Preload("Card.User.Identicon").
		Where("user_id = ?", collector.ID).
		Find(&collectedCards).Error; err != nil {
		return nil, err
	}

	// ③ ドメインモデルに変換（DBアクセスなし）
	var result []domain.CardWithOwner
	for _, cc := range collectedCards {
		blocks, err := parseBlocks(cc.Card.User.Identicon.BlocksData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse blocks: %w", err)
		}
		cardWithOwner := domain.CardWithOwner{
			Card: &domain.Card{
				ID:      domain.CardID(cc.Card.ID),
				OwnerID: domain.UserID(cc.Card.UserID),
			},
			Owner: &domain.User{
				ID:       domain.UserID(cc.Card.User.ID),
				UserName: cc.Card.User.UserName,
				FullName: cc.Card.User.FullName,
				GitHubID: cc.Card.User.GithubID,
				IconURL:  cc.Card.User.IconURL,
				Identicon: domain.Identicon{
					Color:  domain.Color(cc.Card.User.Identicon.Color),
					Blocks: blocks,
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
	if err := r.db.
		Preload("User").
		Preload("User.Identicon").
		First(&dbCard, "id = ?", id).Error; err != nil {
		return nil, err
	}

	blocks, err := parseBlocks(dbCard.User.Identicon.BlocksData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse blocks: %w", err)
	}

	return &domain.CardWithOwner{
		Card: &domain.Card{
			ID:      domain.CardID(dbCard.ID),
			OwnerID: domain.UserID(dbCard.UserID),
		},
		Owner: &domain.User{
			ID:       domain.UserID(dbCard.User.ID),
			UserName: dbCard.User.UserName,
			FullName: dbCard.User.FullName,
			GitHubID: dbCard.User.GithubID,
			IconURL:  dbCard.User.IconURL,
			Identicon: domain.Identicon{
				Color:  domain.Color(dbCard.User.Identicon.Color),
				Blocks: blocks,
			},
		},
	}, nil
}

// FindMyCard はGitHubIDから自分のカードを取得する
func (r *cardRepository) FindMyCard(githubID string) (*domain.CardWithOwner, error) {
	var dbUser database.User
	if err := r.db.
		Preload("Identicon").
		First(&dbUser, "github_id = ?", githubID).Error; err != nil {
		return nil, err
	}

	var dbCard database.Card
	if err := r.db.First(&dbCard, "user_id = ?", dbUser.ID).Error; err != nil {
		return nil, err
	}

	blocks, err := parseBlocks(dbUser.Identicon.BlocksData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse blocks: %w", err)
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
				Color:  domain.Color(dbUser.Identicon.Color),
				Blocks: blocks,
			},
		},
	}, nil
}

// parseBlocks はDBのJSON形式をドメインのBlocks型に変換する
func parseBlocks(data json.RawMessage) (domain.Blocks, error) {
	var blocks domain.Blocks
	if err := json.Unmarshal(data, &blocks); err != nil {
		return domain.Blocks{}, err
	}
	return blocks, nil
}

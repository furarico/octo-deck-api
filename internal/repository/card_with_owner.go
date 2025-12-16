package repository

import (
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
	// GitHubIDからユーザーを特定
	var collector database.User
	if err := r.db.First(&collector, "github_id = ?", githubID).Error; err != nil {
		return nil, err
	}

	// Preloadで関連データを一括取得
	var collectedCards []database.CollectedCard
	if err := r.db.
		Preload("Card").
		Preload("Card.User").
		Preload("Card.User.Identicon").
		Where("user_id = ?", collector.ID).
		Find(&collectedCards).Error; err != nil {
		return nil, err
	}

	// ドメインモデルに変換（DBアクセスなし）
	var result []domain.CardWithOwner
	for _, cc := range collectedCards {
		cardWithOwner := domain.NewCardWithOwner(
			*cc.Card.ToDomain(),
			*cc.Card.User.ToDomain(),
		)
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

	return &domain.CardWithOwner{
		Card:  *dbCard.ToDomain(),
		Owner: *dbCard.User.ToDomain(),
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

	return &domain.CardWithOwner{
		Card:  *dbCard.ToDomain(),
		Owner: *dbUser.ToDomain(),
	}, nil
}

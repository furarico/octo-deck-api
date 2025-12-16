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
func (r *cardRepository) FindAll(githubID string) ([]domain.Card, error) {
	var collectedCards []database.CollectedCard
	if err := r.db.
		Preload("Card").
		Where("collector_github_id = ?", githubID).
		Find(&collectedCards).Error; err != nil {
		return nil, err
	}

	var result []domain.Card
	for _, cc := range collectedCards {
		result = append(result, *cc.Card.ToDomain())
	}

	return result, nil
}

// FindByID は指定されたIDのカードを取得する
func (r *cardRepository) FindByID(id string) (*domain.Card, error) {
	var dbCard database.Card
	if err := r.db.First(&dbCard, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return dbCard.ToDomain(), nil
}

// FindMyCard はGitHubIDから自分のカードを取得する
func (r *cardRepository) FindMyCard(githubID string) (*domain.Card, error) {
	var dbCard database.Card
	if err := r.db.First(&dbCard, "github_id = ?", githubID).Error; err != nil {
		return nil, err
	}

	return dbCard.ToDomain(), nil
}

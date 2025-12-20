package repository

import (
	"context"

	"github.com/furarico/octo-deck-api/internal/database"
	"github.com/furarico/octo-deck-api/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type cardRepository struct {
	db *gorm.DB
}

func NewCardRepository(db *gorm.DB) *cardRepository {
	return &cardRepository{db: db}
}

// FindAll はGitHubIDから自分が集めたカードを全て取得する
func (r *cardRepository) FindAll(ctx context.Context, githubID string) ([]domain.Card, error) {
	var collectedCards []database.CollectedCard
	if err := r.db.WithContext(ctx).
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

// FindByGitHubID はGitHub IDでカードを取得する
func (r *cardRepository) FindByGitHubID(ctx context.Context, githubID string) (*domain.Card, error) {
	var dbCard database.Card
	if err := r.db.WithContext(ctx).First(&dbCard, "github_id = ?", githubID).Error; err != nil {
		return nil, err
	}

	return dbCard.ToDomain(), nil
}

// FindMyCard はGitHubIDから自分のカードを取得する
func (r *cardRepository) FindMyCard(ctx context.Context, githubID string) (*domain.Card, error) {
	var dbCard database.Card
	if err := r.db.WithContext(ctx).First(&dbCard, "github_id = ?", githubID).Error; err != nil {
		return nil, err
	}

	return dbCard.ToDomain(), nil
}

// Create は新しいカードを作成する
func (r *cardRepository) Create(ctx context.Context, card *domain.Card) error {
	dbCard := database.CardFromDomain(card)
	return r.db.WithContext(ctx).Create(dbCard).Error
}

// AddToCollectedCards はカードをデッキに追加する
func (r *cardRepository) AddToCollectedCards(ctx context.Context, collectorGithubID string, cardID domain.CardID) error {
	collectedCard := &database.CollectedCard{
		CollectorGithubID: collectorGithubID,
		CardID:            uuid.UUID(cardID),
	}
	return r.db.WithContext(ctx).Create(collectedCard).Error
}

// RemoveFromCollectedCards はカードをデッキから削除する
func (r *cardRepository) RemoveFromCollectedCards(ctx context.Context, collectorGithubID string, cardID domain.CardID) error {
	return r.db.WithContext(ctx).
		Where("collector_github_id = ? AND card_id = ?", collectorGithubID, uuid.UUID(cardID)).
		Delete(&database.CollectedCard{}).Error
}

// Update はカード情報を更新する
func (r *cardRepository) Update(ctx context.Context, card *domain.Card) error {
	dbCard := database.CardFromDomain(card)
	return r.db.WithContext(ctx).Model(&database.Card{}).Where("id = ?", dbCard.ID).Updates(dbCard).Error
}

// FindAllCardsInDB はデータベース内の全カードを取得する
func (r *cardRepository) FindAllCardsInDB(ctx context.Context) ([]domain.Card, error) {
	var dbCards []database.Card
	if err := r.db.WithContext(ctx).Find(&dbCards).Error; err != nil {
		return nil, err
	}

	result := make([]domain.Card, 0, len(dbCards))
	for _, dbCard := range dbCards {
		result = append(result, *dbCard.ToDomain())
	}

	return result, nil
}

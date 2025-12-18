package repository

import (
	"github.com/furarico/octo-deck-api/internal/database"
	"github.com/furarico/octo-deck-api/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type communityRepository struct {
	db *gorm.DB
}

func NewCommunityRepository(db *gorm.DB) *communityRepository {
	return &communityRepository{db: db}
}

// FindAll はすべてのコミュニティを取得する
func (r *communityRepository) FindAll(githubID string) ([]domain.Community, error) {
	var communities []database.Community
	if err := r.db.
		Model(&database.Community{}).
		Joins("JOIN community_cards cc ON cc.community_id = communities.id").
		Joins("JOIN cards c ON c.id = cc.card_id AND c.github_id = ?", githubID).
		Distinct().
		Find(&communities).Error; err != nil {
		return nil, err
	}

	var result []domain.Community
	for _, community := range communities {
		result = append(result, *community.ToDomain())
	}

	return result, nil
}

// FindByID は指定されたコミュニティIDの情報を取得する
func (r *communityRepository) FindByID(id string) (*domain.Community, error) {
	var community database.Community
	if err := r.db.First(&community, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return community.ToDomain(), nil
}

// FindCards は指定したコミュニティIDのカード一覧を取得する
func (r *communityRepository) FindCards(id string) ([]domain.Card, error) {
	var cards []database.Card
	if err := r.db.
		Joins("JOIN community_cards cc ON cc.card_id = cards.id").
		Where("cc.community_id = ?", id).
		Find(&cards).Error; err != nil {
		return nil, err
	}

	var result []domain.Card
	for _, card := range cards {
		result = append(result, *card.ToDomain())
	}

	return result, nil
}

// Create はコミュニティを作成する
func (r *communityRepository) Create(community *domain.Community) error {
	dbCommunity := &database.Community{
		ID:   uuid.UUID(community.ID),
		Name: community.Name,
	}

	return r.db.Create(dbCommunity).Error
}

// Delete はコミュニティを削除する
func (r *communityRepository) Delete(id string) error {
	return r.db.Delete(&database.Community{}, "id = ?", id).Error
}

// AddCard はコミュニティにカードを追加する
func (r *communityRepository) AddCard(communityID string, cardID string) error {
	communityCard := &database.CommunityCard{
		CommunityID: parseUUID(communityID),
		CardID:      parseUUID(cardID),
	}

	return r.db.Create(communityCard).Error
}

// RemoveCard はコミュニティからカードを削除する
func (r *communityRepository) RemoveCard(communityID string, cardID string) error {
	return r.db.
		Where("community_id = ? AND card_id = ?", communityID, cardID).
		Delete(&database.CommunityCard{}).Error
}

// parseUUID はstringをuuid.UUIDに変換する
func parseUUID(s string) uuid.UUID {
	id, _ := uuid.Parse(s)
	return id
}

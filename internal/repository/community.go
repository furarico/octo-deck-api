package repository

import (
	"context"
	"fmt"

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
func (r *communityRepository) FindAll(ctx context.Context, githubID string) ([]domain.Community, error) {
	var communities []database.Community
	if err := r.db.WithContext(ctx).
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
func (r *communityRepository) FindByID(ctx context.Context, id string) (*domain.Community, error) {
	var community database.Community
	if err := r.db.WithContext(ctx).First(&community, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return community.ToDomain(), nil
}

// FindByIDWithHighlightedCard は指定されたコミュニティIDの情報をHighlightedCard付きで取得する
func (r *communityRepository) FindByIDWithHighlightedCard(ctx context.Context, id string) (*domain.Community, error) {
	var community database.Community
	if err := r.db.WithContext(ctx).
		Preload("BestContributorCard").
		Preload("BestCommitterCard").
		Preload("BestIssuerCard").
		Preload("BestPullRequesterCard").
		Preload("BestReviewerCard").
		First(&community, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return community.ToDomain(), nil
}

// UpdateHighlightedCard はコミュニティのHighlightedCardを更新する
func (r *communityRepository) UpdateHighlightedCard(ctx context.Context, communityID string, highlightedCard *domain.HighlightedCard) error {
	communityUUID, err := parseUUID(communityID)
	if err != nil {
		return fmt.Errorf("invalid community id: %w", err)
	}

	updates := map[string]interface{}{}

	// 各ベストカードのIDを設定（空のカードはnilにする）
	if highlightedCard.BestContributor.GithubID != "" {
		cardID := uuid.UUID(highlightedCard.BestContributor.ID)
		updates["best_contributor_card_id"] = &cardID
	} else {
		updates["best_contributor_card_id"] = nil
	}

	if highlightedCard.BestCommitter.GithubID != "" {
		cardID := uuid.UUID(highlightedCard.BestCommitter.ID)
		updates["best_committer_card_id"] = &cardID
	} else {
		updates["best_committer_card_id"] = nil
	}

	if highlightedCard.BestIssuer.GithubID != "" {
		cardID := uuid.UUID(highlightedCard.BestIssuer.ID)
		updates["best_issuer_card_id"] = &cardID
	} else {
		updates["best_issuer_card_id"] = nil
	}

	if highlightedCard.BestPullRequester.GithubID != "" {
		cardID := uuid.UUID(highlightedCard.BestPullRequester.ID)
		updates["best_pull_requester_card_id"] = &cardID
	} else {
		updates["best_pull_requester_card_id"] = nil
	}

	if highlightedCard.BestReviewer.GithubID != "" {
		cardID := uuid.UUID(highlightedCard.BestReviewer.ID)
		updates["best_reviewer_card_id"] = &cardID
	} else {
		updates["best_reviewer_card_id"] = nil
	}

	return r.db.WithContext(ctx).Model(&database.Community{}).Where("id = ?", communityUUID).Updates(updates).Error
}

// FindCards は指定したコミュニティIDのカード一覧を取得する
func (r *communityRepository) FindCards(ctx context.Context, id string) ([]domain.Card, error) {
	var cards []database.Card
	if err := r.db.WithContext(ctx).
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
func (r *communityRepository) Create(ctx context.Context, community *domain.Community) error {
	dbCommunity := &database.Community{
		ID:        uuid.UUID(community.ID),
		Name:      community.Name,
		StartedAt: community.StartedAt,
		EndedAt:   community.EndedAt,
	}

	return r.db.WithContext(ctx).Create(dbCommunity).Error
}

// Delete はコミュニティを削除する
func (r *communityRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&database.Community{}, "id = ?", id).Error
}

// AddCard はコミュニティにカードを追加する
func (r *communityRepository) AddCard(ctx context.Context, communityID string, cardID string) error {
	communityUUID, err := parseUUID(communityID)
	if err != nil {
		return fmt.Errorf("invalid community id: %w", err)
	}
	cardUUID, err := parseUUID(cardID)
	if err != nil {
		return fmt.Errorf("invalid card id: %w", err)
	}

	communityCard := &database.CommunityCard{
		CommunityID: communityUUID,
		CardID:      cardUUID,
	}

	return r.db.WithContext(ctx).Create(communityCard).Error
}

// RemoveCard はコミュニティからカードを削除する
func (r *communityRepository) RemoveCard(ctx context.Context, communityID string, cardID string) error {
	communityUUID, err := parseUUID(communityID)
	if err != nil {
		return fmt.Errorf("invalid community id: %w", err)
	}
	cardUUID, err := parseUUID(cardID)
	if err != nil {
		return fmt.Errorf("invalid card id: %w", err)
	}

	return r.db.WithContext(ctx).
		Where("community_id = ? AND card_id = ?", communityUUID, cardUUID).
		Delete(&database.CommunityCard{}).Error
}

// parseUUID はstringをuuid.UUIDに変換する
func parseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

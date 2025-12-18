package service

import (
	"fmt"
	"time"

	"github.com/furarico/octo-deck-api/internal/domain"
)

// CommunityRepository はServiceが必要とするRepositoryのインターフェース
type CommunityRepository interface {
	FindAll(githubID string) ([]domain.Community, error)
	FindByID(id string) (*domain.Community, error)
	FindCards(id string) ([]domain.Card, error)
	Create(community *domain.Community) error
	Delete(id string) error
	AddCard(communityID string, cardID string) error
	RemoveCard(communityID string, cardID string) error
}

type CommunityService struct {
	communityRepo CommunityRepository
}

func NewCommunityService(communityRepo CommunityRepository) *CommunityService {
	return &CommunityService{
		communityRepo: communityRepo,
	}
}

// GetAllCommunities はすべてのコミュニティを取得する
func (s *CommunityService) GetAllCommunities(githubID string) ([]domain.Community, error) {
	communities, err := s.communityRepo.FindAll(githubID)
	if err != nil {
		return nil, fmt.Errorf("failed to get all communities: %w", err)
	}

	return communities, nil
}

// GetCommunityByID は指定されたコミュニティIDの情報を取得する
func (s *CommunityService) GetCommunityByID(id string) (*domain.Community, error) {
	community, err := s.communityRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get community by id: %w", err)
	}

	if community == nil {
		return nil, fmt.Errorf("community not found: id=%s", id)
	}

	return community, nil
}

// GetCommunityCards は指定したコミュニティIDのカード一覧を取得する
func (s *CommunityService) GetCommunityCards(id string) ([]domain.Card, error) {
	cards, err := s.communityRepo.FindCards(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get community cards: %w", err)
	}

	return cards, nil
}

// CreateCommunity はコミュニティを作成する
func (s *CommunityService) CreateCommunity(name string, startedAt time.Time, endedAt time.Time, bestContribute domain.BestContribute) (*domain.Community, error) {
	community := domain.NewCommunity(name, startedAt, endedAt, bestContribute)

	if err := s.communityRepo.Create(community); err != nil {
		return nil, fmt.Errorf("failed to create community: %w", err)
	}

	return community, nil
}

// DeleteCommunity はコミュニティを削除する
func (s *CommunityService) DeleteCommunity(id string) error {
	if err := s.communityRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete community: %w", err)
	}

	return nil
}

// AddCardToCommunity はコミュニティにカードを追加する
func (s *CommunityService) AddCardToCommunity(communityID string, cardID string) error {
	if err := s.communityRepo.AddCard(communityID, cardID); err != nil {
		return fmt.Errorf("failed to add card to community: %w", err)
	}

	return nil
}

// RemoveCardFromCommunity はコミュニティからカードを削除する
func (s *CommunityService) RemoveCardFromCommunity(communityID string, cardID string) error {
	if err := s.communityRepo.RemoveCard(communityID, cardID); err != nil {
		return fmt.Errorf("failed to remove card from community: %w", err)
	}

	return nil
}

package service

import (
	"fmt"

	"github.com/furarico/octo-deck-api/internal/domain"
)

// CommunityRepository はServiceが必要とするRepositoryのインターフェース
type CommunityRepository interface {
	FindAll(githubID string) ([]domain.Community, error)
	FindByID(id string) (*domain.Community, error)
	FindCards(id string) ([]domain.Card, error)
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

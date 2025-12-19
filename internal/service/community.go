package service

import (
	"context"
	"fmt"
	"time"

	"github.com/furarico/octo-deck-api/internal/domain"
	"github.com/furarico/octo-deck-api/internal/github"
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

// GetCommunityWithHighlightedCard はコミュニティとHighlightedCardを取得する
func (s *CommunityService) GetCommunityWithHighlightedCard(ctx context.Context, id string, githubClient GitHubClient) (*domain.Community, *domain.HighlightedCard, error) {
	// コミュニティを取得
	community, err := s.communityRepo.FindByID(id)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get community by id: %w", err)
	}
	if community == nil {
		return nil, nil, fmt.Errorf("community not found: id=%s", id)
	}

	// コミュニティのカード一覧を取得
	cards, err := s.communityRepo.FindCards(id)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get community cards: %w", err)
	}

	// カードがない場合は空のHighlightedCardを返す
	if len(cards) == 0 {
		return community, &domain.HighlightedCard{}, nil
	}

	// バッチ処理でGitHub情報を一括取得
	if err := EnrichCardsWithGitHubInfo(ctx, cards, githubClient); err != nil {
		return nil, nil, fmt.Errorf("failed to enrich cards with github info: %w", err)
	}

	// ユーザー名リストを作成
	usernames := make([]string, 0, len(cards))
	cardByUsername := make(map[string]domain.Card)
	for _, card := range cards {
		if card.UserName != "" {
			usernames = append(usernames, card.UserName)
			cardByUsername[card.UserName] = card
		}
	}

	// ユーザーがいない場合は空のHighlightedCardを返す
	if len(usernames) == 0 {
		return community, &domain.HighlightedCard{}, nil
	}

	// GitHub APIで貢献データを取得
	contributions, err := githubClient.GetUsersContributions(ctx, usernames, community.StartedAt, community.EndedAt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get users contributions: %w", err)
	}

	// 各カテゴリのベストユーザーを計算
	highlightedCard := calculateHighlightedCard(contributions, cardByUsername)

	return community, highlightedCard, nil
}

// calculateHighlightedCard は貢献データから各カテゴリのベストユーザーを計算する
func calculateHighlightedCard(contributions []github.UserContributionStats, cardByUsername map[string]domain.Card) *domain.HighlightedCard {
	if len(contributions) == 0 {
		return &domain.HighlightedCard{}
	}

	var bestContributor, bestCommitter, bestIssuer, bestPRer, bestReviewer github.UserContributionStats

	for _, c := range contributions {
		if c.Total > bestContributor.Total {
			bestContributor = c
		}
		if c.Commits > bestCommitter.Commits {
			bestCommitter = c
		}
		if c.Issues > bestIssuer.Issues {
			bestIssuer = c
		}
		if c.PRs > bestPRer.PRs {
			bestPRer = c
		}
		if c.Reviews > bestReviewer.Reviews {
			bestReviewer = c
		}
	}

	return &domain.HighlightedCard{
		BestContributor:   cardByUsername[bestContributor.Login],
		BestCommitter:     cardByUsername[bestCommitter.Login],
		BestIssuer:        cardByUsername[bestIssuer.Login],
		BestPullRequester: cardByUsername[bestPRer.Login],
		BestReviewer:      cardByUsername[bestReviewer.Login],
	}
}

// GetCommunityCards は指定したコミュニティIDのカード一覧を取得し、GitHub情報で補完する
func (s *CommunityService) GetCommunityCards(ctx context.Context, id string, githubClient GitHubClient) ([]domain.Card, error) {
	cards, err := s.communityRepo.FindCards(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get community cards: %w", err)
	}

	// バッチ処理でGitHub情報を一括取得
	if err := EnrichCardsWithGitHubInfo(ctx, cards, githubClient); err != nil {
		return nil, fmt.Errorf("failed to enrich cards with github info: %w", err)
	}

	return cards, nil
}

// CreateCommunityWithPeriod は集計期間を指定してコミュニティを作成する
func (s *CommunityService) CreateCommunityWithPeriod(name string, startDateTime, endDateTime time.Time) (*domain.Community, error) {
	community := domain.NewCommunity(name, startDateTime, endDateTime, domain.HighlightedCard{})

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

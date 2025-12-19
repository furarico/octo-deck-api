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
// 最適化: 統合GraphQLクエリを使用して、ユーザー情報・貢献データ・言語情報を1回のAPI呼び出しで取得
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

	// NodeIDリストを作成し、NodeID -> Cardのインデックスマップを構築
	nodeIDs := make([]string, 0, len(cards))
	cardIndexByNodeID := make(map[string]int)
	for i, card := range cards {
		if card.NodeID != "" {
			nodeIDs = append(nodeIDs, card.NodeID)
			cardIndexByNodeID[card.NodeID] = i
		}
	}

	// NodeIDがない場合は空のHighlightedCardを返す
	if len(nodeIDs) == 0 {
		return community, &domain.HighlightedCard{}, nil
	}

	// 統合GraphQLクエリで全情報を一括取得（ユーザー情報、貢献データ、言語情報）
	// これにより3回のAPI呼び出しが1回に削減される
	usersFullInfo, err := githubClient.GetUsersFullInfoByNodeIDs(ctx, nodeIDs, community.StartedAt, community.EndedAt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get users full info by node ids: %w", err)
	}

	// データがない場合は空のHighlightedCardを返す
	if len(usersFullInfo) == 0 {
		return community, &domain.HighlightedCard{}, nil
	}

	// 各カテゴリのベストユーザーを計算し、カード情報を構築
	highlightedCard := calculateHighlightedCardFromFullInfo(usersFullInfo, cards, cardIndexByNodeID, nodeIDs)

	return community, highlightedCard, nil
}

// calculateHighlightedCardFromFullInfo は統合クエリの結果から各カテゴリのベストユーザーを計算する
func calculateHighlightedCardFromFullInfo(usersFullInfo []github.UserFullInfo, cards []domain.Card, cardIndexByNodeID map[string]int, nodeIDs []string) *domain.HighlightedCard {
	if len(usersFullInfo) == 0 {
		return &domain.HighlightedCard{}
	}

	// login -> UserFullInfoのマップを構築
	userInfoByLogin := make(map[string]github.UserFullInfo)
	for _, info := range usersFullInfo {
		userInfoByLogin[info.Login] = info
	}

	// 各カテゴリのベストユーザーを特定
	var bestContributor, bestCommitter, bestIssuer, bestPRer, bestReviewer github.UserFullInfo

	for _, info := range usersFullInfo {
		if info.Total > bestContributor.Total {
			bestContributor = info
		}
		if info.Commits > bestCommitter.Commits {
			bestCommitter = info
		}
		if info.Issues > bestIssuer.Issues {
			bestIssuer = info
		}
		if info.PRs > bestPRer.PRs {
			bestPRer = info
		}
		if info.Reviews > bestReviewer.Reviews {
			bestReviewer = info
		}
	}

	// login -> nodeID のマッピングを構築
	// usersFullInfoの順序はnodeIDsと同じであることを前提とする
	loginToNodeID := make(map[string]string)
	for i, info := range usersFullInfo {
		if i < len(nodeIDs) {
			loginToNodeID[info.Login] = nodeIDs[i]
		}
	}

	// ベストユーザーのカードを構築するヘルパー関数
	buildCard := func(info github.UserFullInfo) domain.Card {
		if info.Login == "" {
			return domain.Card{}
		}

		// nodeIDからカードを探す
		nodeID, ok := loginToNodeID[info.Login]
		if !ok {
			return domain.Card{}
		}

		cardIdx, ok := cardIndexByNodeID[nodeID]
		if !ok {
			return domain.Card{}
		}

		card := cards[cardIdx]
		// GitHub APIから取得した情報でカードを補完
		card.UserName = info.Login
		card.FullName = info.Name
		card.IconUrl = info.AvatarURL
		card.MostUsedLanguage = domain.Language{
			LanguageName: info.MostUsedLanguage,
			Color:        info.MostUsedLanguageColor,
		}

		return card
	}

	return &domain.HighlightedCard{
		BestContributor:   buildCard(bestContributor),
		BestCommitter:     buildCard(bestCommitter),
		BestIssuer:        buildCard(bestIssuer),
		BestPullRequester: buildCard(bestPRer),
		BestReviewer:      buildCard(bestReviewer),
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

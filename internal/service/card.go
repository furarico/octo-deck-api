package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/furarico/octo-deck-api/internal/domain"
	"gorm.io/gorm"
)

// CardRepository はServiceが必要とするRepositoryのインターフェース
type CardRepository interface {
	FindAll(githubID string) ([]domain.Card, error)
	FindByGitHubID(githubID string) (*domain.Card, error)
	FindMyCard(githubID string) (*domain.Card, error)
	Create(card *domain.Card) error
	Update(card *domain.Card) error
	AddToCollectedCards(collectorGithubID string, cardID domain.CardID) error
	RemoveFromCollectedCards(collectorGithubID string, cardID domain.CardID) error
}

// IdenticonGenerator はServiceが必要とするIdenticon Generatorのインターフェース
type IdenticonGenerator interface {
	Generate(githubID string) (domain.Color, domain.Blocks, error)
}

type CardService struct {
	cardRepo           CardRepository
	identiconGenerator IdenticonGenerator
}

func NewCardService(cardRepo CardRepository, identiconGenerator IdenticonGenerator) *CardService {
	return &CardService{
		cardRepo:           cardRepo,
		identiconGenerator: identiconGenerator,
	}
}

// GetAllCards は自分が集めたカードを全て取得する
func (s *CardService) GetAllCards(ctx context.Context, githubID string, githubClient GitHubClient) ([]domain.Card, error) {
	cards, err := s.cardRepo.FindAll(githubID)
	if err != nil {
		return nil, fmt.Errorf("failed to get all cards: %w", err)
	}

	// バッチ処理でGitHub情報を一括取得
	if err := EnrichCardsWithGitHubInfo(ctx, cards, githubClient); err != nil {
		return nil, fmt.Errorf("failed to enrich cards with github info: %w", err)
	}

	return cards, nil
}

// GetCardByGitHubID は指定されたGitHub IDのカードを取得する
func (s *CardService) GetCardByGitHubID(ctx context.Context, githubID string, githubClient GitHubClient) (*domain.Card, error) {
	card, err := s.cardRepo.FindByGitHubID(githubID)
	if err != nil {
		return nil, fmt.Errorf("failed to get card by github id: %w", err)
	}

	if card == nil {
		return nil, fmt.Errorf("card not found: githubID=%s", githubID)
	}

	// GitHub APIからユーザー情報を取得して補完
	if err := EnrichCardWithGitHubInfo(ctx, card, githubClient); err != nil {
		return nil, err
	}

	return card, nil
}

// GetMyCard は自分のカードを取得する
func (s *CardService) GetMyCard(ctx context.Context, githubID string, githubClient GitHubClient) (*domain.Card, error) {
	card, err := s.cardRepo.FindMyCard(githubID)
	if err != nil {
		return nil, fmt.Errorf("failed to get my card: %w", err)
	}

	if card == nil {
		return nil, fmt.Errorf("my card not found")
	}

	// GitHub APIからユーザー情報を取得して補完
	if err := EnrichCardWithGitHubInfo(ctx, card, githubClient); err != nil {
		return nil, err
	}

	return card, nil
}

// GetOrCreateMyCard は自分のカードを取得し、存在しない場合は新規作成する
func (s *CardService) GetOrCreateMyCard(ctx context.Context, githubID string, nodeID string, githubClient GitHubClient) (*domain.Card, error) {
	card, err := s.cardRepo.FindMyCard(githubID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to get my card: %w", err)
	}

	// カードが存在しない場合は新規作成
	if card == nil || errors.Is(err, gorm.ErrRecordNotFound) {
		// GitHub APIから自分のユーザー情報を取得
		userInfo, err := githubClient.GetAuthenticatedUser(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get authenticated user: %w", err)
		}

		color, blocks, err := s.identiconGenerator.Generate(githubID)
		if err != nil {
			return nil, fmt.Errorf("failed to generate identicon: %w", err)
		}

		// MostUsedLanguageを取得
		langName, langColor, err := githubClient.GetMostUsedLanguage(ctx, userInfo.Login)
		if err != nil {
			return nil, fmt.Errorf("failed to get most used language: %w", err)
		}

		card = domain.NewCard(
			githubID,
			nodeID,
			color,
			blocks,
			domain.Language{LanguageName: langName, Color: langColor},
			userInfo.Login,
			userInfo.Name,
			userInfo.AvatarURL,
		)
		if err := s.cardRepo.Create(card); err != nil {
			return nil, fmt.Errorf("failed to create card: %w", err)
		}
	}

	// GitHub APIからユーザー情報を取得して補完
	if err := EnrichCardWithGitHubInfo(ctx, card, githubClient); err != nil {
		return nil, err
	}

	return card, nil
}

// AddCardToDeck はカードをデッキに追加する
func (s *CardService) AddCardToDeck(ctx context.Context, collectorGithubID string, targetGithubID string, githubClient GitHubClient) (*domain.Card, error) {
	// 追加対象のカードを取得
	card, err := s.cardRepo.FindByGitHubID(targetGithubID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("card not found: githubID=%s", targetGithubID)
		}
		return nil, fmt.Errorf("failed to find card: %w", err)
	}

	// デッキに追加
	if err := s.cardRepo.AddToCollectedCards(collectorGithubID, card.ID); err != nil {
		return nil, fmt.Errorf("failed to add card to deck: %w", err)
	}

	// GitHub APIからユーザー情報を取得して補完
	if err := EnrichCardWithGitHubInfo(ctx, card, githubClient); err != nil {
		return nil, err
	}

	return card, nil
}

// RemoveCardFromDeck はカードをデッキから削除する
func (s *CardService) RemoveCardFromDeck(ctx context.Context, collectorGithubID string, targetGithubID string, githubClient GitHubClient) (*domain.Card, error) {
	// 削除対象のカードを取得
	card, err := s.cardRepo.FindByGitHubID(targetGithubID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("card not found: githubID=%s", targetGithubID)
		}
		return nil, fmt.Errorf("failed to find card: %w", err)
	}

	// デッキから削除
	if err := s.cardRepo.RemoveFromCollectedCards(collectorGithubID, card.ID); err != nil {
		return nil, fmt.Errorf("failed to remove card from deck: %w", err)
	}

	// GitHub APIからユーザー情報を取得して補完
	if err := EnrichCardWithGitHubInfo(ctx, card, githubClient); err != nil {
		return nil, err
	}

	return card, nil
}

// EnrichCardWithGitHubInfo はGitHub APIからユーザー情報を取得してCardに設定する
func EnrichCardWithGitHubInfo(ctx context.Context, card *domain.Card, githubClient GitHubClient) error {
	githubID, err := strconv.ParseInt(card.GithubID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid github id: %w", err)
	}

	userInfo, err := githubClient.GetUserByID(ctx, githubID)
	if err != nil {
		return fmt.Errorf("failed to get github user info: %w", err)
	}

	card.UserName = userInfo.Login
	card.FullName = userInfo.Name
	card.IconUrl = userInfo.AvatarURL

	// MostUsedLanguageを取得して設定
	langName, langColor, err := githubClient.GetMostUsedLanguage(ctx, userInfo.Login)
	if err != nil {
		return fmt.Errorf("failed to get most used language: %w", err)
	}

	card.MostUsedLanguage = domain.Language{
		LanguageName: langName,
		Color:        langColor,
	}

	return nil
}

// EnrichCardsWithGitHubInfo は複数のカードにGitHub情報を一括で設定する（バッチ処理版）
// N+1問題を解消し、並列処理でパフォーマンスを向上させる
func EnrichCardsWithGitHubInfo(ctx context.Context, cards []domain.Card, githubClient GitHubClient) error {
	if len(cards) == 0 {
		return nil
	}

	// GitHub IDを収集
	githubIDs := make([]int64, 0, len(cards))
	idToIndex := make(map[int64][]int) // 同じGitHub IDを持つカードのインデックスを保持

	for i, card := range cards {
		githubID, err := strconv.ParseInt(card.GithubID, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid github id for card %d: %w", i, err)
		}
		if _, exists := idToIndex[githubID]; !exists {
			githubIDs = append(githubIDs, githubID)
		}
		idToIndex[githubID] = append(idToIndex[githubID], i)
	}

	// 一括でユーザー情報を取得（並列処理）
	userInfoMap, err := githubClient.GetUsersByIDs(ctx, githubIDs)
	if err != nil {
		return fmt.Errorf("failed to get users info: %w", err)
	}

	// ログイン名を収集
	logins := make([]string, 0, len(userInfoMap))
	for _, userInfo := range userInfoMap {
		logins = append(logins, userInfo.Login)
	}

	// 一括で言語情報を取得（並列処理）
	langInfoMap, err := githubClient.GetMostUsedLanguages(ctx, logins)
	if err != nil {
		return fmt.Errorf("failed to get languages info: %w", err)
	}

	// カードに情報を設定
	for githubID, indices := range idToIndex {
		userInfo, ok := userInfoMap[githubID]
		if !ok {
			continue
		}

		langInfo := langInfoMap[userInfo.Login]

		for _, idx := range indices {
			cards[idx].UserName = userInfo.Login
			cards[idx].FullName = userInfo.Name
			cards[idx].IconUrl = userInfo.AvatarURL
			cards[idx].MostUsedLanguage = domain.Language{
				LanguageName: langInfo.Name,
				Color:        langInfo.Color,
			}
		}
	}

	return nil
}

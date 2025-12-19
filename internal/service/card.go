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

	// バッチ処理でGitHub情報を補完（N+1問題解消）
	if err := EnrichCardsWithGitHubInfo(ctx, cards, githubClient); err != nil {
		return nil, err
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
func (s *CardService) GetOrCreateMyCard(ctx context.Context, githubID string, githubClient GitHubClient) (*domain.Card, error) {
	card, err := s.cardRepo.FindMyCard(githubID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to get my card: %w", err)
	}

	// カードが存在しない場合は新規作成
	if card == nil || errors.Is(err, gorm.ErrRecordNotFound) {
		color, blocks, err := s.identiconGenerator.Generate(githubID)
		if err != nil {
			return nil, fmt.Errorf("failed to generate identicon: %w", err)
		}

		card = domain.NewCard(githubID, color, blocks, domain.Language{})
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

// EnrichCardWithGitHubInfo はGitHub APIからユーザー情報を取得してCardに設定する（単一カード用）
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

// EnrichCardsWithGitHubInfo はGitHub APIから複数ユーザーの情報を一括取得してCardsに設定する（バッチ版）
// N+1問題を解消するため、2回のAPIコールで全カードの情報を取得する
func EnrichCardsWithGitHubInfo(ctx context.Context, cards []domain.Card, githubClient GitHubClient) error {
	if len(cards) == 0 {
		return nil
	}

	// 1. 全カードのGitHub IDを収集
	ids := make([]int64, 0, len(cards))
	for _, card := range cards {
		githubID, err := strconv.ParseInt(card.GithubID, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid github id: %w", err)
		}
		ids = append(ids, githubID)
	}

	// 2. バッチでユーザー情報を取得
	userInfoMap, err := githubClient.GetUsersByIDs(ctx, ids)
	if err != nil {
		return fmt.Errorf("failed to get users by ids: %w", err)
	}

	// 3. ログイン名リストを作成
	logins := make([]string, 0, len(userInfoMap))
	for _, userInfo := range userInfoMap {
		logins = append(logins, userInfo.Login)
	}

	// 4. バッチで言語情報を取得
	languageMap, err := githubClient.GetUsersLanguages(ctx, logins)
	if err != nil {
		return fmt.Errorf("failed to get users languages: %w", err)
	}

	// 5. 各カードに情報を設定
	for i := range cards {
		githubID, _ := strconv.ParseInt(cards[i].GithubID, 10, 64)
		userInfo, ok := userInfoMap[githubID]
		if !ok {
			continue
		}

		cards[i].UserName = userInfo.Login
		cards[i].FullName = userInfo.Name
		cards[i].IconUrl = userInfo.AvatarURL

		langInfo, ok := languageMap[userInfo.Login]
		if ok {
			cards[i].MostUsedLanguage = domain.Language{
				LanguageName: langInfo.Name,
				Color:        langInfo.Color,
			}
		} else {
			cards[i].MostUsedLanguage = domain.Language{
				LanguageName: "Unknown",
				Color:        "#586069",
			}
		}
	}

	return nil
}

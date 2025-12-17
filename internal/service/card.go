package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/furarico/octo-deck-api/internal/domain"
	"github.com/furarico/octo-deck-api/internal/github"
)

// CardRepository はServiceが必要とするRepositoryのインターフェース
type CardRepository interface {
	FindAll(githubID string) ([]domain.Card, error)
	FindByGitHubID(githubID string) (*domain.Card, error)
	FindMyCard(githubID string) (*domain.Card, error)
}

type CardService struct {
	cardRepo CardRepository
}

func NewCardService(cardRepo CardRepository) *CardService {
	return &CardService{
		cardRepo: cardRepo,
	}
}

// GetAllCards は自分が集めたカードを全て取得する
func (s *CardService) GetAllCards(ctx context.Context, githubID string, githubClient *github.Client) ([]domain.Card, error) {
	cards, err := s.cardRepo.FindAll(githubID)
	if err != nil {
		return nil, fmt.Errorf("failed to get all cards: %w", err)
	}

	// 各カードにGitHub情報を補完
	for i := range cards {
		if err := enrichCardWithGitHubInfo(ctx, &cards[i], githubClient); err != nil {
			return nil, err
		}
	}

	return cards, nil
}

// GetCardByGitHubID は指定されたGitHub IDのカードを取得する
func (s *CardService) GetCardByGitHubID(ctx context.Context, githubID string, githubClient *github.Client) (*domain.Card, error) {
	card, err := s.cardRepo.FindByGitHubID(githubID)
	if err != nil {
		return nil, fmt.Errorf("failed to get card by github id: %w", err)
	}

	if card == nil {
		return nil, fmt.Errorf("card not found: githubID=%s", githubID)
	}

	// GitHub APIからユーザー情報を取得して補完
	if err := enrichCardWithGitHubInfo(ctx, card, githubClient); err != nil {
		return nil, err
	}

	return card, nil
}

// GetMyCard は自分のカードを取得する
func (s *CardService) GetMyCard(ctx context.Context, githubID string, githubClient *github.Client) (*domain.Card, error) {
	card, err := s.cardRepo.FindMyCard(githubID)
	if err != nil {
		return nil, fmt.Errorf("failed to get my card: %w", err)
	}

	if card == nil {
		return nil, fmt.Errorf("my card not found")
	}

	// GitHub APIからユーザー情報を取得して補完
	if err := enrichCardWithGitHubInfo(ctx, card, githubClient); err != nil {
		return nil, err
	}

	return card, nil
}

// enrichCardWithGitHubInfo はGitHub APIからユーザー情報を取得してCardに設定する
func enrichCardWithGitHubInfo(ctx context.Context, card *domain.Card, githubClient *github.Client) error {
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

	return nil
}

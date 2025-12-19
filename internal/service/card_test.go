package service

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/furarico/octo-deck-api/internal/domain"
	"github.com/furarico/octo-deck-api/internal/github"
	"github.com/furarico/octo-deck-api/internal/identicon"
	"github.com/furarico/octo-deck-api/internal/repository"
	"gorm.io/gorm"
)

// テスト用のヘルパー関数: 正常なGitHubClientを返す
func createMockGitHubClient() *github.MockClient {
	return &github.MockClient{
		GetUserByIDFunc: func(ctx context.Context, id int64) (*github.UserInfo, error) {
			return &github.UserInfo{
				ID:        id,
				Login:     "testuser",
				Name:      "Test User",
				AvatarURL: "https://example.com/avatar.png",
			}, nil
		},
		GetMostUsedLanguageFunc: func(ctx context.Context, login string) (string, string, error) {
			return "Go", "#00ADD8", nil
		},
	}
}

// テスト用のヘルパー関数: 正常なカードを返す
func createTestCard(githubID string) *domain.Card {
	return &domain.Card{
		ID:       domain.NewCardID(),
		GithubID: githubID,
		Color:    domain.Color("#000000"),
		Blocks:   domain.Blocks{},
	}
}

// GetAllCards は自分が集めたカードを全て取得する
func TestGetAllCards(t *testing.T) {
	tests := []struct {
		name          string
		githubID      string
		setupRepo     func() *repository.MockCardRepository
		setupGitHub   func() *github.MockClient
		wantErr       bool
		wantErrMsg    string
		wantCardCount int
	}{
		{
			name:     "正常にカード一覧を取得できる",
			githubID: "12345",
			setupRepo: func() *repository.MockCardRepository {
				return &repository.MockCardRepository{
					FindAllFunc: func(githubID string) ([]domain.Card, error) {
						return []domain.Card{
							*createTestCard("12345"),
							*createTestCard("67890"),
						}, nil
					},
				}
			},
			setupGitHub:   createMockGitHubClient,
			wantErr:       false,
			wantCardCount: 2,
		},
		{
			name:     "Repositoryエラーが発生した場合",
			githubID: "12345",
			setupRepo: func() *repository.MockCardRepository {
				return &repository.MockCardRepository{
					FindAllFunc: func(githubID string) ([]domain.Card, error) {
						return nil, fmt.Errorf("database error")
					},
				}
			},
			setupGitHub: createMockGitHubClient,
			wantErr:     true,
			wantErrMsg:  "failed to get all cards",
		},
		{
			name:     "GitHubClientエラーが発生した場合",
			githubID: "12345",
			setupRepo: func() *repository.MockCardRepository {
				return &repository.MockCardRepository{
					FindAllFunc: func(githubID string) ([]domain.Card, error) {
						return []domain.Card{*createTestCard("12345")}, nil
					},
				}
			},
			setupGitHub: func() *github.MockClient {
				return &github.MockClient{
					GetUserByIDFunc: func(ctx context.Context, id int64) (*github.UserInfo, error) {
						return nil, fmt.Errorf("github api error")
					},
				}
			},
			wantErr:    true,
			wantErrMsg: "failed to get github user info",
		},
		{
			name:     "無効なGitHubIDの場合",
			githubID: "12345",
			setupRepo: func() *repository.MockCardRepository {
				return &repository.MockCardRepository{
					FindAllFunc: func(githubID string) ([]domain.Card, error) {
						card := createTestCard("invalid_id")
						return []domain.Card{*card}, nil
					},
				}
			},
			setupGitHub: createMockGitHubClient,
			wantErr:     true,
			wantErrMsg:  "invalid github id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			cardRepo := tt.setupRepo()
			identiconGen := &identicon.MockIdenticonGenerator{}
			githubClient := tt.setupGitHub()

			service := NewCardService(cardRepo, identiconGen)
			cards, err := service.GetAllCards(ctx, tt.githubID, githubClient)

			if tt.wantErr {
				if err == nil {
					t.Errorf("エラーが期待されましたが、エラーが発生しませんでした")
					return
				}
				if tt.wantErrMsg != "" && !contains(err.Error(), tt.wantErrMsg) {
					t.Errorf("エラーメッセージが期待と異なります: 期待=%s, 実際=%s", tt.wantErrMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラーが発生しました: %v", err)
					return
				}
				if len(cards) != tt.wantCardCount {
					t.Errorf("カード数が期待と異なります: 期待=%d, 実際=%d", tt.wantCardCount, len(cards))
				}
			}
		})
	}
}

// GetCardByGitHubID は指定されたGitHub IDのカードを取得する
func TestGetCardByGitHubID(t *testing.T) {
	tests := []struct {
		name        string
		githubID    string
		setupRepo   func() *repository.MockCardRepository
		setupGitHub func() *github.MockClient
		wantErr     bool
		wantErrMsg  string
	}{
		{
			name:     "正常にカードを取得できる",
			githubID: "12345",
			setupRepo: func() *repository.MockCardRepository {
				return &repository.MockCardRepository{
					FindByGitHubIDFunc: func(githubID string) (*domain.Card, error) {
						return createTestCard(githubID), nil
					},
				}
			},
			setupGitHub: createMockGitHubClient,
			wantErr:     false,
		},
		{
			name:     "カードが見つからない場合",
			githubID: "12345",
			setupRepo: func() *repository.MockCardRepository {
				return &repository.MockCardRepository{
					FindByGitHubIDFunc: func(githubID string) (*domain.Card, error) {
						return nil, nil
					},
				}
			},
			setupGitHub: createMockGitHubClient,
			wantErr:     true,
			wantErrMsg:  "card not found",
		},
		{
			name:     "Repositoryエラーが発生した場合",
			githubID: "12345",
			setupRepo: func() *repository.MockCardRepository {
				return &repository.MockCardRepository{
					FindByGitHubIDFunc: func(githubID string) (*domain.Card, error) {
						return nil, fmt.Errorf("database error")
					},
				}
			},
			setupGitHub: createMockGitHubClient,
			wantErr:     true,
			wantErrMsg:  "failed to get card by github id",
		},
		{
			name:     "GitHubClientエラーが発生した場合",
			githubID: "12345",
			setupRepo: func() *repository.MockCardRepository {
				return &repository.MockCardRepository{
					FindByGitHubIDFunc: func(githubID string) (*domain.Card, error) {
						return createTestCard(githubID), nil
					},
				}
			},
			setupGitHub: func() *github.MockClient {
				return &github.MockClient{
					GetUserByIDFunc: func(ctx context.Context, id int64) (*github.UserInfo, error) {
						return nil, fmt.Errorf("github api error")
					},
				}
			},
			wantErr:    true,
			wantErrMsg: "failed to get github user info",
		},
		{
			name:     "無効なGitHubIDの場合",
			githubID: "12345",
			setupRepo: func() *repository.MockCardRepository {
				return &repository.MockCardRepository{
					FindByGitHubIDFunc: func(githubID string) (*domain.Card, error) {
						card := createTestCard("invalid_id")
						return card, nil
					},
				}
			},
			setupGitHub: createMockGitHubClient,
			wantErr:     true,
			wantErrMsg:  "invalid github id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			cardRepo := tt.setupRepo()
			identiconGen := &identicon.MockIdenticonGenerator{}
			githubClient := tt.setupGitHub()

			service := NewCardService(cardRepo, identiconGen)
			card, err := service.GetCardByGitHubID(ctx, tt.githubID, githubClient)

			if tt.wantErr {
				if err == nil {
					t.Errorf("エラーが期待されましたが、エラーが発生しませんでした")
					return
				}
				if tt.wantErrMsg != "" && !contains(err.Error(), tt.wantErrMsg) {
					t.Errorf("エラーメッセージが期待と異なります: 期待=%s, 実際=%s", tt.wantErrMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラーが発生しました: %v", err)
					return
				}
				if card == nil {
					t.Errorf("カードがnilです")
					return
				}
				if card.GithubID != tt.githubID {
					t.Errorf("GitHubIDが期待と異なります: 期待=%s, 実際=%s", tt.githubID, card.GithubID)
				}
			}
		})
	}
}

// GetMyCard は自分のカードを取得する
func TestGetMyCard(t *testing.T) {
	tests := []struct {
		name        string
		githubID    string
		setupRepo   func() *repository.MockCardRepository
		setupGitHub func() *github.MockClient
		wantErr     bool
		wantErrMsg  string
	}{
		{
			name:     "正常に自分のカードを取得できる",
			githubID: "12345",
			setupRepo: func() *repository.MockCardRepository {
				return &repository.MockCardRepository{
					FindMyCardFunc: func(githubID string) (*domain.Card, error) {
						return createTestCard(githubID), nil
					},
				}
			},
			setupGitHub: createMockGitHubClient,
			wantErr:     false,
		},
		{
			name:     "カードが見つからない場合",
			githubID: "12345",
			setupRepo: func() *repository.MockCardRepository {
				return &repository.MockCardRepository{
					FindMyCardFunc: func(githubID string) (*domain.Card, error) {
						return nil, nil
					},
				}
			},
			setupGitHub: createMockGitHubClient,
			wantErr:     true,
			wantErrMsg:  "my card not found",
		},
		{
			name:     "Repositoryエラーが発生した場合",
			githubID: "12345",
			setupRepo: func() *repository.MockCardRepository {
				return &repository.MockCardRepository{
					FindMyCardFunc: func(githubID string) (*domain.Card, error) {
						return nil, fmt.Errorf("database error")
					},
				}
			},
			setupGitHub: createMockGitHubClient,
			wantErr:     true,
			wantErrMsg:  "failed to get my card",
		},
		{
			name:     "GitHubClientエラーが発生した場合",
			githubID: "12345",
			setupRepo: func() *repository.MockCardRepository {
				return &repository.MockCardRepository{
					FindMyCardFunc: func(githubID string) (*domain.Card, error) {
						return createTestCard(githubID), nil
					},
				}
			},
			setupGitHub: func() *github.MockClient {
				return &github.MockClient{
					GetUserByIDFunc: func(ctx context.Context, id int64) (*github.UserInfo, error) {
						return nil, fmt.Errorf("github api error")
					},
				}
			},
			wantErr:    true,
			wantErrMsg: "failed to get github user info",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			cardRepo := tt.setupRepo()
			identiconGen := &identicon.MockIdenticonGenerator{}
			githubClient := tt.setupGitHub()

			service := NewCardService(cardRepo, identiconGen)
			card, err := service.GetMyCard(ctx, tt.githubID, githubClient)

			if tt.wantErr {
				if err == nil {
					t.Errorf("エラーが期待されましたが、エラーが発生しませんでした")
					return
				}
				if tt.wantErrMsg != "" && !contains(err.Error(), tt.wantErrMsg) {
					t.Errorf("エラーメッセージが期待と異なります: 期待=%s, 実際=%s", tt.wantErrMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラーが発生しました: %v", err)
					return
				}
				if card == nil {
					t.Errorf("カードがnilです")
					return
				}
				if card.GithubID != tt.githubID {
					t.Errorf("GitHubIDが期待と異なります: 期待=%s, 実際=%s", tt.githubID, card.GithubID)
				}
			}
		})
	}
}

// GetOrCreateMyCard は自分のカードを取得し、存在しない場合は新規作成する
func TestGetOrCreateMyCard(t *testing.T) {
	tests := []struct {
		name           string
		githubID       string
		setupRepo      func() *repository.MockCardRepository
		setupIdenticon func() *identicon.MockIdenticonGenerator
		setupGitHub    func() *github.MockClient
		wantErr        bool
		wantErrMsg     string
		wantCreated    bool
	}{
		{
			name:     "既存のカードを取得できる",
			githubID: "12345",
			setupRepo: func() *repository.MockCardRepository {
				return &repository.MockCardRepository{
					FindMyCardFunc: func(githubID string) (*domain.Card, error) {
						return createTestCard(githubID), nil
					},
				}
			},
			setupIdenticon: func() *identicon.MockIdenticonGenerator {
				return &identicon.MockIdenticonGenerator{}
			},
			setupGitHub: createMockGitHubClient,
			wantErr:     false,
			wantCreated: false,
		},
		{
			name:     "カードが存在しない場合、新規作成する",
			githubID: "12345",
			setupRepo: func() *repository.MockCardRepository {
				return &repository.MockCardRepository{
					FindMyCardFunc: func(githubID string) (*domain.Card, error) {
						return nil, gorm.ErrRecordNotFound
					},
					CreateFunc: func(card *domain.Card) error {
						return nil
					},
				}
			},
			setupIdenticon: func() *identicon.MockIdenticonGenerator {
				return &identicon.MockIdenticonGenerator{
					GenerateFunc: func(githubID string) (domain.Color, domain.Blocks, error) {
						return domain.Color("#000000"), domain.Blocks{}, nil
					},
				}
			},
			setupGitHub: createMockGitHubClient,
			wantErr:     false,
			wantCreated: true,
		},
		{
			name:     "Repositoryエラーが発生した場合（RecordNotFound以外）",
			githubID: "12345",
			setupRepo: func() *repository.MockCardRepository {
				return &repository.MockCardRepository{
					FindMyCardFunc: func(githubID string) (*domain.Card, error) {
						return nil, fmt.Errorf("database error")
					},
				}
			},
			setupIdenticon: func() *identicon.MockIdenticonGenerator {
				return &identicon.MockIdenticonGenerator{}
			},
			setupGitHub: createMockGitHubClient,
			wantErr:     true,
			wantErrMsg:  "failed to get my card",
		},
		{
			name:     "IdenticonGeneratorエラーが発生した場合",
			githubID: "12345",
			setupRepo: func() *repository.MockCardRepository {
				return &repository.MockCardRepository{
					FindMyCardFunc: func(githubID string) (*domain.Card, error) {
						return nil, gorm.ErrRecordNotFound
					},
				}
			},
			setupIdenticon: func() *identicon.MockIdenticonGenerator {
				return &identicon.MockIdenticonGenerator{
					GenerateFunc: func(githubID string) (domain.Color, domain.Blocks, error) {
						return "", domain.Blocks{}, fmt.Errorf("identicon generation error")
					},
				}
			},
			setupGitHub: createMockGitHubClient,
			wantErr:     true,
			wantErrMsg:  "failed to generate identicon",
		},
		{
			name:     "カード作成時のRepositoryエラー",
			githubID: "12345",
			setupRepo: func() *repository.MockCardRepository {
				return &repository.MockCardRepository{
					FindMyCardFunc: func(githubID string) (*domain.Card, error) {
						return nil, gorm.ErrRecordNotFound
					},
					CreateFunc: func(card *domain.Card) error {
						return fmt.Errorf("create error")
					},
				}
			},
			setupIdenticon: func() *identicon.MockIdenticonGenerator {
				return &identicon.MockIdenticonGenerator{
					GenerateFunc: func(githubID string) (domain.Color, domain.Blocks, error) {
						return domain.Color("#000000"), domain.Blocks{}, nil
					},
				}
			},
			setupGitHub: createMockGitHubClient,
			wantErr:     true,
			wantErrMsg:  "failed to create card",
		},
		{
			name:     "GitHubClientエラーが発生した場合",
			githubID: "12345",
			setupRepo: func() *repository.MockCardRepository {
				return &repository.MockCardRepository{
					FindMyCardFunc: func(githubID string) (*domain.Card, error) {
						return createTestCard(githubID), nil
					},
				}
			},
			setupIdenticon: func() *identicon.MockIdenticonGenerator {
				return &identicon.MockIdenticonGenerator{}
			},
			setupGitHub: func() *github.MockClient {
				return &github.MockClient{
					GetUserByIDFunc: func(ctx context.Context, id int64) (*github.UserInfo, error) {
						return nil, fmt.Errorf("github api error")
					},
				}
			},
			wantErr:    true,
			wantErrMsg: "failed to get github user info",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			cardRepo := tt.setupRepo()
			identiconGen := tt.setupIdenticon()
			githubClient := tt.setupGitHub()

			service := NewCardService(cardRepo, identiconGen)
			card, err := service.GetOrCreateMyCard(ctx, tt.githubID, githubClient)

			if tt.wantErr {
				if err == nil {
					t.Errorf("エラーが期待されましたが、エラーが発生しませんでした")
					return
				}
				if tt.wantErrMsg != "" && !contains(err.Error(), tt.wantErrMsg) {
					t.Errorf("エラーメッセージが期待と異なります: 期待=%s, 実際=%s", tt.wantErrMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラーが発生しました: %v", err)
					return
				}
				if card == nil {
					t.Errorf("カードがnilです")
					return
				}
				if card.GithubID != tt.githubID {
					t.Errorf("GitHubIDが期待と異なります: 期待=%s, 実際=%s", tt.githubID, card.GithubID)
				}
			}
		})
	}
}

// AddCardToDeck はカードをデッキに追加する
func TestAddCardToDeck(t *testing.T) {
	tests := []struct {
		name              string
		collectorGithubID string
		targetGithubID    string
		setupRepo         func() *repository.MockCardRepository
		setupGitHub       func() *github.MockClient
		wantErr           bool
		wantErrMsg        string
	}{
		{
			name:              "正常にカードをデッキに追加できる",
			collectorGithubID: "11111",
			targetGithubID:    "12345",
			setupRepo: func() *repository.MockCardRepository {
				return &repository.MockCardRepository{
					FindByGitHubIDFunc: func(githubID string) (*domain.Card, error) {
						return createTestCard(githubID), nil
					},
					AddToCollectedCardsFunc: func(collectorGithubID string, cardID domain.CardID) error {
						return nil
					},
				}
			},
			setupGitHub: createMockGitHubClient,
			wantErr:     false,
		},
		{
			name:              "カードが見つからない場合",
			collectorGithubID: "11111",
			targetGithubID:    "12345",
			setupRepo: func() *repository.MockCardRepository {
				return &repository.MockCardRepository{
					FindByGitHubIDFunc: func(githubID string) (*domain.Card, error) {
						return nil, gorm.ErrRecordNotFound
					},
				}
			},
			setupGitHub: createMockGitHubClient,
			wantErr:     true,
			wantErrMsg:  "card not found",
		},
		{
			name:              "Repositoryエラーが発生した場合（RecordNotFound以外）",
			collectorGithubID: "11111",
			targetGithubID:    "12345",
			setupRepo: func() *repository.MockCardRepository {
				return &repository.MockCardRepository{
					FindByGitHubIDFunc: func(githubID string) (*domain.Card, error) {
						return nil, fmt.Errorf("database error")
					},
				}
			},
			setupGitHub: createMockGitHubClient,
			wantErr:     true,
			wantErrMsg:  "failed to find card",
		},
		{
			name:              "デッキ追加時のRepositoryエラー",
			collectorGithubID: "11111",
			targetGithubID:    "12345",
			setupRepo: func() *repository.MockCardRepository {
				return &repository.MockCardRepository{
					FindByGitHubIDFunc: func(githubID string) (*domain.Card, error) {
						return createTestCard(githubID), nil
					},
					AddToCollectedCardsFunc: func(collectorGithubID string, cardID domain.CardID) error {
						return fmt.Errorf("add error")
					},
				}
			},
			setupGitHub: createMockGitHubClient,
			wantErr:     true,
			wantErrMsg:  "failed to add card to deck",
		},
		{
			name:              "GitHubClientエラーが発生した場合",
			collectorGithubID: "11111",
			targetGithubID:    "12345",
			setupRepo: func() *repository.MockCardRepository {
				return &repository.MockCardRepository{
					FindByGitHubIDFunc: func(githubID string) (*domain.Card, error) {
						return createTestCard(githubID), nil
					},
					AddToCollectedCardsFunc: func(collectorGithubID string, cardID domain.CardID) error {
						return nil
					},
				}
			},
			setupGitHub: func() *github.MockClient {
				return &github.MockClient{
					GetUserByIDFunc: func(ctx context.Context, id int64) (*github.UserInfo, error) {
						return nil, fmt.Errorf("github api error")
					},
				}
			},
			wantErr:    true,
			wantErrMsg: "failed to get github user info",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			cardRepo := tt.setupRepo()
			identiconGen := &identicon.MockIdenticonGenerator{}
			githubClient := tt.setupGitHub()

			service := NewCardService(cardRepo, identiconGen)
			card, err := service.AddCardToDeck(ctx, tt.collectorGithubID, tt.targetGithubID, githubClient)

			if tt.wantErr {
				if err == nil {
					t.Errorf("エラーが期待されましたが、エラーが発生しませんでした")
					return
				}
				if tt.wantErrMsg != "" && !contains(err.Error(), tt.wantErrMsg) {
					t.Errorf("エラーメッセージが期待と異なります: 期待=%s, 実際=%s", tt.wantErrMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラーが発生しました: %v", err)
					return
				}
				if card == nil {
					t.Errorf("カードがnilです")
					return
				}
				if card.GithubID != tt.targetGithubID {
					t.Errorf("GitHubIDが期待と異なります: 期待=%s, 実際=%s", tt.targetGithubID, card.GithubID)
				}
			}
		})
	}
}

// RemoveCardFromDeck はカードをデッキから削除する
func TestRemoveCardFromDeck(t *testing.T) {
	tests := []struct {
		name              string
		collectorGithubID string
		targetGithubID    string
		setupRepo         func() *repository.MockCardRepository
		setupGitHub       func() *github.MockClient
		wantErr           bool
		wantErrMsg        string
	}{
		{
			name:              "正常にカードをデッキから削除できる",
			collectorGithubID: "11111",
			targetGithubID:    "12345",
			setupRepo: func() *repository.MockCardRepository {
				return &repository.MockCardRepository{
					FindByGitHubIDFunc: func(githubID string) (*domain.Card, error) {
						return createTestCard(githubID), nil
					},
					RemoveFromCollectedCardsFunc: func(collectorGithubID string, cardID domain.CardID) error {
						return nil
					},
				}
			},
			setupGitHub: createMockGitHubClient,
			wantErr:     false,
		},
		{
			name:              "カードが見つからない場合",
			collectorGithubID: "11111",
			targetGithubID:    "12345",
			setupRepo: func() *repository.MockCardRepository {
				return &repository.MockCardRepository{
					FindByGitHubIDFunc: func(githubID string) (*domain.Card, error) {
						return nil, gorm.ErrRecordNotFound
					},
				}
			},
			setupGitHub: createMockGitHubClient,
			wantErr:     true,
			wantErrMsg:  "card not found",
		},
		{
			name:              "Repositoryエラーが発生した場合（RecordNotFound以外）",
			collectorGithubID: "11111",
			targetGithubID:    "12345",
			setupRepo: func() *repository.MockCardRepository {
				return &repository.MockCardRepository{
					FindByGitHubIDFunc: func(githubID string) (*domain.Card, error) {
						return nil, fmt.Errorf("database error")
					},
				}
			},
			setupGitHub: createMockGitHubClient,
			wantErr:     true,
			wantErrMsg:  "failed to find card",
		},
		{
			name:              "デッキ削除時のRepositoryエラー",
			collectorGithubID: "11111",
			targetGithubID:    "12345",
			setupRepo: func() *repository.MockCardRepository {
				return &repository.MockCardRepository{
					FindByGitHubIDFunc: func(githubID string) (*domain.Card, error) {
						return createTestCard(githubID), nil
					},
					RemoveFromCollectedCardsFunc: func(collectorGithubID string, cardID domain.CardID) error {
						return fmt.Errorf("remove error")
					},
				}
			},
			setupGitHub: createMockGitHubClient,
			wantErr:     true,
			wantErrMsg:  "failed to remove card from deck",
		},
		{
			name:              "GitHubClientエラーが発生した場合",
			collectorGithubID: "11111",
			targetGithubID:    "12345",
			setupRepo: func() *repository.MockCardRepository {
				return &repository.MockCardRepository{
					FindByGitHubIDFunc: func(githubID string) (*domain.Card, error) {
						return createTestCard(githubID), nil
					},
					RemoveFromCollectedCardsFunc: func(collectorGithubID string, cardID domain.CardID) error {
						return nil
					},
				}
			},
			setupGitHub: func() *github.MockClient {
				return &github.MockClient{
					GetUserByIDFunc: func(ctx context.Context, id int64) (*github.UserInfo, error) {
						return nil, fmt.Errorf("github api error")
					},
				}
			},
			wantErr:    true,
			wantErrMsg: "failed to get github user info",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			cardRepo := tt.setupRepo()
			identiconGen := &identicon.MockIdenticonGenerator{}
			githubClient := tt.setupGitHub()

			service := NewCardService(cardRepo, identiconGen)
			card, err := service.RemoveCardFromDeck(ctx, tt.collectorGithubID, tt.targetGithubID, githubClient)

			if tt.wantErr {
				if err == nil {
					t.Errorf("エラーが期待されましたが、エラーが発生しませんでした")
					return
				}
				if tt.wantErrMsg != "" && !contains(err.Error(), tt.wantErrMsg) {
					t.Errorf("エラーメッセージが期待と異なります: 期待=%s, 実際=%s", tt.wantErrMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラーが発生しました: %v", err)
					return
				}
				if card == nil {
					t.Errorf("カードがnilです")
					return
				}
				if card.GithubID != tt.targetGithubID {
					t.Errorf("GitHubIDが期待と異なります: 期待=%s, 実際=%s", tt.targetGithubID, card.GithubID)
				}
			}
		})
	}
}

// ヘルパー関数: 文字列が部分文字列を含むかチェック
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

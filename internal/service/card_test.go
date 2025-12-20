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
		NodeID:   "U_" + githubID, // テスト用のNodeID
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
			wantErr:    true,
			wantErrMsg: "failed to get all cards",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cardRepo := tt.setupRepo()
			identiconGen := &identicon.MockIdenticonGenerator{}

			service := NewCardService(cardRepo, identiconGen)
			cards, err := service.GetAllCards(tt.githubID)

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
						return nil, gorm.ErrRecordNotFound
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
			setupGitHub: func() *github.MockClient {
				return &github.MockClient{
					GetAuthenticatedUserFunc: func(ctx context.Context) (*github.UserInfo, error) {
						return nil, fmt.Errorf("github api error")
					},
				}
			},
			wantErr:    true,
			wantErrMsg: "failed to get authenticated user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			cardRepo := tt.setupRepo()
			identiconGen := tt.setupIdenticon()
			githubClient := tt.setupGitHub()

			service := NewCardService(cardRepo, identiconGen)
			card, err := service.GetOrCreateMyCard(ctx, tt.githubID, "MDQ6VXNlcjEyMzQ1", githubClient)

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

// EnrichCardWithGitHubInfo はGitHub APIからユーザー情報を取得してCardに設定する
func TestEnrichCardWithGitHubInfo(t *testing.T) {
	tests := []struct {
		name        string
		card        *domain.Card
		setupGitHub func() *github.MockClient
		wantErr     bool
		wantErrMsg  string
		validate    func(t *testing.T, card *domain.Card)
	}{
		{
			name: "正常にGitHub情報を取得してCardに設定できる",
			card: &domain.Card{
				ID:       domain.NewCardID(),
				GithubID: "12345",
				NodeID:   "U_12345",
				Color:    domain.Color("#000000"),
				Blocks:   domain.Blocks{},
			},
			setupGitHub: func() *github.MockClient {
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
			},
			wantErr: false,
			validate: func(t *testing.T, card *domain.Card) {
				if card.UserName != "testuser" {
					t.Errorf("UserName = %v, want testuser", card.UserName)
				}
				if card.FullName != "Test User" {
					t.Errorf("FullName = %v, want Test User", card.FullName)
				}
				if card.IconUrl != "https://example.com/avatar.png" {
					t.Errorf("IconUrl = %v, want https://example.com/avatar.png", card.IconUrl)
				}
				if card.MostUsedLanguage.LanguageName != "Go" {
					t.Errorf("MostUsedLanguage.LanguageName = %v, want Go", card.MostUsedLanguage.LanguageName)
				}
				if card.MostUsedLanguage.Color != "#00ADD8" {
					t.Errorf("MostUsedLanguage.Color = %v, want #00ADD8", card.MostUsedLanguage.Color)
				}
			},
		},
		{
			name: "無効なGitHub IDの場合はエラー",
			card: &domain.Card{
				ID:       domain.NewCardID(),
				GithubID: "invalid_id",
				NodeID:   "U_invalid",
				Color:    domain.Color("#000000"),
				Blocks:   domain.Blocks{},
			},
			setupGitHub: func() *github.MockClient {
				return &github.MockClient{}
			},
			wantErr:    true,
			wantErrMsg: "invalid github id",
		},
		{
			name: "GitHub APIエラーの場合はエラー",
			card: &domain.Card{
				ID:       domain.NewCardID(),
				GithubID: "12345",
				NodeID:   "U_12345",
				Color:    domain.Color("#000000"),
				Blocks:   domain.Blocks{},
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
			name: "言語情報取得エラーの場合はエラー",
			card: &domain.Card{
				ID:       domain.NewCardID(),
				GithubID: "12345",
				NodeID:   "U_12345",
				Color:    domain.Color("#000000"),
				Blocks:   domain.Blocks{},
			},
			setupGitHub: func() *github.MockClient {
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
						return "", "", fmt.Errorf("language api error")
					},
				}
			},
			wantErr:    true,
			wantErrMsg: "failed to get most used language",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			githubClient := tt.setupGitHub()
			err := EnrichCardWithGitHubInfo(ctx, tt.card, githubClient)

			if tt.wantErr {
				if err == nil {
					t.Errorf("EnrichCardWithGitHubInfo() error = nil, want error")
					return
				}
				if tt.wantErrMsg != "" && !contains(err.Error(), tt.wantErrMsg) {
					t.Errorf("EnrichCardWithGitHubInfo() error = %v, want error containing %v", err.Error(), tt.wantErrMsg)
				}
			} else {
				if err != nil {
					t.Errorf("EnrichCardWithGitHubInfo() error = %v, want nil", err)
					return
				}
				if tt.validate != nil {
					tt.validate(t, tt.card)
				}
			}
		})
	}
}

// EnrichCardsWithGitHubInfo は複数のカードにGitHub情報を一括で設定する（バッチ処理版）
// N+1問題を解消し、並列処理でパフォーマンスを向上させる
func TestEnrichCardsWithGitHubInfo(t *testing.T) {
	tests := []struct {
		name        string
		cards       []domain.Card
		setupGitHub func() *github.MockClient
		wantErr     bool
		wantErrMsg  string
		validate    func(t *testing.T, cards []domain.Card)
	}{
		{
			name: "正常に複数カードの一括処理が動作する",
			cards: []domain.Card{
				{
					ID:       domain.NewCardID(),
					GithubID: "12345",
					NodeID:   "U_12345",
					Color:    domain.Color("#000000"),
					Blocks:   domain.Blocks{},
				},
				{
					ID:       domain.NewCardID(),
					GithubID: "67890",
					NodeID:   "U_67890",
					Color:    domain.Color("#000000"),
					Blocks:   domain.Blocks{},
				},
			},
			setupGitHub: func() *github.MockClient {
				return &github.MockClient{
					GetUsersByIDsFunc: func(ctx context.Context, ids []int64) (map[int64]*github.UserInfo, error) {
						result := make(map[int64]*github.UserInfo)
						for _, id := range ids {
							result[id] = &github.UserInfo{
								ID:        id,
								Login:     fmt.Sprintf("user%d", id),
								Name:      fmt.Sprintf("User %d", id),
								AvatarURL: fmt.Sprintf("https://example.com/user%d.png", id),
							}
						}
						return result, nil
					},
					GetMostUsedLanguagesFunc: func(ctx context.Context, logins []string) (map[string]github.LanguageInfo, error) {
						result := make(map[string]github.LanguageInfo)
						for _, login := range logins {
							result[login] = github.LanguageInfo{
								Name:  "Go",
								Color: "#00ADD8",
							}
						}
						return result, nil
					},
				}
			},
			wantErr: false,
			validate: func(t *testing.T, cards []domain.Card) {
				if len(cards) != 2 {
					t.Errorf("cards length = %v, want 2", len(cards))
				}
				if cards[0].UserName == "" {
					t.Errorf("cards[0].UserName is empty")
				}
				if cards[1].UserName == "" {
					t.Errorf("cards[1].UserName is empty")
				}
			},
		},
		{
			name:  "空のカードリストの場合は早期リターン",
			cards: []domain.Card{},
			setupGitHub: func() *github.MockClient {
				return &github.MockClient{}
			},
			wantErr: false,
			validate: func(t *testing.T, cards []domain.Card) {
				if len(cards) != 0 {
					t.Errorf("cards length = %v, want 0", len(cards))
				}
			},
		},
		{
			name: "無効なGitHub IDが含まれる場合はエラー",
			cards: []domain.Card{
				{
					ID:       domain.NewCardID(),
					GithubID: "invalid_id",
					NodeID:   "U_invalid",
					Color:    domain.Color("#000000"),
					Blocks:   domain.Blocks{},
				},
			},
			setupGitHub: func() *github.MockClient {
				return &github.MockClient{}
			},
			wantErr:    true,
			wantErrMsg: "invalid github id for card",
		},
		{
			name: "GitHub APIエラーの場合はエラー",
			cards: []domain.Card{
				{
					ID:       domain.NewCardID(),
					GithubID: "12345",
					NodeID:   "U_12345",
					Color:    domain.Color("#000000"),
					Blocks:   domain.Blocks{},
				},
			},
			setupGitHub: func() *github.MockClient {
				return &github.MockClient{
					GetUsersByIDsFunc: func(ctx context.Context, ids []int64) (map[int64]*github.UserInfo, error) {
						return nil, fmt.Errorf("github api error")
					},
				}
			},
			wantErr:    true,
			wantErrMsg: "failed to get users info",
		},
		{
			name: "言語情報取得エラーの場合はエラー",
			cards: []domain.Card{
				{
					ID:       domain.NewCardID(),
					GithubID: "12345",
					NodeID:   "U_12345",
					Color:    domain.Color("#000000"),
					Blocks:   domain.Blocks{},
				},
			},
			setupGitHub: func() *github.MockClient {
				return &github.MockClient{
					GetUsersByIDsFunc: func(ctx context.Context, ids []int64) (map[int64]*github.UserInfo, error) {
						return map[int64]*github.UserInfo{
							12345: {
								ID:        12345,
								Login:     "testuser",
								Name:      "Test User",
								AvatarURL: "https://example.com/avatar.png",
							},
						}, nil
					},
					GetMostUsedLanguagesFunc: func(ctx context.Context, logins []string) (map[string]github.LanguageInfo, error) {
						return nil, fmt.Errorf("language api error")
					},
				}
			},
			wantErr:    true,
			wantErrMsg: "failed to get languages info",
		},
		{
			name: "同じGitHub IDを持つカードが複数ある場合も正常に処理できる",
			cards: []domain.Card{
				{
					ID:       domain.NewCardID(),
					GithubID: "12345",
					NodeID:   "U_12345",
					Color:    domain.Color("#000000"),
					Blocks:   domain.Blocks{},
				},
				{
					ID:       domain.NewCardID(),
					GithubID: "12345",
					NodeID:   "U_12345",
					Color:    domain.Color("#000000"),
					Blocks:   domain.Blocks{},
				},
			},
			setupGitHub: func() *github.MockClient {
				return &github.MockClient{
					GetUsersByIDsFunc: func(ctx context.Context, ids []int64) (map[int64]*github.UserInfo, error) {
						return map[int64]*github.UserInfo{
							12345: {
								ID:        12345,
								Login:     "testuser",
								Name:      "Test User",
								AvatarURL: "https://example.com/avatar.png",
							},
						}, nil
					},
					GetMostUsedLanguagesFunc: func(ctx context.Context, logins []string) (map[string]github.LanguageInfo, error) {
						return map[string]github.LanguageInfo{
							"testuser": {
								Name:  "Go",
								Color: "#00ADD8",
							},
						}, nil
					},
				}
			},
			wantErr: false,
			validate: func(t *testing.T, cards []domain.Card) {
				if len(cards) != 2 {
					t.Errorf("cards length = %v, want 2", len(cards))
				}
				if cards[0].UserName != "testuser" {
					t.Errorf("cards[0].UserName = %v, want testuser", cards[0].UserName)
				}
				if cards[1].UserName != "testuser" {
					t.Errorf("cards[1].UserName = %v, want testuser", cards[1].UserName)
				}
			},
		},
		{
			name: "ユーザー情報が見つからない場合はスキップされる",
			cards: []domain.Card{
				{
					ID:       domain.NewCardID(),
					GithubID: "12345",
					NodeID:   "U_12345",
					Color:    domain.Color("#000000"),
					Blocks:   domain.Blocks{},
				},
			},
			setupGitHub: func() *github.MockClient {
				return &github.MockClient{
					GetUsersByIDsFunc: func(ctx context.Context, ids []int64) (map[int64]*github.UserInfo, error) {
						// 空のマップを返す（ユーザー情報が見つからない）
						return map[int64]*github.UserInfo{}, nil
					},
					GetMostUsedLanguagesFunc: func(ctx context.Context, logins []string) (map[string]github.LanguageInfo, error) {
						return map[string]github.LanguageInfo{}, nil
					},
				}
			},
			wantErr: false,
			validate: func(t *testing.T, cards []domain.Card) {
				// ユーザー情報が見つからない場合は、カードの情報は更新されない
				if cards[0].UserName != "" {
					t.Errorf("cards[0].UserName = %v, want empty", cards[0].UserName)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			githubClient := tt.setupGitHub()
			err := EnrichCardsWithGitHubInfo(ctx, tt.cards, githubClient)

			if tt.wantErr {
				if err == nil {
					t.Errorf("EnrichCardsWithGitHubInfo() error = nil, want error")
					return
				}
				if tt.wantErrMsg != "" && !contains(err.Error(), tt.wantErrMsg) {
					t.Errorf("EnrichCardsWithGitHubInfo() error = %v, want error containing %v", err.Error(), tt.wantErrMsg)
				}
			} else {
				if err != nil {
					t.Errorf("EnrichCardsWithGitHubInfo() error = %v, want nil", err)
					return
				}
				if tt.validate != nil {
					tt.validate(t, tt.cards)
				}
			}
		})
	}
}

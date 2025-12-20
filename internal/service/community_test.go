package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/furarico/octo-deck-api/internal/domain"
	"github.com/furarico/octo-deck-api/internal/github"
	"github.com/furarico/octo-deck-api/internal/repository"
)

// テスト用のヘルパー関数: 正常なコミュニティを返す
func createTestCommunity(name string) *domain.Community {
	now := time.Now()
	return &domain.Community{
		ID:        domain.NewCommunityID(),
		Name:      name,
		StartedAt: now.AddDate(0, 0, -7),
		EndedAt:   now,
	}
}

// GetAllCommunities はすべてのコミュニティを取得する
func TestGetAllCommunities(t *testing.T) {
	tests := []struct {
		name       string
		githubID   string
		setupRepo  func() *repository.MockCommunityRepository
		wantErr    bool
		wantErrMsg string
		wantCount  int
	}{
		{
			name:     "正常にコミュニティ一覧を取得できる",
			githubID: "12345",
			setupRepo: func() *repository.MockCommunityRepository {
				return &repository.MockCommunityRepository{
					FindAllFunc: func(githubID string) ([]domain.Community, error) {
						return []domain.Community{
							*createTestCommunity("Community 1"),
							*createTestCommunity("Community 2"),
						}, nil
					},
				}
			},
			wantErr:   false,
			wantCount: 2,
		},
		{
			name:     "コミュニティが存在しない場合",
			githubID: "12345",
			setupRepo: func() *repository.MockCommunityRepository {
				return &repository.MockCommunityRepository{
					FindAllFunc: func(githubID string) ([]domain.Community, error) {
						return []domain.Community{}, nil
					},
				}
			},
			wantErr:   false,
			wantCount: 0,
		},
		{
			name:     "Repositoryエラーが発生した場合",
			githubID: "12345",
			setupRepo: func() *repository.MockCommunityRepository {
				return &repository.MockCommunityRepository{
					FindAllFunc: func(githubID string) ([]domain.Community, error) {
						return nil, fmt.Errorf("database error")
					},
				}
			},
			wantErr:    true,
			wantErrMsg: "failed to get all communities",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			communityRepo := tt.setupRepo()
			cardRepo := &repository.MockCardRepository{}
			service := NewCommunityService(communityRepo, cardRepo)
			communities, err := service.GetAllCommunities(tt.githubID)

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
				if len(communities) != tt.wantCount {
					t.Errorf("コミュニティ数が期待と異なります: 期待=%d, 実際=%d", tt.wantCount, len(communities))
				}
			}
		})
	}
}

// GetCommunityByID は指定されたコミュニティIDの情報を取得する
func TestGetCommunityByID(t *testing.T) {
	tests := []struct {
		name        string
		communityID string
		setupRepo   func() *repository.MockCommunityRepository
		wantErr     bool
		wantErrMsg  string
	}{
		{
			name:        "正常にコミュニティを取得できる",
			communityID: "test-community-id",
			setupRepo: func() *repository.MockCommunityRepository {
				return &repository.MockCommunityRepository{
					FindByIDFunc: func(id string) (*domain.Community, error) {
						return createTestCommunity("Test Community"), nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:        "コミュニティが見つからない場合",
			communityID: "test-community-id",
			setupRepo: func() *repository.MockCommunityRepository {
				return &repository.MockCommunityRepository{
					FindByIDFunc: func(id string) (*domain.Community, error) {
						return nil, nil
					},
				}
			},
			wantErr:    true,
			wantErrMsg: "community not found",
		},
		{
			name:        "Repositoryエラーが発生した場合",
			communityID: "test-community-id",
			setupRepo: func() *repository.MockCommunityRepository {
				return &repository.MockCommunityRepository{
					FindByIDFunc: func(id string) (*domain.Community, error) {
						return nil, fmt.Errorf("database error")
					},
				}
			},
			wantErr:    true,
			wantErrMsg: "failed to get community by id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			communityRepo := tt.setupRepo()
			cardRepo := &repository.MockCardRepository{}
			service := NewCommunityService(communityRepo, cardRepo)
			community, err := service.GetCommunityByID(tt.communityID)

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
				if community == nil {
					t.Errorf("コミュニティがnilです")
					return
				}
			}
		})
	}
}

// GetCommunityWithHighlightedCard はコミュニティとHighlightedCardをデータベースから取得する
func TestGetCommunityWithHighlightedCard(t *testing.T) {
	tests := []struct {
		name        string
		communityID string
		setupRepo   func() *repository.MockCommunityRepository
		wantErr     bool
		wantErrMsg  string
	}{
		{
			name:        "正常にコミュニティとHighlightedCardを取得できる",
			communityID: "test-community-id",
			setupRepo: func() *repository.MockCommunityRepository {
				return &repository.MockCommunityRepository{
					FindByIDWithHighlightedCardFunc: func(id string) (*domain.Community, error) {
						community := createTestCommunity("Test Community")
						community.HighlightedCard = domain.HighlightedCard{
							BestContributor: *createTestCard("12345"),
						}
						return community, nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:        "コミュニティが見つからない場合",
			communityID: "test-community-id",
			setupRepo: func() *repository.MockCommunityRepository {
				return &repository.MockCommunityRepository{
					FindByIDWithHighlightedCardFunc: func(id string) (*domain.Community, error) {
						return nil, nil
					},
				}
			},
			wantErr:    true,
			wantErrMsg: "community not found",
		},
		{
			name:        "Repositoryエラーが発生した場合",
			communityID: "test-community-id",
			setupRepo: func() *repository.MockCommunityRepository {
				return &repository.MockCommunityRepository{
					FindByIDWithHighlightedCardFunc: func(id string) (*domain.Community, error) {
						return nil, fmt.Errorf("database error")
					},
				}
			},
			wantErr:    true,
			wantErrMsg: "failed to get community by id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			communityRepo := tt.setupRepo()
			cardRepo := &repository.MockCardRepository{}
			service := NewCommunityService(communityRepo, cardRepo)
			community, highlightedCard, err := service.GetCommunityWithHighlightedCard(tt.communityID)

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
				if community == nil {
					t.Errorf("コミュニティがnilです")
					return
				}
				if highlightedCard == nil {
					t.Errorf("HighlightedCardがnilです")
					return
				}
			}
		})
	}
}

// GetCommunityCards は指定したコミュニティIDのカード一覧を取得する
func TestGetCommunityCards(t *testing.T) {
	tests := []struct {
		name        string
		communityID string
		setupRepo   func() *repository.MockCommunityRepository
		wantErr     bool
		wantErrMsg  string
		wantCount   int
	}{
		{
			name:        "正常にカード一覧を取得できる",
			communityID: "test-community-id",
			setupRepo: func() *repository.MockCommunityRepository {
				return &repository.MockCommunityRepository{
					FindCardsFunc: func(id string) ([]domain.Card, error) {
						return []domain.Card{
							*createTestCard("12345"),
							*createTestCard("67890"),
						}, nil
					},
				}
			},
			wantErr:   false,
			wantCount: 2,
		},
		{
			name:        "カードが存在しない場合",
			communityID: "test-community-id",
			setupRepo: func() *repository.MockCommunityRepository {
				return &repository.MockCommunityRepository{
					FindCardsFunc: func(id string) ([]domain.Card, error) {
						return []domain.Card{}, nil
					},
				}
			},
			wantErr:   false,
			wantCount: 0,
		},
		{
			name:        "Repositoryエラーが発生した場合",
			communityID: "test-community-id",
			setupRepo: func() *repository.MockCommunityRepository {
				return &repository.MockCommunityRepository{
					FindCardsFunc: func(id string) ([]domain.Card, error) {
						return nil, fmt.Errorf("database error")
					},
				}
			},
			wantErr:    true,
			wantErrMsg: "failed to get community cards",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			communityRepo := tt.setupRepo()
			mockGitHub := &github.MockClient{
				GetUserByIDFunc: func(ctx context.Context, id int64) (*github.UserInfo, error) {
					return &github.UserInfo{
						ID:        id,
						Login:     "user",
						Name:      "User",
						AvatarURL: "https://example.com/avatar.png",
					}, nil
				},
				GetMostUsedLanguageFunc: func(ctx context.Context, login string) (string, string, error) {
					return "Go", "#00ADD8", nil
				},
			}
			cardRepo := &repository.MockCardRepository{}
			service := NewCommunityService(communityRepo, cardRepo)
			cards, err := service.GetCommunityCards(ctx, tt.communityID, mockGitHub)

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
				if len(cards) != tt.wantCount {
					t.Errorf("カード数が期待と異なります: 期待=%d, 実際=%d", tt.wantCount, len(cards))
				}
				// fullname / avatar url が補完されていることの簡易チェック
				for _, c := range cards {
					if c.FullName == "" || c.IconUrl == "" {
						t.Errorf("カードの FullName または IconUrl が補完されていません")
					}
				}
			}
		})
	}
}

// CreateCommunityWithPeriod はコミュニティを作成する
func TestCreateCommunityWithPeriod(t *testing.T) {
	now := time.Now()
	startDateTime := now.AddDate(0, 0, -7)
	endDateTime := now

	tests := []struct {
		name          string
		communityName string
		startDateTime time.Time
		endDateTime   time.Time
		setupRepo     func() *repository.MockCommunityRepository
		wantErr       bool
		wantErrMsg    string
	}{
		{
			name:          "正常にコミュニティを作成できる",
			communityName: "Test Community",
			startDateTime: startDateTime,
			endDateTime:   endDateTime,
			setupRepo: func() *repository.MockCommunityRepository {
				return &repository.MockCommunityRepository{
					CreateFunc: func(community *domain.Community) error {
						return nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:          "Repositoryエラーが発生した場合",
			communityName: "Test Community",
			startDateTime: startDateTime,
			endDateTime:   endDateTime,
			setupRepo: func() *repository.MockCommunityRepository {
				return &repository.MockCommunityRepository{
					CreateFunc: func(community *domain.Community) error {
						return fmt.Errorf("database error")
					},
				}
			},
			wantErr:    true,
			wantErrMsg: "failed to create community",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			communityRepo := tt.setupRepo()
			cardRepo := &repository.MockCardRepository{}
			service := NewCommunityService(communityRepo, cardRepo)
			community, err := service.CreateCommunityWithPeriod(tt.communityName, tt.startDateTime, tt.endDateTime)

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
				if community == nil {
					t.Errorf("コミュニティがnilです")
					return
				}
				if community.Name != tt.communityName {
					t.Errorf("コミュニティ名が期待と異なります: 期待=%s, 実際=%s", tt.communityName, community.Name)
				}
			}
		})
	}
}

// DeleteCommunity はコミュニティを削除する
func TestDeleteCommunity(t *testing.T) {
	tests := []struct {
		name        string
		communityID string
		setupRepo   func() *repository.MockCommunityRepository
		wantErr     bool
		wantErrMsg  string
	}{
		{
			name:        "正常にコミュニティを削除できる",
			communityID: "test-community-id",
			setupRepo: func() *repository.MockCommunityRepository {
				return &repository.MockCommunityRepository{
					DeleteFunc: func(id string) error {
						return nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:        "Repositoryエラーが発生した場合",
			communityID: "test-community-id",
			setupRepo: func() *repository.MockCommunityRepository {
				return &repository.MockCommunityRepository{
					DeleteFunc: func(id string) error {
						return fmt.Errorf("database error")
					},
				}
			},
			wantErr:    true,
			wantErrMsg: "failed to delete community",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			communityRepo := tt.setupRepo()
			cardRepo := &repository.MockCardRepository{}
			service := NewCommunityService(communityRepo, cardRepo)
			err := service.DeleteCommunity(tt.communityID)

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
			}
		})
	}
}

// AddCardToCommunity はコミュニティにカードを追加する
func TestAddCardToCommunity(t *testing.T) {
	tests := []struct {
		name        string
		communityID string
		cardID      string
		setupRepo   func() *repository.MockCommunityRepository
		wantErr     bool
		wantErrMsg  string
	}{
		{
			name:        "正常にカードを追加できる",
			communityID: "test-community-id",
			cardID:      "test-card-id",
			setupRepo: func() *repository.MockCommunityRepository {
				return &repository.MockCommunityRepository{
					AddCardFunc: func(communityID string, cardID string) error {
						return nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:        "Repositoryエラーが発生した場合",
			communityID: "test-community-id",
			cardID:      "test-card-id",
			setupRepo: func() *repository.MockCommunityRepository {
				return &repository.MockCommunityRepository{
					AddCardFunc: func(communityID string, cardID string) error {
						return fmt.Errorf("database error")
					},
				}
			},
			wantErr:    true,
			wantErrMsg: "failed to add card to community",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			communityRepo := tt.setupRepo()
			cardRepo := &repository.MockCardRepository{}
			service := NewCommunityService(communityRepo, cardRepo)
			err := service.AddCardToCommunity(tt.communityID, tt.cardID)

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
			}
		})
	}
}

// RemoveCardFromCommunity はコミュニティからカードを削除する
func TestRemoveCardFromCommunity(t *testing.T) {
	tests := []struct {
		name        string
		communityID string
		cardID      string
		setupRepo   func() *repository.MockCommunityRepository
		wantErr     bool
		wantErrMsg  string
	}{
		{
			name:        "正常にカードを削除できる",
			communityID: "test-community-id",
			cardID:      "test-card-id",
			setupRepo: func() *repository.MockCommunityRepository {
				return &repository.MockCommunityRepository{
					RemoveCardFunc: func(communityID string, cardID string) error {
						return nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:        "Repositoryエラーが発生した場合",
			communityID: "test-community-id",
			cardID:      "test-card-id",
			setupRepo: func() *repository.MockCommunityRepository {
				return &repository.MockCommunityRepository{
					RemoveCardFunc: func(communityID string, cardID string) error {
						return fmt.Errorf("database error")
					},
				}
			},
			wantErr:    true,
			wantErrMsg: "failed to remove card from community",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			communityRepo := tt.setupRepo()
			cardRepo := &repository.MockCardRepository{}
			service := NewCommunityService(communityRepo, cardRepo)
			err := service.RemoveCardFromCommunity(tt.communityID, tt.cardID)

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
			}
		})
	}
}

// テスト用のヘルパー関数: UserFullInfoを作成する
func createTestUserFullInfo(login string, total, commits, issues, prs, reviews int) github.UserFullInfo {
	return github.UserFullInfo{
		Login:                 login,
		Name:                  "Test " + login,
		AvatarURL:             "https://example.com/" + login + ".png",
		Total:                 total,
		Commits:               commits,
		Issues:                issues,
		PRs:                   prs,
		Reviews:               reviews,
		MostUsedLanguage:      "Go",
		MostUsedLanguageColor: "#00ADD8",
	}
}

// RefreshHighlightedCard はGitHub APIを呼び出してHighlightedCardを再計算し、データベースに保存する
func TestRefreshHighlightedCard(t *testing.T) {
	now := time.Now()
	startDateTime := now.AddDate(0, 0, -7)
	endDateTime := now

	tests := []struct {
		name          string
		communityID   string
		setupRepo     func() *repository.MockCommunityRepository
		setupCardRepo func() *repository.MockCardRepository
		setupGitHub   func() *github.MockClient
		wantErr       bool
		wantErrMsg    string
		validate      func(t *testing.T, community *domain.Community, highlightedCard *domain.HighlightedCard)
	}{
		{
			name:        "正常にHighlightedCardが更新される",
			communityID: "test-community-id",
			setupRepo: func() *repository.MockCommunityRepository {
				community := createTestCommunity("Test Community")
				community.StartedAt = startDateTime
				community.EndedAt = endDateTime
				var savedHighlightedCard *domain.HighlightedCard
				return &repository.MockCommunityRepository{
					FindByIDFunc: func(id string) (*domain.Community, error) {
						return community, nil
					},
					FindCardsFunc: func(id string) ([]domain.Card, error) {
						return []domain.Card{
							*createTestCard("12345"),
							*createTestCard("67890"),
							*createTestCard("11111"),
						}, nil
					},
					UpdateHighlightedCardFunc: func(communityID string, highlightedCard *domain.HighlightedCard) error {
						savedHighlightedCard = highlightedCard
						return nil
					},
					FindByIDWithHighlightedCardFunc: func(id string) (*domain.Community, error) {
						updatedCommunity := createTestCommunity("Test Community")
						updatedCommunity.StartedAt = startDateTime
						updatedCommunity.EndedAt = endDateTime
						if savedHighlightedCard != nil {
							updatedCommunity.HighlightedCard = *savedHighlightedCard
						}
						return updatedCommunity, nil
					},
				}
			},
			setupCardRepo: func() *repository.MockCardRepository {
				return &repository.MockCardRepository{
					UpdateFunc: func(card *domain.Card) error {
						return nil
					},
				}
			},
			setupGitHub: func() *github.MockClient {
				return &github.MockClient{
					GetUsersFullInfoByNodeIDsFunc: func(ctx context.Context, nodeIDs []string, from, to time.Time) ([]github.UserFullInfo, error) {
						// 各カテゴリで異なるユーザーがベストになるように設定
						return []github.UserFullInfo{
							createTestUserFullInfo("user1", 100, 50, 10, 20, 20), // BestContributor (Total=100)
							createTestUserFullInfo("user2", 80, 60, 5, 10, 5),    // BestCommitter (Commits=60)
							createTestUserFullInfo("user3", 70, 30, 30, 5, 5),    // BestIssuer (Issues=30)
						}, nil
					},
				}
			},
			wantErr: false,
			validate: func(t *testing.T, community *domain.Community, highlightedCard *domain.HighlightedCard) {
				if community == nil {
					t.Errorf("コミュニティがnilです")
					return
				}
				if highlightedCard == nil {
					t.Errorf("HighlightedCardがnilです")
					return
				}
				// BestContributorが正しく設定されているか確認
				if highlightedCard.BestContributor.GithubID == "" {
					t.Errorf("BestContributorが設定されていません")
				}
			},
		},
		{
			name:        "コミュニティが見つからない場合",
			communityID: "test-community-id",
			setupRepo: func() *repository.MockCommunityRepository {
				return &repository.MockCommunityRepository{
					FindByIDFunc: func(id string) (*domain.Community, error) {
						return nil, nil
					},
				}
			},
			setupCardRepo: func() *repository.MockCardRepository {
				return &repository.MockCardRepository{}
			},
			setupGitHub: func() *github.MockClient {
				return &github.MockClient{}
			},
			wantErr:    true,
			wantErrMsg: "community not found",
		},
		{
			name:        "コミュニティ取得時のRepositoryエラー",
			communityID: "test-community-id",
			setupRepo: func() *repository.MockCommunityRepository {
				return &repository.MockCommunityRepository{
					FindByIDFunc: func(id string) (*domain.Community, error) {
						return nil, fmt.Errorf("database error")
					},
				}
			},
			setupCardRepo: func() *repository.MockCardRepository {
				return &repository.MockCardRepository{}
			},
			setupGitHub: func() *github.MockClient {
				return &github.MockClient{}
			},
			wantErr:    true,
			wantErrMsg: "failed to get community by id",
		},
		{
			name:        "カード一覧取得時のRepositoryエラー",
			communityID: "test-community-id",
			setupRepo: func() *repository.MockCommunityRepository {
				community := createTestCommunity("Test Community")
				community.StartedAt = startDateTime
				community.EndedAt = endDateTime
				return &repository.MockCommunityRepository{
					FindByIDFunc: func(id string) (*domain.Community, error) {
						return community, nil
					},
					FindCardsFunc: func(id string) ([]domain.Card, error) {
						return nil, fmt.Errorf("database error")
					},
				}
			},
			setupCardRepo: func() *repository.MockCardRepository {
				return &repository.MockCardRepository{}
			},
			setupGitHub: func() *github.MockClient {
				return &github.MockClient{}
			},
			wantErr:    true,
			wantErrMsg: "failed to get community cards",
		},
		{
			name:        "カードが0件の場合（空のHighlightedCardを保存）",
			communityID: "test-community-id",
			setupRepo: func() *repository.MockCommunityRepository {
				community := createTestCommunity("Test Community")
				community.StartedAt = startDateTime
				community.EndedAt = endDateTime
				return &repository.MockCommunityRepository{
					FindByIDFunc: func(id string) (*domain.Community, error) {
						return community, nil
					},
					FindCardsFunc: func(id string) ([]domain.Card, error) {
						return []domain.Card{}, nil
					},
					UpdateHighlightedCardFunc: func(communityID string, highlightedCard *domain.HighlightedCard) error {
						return nil
					},
				}
			},
			setupCardRepo: func() *repository.MockCardRepository {
				return &repository.MockCardRepository{}
			},
			setupGitHub: func() *github.MockClient {
				return &github.MockClient{}
			},
			wantErr: false,
			validate: func(t *testing.T, community *domain.Community, highlightedCard *domain.HighlightedCard) {
				if highlightedCard == nil {
					t.Errorf("HighlightedCardがnilです")
					return
				}
				// 空のHighlightedCardが返されることを確認
				if highlightedCard.BestContributor.GithubID != "" {
					t.Errorf("空のHighlightedCardが期待されましたが、BestContributorが設定されています")
				}
			},
		},
		{
			name:        "NodeIDがない場合（空のHighlightedCardを保存）",
			communityID: "test-community-id",
			setupRepo: func() *repository.MockCommunityRepository {
				community := createTestCommunity("Test Community")
				community.StartedAt = startDateTime
				community.EndedAt = endDateTime
				return &repository.MockCommunityRepository{
					FindByIDFunc: func(id string) (*domain.Community, error) {
						return community, nil
					},
					FindCardsFunc: func(id string) ([]domain.Card, error) {
						card1 := createTestCard("12345")
						card1.NodeID = "" // NodeIDを空にする
						return []domain.Card{*card1}, nil
					},
					UpdateHighlightedCardFunc: func(communityID string, highlightedCard *domain.HighlightedCard) error {
						return nil
					},
				}
			},
			setupCardRepo: func() *repository.MockCardRepository {
				return &repository.MockCardRepository{}
			},
			setupGitHub: func() *github.MockClient {
				return &github.MockClient{}
			},
			wantErr: false,
			validate: func(t *testing.T, community *domain.Community, highlightedCard *domain.HighlightedCard) {
				if highlightedCard == nil {
					t.Errorf("HighlightedCardがnilです")
					return
				}
				// 空のHighlightedCardが返されることを確認
				if highlightedCard.BestContributor.GithubID != "" {
					t.Errorf("空のHighlightedCardが期待されましたが、BestContributorが設定されています")
				}
			},
		},
		{
			name:        "GitHub API呼び出し時のエラー",
			communityID: "test-community-id",
			setupRepo: func() *repository.MockCommunityRepository {
				community := createTestCommunity("Test Community")
				community.StartedAt = startDateTime
				community.EndedAt = endDateTime
				return &repository.MockCommunityRepository{
					FindByIDFunc: func(id string) (*domain.Community, error) {
						return community, nil
					},
					FindCardsFunc: func(id string) ([]domain.Card, error) {
						return []domain.Card{
							*createTestCard("12345"),
						}, nil
					},
				}
			},
			setupCardRepo: func() *repository.MockCardRepository {
				return &repository.MockCardRepository{}
			},
			setupGitHub: func() *github.MockClient {
				return &github.MockClient{
					GetUsersFullInfoByNodeIDsFunc: func(ctx context.Context, nodeIDs []string, from, to time.Time) ([]github.UserFullInfo, error) {
						return nil, fmt.Errorf("github api error")
					},
				}
			},
			wantErr:    true,
			wantErrMsg: "failed to get users full info by node ids",
		},
		{
			name:        "usersFullInfoが空の場合（空のHighlightedCardを保存）",
			communityID: "test-community-id",
			setupRepo: func() *repository.MockCommunityRepository {
				community := createTestCommunity("Test Community")
				community.StartedAt = startDateTime
				community.EndedAt = endDateTime
				return &repository.MockCommunityRepository{
					FindByIDFunc: func(id string) (*domain.Community, error) {
						return community, nil
					},
					FindCardsFunc: func(id string) ([]domain.Card, error) {
						return []domain.Card{
							*createTestCard("12345"),
						}, nil
					},
					UpdateHighlightedCardFunc: func(communityID string, highlightedCard *domain.HighlightedCard) error {
						return nil
					},
				}
			},
			setupCardRepo: func() *repository.MockCardRepository {
				return &repository.MockCardRepository{}
			},
			setupGitHub: func() *github.MockClient {
				return &github.MockClient{
					GetUsersFullInfoByNodeIDsFunc: func(ctx context.Context, nodeIDs []string, from, to time.Time) ([]github.UserFullInfo, error) {
						return []github.UserFullInfo{}, nil
					},
				}
			},
			wantErr: false,
			validate: func(t *testing.T, community *domain.Community, highlightedCard *domain.HighlightedCard) {
				if highlightedCard == nil {
					t.Errorf("HighlightedCardがnilです")
					return
				}
				// 空のHighlightedCardが返されることを確認
				if highlightedCard.BestContributor.GithubID != "" {
					t.Errorf("空のHighlightedCardが期待されましたが、BestContributorが設定されています")
				}
			},
		},
		{
			name:        "カード更新時のRepositoryエラー",
			communityID: "test-community-id",
			setupRepo: func() *repository.MockCommunityRepository {
				community := createTestCommunity("Test Community")
				community.StartedAt = startDateTime
				community.EndedAt = endDateTime
				return &repository.MockCommunityRepository{
					FindByIDFunc: func(id string) (*domain.Community, error) {
						return community, nil
					},
					FindCardsFunc: func(id string) ([]domain.Card, error) {
						return []domain.Card{
							*createTestCard("12345"),
						}, nil
					},
				}
			},
			setupCardRepo: func() *repository.MockCardRepository {
				return &repository.MockCardRepository{
					UpdateFunc: func(card *domain.Card) error {
						return fmt.Errorf("update error")
					},
				}
			},
			setupGitHub: func() *github.MockClient {
				return &github.MockClient{
					GetUsersFullInfoByNodeIDsFunc: func(ctx context.Context, nodeIDs []string, from, to time.Time) ([]github.UserFullInfo, error) {
						return []github.UserFullInfo{
							createTestUserFullInfo("user1", 100, 50, 10, 20, 20),
						}, nil
					},
				}
			},
			wantErr:    true,
			wantErrMsg: "failed to update card",
		},
		{
			name:        "HighlightedCard更新時のRepositoryエラー",
			communityID: "test-community-id",
			setupRepo: func() *repository.MockCommunityRepository {
				community := createTestCommunity("Test Community")
				community.StartedAt = startDateTime
				community.EndedAt = endDateTime
				return &repository.MockCommunityRepository{
					FindByIDFunc: func(id string) (*domain.Community, error) {
						return community, nil
					},
					FindCardsFunc: func(id string) ([]domain.Card, error) {
						return []domain.Card{
							*createTestCard("12345"),
						}, nil
					},
					UpdateHighlightedCardFunc: func(communityID string, highlightedCard *domain.HighlightedCard) error {
						return fmt.Errorf("update highlighted card error")
					},
				}
			},
			setupCardRepo: func() *repository.MockCardRepository {
				return &repository.MockCardRepository{
					UpdateFunc: func(card *domain.Card) error {
						return nil
					},
				}
			},
			setupGitHub: func() *github.MockClient {
				return &github.MockClient{
					GetUsersFullInfoByNodeIDsFunc: func(ctx context.Context, nodeIDs []string, from, to time.Time) ([]github.UserFullInfo, error) {
						return []github.UserFullInfo{
							createTestUserFullInfo("user1", 100, 50, 10, 20, 20),
						}, nil
					},
				}
			},
			wantErr:    true,
			wantErrMsg: "failed to update highlighted card",
		},
		{
			name:        "更新後のコミュニティ取得時のRepositoryエラー",
			communityID: "test-community-id",
			setupRepo: func() *repository.MockCommunityRepository {
				community := createTestCommunity("Test Community")
				community.StartedAt = startDateTime
				community.EndedAt = endDateTime
				return &repository.MockCommunityRepository{
					FindByIDFunc: func(id string) (*domain.Community, error) {
						return community, nil
					},
					FindCardsFunc: func(id string) ([]domain.Card, error) {
						return []domain.Card{
							*createTestCard("12345"),
						}, nil
					},
					UpdateHighlightedCardFunc: func(communityID string, highlightedCard *domain.HighlightedCard) error {
						return nil
					},
					FindByIDWithHighlightedCardFunc: func(id string) (*domain.Community, error) {
						return nil, fmt.Errorf("get updated community error")
					},
				}
			},
			setupCardRepo: func() *repository.MockCardRepository {
				return &repository.MockCardRepository{
					UpdateFunc: func(card *domain.Card) error {
						return nil
					},
				}
			},
			setupGitHub: func() *github.MockClient {
				return &github.MockClient{
					GetUsersFullInfoByNodeIDsFunc: func(ctx context.Context, nodeIDs []string, from, to time.Time) ([]github.UserFullInfo, error) {
						return []github.UserFullInfo{
							createTestUserFullInfo("user1", 100, 50, 10, 20, 20),
						}, nil
					},
				}
			},
			wantErr:    true,
			wantErrMsg: "failed to get updated community",
		},
		{
			name:        "一部のカードにNodeIDがない場合",
			communityID: "test-community-id",
			setupRepo: func() *repository.MockCommunityRepository {
				community := createTestCommunity("Test Community")
				community.StartedAt = startDateTime
				community.EndedAt = endDateTime
				var savedHighlightedCard *domain.HighlightedCard
				return &repository.MockCommunityRepository{
					FindByIDFunc: func(id string) (*domain.Community, error) {
						return community, nil
					},
					FindCardsFunc: func(id string) ([]domain.Card, error) {
						card1 := createTestCard("12345")
						card2 := createTestCard("67890")
						card2.NodeID = "" // NodeIDを空にする
						return []domain.Card{*card1, *card2}, nil
					},
					UpdateHighlightedCardFunc: func(communityID string, highlightedCard *domain.HighlightedCard) error {
						savedHighlightedCard = highlightedCard
						return nil
					},
					FindByIDWithHighlightedCardFunc: func(id string) (*domain.Community, error) {
						updatedCommunity := createTestCommunity("Test Community")
						updatedCommunity.StartedAt = startDateTime
						updatedCommunity.EndedAt = endDateTime
						if savedHighlightedCard != nil {
							updatedCommunity.HighlightedCard = *savedHighlightedCard
						}
						return updatedCommunity, nil
					},
				}
			},
			setupCardRepo: func() *repository.MockCardRepository {
				return &repository.MockCardRepository{
					UpdateFunc: func(card *domain.Card) error {
						return nil
					},
				}
			},
			setupGitHub: func() *github.MockClient {
				return &github.MockClient{
					GetUsersFullInfoByNodeIDsFunc: func(ctx context.Context, nodeIDs []string, from, to time.Time) ([]github.UserFullInfo, error) {
						// NodeIDがあるカード1つ分の情報のみ返す
						return []github.UserFullInfo{
							createTestUserFullInfo("user1", 100, 50, 10, 20, 20),
						}, nil
					},
				}
			},
			wantErr: false,
			validate: func(t *testing.T, community *domain.Community, highlightedCard *domain.HighlightedCard) {
				if highlightedCard == nil {
					t.Errorf("HighlightedCardがnilです")
					return
				}
				// NodeIDがあるカードのみが処理されることを確認
				if highlightedCard.BestContributor.GithubID == "" {
					t.Errorf("BestContributorが設定されていません")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			communityRepo := tt.setupRepo()
			cardRepo := tt.setupCardRepo()
			githubClient := tt.setupGitHub()
			service := NewCommunityService(communityRepo, cardRepo)

			community, highlightedCard, err := service.RefreshHighlightedCard(ctx, tt.communityID, githubClient)

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
				if community == nil {
					t.Errorf("コミュニティがnilです")
					return
				}
				if highlightedCard == nil {
					t.Errorf("HighlightedCardがnilです")
					return
				}
				if tt.validate != nil {
					tt.validate(t, community, highlightedCard)
				}
			}
		})
	}
}

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/furarico/octo-deck-api/internal/database"
	"github.com/furarico/octo-deck-api/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// テスト用のコミュニティを作成するヘルパー関数
func createTestCommunity(name string) *domain.Community {
	now := time.Now()
	return domain.NewCommunity(
		name,
		now,
		now.Add(30*24*time.Hour), // 30日後
		domain.HighlightedCard{},
	)
}

// CommunityRepositoryのFindAllメソッドをテスト
func TestCommunityRepository_FindAll(t *testing.T) {
	db := SetupTestDB(t)

	tests := []struct {
		name               string
		setup              func(db *gorm.DB) string // githubIDを返す
		wantCommunityCount int
		wantErr            bool
	}{
		{
			name: "自分のカードが所属するコミュニティを全て取得できる",
			setup: func(db *gorm.DB) string {
				githubID := "user123"

				// カードを作成
				card := createTestCard(githubID, "U_user123")
				dbCard := database.CardFromDomain(card)
				db.Create(dbCard)

				// コミュニティを2つ作成
				community1 := createTestCommunity("Community 1")
				community2 := createTestCommunity("Community 2")
				dbCommunity1 := &database.Community{
					ID:        uuid.UUID(community1.ID),
					Name:      community1.Name,
					StartedAt: community1.StartedAt,
					EndedAt:   community1.EndedAt,
				}
				dbCommunity2 := &database.Community{
					ID:        uuid.UUID(community2.ID),
					Name:      community2.Name,
					StartedAt: community2.StartedAt,
					EndedAt:   community2.EndedAt,
				}
				db.Create(dbCommunity1)
				db.Create(dbCommunity2)

				// 両方のコミュニティにカードを追加
				db.Create(&database.CommunityCard{
					CommunityID: dbCommunity1.ID,
					CardID:      dbCard.ID,
				})
				db.Create(&database.CommunityCard{
					CommunityID: dbCommunity2.ID,
					CardID:      dbCard.ID,
				})

				return githubID
			},
			wantCommunityCount: 2,
			wantErr:            false,
		},
		{
			name: "所属するコミュニティがない場合は空のスライスを返す",
			setup: func(db *gorm.DB) string {
				return "nocommunity_user"
			},
			wantCommunityCount: 0,
			wantErr:            false,
		},
		{
			name: "他のユーザーのコミュニティは取得しない",
			setup: func(db *gorm.DB) string {
				myGithubID := "myuser"
				otherGithubID := "otheruser"

				// 自分のカードを作成
				myCard := createTestCard(myGithubID, "U_myuser")
				dbMyCard := database.CardFromDomain(myCard)
				db.Create(dbMyCard)

				// 他のユーザーのカードを作成
				otherCard := createTestCard(otherGithubID, "U_otheruser")
				dbOtherCard := database.CardFromDomain(otherCard)
				db.Create(dbOtherCard)

				// コミュニティを2つ作成
				myCommunity := createTestCommunity("My Community")
				otherCommunity := createTestCommunity("Other Community")
				dbMyCommunity := &database.Community{
					ID:        uuid.UUID(myCommunity.ID),
					Name:      myCommunity.Name,
					StartedAt: myCommunity.StartedAt,
					EndedAt:   myCommunity.EndedAt,
				}
				dbOtherCommunity := &database.Community{
					ID:        uuid.UUID(otherCommunity.ID),
					Name:      otherCommunity.Name,
					StartedAt: otherCommunity.StartedAt,
					EndedAt:   otherCommunity.EndedAt,
				}
				db.Create(dbMyCommunity)
				db.Create(dbOtherCommunity)

				// 自分のコミュニティに自分のカードを追加
				db.Create(&database.CommunityCard{
					CommunityID: dbMyCommunity.ID,
					CardID:      dbMyCard.ID,
				})

				// 他のコミュニティに他のユーザーのカードを追加
				db.Create(&database.CommunityCard{
					CommunityID: dbOtherCommunity.ID,
					CardID:      dbOtherCard.ID,
				})

				return myGithubID
			},
			wantCommunityCount: 1,
			wantErr:            false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CleanupTestData(t, db)
			githubID := tt.setup(db)
			ctx := context.Background()

			repo := NewCommunityRepository(db)
			communities, err := repo.FindAll(ctx, githubID)

			if (err != nil) != tt.wantErr {
				t.Errorf("FindAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(communities) != tt.wantCommunityCount {
				t.Errorf("FindAll() returned %d communities, want %d", len(communities), tt.wantCommunityCount)
			}
		})
	}
}

// CommunityRepositoryのFindByIDメソッドをテスト
func TestCommunityRepository_FindByID(t *testing.T) {
	db := SetupTestDB(t)

	tests := []struct {
		name    string
		setup   func(db *gorm.DB) string // communityIDを返す
		wantErr bool
	}{
		{
			name: "存在するIDでコミュニティを取得できる",
			setup: func(db *gorm.DB) string {
				community := createTestCommunity("Test Community")
				dbCommunity := &database.Community{
					ID:        uuid.UUID(community.ID),
					Name:      community.Name,
					StartedAt: community.StartedAt,
					EndedAt:   community.EndedAt,
				}
				db.Create(dbCommunity)
				return dbCommunity.ID.String()
			},
			wantErr: false,
		},
		{
			name: "存在しないIDの場合エラーになる",
			setup: func(db *gorm.DB) string {
				return uuid.New().String()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CleanupTestData(t, db)
			communityID := tt.setup(db)
			ctx := context.Background()

			repo := NewCommunityRepository(db)
			community, err := repo.FindByID(ctx, communityID)

			if (err != nil) != tt.wantErr {
				t.Errorf("FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && community != nil {
				if uuid.UUID(community.ID).String() != communityID {
					t.Errorf("ID = %v, want %v", uuid.UUID(community.ID).String(), communityID)
				}
			}
		})
	}
}

// CommunityRepositoryのFindByIDWithHighlightedCardメソッドをテスト
func TestCommunityRepository_FindByIDWithHighlightedCard(t *testing.T) {
	db := SetupTestDB(t)

	tests := []struct {
		name                 string
		setup                func(db *gorm.DB) string // communityIDを返す
		wantErr              bool
		wantHighlightedCards bool
	}{
		{
			name: "HighlightedCard付きでコミュニティを取得できる",
			setup: func(db *gorm.DB) string {
				// ベストカードを作成
				card := createTestCard("bestuser", "U_bestuser")
				dbCard := database.CardFromDomain(card)
				db.Create(dbCard)

				// コミュニティを作成（HighlightedCardを設定）
				community := createTestCommunity("Community with Highlighted")
				dbCommunity := &database.Community{
					ID:                    uuid.UUID(community.ID),
					Name:                  community.Name,
					StartedAt:             community.StartedAt,
					EndedAt:               community.EndedAt,
					BestContributorCardID: &dbCard.ID,
					BestCommitterCardID:   &dbCard.ID,
				}
				db.Create(dbCommunity)
				return dbCommunity.ID.String()
			},
			wantErr:              false,
			wantHighlightedCards: true,
		},
		{
			name: "HighlightedCardがないコミュニティも取得できる",
			setup: func(db *gorm.DB) string {
				community := createTestCommunity("Community without Highlighted")
				dbCommunity := &database.Community{
					ID:        uuid.UUID(community.ID),
					Name:      community.Name,
					StartedAt: community.StartedAt,
					EndedAt:   community.EndedAt,
				}
				db.Create(dbCommunity)
				return dbCommunity.ID.String()
			},
			wantErr:              false,
			wantHighlightedCards: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CleanupTestData(t, db)
			communityID := tt.setup(db)
			ctx := context.Background()

			repo := NewCommunityRepository(db)
			community, err := repo.FindByIDWithHighlightedCard(ctx, communityID)

			if (err != nil) != tt.wantErr {
				t.Errorf("FindByIDWithHighlightedCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && community != nil {
				if tt.wantHighlightedCards {
					if community.HighlightedCard.BestContributor.GithubID == "" {
						t.Error("BestContributorが設定されていません")
					}
					if community.HighlightedCard.BestCommitter.GithubID == "" {
						t.Error("BestCommitterが設定されていません")
					}
				}
			}
		})
	}
}

// CommunityRepositoryのUpdateHighlightedCardメソッドをテスト
func TestCommunityRepository_UpdateHighlightedCard(t *testing.T) {
	db := SetupTestDB(t)

	tests := []struct {
		name    string
		setup   func(db *gorm.DB) (string, *domain.HighlightedCard)
		wantErr bool
	}{
		{
			name: "HighlightedCardを更新できる",
			setup: func(db *gorm.DB) (string, *domain.HighlightedCard) {
				// カードを作成
				card := createTestCard("highlighteduser", "U_highlighteduser")
				dbCard := database.CardFromDomain(card)
				db.Create(dbCard)

				// コミュニティを作成
				community := createTestCommunity("Test Community")
				dbCommunity := &database.Community{
					ID:        uuid.UUID(community.ID),
					Name:      community.Name,
					StartedAt: community.StartedAt,
					EndedAt:   community.EndedAt,
				}
				db.Create(dbCommunity)

				// HighlightedCardを作成
				domainCard := dbCard.ToDomain()
				highlightedCard := domain.NewHighlightedCard(
					*domainCard, // BestCommitter
					*domainCard, // BestContributor
					*domainCard, // BestIssuer
					*domainCard, // BestPullRequester
					*domainCard, // BestReviewer
				)

				return dbCommunity.ID.String(), highlightedCard
			},
			wantErr: false,
		},
		{
			name: "無効なコミュニティIDの場合エラーになる",
			setup: func(db *gorm.DB) (string, *domain.HighlightedCard) {
				return "invalid-uuid", &domain.HighlightedCard{}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CleanupTestData(t, db)
			communityID, highlightedCard := tt.setup(db)
			ctx := context.Background()

			repo := NewCommunityRepository(db)
			err := repo.UpdateHighlightedCard(ctx, communityID, highlightedCard)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateHighlightedCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// 更新されたことを確認
				var dbCommunity database.Community
				if err := db.Preload("BestContributorCard").First(&dbCommunity, "id = ?", communityID).Error; err != nil {
					t.Errorf("コミュニティが見つかりません: %v", err)
					return
				}

				if highlightedCard.BestContributor.GithubID != "" && dbCommunity.BestContributorCardID == nil {
					t.Error("BestContributorCardIDが更新されていません")
				}
			}
		})
	}
}

// CommunityRepositoryのFindCardsメソッドをテスト
func TestCommunityRepository_FindCards(t *testing.T) {
	db := SetupTestDB(t)

	tests := []struct {
		name          string
		setup         func(db *gorm.DB) string // communityIDを返す
		wantCardCount int
		wantErr       bool
	}{
		{
			name: "コミュニティに所属するカードを全て取得できる",
			setup: func(db *gorm.DB) string {
				// カードを2枚作成
				card1 := createTestCard("card1", "U_card1")
				card2 := createTestCard("card2", "U_card2")
				dbCard1 := database.CardFromDomain(card1)
				dbCard2 := database.CardFromDomain(card2)
				db.Create(dbCard1)
				db.Create(dbCard2)

				// コミュニティを作成
				community := createTestCommunity("Test Community")
				dbCommunity := &database.Community{
					ID:        uuid.UUID(community.ID),
					Name:      community.Name,
					StartedAt: community.StartedAt,
					EndedAt:   community.EndedAt,
				}
				db.Create(dbCommunity)

				// コミュニティにカードを追加
				db.Create(&database.CommunityCard{
					CommunityID: dbCommunity.ID,
					CardID:      dbCard1.ID,
				})
				db.Create(&database.CommunityCard{
					CommunityID: dbCommunity.ID,
					CardID:      dbCard2.ID,
				})

				return dbCommunity.ID.String()
			},
			wantCardCount: 2,
			wantErr:       false,
		},
		{
			name: "カードがないコミュニティは空のスライスを返す",
			setup: func(db *gorm.DB) string {
				community := createTestCommunity("Empty Community")
				dbCommunity := &database.Community{
					ID:        uuid.UUID(community.ID),
					Name:      community.Name,
					StartedAt: community.StartedAt,
					EndedAt:   community.EndedAt,
				}
				db.Create(dbCommunity)
				return dbCommunity.ID.String()
			},
			wantCardCount: 0,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CleanupTestData(t, db)
			communityID := tt.setup(db)
			ctx := context.Background()

			repo := NewCommunityRepository(db)
			cards, err := repo.FindCards(ctx, communityID)

			if (err != nil) != tt.wantErr {
				t.Errorf("FindCards() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(cards) != tt.wantCardCount {
				t.Errorf("FindCards() returned %d cards, want %d", len(cards), tt.wantCardCount)
			}
		})
	}
}

// CommunityRepositoryのCreateメソッドをテスト
func TestCommunityRepository_Create(t *testing.T) {
	db := SetupTestDB(t)

	tests := []struct {
		name      string
		community *domain.Community
		wantErr   bool
	}{
		{
			name:      "正常にコミュニティを作成できる",
			community: createTestCommunity("New Community"),
			wantErr:   false,
		},
		{
			name:      "別の名前でコミュニティを作成できる",
			community: createTestCommunity("Another Community"),
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CleanupTestData(t, db)
			ctx := context.Background()

			repo := NewCommunityRepository(db)
			err := repo.Create(ctx, tt.community)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// 作成されたコミュニティを確認
				var dbCommunity database.Community
				if err := db.First(&dbCommunity, "id = ?", uuid.UUID(tt.community.ID)).Error; err != nil {
					t.Errorf("作成されたコミュニティが見つかりません: %v", err)
					return
				}

				if dbCommunity.Name != tt.community.Name {
					t.Errorf("Name = %v, want %v", dbCommunity.Name, tt.community.Name)
				}
			}
		})
	}
}

// CommunityRepositoryのDeleteメソッドをテスト
func TestCommunityRepository_Delete(t *testing.T) {
	db := SetupTestDB(t)

	tests := []struct {
		name    string
		setup   func(db *gorm.DB) string
		wantErr bool
	}{
		{
			name: "コミュニティを削除できる",
			setup: func(db *gorm.DB) string {
				community := createTestCommunity("Delete Target")
				dbCommunity := &database.Community{
					ID:        uuid.UUID(community.ID),
					Name:      community.Name,
					StartedAt: community.StartedAt,
					EndedAt:   community.EndedAt,
				}
				db.Create(dbCommunity)
				return dbCommunity.ID.String()
			},
			wantErr: false,
		},
		{
			name: "存在しないコミュニティを削除してもエラーにならない",
			setup: func(db *gorm.DB) string {
				return uuid.New().String()
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CleanupTestData(t, db)
			communityID := tt.setup(db)
			ctx := context.Background()

			repo := NewCommunityRepository(db)
			err := repo.Delete(ctx, communityID)

			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// 削除されたことを確認
				var count int64
				db.Model(&database.Community{}).Where("id = ?", communityID).Count(&count)
				if count != 0 {
					t.Error("コミュニティが削除されていません")
				}
			}
		})
	}
}

// CommunityRepositoryのAddCardメソッドをテスト
func TestCommunityRepository_AddCard(t *testing.T) {
	db := SetupTestDB(t)

	tests := []struct {
		name    string
		setup   func(db *gorm.DB) (string, string) // communityID, cardIDを返す
		wantErr bool
	}{
		{
			name: "コミュニティにカードを追加できる",
			setup: func(db *gorm.DB) (string, string) {
				// カードを作成
				card := createTestCard("addcard", "U_addcard")
				dbCard := database.CardFromDomain(card)
				db.Create(dbCard)

				// コミュニティを作成
				community := createTestCommunity("Test Community")
				dbCommunity := &database.Community{
					ID:        uuid.UUID(community.ID),
					Name:      community.Name,
					StartedAt: community.StartedAt,
					EndedAt:   community.EndedAt,
				}
				db.Create(dbCommunity)

				return dbCommunity.ID.String(), dbCard.ID.String()
			},
			wantErr: false,
		},
		{
			name: "無効なコミュニティIDの場合エラーになる",
			setup: func(db *gorm.DB) (string, string) {
				return "invalid-uuid", uuid.New().String()
			},
			wantErr: true,
		},
		{
			name: "無効なカードIDの場合エラーになる",
			setup: func(db *gorm.DB) (string, string) {
				return uuid.New().String(), "invalid-uuid"
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CleanupTestData(t, db)
			communityID, cardID := tt.setup(db)
			ctx := context.Background()

			repo := NewCommunityRepository(db)
			err := repo.AddCard(ctx, communityID, cardID)

			if (err != nil) != tt.wantErr {
				t.Errorf("AddCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// 追加されたことを確認
				var count int64
				db.Model(&database.CommunityCard{}).
					Where("community_id = ? AND card_id = ?", communityID, cardID).
					Count(&count)

				if count != 1 {
					t.Error("コミュニティにカードが追加されていません")
				}
			}
		})
	}
}

// CommunityRepositoryのRemoveCardメソッドをテスト
func TestCommunityRepository_RemoveCard(t *testing.T) {
	db := SetupTestDB(t)

	tests := []struct {
		name    string
		setup   func(db *gorm.DB) (string, string) // communityID, cardIDを返す
		wantErr bool
	}{
		{
			name: "コミュニティからカードを削除できる",
			setup: func(db *gorm.DB) (string, string) {
				// カードを作成
				card := createTestCard("removecard", "U_removecard")
				dbCard := database.CardFromDomain(card)
				db.Create(dbCard)

				// コミュニティを作成
				community := createTestCommunity("Test Community")
				dbCommunity := &database.Community{
					ID:        uuid.UUID(community.ID),
					Name:      community.Name,
					StartedAt: community.StartedAt,
					EndedAt:   community.EndedAt,
				}
				db.Create(dbCommunity)

				// コミュニティにカードを追加
				db.Create(&database.CommunityCard{
					CommunityID: dbCommunity.ID,
					CardID:      dbCard.ID,
				})

				return dbCommunity.ID.String(), dbCard.ID.String()
			},
			wantErr: false,
		},
		{
			name: "存在しないカードを削除してもエラーにならない",
			setup: func(db *gorm.DB) (string, string) {
				community := createTestCommunity("Empty Community")
				dbCommunity := &database.Community{
					ID:        uuid.UUID(community.ID),
					Name:      community.Name,
					StartedAt: community.StartedAt,
					EndedAt:   community.EndedAt,
				}
				db.Create(dbCommunity)

				return dbCommunity.ID.String(), uuid.New().String()
			},
			wantErr: false,
		},
		{
			name: "無効なコミュニティIDの場合エラーになる",
			setup: func(db *gorm.DB) (string, string) {
				return "invalid-uuid", uuid.New().String()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CleanupTestData(t, db)
			communityID, cardID := tt.setup(db)
			ctx := context.Background()

			repo := NewCommunityRepository(db)
			err := repo.RemoveCard(ctx, communityID, cardID)

			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// 削除されたことを確認
				var count int64
				db.Model(&database.CommunityCard{}).
					Where("community_id = ? AND card_id = ?", communityID, cardID).
					Count(&count)

				if count != 0 {
					t.Error("コミュニティからカードが削除されていません")
				}
			}
		})
	}
}

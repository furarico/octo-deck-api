package repository

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/furarico/octo-deck-api/internal/database"
	"github.com/furarico/octo-deck-api/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// テスト用のカードを作成するヘルパー関数
func createTestCard(githubID, nodeID string) *domain.Card {
	return &domain.Card{
		ID:       domain.NewCardID(),
		GithubID: githubID,
		NodeID:   nodeID,
		Color:    domain.Color("#FF5733"),
		Blocks: domain.Blocks{
			{true, false, true, false, true},
			{false, true, false, true, false},
			{true, false, true, false, true},
			{false, true, false, true, false},
			{true, false, true, false, true},
		},
		UserName: "testuser",
		FullName: "Test User",
		IconUrl:  "https://example.com/icon.png",
		MostUsedLanguage: domain.Language{
			LanguageName: "Go",
			Color:        "#00ADD8",
		},
	}
}

// CardRepositoryのCreateメソッドをテスト
func TestCardRepository_Create(t *testing.T) {
	db := SetupTestDB(t)

	tests := []struct {
		name    string
		card    *domain.Card
		wantErr bool
	}{
		{
			name:    "正常にカードを作成できる",
			card:    createTestCard("12345", "U_12345"),
			wantErr: false,
		},
		{
			name:    "別のGitHubIDでカードを作成できる",
			card:    createTestCard("67890", "U_67890"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CleanupTestData(t, db)
			ctx := context.Background()

			repo := NewCardRepository(db)
			err := repo.Create(ctx, tt.card)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// 作成されたカードを確認
				var dbCard database.Card
				if err := db.First(&dbCard, "github_id = ?", tt.card.GithubID).Error; err != nil {
					t.Errorf("作成されたカードが見つかりません: %v", err)
					return
				}

				if dbCard.GithubID != tt.card.GithubID {
					t.Errorf("GithubID = %v, want %v", dbCard.GithubID, tt.card.GithubID)
				}
				if dbCard.NodeID != tt.card.NodeID {
					t.Errorf("NodeID = %v, want %v", dbCard.NodeID, tt.card.NodeID)
				}
				if dbCard.UserName != tt.card.UserName {
					t.Errorf("UserName = %v, want %v", dbCard.UserName, tt.card.UserName)
				}
			}
		})
	}
}

// CardRepositoryのFindByGitHubIDメソッドをテスト
func TestCardRepository_FindByGitHubID(t *testing.T) {
	db := SetupTestDB(t)

	tests := []struct {
		name      string
		setup     func(db *gorm.DB)
		githubID  string
		wantErr   bool
		wantMatch bool
	}{
		{
			name: "存在するGitHubIDでカードを取得できる",
			setup: func(db *gorm.DB) {
				card := createTestCard("12345", "U_12345")
				dbCard := database.CardFromDomain(card)
				db.Create(dbCard)
			},
			githubID:  "12345",
			wantErr:   false,
			wantMatch: true,
		},
		{
			name: "存在しないGitHubIDの場合エラーになる",
			setup: func(db *gorm.DB) {
				// データを作成しない
			},
			githubID:  "nonexistent",
			wantErr:   true,
			wantMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CleanupTestData(t, db)
			tt.setup(db)
			ctx := context.Background()

			repo := NewCardRepository(db)
			card, err := repo.FindByGitHubID(ctx, tt.githubID)

			if (err != nil) != tt.wantErr {
				t.Errorf("FindByGitHubID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantMatch && card != nil {
				if card.GithubID != tt.githubID {
					t.Errorf("GithubID = %v, want %v", card.GithubID, tt.githubID)
				}
			}
		})
	}
}

// CardRepositoryのFindMyCardメソッドをテスト
func TestCardRepository_FindMyCard(t *testing.T) {
	db := SetupTestDB(t)

	tests := []struct {
		name     string
		setup    func(db *gorm.DB)
		githubID string
		wantErr  bool
	}{
		{
			name: "自分のカードを取得できる",
			setup: func(db *gorm.DB) {
				card := createTestCard("mycard123", "U_mycard123")
				dbCard := database.CardFromDomain(card)
				db.Create(dbCard)
			},
			githubID: "mycard123",
			wantErr:  false,
		},
		{
			name: "カードが存在しない場合エラーになる",
			setup: func(db *gorm.DB) {
				// データを作成しない
			},
			githubID: "nocard",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CleanupTestData(t, db)
			tt.setup(db)
			ctx := context.Background()

			repo := NewCardRepository(db)
			card, err := repo.FindMyCard(ctx, tt.githubID)

			if (err != nil) != tt.wantErr {
				t.Errorf("FindMyCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && card != nil {
				if card.GithubID != tt.githubID {
					t.Errorf("GithubID = %v, want %v", card.GithubID, tt.githubID)
				}
			}
		})
	}
}

// CardRepositoryのFindAllメソッドをテスト
func TestCardRepository_FindAll(t *testing.T) {
	db := SetupTestDB(t)

	tests := []struct {
		name          string
		setup         func(db *gorm.DB) string // collectorGithubIDを返す
		wantCardCount int
		wantErr       bool
	}{
		{
			name: "コレクションにあるカードを全て取得できる",
			setup: func(db *gorm.DB) string {
				collectorID := "collector123"

				// カードを2枚作成
				card1 := createTestCard("card1", "U_card1")
				card2 := createTestCard("card2", "U_card2")
				dbCard1 := database.CardFromDomain(card1)
				dbCard2 := database.CardFromDomain(card2)
				db.Create(dbCard1)
				db.Create(dbCard2)

				// コレクションに追加
				db.Create(&database.CollectedCard{
					ID:                uuid.New(),
					CollectorGithubID: collectorID,
					CardID:            dbCard1.ID,
				})
				db.Create(&database.CollectedCard{
					ID:                uuid.New(),
					CollectorGithubID: collectorID,
					CardID:            dbCard2.ID,
				})

				return collectorID
			},
			wantCardCount: 2,
			wantErr:       false,
		},
		{
			name: "コレクションが空の場合は空のスライスを返す",
			setup: func(db *gorm.DB) string {
				return "emptycollector"
			},
			wantCardCount: 0,
			wantErr:       false,
		},
		{
			name: "他のユーザーのコレクションは取得しない",
			setup: func(db *gorm.DB) string {
				myCollectorID := "mycollector"
				otherCollectorID := "othercollector"

				// カードを2枚作成
				card1 := createTestCard("mycard", "U_mycard")
				card2 := createTestCard("othercard", "U_othercard")
				dbCard1 := database.CardFromDomain(card1)
				dbCard2 := database.CardFromDomain(card2)
				db.Create(dbCard1)
				db.Create(dbCard2)

				// 自分のコレクションに1枚追加
				db.Create(&database.CollectedCard{
					ID:                uuid.New(),
					CollectorGithubID: myCollectorID,
					CardID:            dbCard1.ID,
				})

				// 他人のコレクションに1枚追加
				db.Create(&database.CollectedCard{
					ID:                uuid.New(),
					CollectorGithubID: otherCollectorID,
					CardID:            dbCard2.ID,
				})

				return myCollectorID
			},
			wantCardCount: 1,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CleanupTestData(t, db)
			collectorID := tt.setup(db)
			ctx := context.Background()

			repo := NewCardRepository(db)
			cards, err := repo.FindAll(ctx, collectorID)

			if (err != nil) != tt.wantErr {
				t.Errorf("FindAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(cards) != tt.wantCardCount {
				t.Errorf("FindAll() returned %d cards, want %d", len(cards), tt.wantCardCount)
			}
		})
	}
}

// CardRepositoryのUpdateメソッドをテスト
func TestCardRepository_Update(t *testing.T) {
	db := SetupTestDB(t)

	tests := []struct {
		name       string
		setup      func(db *gorm.DB) *domain.Card
		updateFunc func(card *domain.Card)
		wantErr    bool
	}{
		{
			name: "カード情報を更新できる",
			setup: func(db *gorm.DB) *domain.Card {
				card := createTestCard("updatetest", "U_updatetest")
				dbCard := database.CardFromDomain(card)
				db.Create(dbCard)
				card.ID = domain.CardID(dbCard.ID)
				return card
			},
			updateFunc: func(card *domain.Card) {
				card.UserName = "updateduser"
				card.FullName = "Updated User"
				card.Color = domain.Color("#00FF00")
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CleanupTestData(t, db)
			card := tt.setup(db)
			tt.updateFunc(card)
			ctx := context.Background()

			repo := NewCardRepository(db)
			err := repo.Update(ctx, card)

			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// 更新されたカードを確認
				var dbCard database.Card
				if err := db.First(&dbCard, "id = ?", uuid.UUID(card.ID)).Error; err != nil {
					t.Errorf("更新されたカードが見つかりません: %v", err)
					return
				}

				if dbCard.UserName != card.UserName {
					t.Errorf("UserName = %v, want %v", dbCard.UserName, card.UserName)
				}
				if dbCard.FullName != card.FullName {
					t.Errorf("FullName = %v, want %v", dbCard.FullName, card.FullName)
				}
			}
		})
	}
}

// CardRepositoryのAddToCollectedCardsメソッドをテスト
func TestCardRepository_AddToCollectedCards(t *testing.T) {
	db := SetupTestDB(t)

	tests := []struct {
		name    string
		setup   func(db *gorm.DB) (string, domain.CardID)
		wantErr bool
	}{
		{
			name: "カードをコレクションに追加できる",
			setup: func(db *gorm.DB) (string, domain.CardID) {
				card := createTestCard("addtest", "U_addtest")
				dbCard := database.CardFromDomain(card)
				db.Create(dbCard)
				return "collector123", domain.CardID(dbCard.ID)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CleanupTestData(t, db)
			collectorID, cardID := tt.setup(db)
			ctx := context.Background()

			repo := NewCardRepository(db)
			err := repo.AddToCollectedCards(ctx, collectorID, cardID)

			if (err != nil) != tt.wantErr {
				t.Errorf("AddToCollectedCards() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// 追加されたことを確認
				var count int64
				db.Model(&database.CollectedCard{}).
					Where("collector_github_id = ? AND card_id = ?", collectorID, uuid.UUID(cardID)).
					Count(&count)

				if count != 1 {
					t.Errorf("コレクションにカードが追加されていません")
				}
			}
		})
	}
}

// CardRepositoryのRemoveFromCollectedCardsメソッドをテスト
func TestCardRepository_RemoveFromCollectedCards(t *testing.T) {
	db := SetupTestDB(t)

	tests := []struct {
		name    string
		setup   func(db *gorm.DB) (string, domain.CardID)
		wantErr bool
	}{
		{
			name: "カードをコレクションから削除できる",
			setup: func(db *gorm.DB) (string, domain.CardID) {
				card := createTestCard("removetest", "U_removetest")
				dbCard := database.CardFromDomain(card)
				db.Create(dbCard)

				collectorID := "collector123"
				db.Create(&database.CollectedCard{
					ID:                uuid.New(),
					CollectorGithubID: collectorID,
					CardID:            dbCard.ID,
				})

				return collectorID, domain.CardID(dbCard.ID)
			},
			wantErr: false,
		},
		{
			name: "存在しないカードを削除してもエラーにならない",
			setup: func(db *gorm.DB) (string, domain.CardID) {
				return "nonexistent", domain.NewCardID()
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CleanupTestData(t, db)
			collectorID, cardID := tt.setup(db)
			ctx := context.Background()

			repo := NewCardRepository(db)
			err := repo.RemoveFromCollectedCards(ctx, collectorID, cardID)

			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveFromCollectedCards() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// 削除されたことを確認
				var count int64
				db.Model(&database.CollectedCard{}).
					Where("collector_github_id = ? AND card_id = ?", collectorID, uuid.UUID(cardID)).
					Count(&count)

				if count != 0 {
					t.Errorf("コレクションからカードが削除されていません")
				}
			}
		})
	}
}

// CardRepositoryのFindAllCardsInDBメソッドをテスト
func TestCardRepository_FindAllCardsInDB(t *testing.T) {
	db := SetupTestDB(t)

	tests := []struct {
		name          string
		setup         func(db *gorm.DB)
		wantCardCount int
		wantErr       bool
	}{
		{
			name: "データベース内の全カードを取得できる",
			setup: func(db *gorm.DB) {
				// カードを3枚作成
				card1 := createTestCard("card1", "U_card1")
				card2 := createTestCard("card2", "U_card2")
				card3 := createTestCard("card3", "U_card3")
				db.Create(database.CardFromDomain(card1))
				db.Create(database.CardFromDomain(card2))
				db.Create(database.CardFromDomain(card3))
			},
			wantCardCount: 3,
			wantErr:       false,
		},
		{
			name: "カードがない場合は空のスライスを返す",
			setup: func(db *gorm.DB) {
				// データを作成しない
			},
			wantCardCount: 0,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CleanupTestData(t, db)
			tt.setup(db)
			ctx := context.Background()

			repo := NewCardRepository(db)
			cards, err := repo.FindAllCardsInDB(ctx)

			if (err != nil) != tt.wantErr {
				t.Errorf("FindAllCardsInDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(cards) != tt.wantCardCount {
				t.Errorf("FindAllCardsInDB() returned %d cards, want %d", len(cards), tt.wantCardCount)
			}
		})
	}
}

// CardRepositoryのPostgres特定の機能をテスト
func TestCardRepository_PostgresSpecificFeatures(t *testing.T) {
	db := SetupTestDB(t)

	t.Run("UUID自動生成が動作する", func(t *testing.T) {
		CleanupTestData(t, db)

		// IDを指定せずにカードを作成
		card := &database.Card{
			GithubID:   "uuidtest",
			NodeID:     "U_uuidtest",
			Color:      "#FF0000",
			BlocksData: json.RawMessage(`[]`),
		}
		if err := db.Create(card).Error; err != nil {
			t.Fatalf("カード作成に失敗: %v", err)
		}

		if card.ID == uuid.Nil {
			t.Error("UUIDが自動生成されていません")
		}
	})

	t.Run("JSONB型のBlocksDataが正しく保存・取得される", func(t *testing.T) {
		CleanupTestData(t, db)

		blocks := domain.Blocks{
			{true, false, true, false, true},
			{false, true, false, true, false},
			{true, false, true, false, true},
			{false, true, false, true, false},
			{true, false, true, false, true},
		}
		blocksData, _ := json.Marshal(blocks)

		card := &database.Card{
			GithubID:   "jsonbtest",
			NodeID:     "U_jsonbtest",
			Color:      "#00FF00",
			BlocksData: blocksData,
		}
		if err := db.Create(card).Error; err != nil {
			t.Fatalf("カード作成に失敗: %v", err)
		}

		var retrieved database.Card
		if err := db.First(&retrieved, "github_id = ?", "jsonbtest").Error; err != nil {
			t.Fatalf("カード取得に失敗: %v", err)
		}

		var retrievedBlocks domain.Blocks
		if err := json.Unmarshal(retrieved.BlocksData, &retrievedBlocks); err != nil {
			t.Fatalf("JSONBのアンマーシャルに失敗: %v", err)
		}

		// [5][5]bool型なので行数を確認
		if len(retrievedBlocks) != 5 {
			t.Errorf("Blocks行数 = %d, want 5", len(retrievedBlocks))
		}
		// 最初の行の最初の値を確認
		if retrievedBlocks[0][0] != true {
			t.Errorf("Block[0][0] = %v, want true", retrievedBlocks[0][0])
		}
	})
}

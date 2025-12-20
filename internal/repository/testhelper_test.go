package repository

import (
	"context"
	"testing"
	"time"

	"github.com/furarico/octo-deck-api/internal/database"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupTestDB はテスト用のPostgreSQLコンテナを起動し、GORMのDB接続を返す
// テスト終了時に自動的にコンテナがクリーンアップされる
func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	ctx := context.Background()

	// PostgreSQLコンテナを起動
	pgContainer, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	// テスト終了時にコンテナをクリーンアップ
	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Logf("failed to terminate postgres container: %v", err)
		}
	})

	// 接続文字列を取得
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	// GORMで接続
	db, err := gorm.Open(gormpostgres.Open(connStr), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	// マイグレーションを実行
	if err := database.AutoMigrate(db); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	return db
}

// CleanupTestData はテストデータをクリーンアップする
func CleanupTestData(t *testing.T, db *gorm.DB) {
	t.Helper()

	// 外部キー制約を考慮して削除順序を指定
	tables := []string{"collected_cards", "community_cards", "communities", "cards"}
	for _, table := range tables {
		if err := db.Exec("TRUNCATE TABLE " + table + " CASCADE").Error; err != nil {
			t.Logf("failed to truncate table %s: %v", table, err)
		}
	}
}

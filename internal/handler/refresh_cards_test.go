package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	api "github.com/furarico/octo-deck-api/generated"
	"github.com/furarico/octo-deck-api/internal/domain"
	"github.com/furarico/octo-deck-api/internal/service"
	"github.com/gin-gonic/gin"
)

// 全カードを更新するテスト
func TestRefreshAllCards(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		setupMock   func() *service.MockCardService
		setupContext func(c *gin.Context)
		wantCode    int
		validate    func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常に全カードを更新できる",
			setupMock: func() *service.MockCardService {
				return &service.MockCardService{
					RefreshAllCardsFunc: func(ctx context.Context, githubClient service.GitHubClient) ([]domain.Card, error) {
						return []domain.Card{
							{
								ID:       domain.NewCardID(),
								GithubID: "12345",
								UserName: "user1",
								FullName: "User One",
								IconUrl:  "https://example.com/user1.png",
								Color:    domain.Color("#111111"),
								Blocks:   domain.Blocks{},
								MostUsedLanguage: domain.Language{
									LanguageName: "Go",
									Color:        "#00ADD8",
								},
							},
							{
								ID:       domain.NewCardID(),
								GithubID: "67890",
								UserName: "user2",
								FullName: "User Two",
								IconUrl:  "https://example.com/user2.png",
								Color:    domain.Color("#222222"),
								Blocks:   domain.Blocks{},
								MostUsedLanguage: domain.Language{
									LanguageName: "TypeScript",
									Color:        "#3178C6",
								},
							},
						}, nil
					},
				}
			},
			wantCode: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Card []api.Card `json:"card"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("JSONパースに失敗しました: %v", err)
				}

				// カード数をチェック
				if len(response.Card) != 2 {
					t.Errorf("カード数が違う: 期待=2, 実際=%d", len(response.Card))
				}

				// 最初のカードの内容をチェック
				if response.Card[0].GithubId != "12345" {
					t.Errorf("最初のカードのGithubIDが違う: 期待=12345, 実際=%s", response.Card[0].GithubId)
				}
				if response.Card[0].UserName != "user1" {
					t.Errorf("最初のカードのUserNameが違う: 期待=user1, 実際=%s", response.Card[0].UserName)
				}
			},
		},
		{
			name: "空の結果を正常に返せる",
			setupMock: func() *service.MockCardService {
				return &service.MockCardService{
					RefreshAllCardsFunc: func(ctx context.Context, githubClient service.GitHubClient) ([]domain.Card, error) {
						return []domain.Card{}, nil
					},
				}
			},
			wantCode: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Card []api.Card `json:"card"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("JSONパースに失敗しました: %v", err)
				}

				// カード数をチェック
				if len(response.Card) != 0 {
					t.Errorf("カード数が違う: 期待=0, 実際=%d", len(response.Card))
				}
			},
		},
		{
			name: "RefreshAllCardsでエラーが発生した場合",
			setupMock: func() *service.MockCardService {
				return &service.MockCardService{
					RefreshAllCardsFunc: func(ctx context.Context, githubClient service.GitHubClient) ([]domain.Card, error) {
						return nil, fmt.Errorf("failed to refresh all cards")
					},
				}
			},
			wantCode: http.StatusInternalServerError,
			validate: nil,
		},
		{
			name: "GitHub Clientが取得できない場合",
			setupMock: func() *service.MockCardService {
				return &service.MockCardService{
					RefreshAllCardsFunc: func(ctx context.Context, githubClient service.GitHubClient) ([]domain.Card, error) {
						return []domain.Card{}, nil
					},
				}
			},
			wantCode: http.StatusInternalServerError,
			setupContext: func(c *gin.Context) {
				// GitHub Clientを設定しない（エラーになる）
				ctx := c.Request.Context()
				ctx = context.WithValue(ctx, GitHubIDKey, "test_user")
				ctx = context.WithValue(ctx, GitHubNodeIDKey, "MDQ6VXNlcjEyMzQ1")
				c.Request = c.Request.WithContext(ctx)
			},
			validate: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			mockService := tt.setupMock()
			cardHandler := NewCardHandler(mockService)
			router := gin.Default()

			// setTestContextまたはカスタムコンテキスト設定を使用
			if tt.setupContext != nil {
				router.Use(tt.setupContext)
			} else {
				router.Use(setTestContext)
			}

			strictHandler := api.NewStrictHandler(cardHandler, nil)
			api.RegisterHandlers(router, strictHandler)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", "/cards/refresh", nil)
			router.ServeHTTP(w, req)

			if w.Code != tt.wantCode {
				t.Errorf("ステータスコードが違う: 期待=%d, 実際=%d", tt.wantCode, w.Code)
			}

			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}

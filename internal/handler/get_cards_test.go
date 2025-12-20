package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	api "github.com/furarico/octo-deck-api/generated"
	"github.com/furarico/octo-deck-api/internal/domain"
	"github.com/furarico/octo-deck-api/internal/service"
	"github.com/gin-gonic/gin"
)

// カード一覧を取得するテスト
func TestGetCards(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name      string
		setupMock func() *service.MockCardService
		wantCode  int
		validate  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常にカード一覧を取得できる",
			setupMock: func() *service.MockCardService {
				return &service.MockCardService{
					GetAllCardsFunc: func(githubID string) ([]domain.Card, error) {
						return []domain.Card{
							{
								ID:       domain.NewCardID(),
								GithubID: "user1",
								UserName: "user1",
								FullName: "User One",
								IconUrl:  "https://example.com/user1.png",
								Color:    "#111111",
								Blocks:   domain.Blocks{},
							},
							{
								ID:       domain.NewCardID(),
								GithubID: "user2",
								UserName: "user2",
								FullName: "User Two",
								IconUrl:  "https://example.com/user2.png",
								Color:    "#222222",
								Blocks:   domain.Blocks{},
							},
							{
								ID:       domain.NewCardID(),
								GithubID: "user3",
								UserName: "user3",
								FullName: "User Three",
								IconUrl:  "https://example.com/user3.png",
								Color:    "#333333",
								Blocks:   domain.Blocks{},
							},
						}, nil
					},
				}
			},
			wantCode: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Cards []api.Card `json:"cards"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("JSONパースに失敗しました: %v", err)
				}

				// カード数をチェック
				if len(response.Cards) != 3 {
					t.Errorf("カード数が違う: 期待=3, 実際=%d", len(response.Cards))
				}
			},
		},
		{
			name: "空の結果を正常に返せる",
			setupMock: func() *service.MockCardService {
				return &service.MockCardService{
					GetAllCardsFunc: func(githubID string) ([]domain.Card, error) {
						return []domain.Card{}, nil
					},
				}
			},
			wantCode: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Cards []api.Card `json:"cards"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("JSONパースに失敗しました: %v", err)
				}

				// カード数をチェック
				if len(response.Cards) != 0 {
					t.Errorf("カード数が違う: 期待=0, 実際=%d", len(response.Cards))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			mockService := tt.setupMock()
			cardHandler := NewCardHandler(mockService)
			router := gin.Default()
			// context.Context に値を設定するミドルウェア
			router.Use(setTestContext)
			strictHandler := api.NewStrictHandler(cardHandler, nil)
			api.RegisterHandlers(router, strictHandler)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/cards", nil)
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

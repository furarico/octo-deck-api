package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	api "github.com/furarico/octo-deck-api/generated"
	"github.com/furarico/octo-deck-api/internal/repository"
	"github.com/furarico/octo-deck-api/internal/service"
	"github.com/gin-gonic/gin"
)

// カード一覧を取得するテスト
func TestGetCards(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name      string
		setupMock func() *repository.MockCardRepository
		wantCode  int
		validate  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常にカード一覧を取得できる",
			setupMock: func() *repository.MockCardRepository {
				return repository.NewMockCardRepository()
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
			setupMock: func() *repository.MockCardRepository {
				return repository.NewMockCardRepository()
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			repo := tt.setupMock()
			cardService := service.NewCardService(repo)
			cardHandler := NewHandler(cardService)
			router := gin.Default()
			api.RegisterHandlers(router, cardHandler)

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

// 自分のカードを取得するテスト
func TestGetMyCard(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name      string
		setupMock func() *repository.MockCardRepository
		wantCode  int
		validate  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常に自分のカードを取得できる",
			setupMock: func() *repository.MockCardRepository {
				return repository.NewMockCardRepository()
			},
			wantCode: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response api.Card
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("JSONパースに失敗しました: %v", err)
				}
			},
		},
		{
			name: "空の結果を正常に返せる",
			setupMock: func() *repository.MockCardRepository {
				return repository.NewMockCardRepository()
			},
			wantCode: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response api.Card
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("JSONパースに失敗しました: %v", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			repo := tt.setupMock()
			cardService := service.NewCardService(repo)
			cardHandler := NewHandler(cardService)
			router := gin.Default()
			api.RegisterHandlers(router, cardHandler)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/cards/me", nil)
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

// 特定のカードを取得するテスト
func TestGetCard(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name      string
		setupMock func() *repository.MockCardRepository
		wantCode  int
		validate  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常に特定のカードを取得できる",
			setupMock: func() *repository.MockCardRepository {
				return repository.NewMockCardRepository()
			},
			wantCode: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response api.Card
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("JSONパースに失敗しました: %v", err)
				}
			},
		},
		{
			name: "空の結果を正常に返せる",
			setupMock: func() *repository.MockCardRepository {
				return repository.NewMockCardRepository()
			},
			wantCode: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response api.Card
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("JSONパースに失敗しました: %v", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			repo := tt.setupMock()
			cardService := service.NewCardService(repo)
			cardHandler := NewHandler(cardService)
			router := gin.Default()
			api.RegisterHandlers(router, cardHandler)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/cards/1", nil)
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

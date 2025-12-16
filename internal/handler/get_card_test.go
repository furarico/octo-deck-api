package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	api "github.com/furarico/octo-deck-api/generated"
	"github.com/furarico/octo-deck-api/internal/domain"
	"github.com/furarico/octo-deck-api/internal/repository"
	"github.com/furarico/octo-deck-api/internal/service"
	"github.com/gin-gonic/gin"
)

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
				return &repository.MockCardRepository{
					FindByIDFunc: func(id string) (*domain.Card, error) {
						return &domain.Card{
							ID:       domain.NewCardID(),
							GithubID: "john_doe",
							Color:    domain.Color("#000000"),
							Blocks:   domain.Blocks{},
						}, nil
					},
				}
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
			name: "カードが見つからない場合はエラーを返す",
			setupMock: func() *repository.MockCardRepository {
				return &repository.MockCardRepository{
					FindByIDFunc: func(id string) (*domain.Card, error) {
						return nil, fmt.Errorf("card not found: id=%s", id)
					},
				}
			},
			wantCode: http.StatusNotFound,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("JSONパースに失敗しました: %v", err)
				}
				if _, ok := response["error"]; !ok {
					t.Errorf("エラーメッセージがありません")
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

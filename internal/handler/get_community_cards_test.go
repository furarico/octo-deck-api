package handler

import (
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

// コミュニティのカード一覧取得のテスト
func TestGetCommunityCards(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name      string
		setupMock func() *service.MockCommunityService
		wantCode  int
		validate  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常にコミュニティのカード一覧を取得できる",
			setupMock: func() *service.MockCommunityService {
				return &service.MockCommunityService{
					GetCommunityCardsFunc: func(id string) ([]domain.Card, error) {
						return []domain.Card{
							{
								GithubID: "user1",
								UserName: "user1",
								FullName: "User One",
								IconUrl:  "https://example.com/user1.png",
								Color:    domain.Color("#000000"),
								Blocks:   domain.Blocks{},
							},
							{
								GithubID: "user2",
								UserName: "user2",
								FullName: "User Two",
								IconUrl:  "https://example.com/user2.png",
								Color:    domain.Color("#FFFFFF"),
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
				if len(response.Cards) != 2 {
					t.Errorf("カードの数が違う: 期待=2, 実際=%d", len(response.Cards))
				}
			},
		},
		{
			name: "カードが空の場合",
			setupMock: func() *service.MockCommunityService {
				return &service.MockCommunityService{
					GetCommunityCardsFunc: func(id string) ([]domain.Card, error) {
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
				if len(response.Cards) != 0 {
					t.Errorf("カードの数が違う: 期待=0, 実際=%d", len(response.Cards))
				}
			},
		},
		{
			name: "サービスでエラーが発生した場合",
			setupMock: func() *service.MockCommunityService {
				return &service.MockCommunityService{
					GetCommunityCardsFunc: func(id string) ([]domain.Card, error) {
						return nil, fmt.Errorf("database error")
					},
				}
			},
			wantCode: http.StatusInternalServerError,
			validate: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			mockService := tt.setupMock()
			communityHandler := NewCommunityHandler(mockService)
			router := gin.Default()
			router.Use(setTestContext)
			strictHandler := api.NewStrictHandler(communityHandler, nil)
			api.RegisterHandlers(router, strictHandler)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/communities/test-id/cards", nil)
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

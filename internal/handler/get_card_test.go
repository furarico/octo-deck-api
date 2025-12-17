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
	"github.com/furarico/octo-deck-api/internal/github"
	"github.com/furarico/octo-deck-api/internal/service"
	"github.com/gin-gonic/gin"
)

// 特定のカードを取得するテスト
func TestGetCard(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name      string
		setupMock func() *service.MockCardService
		wantCode  int
		validate  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常に特定のカードを取得できる",
			setupMock: func() *service.MockCardService {
				return &service.MockCardService{
					GetCardByGitHubIDFunc: func(ctx context.Context, githubID string, githubClient *github.Client) (*domain.Card, error) {
						return &domain.Card{
							ID:       domain.NewCardID(),
							GithubID: "john_doe",
							UserName: "john_doe",
							FullName: "John Doe",
							IconUrl:  "https://example.com/avatar.png",
							Color:    domain.Color("#000000"),
							Blocks:   domain.Blocks{},
						}, nil
					},
				}
			},
			wantCode: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Card api.Card `json:"card"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("JSONパースに失敗しました: %v", err)
				}
			},
		},
		{
			name: "カードが見つからない場合はエラーを返す",
			setupMock: func() *service.MockCardService {
				return &service.MockCardService{
					GetCardByGitHubIDFunc: func(ctx context.Context, githubID string, githubClient *github.Client) (*domain.Card, error) {
						return nil, fmt.Errorf("card not found: githubID=%s", githubID)
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
			mockService := tt.setupMock()
			cardHandler := NewCardHandler(mockService)
			router := gin.Default()
			// github_client をコンテキストに設定するミドルウェア
			router.Use(func(c *gin.Context) {
				c.Set("github_client", (*github.Client)(nil))
				c.Next()
			})
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

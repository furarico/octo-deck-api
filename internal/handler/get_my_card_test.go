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

// 自分のカードを取得するテスト
func TestGetMyCard(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name      string
		setupMock func() *service.MockCardService
		wantCode  int
		validate  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常に自分のカードを取得できる",
			setupMock: func() *service.MockCardService {
				return &service.MockCardService{
					GetMyCardFunc: func(ctx context.Context, githubID string, githubClient *github.Client) (*domain.Card, error) {
						return &domain.Card{
							ID:       domain.NewCardID(),
							GithubID: "my_user",
							UserName: "my_user",
							FullName: "My User",
							IconUrl:  "https://example.com/myuser.png",
							Color:    domain.Color("#abcdef"),
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
				if response.UserName != "my_user" {
					t.Errorf("UserNameが違う: 期待=my_user, 実際=%s", response.UserName)
				}
			},
		},
		{
			name: "カードが見つからない場合はエラーを返す",
			setupMock: func() *service.MockCardService {
				return &service.MockCardService{
					GetMyCardFunc: func(ctx context.Context, githubID string, githubClient *github.Client) (*domain.Card, error) {
						return nil, fmt.Errorf("my card not found")
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
			router.Use(func(c *gin.Context) {
				c.Set("github_id", "test_user")
				c.Set("github_client", (*github.Client)(nil))
				c.Next()
			})
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

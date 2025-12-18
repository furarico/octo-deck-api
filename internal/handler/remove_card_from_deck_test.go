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

// カードをデッキから削除するテスト
func TestRemoveCardFromDeck(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name      string
		setupMock func() *service.MockCardService
		githubID  string
		wantCode  int
		validate  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常にカードをデッキから削除できる",
			setupMock: func() *service.MockCardService {
				return &service.MockCardService{
					RemoveCardFromDeckFunc: func(ctx context.Context, collectorGithubID string, targetGithubID string, githubClient *github.Client) (*domain.Card, error) {
						return &domain.Card{
							GithubID: targetGithubID,
							UserName: "target_user",
							FullName: "Target User",
							IconUrl:  "https://example.com/target.png",
							Color:    domain.Color("#000000"),
							Blocks:   domain.Blocks{},
						}, nil
					},
				}
			},
			githubID: "target_user",
			wantCode: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Card api.Card `json:"card"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("JSONパースに失敗しました: %v", err)
				}
				if response.Card.GithubId != "target_user" {
					t.Errorf("GithubIDが違う: 期待=target_user, 実際=%s", response.Card.GithubId)
				}
			},
		},
		{
			name: "カードの削除に失敗した場合",
			setupMock: func() *service.MockCardService {
				return &service.MockCardService{
					RemoveCardFromDeckFunc: func(ctx context.Context, collectorGithubID string, targetGithubID string, githubClient *github.Client) (*domain.Card, error) {
						return nil, fmt.Errorf("card not found in deck")
					},
				}
			},
			githubID: "target_user",
			wantCode: http.StatusInternalServerError,
			validate: nil,
		},
		{
			name: "存在しないカードを削除しようとした場合",
			setupMock: func() *service.MockCardService {
				return &service.MockCardService{
					RemoveCardFromDeckFunc: func(ctx context.Context, collectorGithubID string, targetGithubID string, githubClient *github.Client) (*domain.Card, error) {
						return nil, fmt.Errorf("card does not exist")
					},
				}
			},
			githubID: "non_existent_user",
			wantCode: http.StatusInternalServerError,
			validate: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			mockService := tt.setupMock()
			cardHandler := NewCardHandler(mockService)
			router := gin.Default()
			router.Use(setTestContext)
			strictHandler := api.NewStrictHandler(cardHandler, nil)
			api.RegisterHandlers(router, strictHandler)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", "/cards/"+tt.githubID, nil)
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

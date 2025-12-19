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

// コミュニティからカードを削除するテスト
func TestRemoveCardFromCommunity(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name               string
		setupCardMock      func() *service.MockCardService
		setupCommunityMock func() *service.MockCommunityService
		wantCode           int
		validate           func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常にコミュニティからカードを削除できる",
			setupCardMock: func() *service.MockCardService {
				return &service.MockCardService{
					GetMyCardFunc: func(ctx context.Context, githubID string, githubClient service.GitHubClient) (*domain.Card, error) {
						return &domain.Card{
							ID:       domain.NewCardID(),
							GithubID: "test_user",
							UserName: "test_user",
							FullName: "Test User",
							IconUrl:  "https://example.com/test.png",
							Color:    domain.Color("#000000"),
							Blocks:   domain.Blocks{},
						}, nil
					},
				}
			},
			setupCommunityMock: func() *service.MockCommunityService {
				return &service.MockCommunityService{
					RemoveCardFromCommunityFunc: func(communityID string, cardID string) error {
						return nil
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
				if response.Card.GithubId != "test_user" {
					t.Errorf("GithubIDが違う: 期待=test_user, 実際=%s", response.Card.GithubId)
				}
			},
		},
		{
			name: "カードが見つからない場合",
			setupCardMock: func() *service.MockCardService {
				return &service.MockCardService{
					GetMyCardFunc: func(ctx context.Context, githubID string, githubClient service.GitHubClient) (*domain.Card, error) {
						return nil, fmt.Errorf("card not found")
					},
				}
			},
			setupCommunityMock: func() *service.MockCommunityService {
				return &service.MockCommunityService{}
			},
			wantCode: http.StatusInternalServerError,
			validate: nil,
		},
		{
			name: "コミュニティからのカード削除でエラーが発生した場合",
			setupCardMock: func() *service.MockCardService {
				return &service.MockCardService{
					GetMyCardFunc: func(ctx context.Context, githubID string, githubClient service.GitHubClient) (*domain.Card, error) {
						return &domain.Card{
							ID:       domain.NewCardID(),
							GithubID: "test_user",
							UserName: "test_user",
							FullName: "Test User",
							IconUrl:  "https://example.com/test.png",
							Color:    domain.Color("#000000"),
							Blocks:   domain.Blocks{},
						}, nil
					},
				}
			},
			setupCommunityMock: func() *service.MockCommunityService {
				return &service.MockCommunityService{
					RemoveCardFromCommunityFunc: func(communityID string, cardID string) error {
						return fmt.Errorf("database error")
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
			mockCardService := tt.setupCardMock()
			mockCommunityService := tt.setupCommunityMock()
			handler := NewHandler(mockCardService, mockCommunityService, nil)
			router := gin.Default()
			router.Use(setTestContext)
			strictHandler := api.NewStrictHandler(handler, nil)
			api.RegisterHandlers(router, strictHandler)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", "/communities/test-id/cards", nil)
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

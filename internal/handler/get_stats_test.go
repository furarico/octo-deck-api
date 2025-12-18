package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	api "github.com/furarico/octo-deck-api/generated"
	"github.com/furarico/octo-deck-api/internal/github"
	"github.com/furarico/octo-deck-api/internal/service"
	"github.com/gin-gonic/gin"
)

// ユーザーの統計情報取得のテスト
func TestGetUserStats(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name      string
		setupMock func() *service.MockStatsService
		wantCode  int
		validate  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常にユーザーの統計情報を取得できる",
			setupMock: func() *service.MockStatsService {
				return &service.MockStatsService{
					GetUserStatsFunc: func(ctx context.Context, githubID string, githubClient *github.Client) (*github.ContributionStats, error) {
						return &github.ContributionStats{
							Contributions: []github.Contribution{
								{Date: "2024-01-01", Count: 5},
								{Date: "2024-01-02", Count: 3},
							},
						}, nil
					},
				}
			},
			wantCode: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Stats api.UserStats `json:"stats"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("JSONパースに失敗しました: %v", err)
				}
				if len(response.Stats.Contributions) != 2 {
					t.Errorf("コントリビューション数が違う: 期待=2, 実際=%d", len(response.Stats.Contributions))
				}
			},
		},
		{
			name: "統計情報の取得に失敗した場合",
			setupMock: func() *service.MockStatsService {
				return &service.MockStatsService{
					GetUserStatsFunc: func(ctx context.Context, githubID string, githubClient *github.Client) (*github.ContributionStats, error) {
						return nil, fmt.Errorf("GitHub API error")
					},
				}
			},
			wantCode: http.StatusInternalServerError,
			validate: nil,
		},
		{
			name: "統計情報が空の場合",
			setupMock: func() *service.MockStatsService {
				return &service.MockStatsService{
					GetUserStatsFunc: func(ctx context.Context, githubID string, githubClient *github.Client) (*github.ContributionStats, error) {
						return &github.ContributionStats{
							Contributions: []github.Contribution{},
						}, nil
					},
				}
			},
			wantCode: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Stats api.UserStats `json:"stats"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("JSONパースに失敗しました: %v", err)
				}
				if len(response.Stats.Contributions) != 0 {
					t.Errorf("コントリビューション数が違う: 期待=0, 実際=%d", len(response.Stats.Contributions))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			mockService := tt.setupMock()
			statsHandler := NewStatsHandler(mockService)
			router := gin.Default()
			router.Use(setTestContext)
			strictHandler := api.NewStrictHandler(statsHandler, nil)
			api.RegisterHandlers(router, strictHandler)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/stats/123", nil)
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

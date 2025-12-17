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

// コミュニティ一覧取得のテスト
func TestGetCommunities(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name      string
		setupMock func() *service.MockCommunityService
		wantCode  int
		validate  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常にコミュニティ一覧を取得できる",
			setupMock: func() *service.MockCommunityService {
				return &service.MockCommunityService{
					GetAllCommunitiesFunc: func(githubID string) ([]domain.Community, error) {
						return []domain.Community{
							{
								ID:   domain.NewCommunityID(),
								Name: "Test Community 1",
							},
							{
								ID:   domain.NewCommunityID(),
								Name: "Test Community 2",
							},
						}, nil
					},
				}
			},
			wantCode: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Communities []api.Community `json:"communities"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("JSONパースに失敗しました: %v", err)
				}
				if len(response.Communities) != 2 {
					t.Errorf("コミュニティの数が違う: 期待=2, 実際=%d", len(response.Communities))
				}
			},
		},
		{
			name: "コミュニティが空の場合",
			setupMock: func() *service.MockCommunityService {
				return &service.MockCommunityService{
					GetAllCommunitiesFunc: func(githubID string) ([]domain.Community, error) {
						return []domain.Community{}, nil
					},
				}
			},
			wantCode: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Communities []api.Community `json:"communities"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("JSONパースに失敗しました: %v", err)
				}
				if len(response.Communities) != 0 {
					t.Errorf("コミュニティの数が違う: 期待=0, 実際=%d", len(response.Communities))
				}
			},
		},
		{
			name: "サービスでエラーが発生した場合",
			setupMock: func() *service.MockCommunityService {
				return &service.MockCommunityService{
					GetAllCommunitiesFunc: func(githubID string) ([]domain.Community, error) {
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
			req, _ := http.NewRequest("GET", "/communities", nil)
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

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

// 指定したコミュニティ取得のテスト
func TestGetCommunity(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name      string
		setupMock func() *service.MockCommunityService
		wantCode  int
		validate  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常にコミュニティを取得できる",
			setupMock: func() *service.MockCommunityService {
				return &service.MockCommunityService{
					GetCommunityWithHighlightedCardFunc: func(ctx context.Context, id string) (*domain.Community, *domain.HighlightedCard, error) {
						community := &domain.Community{
							ID:   domain.NewCommunityID(),
							Name: "Test Community",
						}
						highlightedCard := &domain.HighlightedCard{}
						return community, highlightedCard, nil
					},
				}
			},
			wantCode: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Community       api.Community       `json:"community"`
					HighlightedCard api.HighlightedCard `json:"highlightedCard"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("JSONパースに失敗しました: %v", err)
				}
				if response.Community.Name != "Test Community" {
					t.Errorf("コミュニティ名が違う: 期待=Test Community, 実際=%s", response.Community.Name)
				}
			},
		},
		{
			name: "コミュニティが見つからない場合",
			setupMock: func() *service.MockCommunityService {
				return &service.MockCommunityService{
					GetCommunityWithHighlightedCardFunc: func(ctx context.Context, id string) (*domain.Community, *domain.HighlightedCard, error) {
						return nil, nil, fmt.Errorf("community not found: id=%s", id)
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
			req, _ := http.NewRequest("GET", "/communities/test-id", nil)
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

package handler

import (
	"bytes"
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

// コミュニティ作成のテスト
func TestCreateCommunity(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name      string
		setupMock func() *service.MockCommunityService
		body      string
		wantCode  int
		validate  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常にコミュニティを作成できる",
			setupMock: func() *service.MockCommunityService {
				return &service.MockCommunityService{
					CreateCommunityFunc: func(name string) (*domain.Community, error) {
						return &domain.Community{
							ID:   domain.NewCommunityID(),
							Name: name,
						}, nil
					},
				}
			},
			body:     `"New Community"`,
			wantCode: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Community api.Community `json:"community"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("JSONパースに失敗しました: %v", err)
				}
				// レスポンスにコミュニティが含まれていることを確認
				if response.Community.Name == "" {
					t.Errorf("コミュニティ名が空です")
				}
			},
		},
		{
			name: "空の名前でエラーを返す",
			setupMock: func() *service.MockCommunityService {
				return &service.MockCommunityService{}
			},
			body:     `""`,
			wantCode: http.StatusInternalServerError,
			validate: nil,
		},
		{
			name: "サービスでエラーが発生した場合",
			setupMock: func() *service.MockCommunityService {
				return &service.MockCommunityService{
					CreateCommunityFunc: func(name string) (*domain.Community, error) {
						return nil, fmt.Errorf("database error")
					},
				}
			},
			body:     `"Test Community"`,
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
			req, _ := http.NewRequest("POST", "/communities", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
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

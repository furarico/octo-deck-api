package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	api "github.com/furarico/octo-deck-api/generated"
	"github.com/furarico/octo-deck-api/internal/domain"
	"github.com/furarico/octo-deck-api/internal/service"
	"github.com/gin-gonic/gin"
)

// コミュニティのHighlightedCardを更新するテスト
func TestRefreshCommunity(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name      string
		setupMock func() *service.MockCommunityService
		wantCode  int
		validate  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常にコミュニティのHighlightedCardを更新できる",
			setupMock: func() *service.MockCommunityService {
				return &service.MockCommunityService{
					RefreshHighlightedCardFunc: func(ctx context.Context, id string, githubClient service.GitHubClient) (*domain.Community, *domain.HighlightedCard, error) {
						community := &domain.Community{
							ID:        domain.NewCommunityID(),
							Name:      "Test Community",
							StartedAt: time.Now(),
							EndedAt:   time.Now().Add(24 * time.Hour),
						}
						highlightedCard := &domain.HighlightedCard{
							BestContributor: domain.Card{
								ID:       domain.NewCardID(),
								GithubID: "contributor1",
								UserName: "contributor1",
								FullName: "Contributor One",
								IconUrl:  "https://example.com/contributor1.png",
								Color:    domain.Color("#111111"),
								Blocks:   domain.Blocks{},
								MostUsedLanguage: domain.Language{
									LanguageName: "Go",
									Color:        "#00ADD8",
								},
							},
							BestCommitter: domain.Card{
								ID:       domain.NewCardID(),
								GithubID: "committer1",
								UserName: "committer1",
								FullName: "Committer One",
								IconUrl:  "https://example.com/committer1.png",
								Color:    domain.Color("#222222"),
								Blocks:   domain.Blocks{},
								MostUsedLanguage: domain.Language{
									LanguageName: "TypeScript",
									Color:        "#3178C6",
								},
							},
							BestIssuer: domain.Card{
								ID:       domain.NewCardID(),
								GithubID: "issuer1",
								UserName: "issuer1",
								FullName: "Issuer One",
								IconUrl:  "https://example.com/issuer1.png",
								Color:    domain.Color("#333333"),
								Blocks:   domain.Blocks{},
								MostUsedLanguage: domain.Language{
									LanguageName: "Python",
									Color:        "#3776AB",
								},
							},
							BestPullRequester: domain.Card{
								ID:       domain.NewCardID(),
								GithubID: "pr1",
								UserName: "pr1",
								FullName: "PR One",
								IconUrl:  "https://example.com/pr1.png",
								Color:    domain.Color("#444444"),
								Blocks:   domain.Blocks{},
								MostUsedLanguage: domain.Language{
									LanguageName: "JavaScript",
									Color:        "#F7DF1E",
								},
							},
							BestReviewer: domain.Card{
								ID:       domain.NewCardID(),
								GithubID: "reviewer1",
								UserName: "reviewer1",
								FullName: "Reviewer One",
								IconUrl:  "https://example.com/reviewer1.png",
								Color:    domain.Color("#555555"),
								Blocks:   domain.Blocks{},
								MostUsedLanguage: domain.Language{
									LanguageName: "Rust",
									Color:        "#000000",
								},
							},
						}
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
				if response.HighlightedCard.BestContributor.GithubId != "contributor1" {
					t.Errorf("BestContributorのGithubIDが違う: 期待=contributor1, 実際=%s", response.HighlightedCard.BestContributor.GithubId)
				}
				if response.HighlightedCard.BestCommitter.GithubId != "committer1" {
					t.Errorf("BestCommitterのGithubIDが違う: 期待=committer1, 実際=%s", response.HighlightedCard.BestCommitter.GithubId)
				}
			},
		},
		{
			name: "RefreshHighlightedCardでエラーが発生した場合",
			setupMock: func() *service.MockCommunityService {
				return &service.MockCommunityService{
					RefreshHighlightedCardFunc: func(ctx context.Context, id string, githubClient service.GitHubClient) (*domain.Community, *domain.HighlightedCard, error) {
						return nil, nil, fmt.Errorf("failed to refresh highlighted card: id=%s", id)
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
			req, _ := http.NewRequest("PUT", "/communities/test-id/refresh", nil)
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

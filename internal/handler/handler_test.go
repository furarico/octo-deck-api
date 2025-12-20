package handler

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/furarico/octo-deck-api/internal/github"
	"github.com/gin-gonic/gin"
)

// gin.Contextからcontext.Contextを取得するテスト
func TestGetRequestContext(t *testing.T) {
	tests := []struct {
		name string
		ctx  context.Context
		want context.Context
	}{
		{
			name: "通常のcontext.Contextをそのまま返す",
			ctx:  context.Background(),
			want: context.Background(),
		},
		{
			name: "gin.ContextからRequest.Contextを取得する",
			ctx: func() context.Context {
				gin.SetMode(gin.TestMode)
				c, _ := gin.CreateTestContext(nil)
				// Requestを設定する必要がある
				req := httptest.NewRequest("GET", "/", nil)
				c.Request = req
				return c
			}(),
			want: nil, // gin.Contextの場合はRequest.Context()を返すが、値は異なる可能性がある
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getRequestContext(tt.ctx)
			// gin.Contextの場合はRequest.Context()を返すが、値は異なる可能性がある
			// 通常のcontext.Contextの場合はそのまま返される
			if tt.want == nil {
				// gin.Contextの場合はRequest.Context()が返されることを確認
				if got == nil {
					t.Errorf("getRequestContext() = nil, want non-nil")
				}
			} else {
				if got != tt.want {
					t.Errorf("getRequestContext() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

// context.ContextからGitHub Clientを取得するテスト
func TestGetGitHubClient(t *testing.T) {
	tests := []struct {
		name       string
		setup      func() context.Context
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "正常にGitHub Clientを取得できる",
			setup: func() context.Context {
				ctx := context.Background()
				mockClient := github.NewMockClient()
				return context.WithValue(ctx, GitHubClientKey, mockClient)
			},
			wantErr: false,
		},
		{
			name: "contextにGitHub Clientが設定されていない場合はエラー",
			setup: func() context.Context {
				return context.Background()
			},
			wantErr:    true,
			wantErrMsg: "github_client not found in context",
		},
		{
			name: "型が一致しない場合はエラー",
			setup: func() context.Context {
				ctx := context.Background()
				return context.WithValue(ctx, GitHubClientKey, "invalid_type")
			},
			wantErr:    true,
			wantErrMsg: "github_client not found in context",
		},
		{
			name: "gin.Contextから正常に取得できる",
			setup: func() context.Context {
				gin.SetMode(gin.TestMode)
				c, _ := gin.CreateTestContext(nil)
				req := httptest.NewRequest("GET", "/", nil)
				mockClient := github.NewMockClient()
				ctx := context.WithValue(req.Context(), GitHubClientKey, mockClient)
				c.Request = req.WithContext(ctx)
				return c
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setup()
			got, err := getGitHubClient(ctx)
			if tt.wantErr {
				if err == nil {
					t.Errorf("getGitHubClient() error = nil, want error")
					return
				}
				if tt.wantErrMsg != "" && err.Error() != tt.wantErrMsg {
					t.Errorf("getGitHubClient() error = %v, want %v", err.Error(), tt.wantErrMsg)
				}
			} else {
				if err != nil {
					t.Errorf("getGitHubClient() error = %v, want nil", err)
					return
				}
				if got == nil {
					t.Errorf("getGitHubClient() = nil, want non-nil")
				}
			}
		})
	}
}

// context.ContextからGitHub IDを取得するテスト
func TestGetGitHubID(t *testing.T) {
	tests := []struct {
		name       string
		setup      func() context.Context
		wantErr    bool
		wantErrMsg string
		wantID     string
	}{
		{
			name: "正常にGitHub IDを取得できる",
			setup: func() context.Context {
				ctx := context.Background()
				return context.WithValue(ctx, GitHubIDKey, "12345")
			},
			wantErr: false,
			wantID:  "12345",
		},
		{
			name: "contextにGitHub IDが設定されていない場合はエラー",
			setup: func() context.Context {
				return context.Background()
			},
			wantErr:    true,
			wantErrMsg: "github_id not found in context",
		},
		{
			name: "空文字列の場合はエラー",
			setup: func() context.Context {
				ctx := context.Background()
				return context.WithValue(ctx, GitHubIDKey, "")
			},
			wantErr:    true,
			wantErrMsg: "github_id not found in context",
		},
		{
			name: "型が一致しない場合はエラー",
			setup: func() context.Context {
				ctx := context.Background()
				return context.WithValue(ctx, GitHubIDKey, 12345)
			},
			wantErr:    true,
			wantErrMsg: "github_id not found in context",
		},
		{
			name: "gin.Contextから正常に取得できる",
			setup: func() context.Context {
				gin.SetMode(gin.TestMode)
				c, _ := gin.CreateTestContext(nil)
				req := httptest.NewRequest("GET", "/", nil)
				ctx := context.WithValue(req.Context(), GitHubIDKey, "12345")
				c.Request = req.WithContext(ctx)
				return c
			},
			wantErr: false,
			wantID:  "12345",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setup()
			got, err := getGitHubID(ctx)
			if tt.wantErr {
				if err == nil {
					t.Errorf("getGitHubID() error = nil, want error")
					return
				}
				if tt.wantErrMsg != "" && err.Error() != tt.wantErrMsg {
					t.Errorf("getGitHubID() error = %v, want %v", err.Error(), tt.wantErrMsg)
				}
			} else {
				if err != nil {
					t.Errorf("getGitHubID() error = %v, want nil", err)
					return
				}
				if got != tt.wantID {
					t.Errorf("getGitHubID() = %v, want %v", got, tt.wantID)
				}
			}
		})
	}
}

// context.ContextからGitHub NodeIDを取得するテスト
func TestGetGitHubNodeID(t *testing.T) {
	tests := []struct {
		name       string
		setup      func() context.Context
		wantErr    bool
		wantErrMsg string
		wantNodeID string
	}{
		{
			name: "正常にGitHub NodeIDを取得できる",
			setup: func() context.Context {
				ctx := context.Background()
				return context.WithValue(ctx, GitHubNodeIDKey, "MDQ6VXNlcjEyMzQ1")
			},
			wantErr:    false,
			wantNodeID: "MDQ6VXNlcjEyMzQ1",
		},
		{
			name: "contextにGitHub NodeIDが設定されていない場合はエラー",
			setup: func() context.Context {
				return context.Background()
			},
			wantErr:    true,
			wantErrMsg: "github_node_id not found in context",
		},
		{
			name: "空文字列の場合はエラー",
			setup: func() context.Context {
				ctx := context.Background()
				return context.WithValue(ctx, GitHubNodeIDKey, "")
			},
			wantErr:    true,
			wantErrMsg: "github_node_id not found in context",
		},
		{
			name: "型が一致しない場合はエラー",
			setup: func() context.Context {
				ctx := context.Background()
				return context.WithValue(ctx, GitHubNodeIDKey, 12345)
			},
			wantErr:    true,
			wantErrMsg: "github_node_id not found in context",
		},
		{
			name: "gin.Contextから正常に取得できる",
			setup: func() context.Context {
				gin.SetMode(gin.TestMode)
				c, _ := gin.CreateTestContext(nil)
				req := httptest.NewRequest("GET", "/", nil)
				ctx := context.WithValue(req.Context(), GitHubNodeIDKey, "MDQ6VXNlcjEyMzQ1")
				c.Request = req.WithContext(ctx)
				return c
			},
			wantErr:    false,
			wantNodeID: "MDQ6VXNlcjEyMzQ1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setup()
			got, err := getGitHubNodeID(ctx)
			if tt.wantErr {
				if err == nil {
					t.Errorf("getGitHubNodeID() error = nil, want error")
					return
				}
				if tt.wantErrMsg != "" && err.Error() != tt.wantErrMsg {
					t.Errorf("getGitHubNodeID() error = %v, want %v", err.Error(), tt.wantErrMsg)
				}
			} else {
				if err != nil {
					t.Errorf("getGitHubNodeID() error = %v, want nil", err)
					return
				}
				if got != tt.wantNodeID {
					t.Errorf("getGitHubNodeID() = %v, want %v", got, tt.wantNodeID)
				}
			}
		})
	}
}

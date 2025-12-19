package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/furarico/octo-deck-api/internal/github"
	"github.com/furarico/octo-deck-api/internal/handler"
	"github.com/gin-gonic/gin"
)

// GitHub Appのアクセストークンを検証し、ユーザー情報をContextにセット
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format, expected 'Bearer <token>'"})
			c.Abort()
			return
		}

		// GitHub APIでユーザー情報を取得
		ghClient := github.NewClient(token)
		user, err := ghClient.GetAuthenticatedUser(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// context.Context にユーザー情報とClientをセット
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, handler.GitHubClientKey, ghClient)
		ctx = context.WithValue(ctx, handler.GitHubIDKey, strconv.FormatInt(user.ID, 10))
		ctx = context.WithValue(ctx, handler.GitHubNodeIDKey, user.NodeID)
		ctx = context.WithValue(ctx, handler.GitHubLoginKey, user.Login)
		ctx = context.WithValue(ctx, handler.GitHubNameKey, user.Name)
		ctx = context.WithValue(ctx, handler.GitHubAvatarKey, user.AvatarURL)

		// 新しいcontextでリクエストを更新
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

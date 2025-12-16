package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/furarico/octo-deck-api/internal/github"
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

		// Contextにユーザー情報をセット
		c.Set("github_id", strconv.FormatInt(user.ID, 10))
		c.Set("github_token", token)
		c.Set("github_login", user.Login)
		c.Set("github_name", user.Name)
		c.Set("github_avatar_url", user.AvatarURL)

		c.Next()
	}
}

package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v80/github"
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
		client := github.NewClient(nil).WithAuthToken(token)
		user, _, err := client.Users.Get(c.Request.Context(), "")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Contextにユーザー情報をセット
		c.Set("github_id", strconv.FormatInt(user.GetID(), 10))
		c.Set("github_token", token)
		c.Set("github_login", user.GetLogin())
		c.Set("github_name", user.GetName())
		c.Set("github_avatar_url", user.GetAvatarURL())

		c.Next()
	}
}

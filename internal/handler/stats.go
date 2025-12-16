package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 自分の統計情報取得
// (GET /stats/me)
func (h *Handler) GetMyStats(c *gin.Context) {
	// TODO: 実装
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "not implemented",
	})
}

// ユーザーの統計情報取得
// (GET /stats/{githubId})
func (h *Handler) GetUserStats(c *gin.Context, githubId string) {
	// TODO: 実装
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "not implemented",
	})
}

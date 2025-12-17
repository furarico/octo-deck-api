package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// コミュニティを作成
// (POST /communities)
func (h *Handler) CreateCommunity(c *gin.Context) {
	// TODO: 実装
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "not implemented",
	})
}

// コミュニティを削除
// (DELETE /communities/{id})
func (h *Handler) DeleteCommunity(c *gin.Context, id string) {
	// TODO: 実装
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "not implemented",
	})
}

// 指定したコミュニティの自分のカードを削除
// (DELETE /communities/{id}/cards)
func (h *Handler) RemoveCardFromCommunity(c *gin.Context, id string) {
	// TODO: 実装
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "not implemented",
	})
}

// 指定したコミュニティに自分のカードを追加
// (POST /communities/{id}/cards)
func (h *Handler) AddCardToCommunity(c *gin.Context, id string) {
	// TODO: 実装
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "not implemented",
	})
}

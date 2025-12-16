package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// コミュニティ一覧取得
// (GET /communities)
func (h *Handler) GetCommunities(c *gin.Context) {
	// TODO: 実装
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "not implemented",
	})
}

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

// 指定したコミュニティ取得
// (GET /communities/{id})
func (h *Handler) GetCommunity(c *gin.Context, id string) {
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

// 指定したコミュニティのカード一覧取得
// (GET /communities/{id}/cards)
func (h *Handler) GetCommunityCards(c *gin.Context, id string) {
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

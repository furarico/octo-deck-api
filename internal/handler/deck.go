package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// カードをデッキに追加
// (POST /cards)
func (h *Handler) AddCardToDeck(c *gin.Context) {
	// TODO: 実装
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "not implemented",
	})
}

// カードをデッキから削除
// (DELETE /cards/{githubId})
func (h *Handler) RemoveCardFromDeck(c *gin.Context, githubId string) {
	// TODO: 実装
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "not implemented",
	})
}

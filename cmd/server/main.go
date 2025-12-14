package main

import (
	"net/http"

	api "github.com/furarico/octo-deck-api/generated"
	"github.com/gin-gonic/gin"
)

type Server struct{}

// GetCards はカード一覧を返すハンドラー
func (s *Server) GetCards(c *gin.Context) {
	// モックデータを返す
	cards := []api.Card{
		{
			Id:       "1",
			UserName: "taro_yamada",
			FullName: "山田太郎",
			IconUrl:  "https://example.com/icons/taro.png",
		},
		{
			Id:       "2",
			UserName: "hanako_tanaka",
			FullName: "田中花子",
			IconUrl:  "https://example.com/icons/hanako.png",
		},
		{
			Id:       "3",
			UserName: "jiro_sato",
			FullName: "佐藤次郎",
			IconUrl:  "https://example.com/icons/jiro.png",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"cards": cards,
	})
}

// GetMyCard は自分のカードを返すハンドラー（未実装）
func (s *Server) GetMyCard(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Not implemented yet",
	})
}

// GetCard は指定されたIDのカードを返すハンドラー（未実装）
func (s *Server) GetCard(c *gin.Context, id string) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Not implemented yet",
	})
}

func main() {
	r := gin.Default()

	server := &Server{}
	api.RegisterHandlers(r, server)

	r.Run(":8080")
}

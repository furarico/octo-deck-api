package handler

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// コミュニティを作成
// (POST /communities)
func (h *Handler) CreateCommunity(c *gin.Context) {
	// リクエストボディからコミュニティ名を取得
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read request body",
		})
		return
	}

	name := string(body)
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "community name is required",
		})
		return
	}

	community, err := h.communityService.CreateCommunity(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, convertCommunityToAPI(*community))
}

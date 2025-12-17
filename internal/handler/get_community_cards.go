package handler

import (
	"net/http"

	api "github.com/furarico/octo-deck-api/generated"
	"github.com/gin-gonic/gin"
)

// 指定したコミュニティのカード一覧取得
// (GET /communities/{id}/cards)
func (h *Handler) GetCommunityCards(c *gin.Context, id string) {
	cards, err := h.communityService.GetCommunityCards(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	cardsAPI := make([]api.Card, len(cards))
	for i, card := range cards {
		cardsAPI[i] = convertCardToAPI(card)
	}

	c.JSON(http.StatusOK, cardsAPI)
}

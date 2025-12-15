package handler

import (
	"net/http"

	api "github.com/furarico/octo-deck-api/generated"
	"github.com/gin-gonic/gin"
)

// (GET /cards)
func (h *Handler) GetCards(c *gin.Context) {
	// TODO: 認証情報からGitHubIDを取得する
	githubID := c.GetString("github_id")

	cards, err := h.cardService.GetAllCards(githubID)
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

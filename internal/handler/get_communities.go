package handler

import (
	"net/http"

	api "github.com/furarico/octo-deck-api/generated"
	"github.com/gin-gonic/gin"
)

// コミュニティ一覧取得
// (GET /communities)
func (h *Handler) GetCommunities(c *gin.Context) {
	githubID := c.GetString("github_id")
	if githubID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "github_id is missing from context",
		})
		return
	}

	communities, err := h.communityService.GetAllCommunities(githubID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	communitiesAPI := make([]api.Community, len(communities))
	for i, community := range communities {
		communitiesAPI[i] = convertCommunityToAPI(community)
	}

	c.JSON(http.StatusOK, communitiesAPI)
}

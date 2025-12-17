package handler

import (
	"github.com/furarico/octo-deck-api/internal/service"
)

type Handler struct {
	cardService      *service.CardService
	communityService *service.CommunityService
}

// ふつうに依存を注入するコンストラクタ
func NewHandler(cardService *service.CardService, communityService *service.CommunityService) *Handler {
	return &Handler{
		cardService:      cardService,
		communityService: communityService,
	}
}

// （任意）カードだけで組み立てたいケース用の簡易コンストラクタ
func NewCardHandler(cardService *service.CardService) *Handler {
	return &Handler{cardService: cardService}
}

// （任意）コミュニティだけで組み立てたいケース用の簡易コンストラクタ
func NewCommunityHandler(communityService *service.CommunityService) *Handler {
	return &Handler{communityService: communityService}
}

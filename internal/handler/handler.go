package handler

import (
	"github.com/furarico/octo-deck-api/internal/service"
)

type Handler struct {
	cardService *service.CardService
}

func NewHandler(cardService *service.CardService) *Handler {
	return &Handler{
		cardService: cardService,
	}
}

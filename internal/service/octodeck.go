package service

import (
	"fmt"

	api "github.com/furarico/octo-deck-api/generated"
	"github.com/furarico/octo-deck-api/internal/repository"
)

type CardService struct {
	cardRepo repository.CardRepository
}

func NewCardService(cardRepo repository.CardRepository) *CardService {
	return &CardService{
		cardRepo: cardRepo,
	}
}

// GetAllCards は全てのカードを取得する
func (s *CardService) GetAllCards() ([]api.Card, error) {
	cards, err := s.cardRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all cards: %w", err)
	}

	return cards, nil
}

// GetCardByID は指定されたIDのカードを取得する
func (s *CardService) GetCardByID(id string) (*api.Card, error) {
	card, err := s.cardRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get card by id: %w", err)
	}

	if card == nil {
		return nil, fmt.Errorf("card not found: id=%s", id)
	}

	return card, nil
}

// GetMyCard は自分のカードを取得する
func (s *CardService) GetMyCard() (*api.Card, error) {
	card, err := s.cardRepo.FindMyCard()
	if err != nil {
		return nil, fmt.Errorf("failed to get my card: %w", err)
	}

	if card == nil {
		return nil, fmt.Errorf("my card not found")
	}

	return card, nil
}

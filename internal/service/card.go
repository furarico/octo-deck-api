package service

import (
	"fmt"

	"github.com/furarico/octo-deck-api/internal/domain"
)

// CardRepository はServiceが必要とするRepositoryのインターフェース
type CardRepository interface {
	FindAll() ([]domain.CardWithOwner, error)
	FindByID(id string) (*domain.CardWithOwner, error)
	FindMyCard() (*domain.CardWithOwner, error)
}

type CardService struct {
	cardRepo CardRepository
}

func NewCardService(cardRepo CardRepository) *CardService {
	return &CardService{
		cardRepo: cardRepo,
	}
}

// GetAllCards は全てのカードを取得する
func (s *CardService) GetAllCards() ([]domain.CardWithOwner, error) {
	cards, err := s.cardRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all cards: %w", err)
	}

	return cards, nil
}

// GetCardByID は指定されたIDのカードを取得する
func (s *CardService) GetCardByID(id string) (*domain.CardWithOwner, error) {
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
func (s *CardService) GetMyCard() (*domain.CardWithOwner, error) {
	card, err := s.cardRepo.FindMyCard()
	if err != nil {
		return nil, fmt.Errorf("failed to get my card: %w", err)
	}

	if card == nil {
		return nil, fmt.Errorf("my card not found")
	}

	return card, nil
}

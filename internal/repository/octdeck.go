package repository

import (
	api "github.com/furarico/octo-deck-api/generated"
)

// MockCardRepository は service.CardRepository インターフェースを実装する
type MockCardRepository struct{}

func NewMockCardRepository() *MockCardRepository {
	return &MockCardRepository{}
}

// getMockCards はモックカードデータを返す
func (r *MockCardRepository) getMockCards() []api.Card {
	return []api.Card{
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
}

// FindAll は全てのカードを取得する
func (r *MockCardRepository) FindAll() ([]api.Card, error) {
	cards := r.getMockCards()

	if len(cards) == 0 {
		return nil, nil
	}

	return cards, nil
}

// FindByID は指定されたIDのカードを取得する
func (r *MockCardRepository) FindByID(id string) (*api.Card, error) {
	cards := r.getMockCards()

	for _, card := range cards {
		if card.Id == id {
			return &card, nil
		}
	}

	return nil, nil
}

// FindMyCard は自分のカードを取得する
func (r *MockCardRepository) FindMyCard() (*api.Card, error) {
	cards := r.getMockCards()

	if len(cards) == 0 {
		return nil, nil
	}

	return &cards[0], nil
}

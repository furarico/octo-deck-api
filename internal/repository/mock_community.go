package repository

import "github.com/furarico/octo-deck-api/internal/domain"

type MockCommunityRepository struct {
	FindAllFunc     func(githubID string) ([]domain.Community, error)
	FindByIDFunc    func(id string) (*domain.Community, error)
	FindCardsFunc   func(id string) ([]domain.Card, error)
	CreateFunc      func(community *domain.Community) error
	DeleteFunc      func(id string) error
	AddCardFunc     func(communityID string, cardID string) error
	RemoveCardFunc  func(communityID string, cardID string) error
}

func NewMockCommunityRepository() *MockCommunityRepository {
	return &MockCommunityRepository{}
}

// FindAll はすべてのコミュニティを取得する
func (r *MockCommunityRepository) FindAll(githubID string) ([]domain.Community, error) {
	if r.FindAllFunc != nil {
		return r.FindAllFunc(githubID)
	}
	return []domain.Community{}, nil
}

// FindByID は指定されたIDのコミュニティを取得する
func (r *MockCommunityRepository) FindByID(id string) (*domain.Community, error) {
	if r.FindByIDFunc != nil {
		return r.FindByIDFunc(id)
	}
	return nil, nil
}

// FindCards は指定したコミュニティIDのカード一覧を取得する
func (r *MockCommunityRepository) FindCards(id string) ([]domain.Card, error) {
	if r.FindCardsFunc != nil {
		return r.FindCardsFunc(id)
	}
	return []domain.Card{}, nil
}

// Create はコミュニティを作成する
func (r *MockCommunityRepository) Create(community *domain.Community) error {
	if r.CreateFunc != nil {
		return r.CreateFunc(community)
	}
	return nil
}

// Delete はコミュニティを削除する
func (r *MockCommunityRepository) Delete(id string) error {
	if r.DeleteFunc != nil {
		return r.DeleteFunc(id)
	}
	return nil
}

// AddCard はコミュニティにカードを追加する
func (r *MockCommunityRepository) AddCard(communityID string, cardID string) error {
	if r.AddCardFunc != nil {
		return r.AddCardFunc(communityID, cardID)
	}
	return nil
}

// RemoveCard はコミュニティからカードを削除する
func (r *MockCommunityRepository) RemoveCard(communityID string, cardID string) error {
	if r.RemoveCardFunc != nil {
		return r.RemoveCardFunc(communityID, cardID)
	}
	return nil
}

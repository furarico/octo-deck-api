package domain

// CardWithOwner はカードと所有者情報を組み合わせた集約
type CardWithOwner struct {
	Card  *Card
	Owner *User
}

func NewCardWithOwner(card *Card, owner *User) *CardWithOwner {
	return &CardWithOwner{
		Card:  card,
		Owner: owner,
	}
}

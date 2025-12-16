package database

import (
	"time"

	"github.com/furarico/octo-deck-api/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserName  string    `gorm:"not null"`
	FullName  string    `gorm:"not null"`
	GithubID  string    `gorm:"not null"`
	IconURL   string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	Identicon Identicon `gorm:"foreignKey:UserID"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (u *User) ToDomain() *domain.User {
	return &domain.User{
		ID:        domain.UserID(u.ID),
		UserName:  u.UserName,
		FullName:  u.FullName,
		GitHubID:  u.GithubID,
		IconURL:   u.IconURL,
		Identicon: *u.Identicon.ToDomain(),
	}
}

package domain

import "github.com/google/uuid"

type UserID uuid.UUID

func NewUserID() UserID {
	return UserID(uuid.New())
}

type User struct {
	ID        UserID
	UserName  string
	FullName  string
	GitHubID  string
	IconURL   string
	Identicon Identicon
}

func NewUser(userName, fullName, githubID, iconURL string) *User {
	return &User{
		ID:       NewUserID(),
		UserName: userName,
		FullName: fullName,
		GitHubID: githubID,
		IconURL:  iconURL,
	}
}

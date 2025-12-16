package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v80/github"
)

type UserInfo struct {
	Login     string
	Name      string
	AvatarURL string
}

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) GetUserByID(ctx context.Context, token string, id int64) (*UserInfo, error) {
	client := github.NewClient(nil).WithAuthToken(token)
	user, _, err := client.Users.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	name := user.GetLogin()
	if user.Name != nil && *user.Name != "" {
		name = *user.Name
	}

	return &UserInfo{
		Login:     user.GetLogin(),
		Name:      name,
		AvatarURL: user.GetAvatarURL(),
	}, nil
}

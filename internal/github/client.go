package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v80/github"
)

type UserInfo struct {
	ID        int64
	Login     string
	Name      string
	AvatarURL string
}

type Client struct {
	client *github.Client
}

// トークンで認証されたGitHub API Clientを生成する
func NewClient(token string) *Client {
	return &Client{
		client: github.NewClient(nil).WithAuthToken(token),
	}
}

// 認証されたユーザー自身の情報を取得する
func (c *Client) GetAuthenticatedUser(ctx context.Context) (*UserInfo, error) {
	user, _, err := c.client.Users.Get(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get authenticated user: %w", err)
	}

	return c.toUserInfo(user), nil
}

// 指定されたIDのユーザー情報を取得する
func (c *Client) GetUserByID(ctx context.Context, id int64) (*UserInfo, error) {
	user, _, err := c.client.Users.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return c.toUserInfo(user), nil
}

// github.Userをドメインオブジェクトに変換する
func (c *Client) toUserInfo(user *github.User) *UserInfo {
	name := user.GetLogin()
	if user.Name != nil && *user.Name != "" {
		name = *user.Name
	}

	return &UserInfo{
		ID:        user.GetID(),
		Login:     user.GetLogin(),
		Name:      name,
		AvatarURL: user.GetAvatarURL(),
	}
}

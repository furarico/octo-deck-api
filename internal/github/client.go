package github

import (
	"github.com/google/go-github/v80/github"
)

type Client struct {
	client *github.Client
}

// トークンで認証されたGitHub API Clientを生成する
func NewClient(token string) *Client {
	return &Client{
		client: github.NewClient(nil).WithAuthToken(token),
	}
}

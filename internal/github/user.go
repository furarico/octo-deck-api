package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v80/github"
)

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

// GetUsersByIDs は複数のGitHub IDからユーザー情報を一括取得する（GraphQL使用）
func (c *Client) GetUsersByIDs(ctx context.Context, ids []int64) (map[int64]*UserInfo, error) {
	if len(ids) == 0 {
		return make(map[int64]*UserInfo), nil
	}

	// まずREST APIで各IDのログイン名を取得（GraphQLはIDでの直接検索をサポートしていないため）
	// 並列処理で高速化
	type result struct {
		id   int64
		info *UserInfo
		err  error
	}

	results := make(chan result, len(ids))

	for _, id := range ids {
		go func(id int64) {
			info, err := c.GetUserByID(ctx, id)
			results <- result{id: id, info: info, err: err}
		}(id)
	}

	userMap := make(map[int64]*UserInfo)
	var firstErr error

	for range ids {
		r := <-results
		if r.err != nil {
			if firstErr == nil {
				firstErr = r.err
			}
			continue
		}
		userMap[r.id] = r.info
	}

	if firstErr != nil && len(userMap) == 0 {
		return nil, firstErr
	}

	return userMap, nil
}

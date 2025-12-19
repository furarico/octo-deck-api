package github

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/go-github/v80/github"
)

const (
	// maxConcurrentRequests は並列API呼び出しの最大同時実行数
	maxConcurrentRequests = 100
	// defaultLanguageColor は言語が不明な場合のデフォルト色
	defaultLanguageColor = "#586069"
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

// GetUsersByIDs は複数のGitHub IDからユーザー情報を一括取得する
// 並列処理で高速化しつつ、同時実行数を制限してレート制限を回避する
func (c *Client) GetUsersByIDs(ctx context.Context, ids []int64) (map[int64]*UserInfo, error) {
	if len(ids) == 0 {
		return make(map[int64]*UserInfo), nil
	}

	type result struct {
		id   int64
		info *UserInfo
		err  error
	}

	results := make(chan result, len(ids))
	sem := make(chan struct{}, maxConcurrentRequests) // 同時実行数を制限

	var wg sync.WaitGroup
	for _, id := range ids {
		wg.Add(1)
		go func(id int64) {
			defer wg.Done()

			// コンテキストがキャンセルされていたら早期リターン
			select {
			case <-ctx.Done():
				results <- result{id: id, err: ctx.Err()}
				return
			case sem <- struct{}{}: // セマフォを取得
				defer func() { <-sem }()
			}

			info, err := c.GetUserByID(ctx, id)
			results <- result{id: id, info: info, err: err}
		}(id)
	}

	// 全てのgoroutineが完了したらチャネルを閉じる
	go func() {
		wg.Wait()
		close(results)
	}()

	userMap := make(map[int64]*UserInfo)
	var firstErr error

	for r := range results {
		if r.err != nil {
			if firstErr == nil {
				firstErr = r.err
			}
			continue
		}
		userMap[r.id] = r.info
	}

	if firstErr != nil && len(userMap) == 0 {
		return nil, fmt.Errorf("failed to get users by IDs: %w", firstErr)
	}

	return userMap, nil
}

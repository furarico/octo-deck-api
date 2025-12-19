package github

import (
	"context"
	"fmt"
	"strings"

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

// GetUsersByIDs は複数のユーザー情報をバッチで取得する（GraphQL使用）
func (c *Client) GetUsersByIDs(ctx context.Context, ids []int64) (map[int64]*UserInfo, error) {
	if len(ids) == 0 {
		return make(map[int64]*UserInfo), nil
	}

	// まずREST APIで各IDのログイン名を取得（GraphQLはIDでの検索が難しいため）
	// 並行処理で効率化
	type result struct {
		id   int64
		user *UserInfo
		err  error
	}

	results := make(chan result, len(ids))
	for _, id := range ids {
		go func(userID int64) {
			user, err := c.GetUserByID(ctx, userID)
			results <- result{id: userID, user: user, err: err}
		}(id)
	}

	userMap := make(map[int64]*UserInfo)
	var firstErr error
	for i := 0; i < len(ids); i++ {
		r := <-results
		if r.err != nil {
			if firstErr == nil {
				firstErr = r.err
			}
			continue
		}
		userMap[r.id] = r.user
	}

	if firstErr != nil && len(userMap) == 0 {
		return nil, firstErr
	}

	return userMap, nil
}

// GetUsersByLogins は複数のユーザー情報をログイン名で一括取得する（GraphQL使用）
func (c *Client) GetUsersByLogins(ctx context.Context, logins []string) (map[string]*UserInfo, error) {
	if len(logins) == 0 {
		return make(map[string]*UserInfo), nil
	}

	// GraphQL search APIを使用して一括取得
	userQuery := strings.Join(func() []string {
		result := make([]string, len(logins))
		for i, login := range logins {
			result[i] = "user:" + login
		}
		return result
	}(), " ")

	query := `
		query($q: String!) {
			search(query: $q, type: USER, first: 100) {
				nodes {
					... on User {
						databaseId
						login
						name
						avatarUrl
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"q": userQuery,
	}

	var result struct {
		Search struct {
			Nodes []struct {
				DatabaseId int64  `json:"databaseId"`
				Login      string `json:"login"`
				Name       string `json:"name"`
				AvatarUrl  string `json:"avatarUrl"`
			} `json:"nodes"`
		} `json:"search"`
	}

	if err := c.executeGraphQL(ctx, query, variables, &result); err != nil {
		return nil, fmt.Errorf("failed to get users by logins: %w", err)
	}

	userMap := make(map[string]*UserInfo)
	for _, node := range result.Search.Nodes {
		name := node.Login
		if node.Name != "" {
			name = node.Name
		}
		userMap[node.Login] = &UserInfo{
			ID:        node.DatabaseId,
			Login:     node.Login,
			Name:      name,
			AvatarURL: node.AvatarUrl,
		}
	}

	return userMap, nil
}

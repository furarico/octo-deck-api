package github

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

// GitHubIDからコントリビューション統計を取得する
func (c *Client) GetContributionStats(ctx context.Context, githubID int64) (*ContributionStats, error) {
	// GitHubIDからユーザー情報を取得してログイン名を取得
	userInfo, err := c.GetUserByID(ctx, githubID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// GraphQLクエリ（デフォルトで過去1年間）
	query := `
		query($login: String!) {
			user(login: $login) {
				contributionsCollection {
					contributionCalendar {
						weeks {
							contributionDays {
								date
								contributionCount
							}
						}
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"login": userInfo.Login,
	}

	// GraphQLリクエストの実行
	var result struct {
		User struct {
			ContributionsCollection struct {
				ContributionCalendar struct {
					Weeks []struct {
						ContributionDays []struct {
							Date              string `json:"date"`
							ContributionCount int    `json:"contributionCount"`
						} `json:"contributionDays"`
					} `json:"weeks"`
				} `json:"contributionCalendar"`
			} `json:"contributionsCollection"`
		} `json:"user"`
	}

	if err := c.executeGraphQL(ctx, query, variables, &result); err != nil {
		return nil, err
	}

	// 日ごとのコントリビューションデータを平坦化
	var contributions []Contribution
	for _, week := range result.User.ContributionsCollection.ContributionCalendar.Weeks {
		for _, day := range week.ContributionDays {
			contributions = append(contributions, Contribution{
				Date:  day.Date,
				Count: day.ContributionCount,
			})
		}
	}

	stats := &ContributionStats{
		Contributions: contributions,
	}

	return stats, nil
}

// GetUsersContributions は複数ユーザーの貢献データを取得する
// GitHub GraphQL APIのリソース制限を回避するため、バッチ処理で取得する
// 各バッチは並列で実行され、パフォーマンスが向上する
func (c *Client) GetUsersContributions(ctx context.Context, usernames []string, from, to time.Time) ([]UserContributionStats, error) {
	if len(usernames) == 0 {
		return []UserContributionStats{}, nil
	}

	const batchSize = 15 // 一度に処理するユーザー数（リソース制限回避のため）

	// バッチ数を計算
	numBatches := (len(usernames) + batchSize - 1) / batchSize

	// バッチ結果を格納する構造体
	type batchResult struct {
		index int
		stats []UserContributionStats
		err   error
	}

	results := make(chan batchResult, numBatches)
	var wg sync.WaitGroup

	// 各バッチを並列で実行
	for i := 0; i < len(usernames); i += batchSize {
		wg.Add(1)
		batchIndex := i / batchSize
		start := i
		end := i + batchSize
		if end > len(usernames) {
			end = len(usernames)
		}
		batch := usernames[start:end]

		go func(index int, batch []string, start, end int) {
			defer wg.Done()
			stats, err := c.getUsersContributionsBatch(ctx, batch, from, to)
			if err != nil {
				results <- batchResult{
					index: index,
					err:   fmt.Errorf("failed to get users contributions (batch %d-%d): %w", start, end, err),
				}
				return
			}
			results <- batchResult{index: index, stats: stats}
		}(batchIndex, batch, start, end)
	}

	// 全バッチの完了を待ってチャネルをクローズ
	go func() {
		wg.Wait()
		close(results)
	}()

	// 結果を収集（順序を保持するためにインデックスでソート）
	batchResults := make([]batchResult, numBatches)
	for result := range results {
		if result.err != nil {
			return nil, result.err
		}
		batchResults[result.index] = result
	}

	// 結果を順序通りに結合
	var allStats []UserContributionStats
	for _, br := range batchResults {
		allStats = append(allStats, br.stats...)
	}

	return allStats, nil
}

// getUsersContributionsBatch は指定されたユーザーの貢献データを一括取得する（内部用）
func (c *Client) getUsersContributionsBatch(ctx context.Context, usernames []string, from, to time.Time) ([]UserContributionStats, error) {
	// ユーザー名をクエリ形式に変換 (例: "user:octocat user:torvalds")
	userQuery := strings.Join(func() []string {
		result := make([]string, len(usernames))
		for i, u := range usernames {
			result[i] = "user:" + u
		}
		return result
	}(), " ")

	query := `
		query($q: String!, $from: DateTime!, $to: DateTime!) {
			search(query: $q, type: USER, first: 100) {
				nodes {
					... on User {
						login
						contributionsCollection(from: $from, to: $to) {
							total: contributionCalendar {
								totalContributions
							}
							commits: totalCommitContributions
							issues: totalIssueContributions
							prs: totalPullRequestContributions
							reviews: totalPullRequestReviewContributions
						}
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"q":    userQuery,
		"from": from.Format(time.RFC3339),
		"to":   to.Format(time.RFC3339),
	}

	var result struct {
		Search struct {
			Nodes []struct {
				Login                   string `json:"login"`
				ContributionsCollection struct {
					Total struct {
						TotalContributions int `json:"totalContributions"`
					} `json:"total"`
					Commits int `json:"commits"`
					Issues  int `json:"issues"`
					PRs     int `json:"prs"`
					Reviews int `json:"reviews"`
				} `json:"contributionsCollection"`
			} `json:"nodes"`
		} `json:"search"`
	}

	if err := c.executeGraphQL(ctx, query, variables, &result); err != nil {
		return nil, err
	}

	stats := make([]UserContributionStats, 0, len(result.Search.Nodes))
	for _, node := range result.Search.Nodes {
		stats = append(stats, UserContributionStats{
			Login:   node.Login,
			Total:   node.ContributionsCollection.Total.TotalContributions,
			Commits: node.ContributionsCollection.Commits,
			Issues:  node.ContributionsCollection.Issues,
			PRs:     node.ContributionsCollection.PRs,
			Reviews: node.ContributionsCollection.Reviews,
		})
	}

	return stats, nil
}

// GetContributionsByNodeIDs はNodeIDを使って複数ユーザーの貢献データを一括取得する
// usernameではなくNodeIDを直接使用することで、事前のユーザー情報取得が不要になる
func (c *Client) GetContributionsByNodeIDs(ctx context.Context, nodeIDs []string, from, to time.Time) ([]UserContributionStats, error) {
	if len(nodeIDs) == 0 {
		return []UserContributionStats{}, nil
	}

	const batchSize = 20 // nodes APIは一度に最大100件まで取得可能

	// バッチ数を計算
	numBatches := (len(nodeIDs) + batchSize - 1) / batchSize

	// バッチ結果を格納する構造体
	type batchResult struct {
		index int
		stats []UserContributionStats
		err   error
	}

	results := make(chan batchResult, numBatches)
	var wg sync.WaitGroup

	// 各バッチを並列で実行
	for i := 0; i < len(nodeIDs); i += batchSize {
		wg.Add(1)
		batchIndex := i / batchSize
		start := i
		end := i + batchSize
		if end > len(nodeIDs) {
			end = len(nodeIDs)
		}
		batch := nodeIDs[start:end]

		go func(index int, batch []string, start, end int) {
			defer wg.Done()
			stats, err := c.getContributionsByNodeIDsBatch(ctx, batch, from, to)
			if err != nil {
				results <- batchResult{
					index: index,
					err:   fmt.Errorf("failed to get contributions by node ids (batch %d-%d): %w", start, end, err),
				}
				return
			}
			results <- batchResult{index: index, stats: stats}
		}(batchIndex, batch, start, end)
	}

	// 全バッチの完了を待ってチャネルをクローズ
	go func() {
		wg.Wait()
		close(results)
	}()

	// 結果を収集（順序を保持するためにインデックスでソート）
	batchResults := make([]batchResult, numBatches)
	for result := range results {
		if result.err != nil {
			return nil, result.err
		}
		batchResults[result.index] = result
	}

	// 結果を順序通りに結合
	var allStats []UserContributionStats
	for _, br := range batchResults {
		allStats = append(allStats, br.stats...)
	}

	return allStats, nil
}

// GetUsersFullInfoByNodeIDs はNodeIDを使ってユーザーの全情報を一括取得する
// - ユーザー基本情報（login, name, avatarUrl）
// - 貢献データ（commits, issues, PRs, reviews）
// - 言語情報（最も使用している言語）
// 3回のAPI呼び出しを1回に統合することで、レイテンシを大幅に削減する
func (c *Client) GetUsersFullInfoByNodeIDs(ctx context.Context, nodeIDs []string, from, to time.Time) ([]UserFullInfo, error) {
	if len(nodeIDs) == 0 {
		return []UserFullInfo{}, nil
	}

	const batchSize = 10 // 言語情報も含むため、バッチサイズを小さくしてリソース制限を回避

	// バッチ数を計算
	numBatches := (len(nodeIDs) + batchSize - 1) / batchSize

	// バッチ結果を格納する構造体
	type batchResult struct {
		index int
		infos []UserFullInfo
		err   error
	}

	results := make(chan batchResult, numBatches)
	var wg sync.WaitGroup

	// 各バッチを並列で実行
	for i := 0; i < len(nodeIDs); i += batchSize {
		wg.Add(1)
		batchIndex := i / batchSize
		start := i
		end := i + batchSize
		if end > len(nodeIDs) {
			end = len(nodeIDs)
		}
		batch := nodeIDs[start:end]

		go func(index int, batch []string, start, end int) {
			defer wg.Done()
			infos, err := c.getUsersFullInfoByNodeIDsBatch(ctx, batch, from, to)
			if err != nil {
				results <- batchResult{
					index: index,
					err:   fmt.Errorf("failed to get users full info by node ids (batch %d-%d): %w", start, end, err),
				}
				return
			}
			results <- batchResult{index: index, infos: infos}
		}(batchIndex, batch, start, end)
	}

	// 全バッチの完了を待ってチャネルをクローズ
	go func() {
		wg.Wait()
		close(results)
	}()

	// 結果を収集（順序を保持するためにインデックスでソート）
	batchResults := make([]batchResult, numBatches)
	for result := range results {
		if result.err != nil {
			return nil, result.err
		}
		batchResults[result.index] = result
	}

	// 結果を順序通りに結合
	var allInfos []UserFullInfo
	for _, br := range batchResults {
		allInfos = append(allInfos, br.infos...)
	}

	return allInfos, nil
}

// getUsersFullInfoByNodeIDsBatch はNodeIDを使ってユーザーの全情報を一括取得する（内部用）
func (c *Client) getUsersFullInfoByNodeIDsBatch(ctx context.Context, nodeIDs []string, from, to time.Time) ([]UserFullInfo, error) {
	query := `
		query ($ids: [ID!]!, $from: DateTime!, $to: DateTime!) {
			nodes(ids: $ids) {
				... on User {
					login
					name
					avatarUrl
					contributionsCollection(from: $from, to: $to) {
						total: contributionCalendar {
							totalContributions
						}
						commits: totalCommitContributions
						issues: totalIssueContributions
						prs: totalPullRequestContributions
						reviews: totalPullRequestReviewContributions
					}
					repositories(first: 100, ownerAffiliations: OWNER, isFork: false, privacy: PUBLIC) {
						nodes {
							languages(first: 20) {
								edges {
									size
									node {
										name
									}
								}
							}
						}
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"ids":  nodeIDs,
		"from": from.Format(time.RFC3339),
		"to":   to.Format(time.RFC3339),
	}

	var result struct {
		Nodes []struct {
			Login                   string `json:"login"`
			Name                    string `json:"name"`
			AvatarUrl               string `json:"avatarUrl"`
			ContributionsCollection struct {
				Total struct {
					TotalContributions int `json:"totalContributions"`
				} `json:"total"`
				Commits int `json:"commits"`
				Issues  int `json:"issues"`
				PRs     int `json:"prs"`
				Reviews int `json:"reviews"`
			} `json:"contributionsCollection"`
			Repositories struct {
				Nodes []struct {
					Languages struct {
						Edges []struct {
							Size int `json:"size"`
							Node struct {
								Name string `json:"name"`
							} `json:"node"`
						} `json:"edges"`
					} `json:"languages"`
				} `json:"nodes"`
			} `json:"repositories"`
		} `json:"nodes"`
	}

	if err := c.executeGraphQL(ctx, query, variables, &result); err != nil {
		return nil, err
	}

	infos := make([]UserFullInfo, 0, len(result.Nodes))
	for _, node := range result.Nodes {
		// loginが空の場合はスキップ（Userではないノードの可能性）
		if node.Login == "" {
			continue
		}

		// 最も使用している言語を集計
		languageStats := make(map[string]int)
		for _, repo := range node.Repositories.Nodes {
			for _, edge := range repo.Languages.Edges {
				languageStats[edge.Node.Name] += edge.Size
			}
		}

		// 最大の言語を見つける
		mostUsedLanguage := "Unknown"
		maxSize := 0
		for lang, size := range languageStats {
			if size > maxSize {
				maxSize = size
				mostUsedLanguage = lang
			}
		}

		// 名前がない場合はloginを使用
		name := node.Name
		if name == "" {
			name = node.Login
		}

		infos = append(infos, UserFullInfo{
			Login:                 node.Login,
			Name:                  name,
			AvatarURL:             node.AvatarUrl,
			Total:                 node.ContributionsCollection.Total.TotalContributions,
			Commits:               node.ContributionsCollection.Commits,
			Issues:                node.ContributionsCollection.Issues,
			PRs:                   node.ContributionsCollection.PRs,
			Reviews:               node.ContributionsCollection.Reviews,
			MostUsedLanguage:      mostUsedLanguage,
			MostUsedLanguageColor: GetLanguageColor(mostUsedLanguage),
		})
	}

	return infos, nil
}

// getContributionsByNodeIDsBatch はNodeIDを使って貢献データを一括取得する（内部用）
func (c *Client) getContributionsByNodeIDsBatch(ctx context.Context, nodeIDs []string, from, to time.Time) ([]UserContributionStats, error) {
	query := `
		query ($ids: [ID!]!, $from: DateTime!, $to: DateTime!) {
			nodes(ids: $ids) {
				... on User {
					login
					contributionsCollection(from: $from, to: $to) {
						total: contributionCalendar {
							totalContributions
						}
						commits: totalCommitContributions
						issues: totalIssueContributions
						prs: totalPullRequestContributions
						reviews: totalPullRequestReviewContributions
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"ids":  nodeIDs,
		"from": from.Format(time.RFC3339),
		"to":   to.Format(time.RFC3339),
	}

	var result struct {
		Nodes []struct {
			Login                   string `json:"login"`
			ContributionsCollection struct {
				Total struct {
					TotalContributions int `json:"totalContributions"`
				} `json:"total"`
				Commits int `json:"commits"`
				Issues  int `json:"issues"`
				PRs     int `json:"prs"`
				Reviews int `json:"reviews"`
			} `json:"contributionsCollection"`
		} `json:"nodes"`
	}

	if err := c.executeGraphQL(ctx, query, variables, &result); err != nil {
		return nil, err
	}

	stats := make([]UserContributionStats, 0, len(result.Nodes))
	for _, node := range result.Nodes {
		// loginが空の場合はスキップ（Userではないノードの可能性）
		if node.Login == "" {
			continue
		}
		stats = append(stats, UserContributionStats{
			Login:   node.Login,
			Total:   node.ContributionsCollection.Total.TotalContributions,
			Commits: node.ContributionsCollection.Commits,
			Issues:  node.ContributionsCollection.Issues,
			PRs:     node.ContributionsCollection.PRs,
			Reviews: node.ContributionsCollection.Reviews,
		})
	}

	return stats, nil
}

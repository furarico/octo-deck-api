package github

import (
	"context"
	"fmt"
	"sync"
	"time"
)

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

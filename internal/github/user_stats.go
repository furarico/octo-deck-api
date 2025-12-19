package github

import (
	"context"
	"fmt"
	"sync"
)

// GitHubIDからコントリビューション統計を取得する
func (c *Client) GetUserStats(ctx context.Context, githubID int64) (*UserStats, error) {
	// GitHubIDからユーザー情報を取得してログイン名を取得
	userInfo, err := c.GetUserByID(ctx, githubID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// GraphQLクエリ（デフォルトで過去1年間）
	query := `
		query($login: String!) {
			user(login: $login) {
				# 1. コントリビューション関連の集計
				contributionsCollection {
					# 過去1年間のトータルと日毎のデータ
					contributionCalendar {
						totalContributions
						weeks {
							contributionDays {
								date
								contributionCount
							}
						}
					}
					# コントリビューションの内訳
					totalCommitContributions
					totalIssueContributions
					totalPullRequestContributions
					totalPullRequestReviewContributions
				}
				# 2. 言語統計のためのリポジトリ情報
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
	`

	variables := map[string]interface{}{
		"login": userInfo.Login,
	}

	// GraphQLリクエストの実行
	var result struct {
		User struct {
			ContributionsCollection struct {
				ContributionCalendar struct {
					TotalContributions int `json:"totalContributions"`
					Weeks              []struct {
						ContributionDays []struct {
							Date              string `json:"date"`
							ContributionCount int    `json:"contributionCount"`
						} `json:"contributionDays"`
					} `json:"weeks"`
				} `json:"contributionCalendar"`
				TotalCommitContributions            int `json:"totalCommitContributions"`
				TotalIssueContributions             int `json:"totalIssueContributions"`
				TotalPullRequestContributions       int `json:"totalPullRequestContributions"`
				TotalPullRequestReviewContributions int `json:"totalPullRequestReviewContributions"`
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

	// 最も使用している言語を集計
	languageStats := make(map[string]int)
	for _, repo := range result.User.Repositories.Nodes {
		for _, edge := range repo.Languages.Edges {
			languageStats[edge.Node.Name] += edge.Size
		}
	}

	// 最大の言語を見つける
	mostUsedLanguage := ""
	maxSize := 0
	for lang, size := range languageStats {
		if size > maxSize {
			maxSize = size
			mostUsedLanguage = lang
		}
	}

	stats := &UserStats{
		Contributions:         contributions,
		MostUsedLanguage:      mostUsedLanguage,
		MostUsedLanguageColor: GetLanguageColor(mostUsedLanguage),
		TotalContribution:     result.User.ContributionsCollection.ContributionCalendar.TotalContributions,
		ContributionDetail: ContributionDetail{
			ReviewCount:      result.User.ContributionsCollection.TotalPullRequestReviewContributions,
			CommitCount:      result.User.ContributionsCollection.TotalCommitContributions,
			IssueCount:       result.User.ContributionsCollection.TotalIssueContributions,
			PullRequestCount: result.User.ContributionsCollection.TotalPullRequestContributions,
		},
	}

	return stats, nil
}

// GetMostUsedLanguage はユーザーの最も使用している言語を取得する
func (c *Client) GetMostUsedLanguage(ctx context.Context, login string) (string, string, error) {
	query := `
		query($login: String!) {
			user(login: $login) {
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
	`

	variables := map[string]interface{}{
		"login": login,
	}

	var result struct {
		User struct {
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
		} `json:"user"`
	}

	if err := c.executeGraphQL(ctx, query, variables, &result); err != nil {
		return "", "", fmt.Errorf("failed to execute GraphQL query: %w", err)
	}

	// 最も使用している言語を集計
	languageStats := make(map[string]int)
	for _, repo := range result.User.Repositories.Nodes {
		for _, edge := range repo.Languages.Edges {
			languageStats[edge.Node.Name] += edge.Size
		}
	}

	// 最大の言語を見つける
	mostUsedLanguage := ""
	maxSize := 0
	for lang, size := range languageStats {
		if size > maxSize {
			maxSize = size
			mostUsedLanguage = lang
		}
	}

	// 言語が見つからない場合はデフォルト値を返す
	if mostUsedLanguage == "" {
		return "Unknown", defaultLanguageColor, nil
	}

	return mostUsedLanguage, GetLanguageColor(mostUsedLanguage), nil
}

// LanguageInfo は言語名と色を保持する構造体
type LanguageInfo struct {
	Name  string // 言語名（例: "Go", "Python"）
	Color string // 言語の表示色（例: "#00ADD8"）
}

// GetMostUsedLanguages は複数ユーザーの最も使用している言語を一括取得する
// 並列処理で高速化しつつ、同時実行数を制限してレート制限を回避する
func (c *Client) GetMostUsedLanguages(ctx context.Context, logins []string) (map[string]LanguageInfo, error) {
	if len(logins) == 0 {
		return make(map[string]LanguageInfo), nil
	}

	type result struct {
		login string
		info  LanguageInfo
		err   error
	}

	results := make(chan result, len(logins))
	sem := make(chan struct{}, maxConcurrentRequests) // 同時実行数を制限

	var wg sync.WaitGroup
	for _, login := range logins {
		wg.Add(1)
		go func(login string) {
			defer wg.Done()

			// コンテキストがキャンセルされていたら早期リターン
			select {
			case <-ctx.Done():
				results <- result{
					login: login,
					info:  LanguageInfo{Name: "Unknown", Color: defaultLanguageColor},
					err:   ctx.Err(),
				}
				return
			case sem <- struct{}{}: // セマフォを取得
				defer func() { <-sem }()
			}

			langName, langColor, err := c.GetMostUsedLanguage(ctx, login)
			results <- result{
				login: login,
				info:  LanguageInfo{Name: langName, Color: langColor},
				err:   err,
			}
		}(login)
	}

	// 全てのgoroutineが完了したらチャネルを閉じる
	go func() {
		wg.Wait()
		close(results)
	}()

	langMap := make(map[string]LanguageInfo)

	for r := range results {
		if r.err != nil {
			// エラーが発生した場合はデフォルト値を設定（部分的なエラーを許容）
			langMap[r.login] = LanguageInfo{Name: "Unknown", Color: defaultLanguageColor}
			continue
		}
		langMap[r.login] = r.info
	}

	return langMap, nil
}

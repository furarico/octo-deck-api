package github

import (
	"context"
	"fmt"
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
		return "Unknown", "#586069", nil
	}

	return mostUsedLanguage, GetLanguageColor(mostUsedLanguage), nil
}

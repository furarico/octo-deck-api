package github

import (
	"context"
	"fmt"
	"strings"
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
func (c *Client) GetUsersContributions(ctx context.Context, usernames []string, from, to time.Time) ([]UserContributionStats, error) {
	if len(usernames) == 0 {
		return []UserContributionStats{}, nil
	}

	const batchSize = 15 // 一度に処理するユーザー数（リソース制限回避のため）
	var allStats []UserContributionStats

	for i := 0; i < len(usernames); i += batchSize {
		end := i + batchSize
		if end > len(usernames) {
			end = len(usernames)
		}
		batch := usernames[i:end]

		stats, err := c.getUsersContributionsBatch(ctx, batch, from, to)
		if err != nil {
			return nil, fmt.Errorf("failed to get users contributions (batch %d-%d): %w", i, end, err)
		}
		allStats = append(allStats, stats...)
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

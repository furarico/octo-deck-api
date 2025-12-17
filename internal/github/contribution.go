package github

import (
	"context"
	"fmt"
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

package github

import (
	"time"

	"github.com/furarico/octo-deck-api/internal/domain"
)

// ToDomainStats converts GitHub UserStats to domain Stats
func (us *UserStats) ToDomainStats() (*domain.Stats, error) {
	contributions := make([]domain.Contribution, len(us.Contributions))
	for i, c := range us.Contributions {
		date, err := time.Parse("2006-01-02", c.Date)
		if err != nil {
			return nil, err
		}
		contributions[i] = *domain.NewContribution(date, c.Count)
	}

	language := domain.NewLanguage(us.MostUsedLanguage, us.MostUsedLanguageColor)

	contributionDetail := domain.NewContributionDetail(
		us.ContributionDetail.ReviewCount,
		us.ContributionDetail.CommitCount,
		us.ContributionDetail.PullRequestCount,
		us.ContributionDetail.IssueCount,
	)

	return domain.NewStats(
		contributions,
		us.TotalContribution,
		*language,
		*contributionDetail,
	), nil
}

package domain

type Stats struct {
	Contributions      []Contribution
	TotalContribution  int
	MostUsedLanguage   Language
	ContributionDetail ContributionDetail
}

func NewStats(contributions []Contribution, totalContribution int, mostUsedLanguage Language, contributionDetail ContributionDetail) *Stats {
	return &Stats{
		Contributions:      contributions,
		TotalContribution:  totalContribution,
		MostUsedLanguage:   mostUsedLanguage,
		ContributionDetail: contributionDetail,
	}
}

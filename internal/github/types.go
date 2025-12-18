package github

type UserInfo struct {
	ID        int64
	Login     string
	Name      string
	AvatarURL string
}

type Contribution struct {
	Date  string
	Count int
}

type ContributionStats struct {
	Contributions []Contribution
}

type UserStats struct {
	Contributions      []Contribution
	MostUsedLanguage   string
	TotalContribution  int
	ContributionDetail ContributionDetail
}

type ContributionDetail struct {
	ReviewCount      int
	CommitCount      int
	IssueCount       int
	PullRequestCount int
}

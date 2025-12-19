package github

type UserInfo struct {
	ID        int64
	NodeID    string
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
	Contributions         []Contribution
	MostUsedLanguage      string
	MostUsedLanguageColor string
	TotalContribution     int
	ContributionDetail    ContributionDetail
}

type ContributionDetail struct {
	ReviewCount      int
	CommitCount      int
	IssueCount       int
	PullRequestCount int
}

// UserContributionStats は複数ユーザーの貢献統計を表す
type UserContributionStats struct {
	Login   string
	Total   int
	Commits int
	Issues  int
	PRs     int
	Reviews int
}

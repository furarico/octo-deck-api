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

// UserFullInfo はユーザーの全情報を保持する（統合クエリ用）
// ユーザー基本情報、貢献データ、言語情報を1回のGraphQLクエリで取得するために使用
type UserFullInfo struct {
	Login                 string
	Name                  string
	AvatarURL             string
	Total                 int
	Commits               int
	Issues                int
	PRs                   int
	Reviews               int
	MostUsedLanguage      string
	MostUsedLanguageColor string
}

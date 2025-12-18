package domain

type ContributionDetail struct {
	ReviewCount      int
	CommitCount      int
	PullRequestCount int
	IssueCount       int
}

func NewContributionDetail(reviewCount int, commitCount int, pullRequestCount int, issueCount int) *ContributionDetail {
	return &ContributionDetail{
		ReviewCount:      reviewCount,
		CommitCount:      commitCount,
		PullRequestCount: pullRequestCount,
		IssueCount:       issueCount,
	}
}

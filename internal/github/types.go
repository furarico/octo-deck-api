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

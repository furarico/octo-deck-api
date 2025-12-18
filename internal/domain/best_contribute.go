package domain

type BestContribute struct {
	BestCommitter     Card
	BestContributor   Card
	BestIssuer        Card
	BestPullRequester Card
	BestReviewer      Card
}

func NewBestContribute(bestCommitter Card, bestContributor Card, bestIssuer Card, bestPullRequester Card, bestReviewer Card) *BestContribute {
	return &BestContribute{
		BestCommitter:     bestCommitter,
		BestContributor:   bestContributor,
		BestIssuer:        bestIssuer,
		BestPullRequester: bestPullRequester,
		BestReviewer:      bestReviewer,
	}
}

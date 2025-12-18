package domain

type HighlightedCard struct {
	BestCommitter     Card
	BestContributor   Card
	BestIssuer        Card
	BestPullRequester Card
	BestReviewer      Card
}

func NewHighlightedCard(bestCommitter Card, bestContributor Card, bestIssuer Card, bestPullRequester Card, bestReviewer Card) *HighlightedCard {
	return &HighlightedCard{
		BestCommitter:     bestCommitter,
		BestContributor:   bestContributor,
		BestIssuer:        bestIssuer,
		BestPullRequester: bestPullRequester,
		BestReviewer:      bestReviewer,
	}
}

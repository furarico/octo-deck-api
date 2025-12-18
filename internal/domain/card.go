package domain

type Color string

type Blocks [5][5]bool

type Card struct {
	GithubID         string
	UserName         string
	FullName         string
	IconUrl          string
	Color            Color
	Blocks           Blocks
	MostUsedLanguage Language
}

func NewCard(githubID string, color Color, blocks Blocks, mostUsedLanguage Language) *Card {
	return &Card{
		GithubID:         githubID,
		Color:            color,
		Blocks:           blocks,
		MostUsedLanguage: mostUsedLanguage,
	}
}

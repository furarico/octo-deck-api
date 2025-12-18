package identicon

import (
	"github.com/furarico/octo-deck-api/internal/domain"
)

type MockIdenticonGenerator struct {
	GenerateFunc func(githubID string) (domain.Color, domain.Blocks, error)
}

func (g *MockIdenticonGenerator) Generate(githubID string) (domain.Color, domain.Blocks, error) {
	if g.GenerateFunc != nil {
		return g.GenerateFunc(githubID)
	}
	return domain.Color("#000000"), domain.Blocks{}, nil
}

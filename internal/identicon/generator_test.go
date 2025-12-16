package identicon

import (
	"testing"

	"github.com/furarico/octo-deck-api/internal/domain"
)

func TestGenerator_Generate(t *testing.T) {
	tests := []struct {
		name           string
		githubID       string
		expectedBlocks domain.Blocks
		expectedColor  domain.Color
	}{
		{
			name:     "GitHub ID 136790650",
			githubID: "136790650",
			expectedBlocks: domain.Blocks{
				{true, true, true, true, true},
				{false, false, true, false, false},
				{true, false, false, false, true},
				{false, true, true, true, false},
				{false, true, true, true, false},
			},
			expectedColor: "#8058d7",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGenerator()
			color, blocks, err := g.Generate(tt.githubID)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// blocksの検証
			for i := 0; i < 5; i++ {
				for j := 0; j < 5; j++ {
					if blocks[i][j] != tt.expectedBlocks[i][j] {
						t.Errorf("blocks[%d][%d] = %v, want %v", i, j, blocks[i][j], tt.expectedBlocks[i][j])
					}
				}
			}

			// colorの検証
			if color != tt.expectedColor {
				t.Errorf("color = %v, want %v", color, tt.expectedColor)
			}
		})
	}
}

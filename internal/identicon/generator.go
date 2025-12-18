package identicon

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/furarico/octo-deck-api/internal/domain"
)

const (
	patternSize    = 5
	patternHashLen = 15
	minHashLength  = 32
)

// Generator はIdenticonを生成する実装
type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) Generate(githubID string) (domain.Color, domain.Blocks, error) {
	hash := md5.Sum([]byte(githubID))
	hashStr := hex.EncodeToString(hash[:])

	pattern, err := createPattern(hashStr)
	if err != nil {
		return "", domain.Blocks{}, fmt.Errorf("failed to create pattern: %w", err)
	}

	color, err := createColor(hashStr)
	if err != nil {
		return "", domain.Blocks{}, fmt.Errorf("failed to create color: %w", err)
	}

	var blocks domain.Blocks
	for i := 0; i < patternSize; i++ {
		for j := 0; j < patternSize; j++ {
			blocks[i][j] = pattern[i][j] == 1
		}
	}

	return color, blocks, nil
}

func createPattern(hashValue string) ([][]int, error) {
	if len(hashValue) < patternHashLen {
		return nil, fmt.Errorf("hash must be at least %d characters", patternHashLen)
	}

	pat := make([][]int, patternSize)
	for i := range pat {
		pat[i] = make([]int, patternSize)
	}

	for i := 0; i < patternHashLen; i++ {
		val, err := strconv.ParseInt(string(hashValue[i]), 16, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid hex character at position %d: %w", i, err)
		}

		bit := 1 - int(val%2)
		row := i % patternSize
		col := i / patternSize

		switch col {
		case 0:
			pat[row][2] = bit
		case 1:
			pat[row][1], pat[row][3] = bit, bit
		case 2:
			pat[row][0], pat[row][4] = bit, bit
		}
	}

	return pat, nil
}

func hslToRGB(h, s, l float64) (r, g, b float64) {
	var max, min float64
	if l < 50 {
		max = 2.55 * (l + l*(s/100))
		min = 2.55 * (l - l*(s/100))
	} else {
		max = 2.55 * (l + (100-l)*(s/100))
		min = 2.55 * (l - (100-l)*(s/100))
	}

	switch {
	case h < 60:
		return max, (h/60)*(max-min) + min, min
	case h < 120:
		return ((120-h)/60)*(max-min) + min, max, min
	case h < 180:
		return min, max, ((h-120)/60)*(max-min) + min
	case h < 240:
		return min, ((240-h)/60)*(max-min) + min, max
	case h < 300:
		return ((h-240)/60)*(max-min) + min, min, max
	default:
		return max, min, ((360-h)/60)*(max-min) + min
	}
}

func createColor(hashValue string) (domain.Color, error) {
	if len(hashValue) < minHashLength {
		return "", fmt.Errorf("hash must be at least %d characters", minHashLength)
	}

	hueInt, err := strconv.ParseInt(hashValue[25:28], 16, 64)
	if err != nil {
		return "", fmt.Errorf("invalid hue hex: %w", err)
	}
	hue := float64(hueInt) / 4095 * 360

	satInt, err := strconv.ParseInt(hashValue[28:30], 16, 64)
	if err != nil {
		return "", fmt.Errorf("invalid saturation hex: %w", err)
	}
	sat := 65 - float64(satInt)/255*20

	lumInt, err := strconv.ParseInt(hashValue[30:32], 16, 64)
	if err != nil {
		return "", fmt.Errorf("invalid luminosity hex: %w", err)
	}
	lum := 75 - float64(lumInt)/255*20

	r, g, b := hslToRGB(hue, sat, lum)

	colorHex := fmt.Sprintf("#%02x%02x%02x", int(r), int(g), int(b))
	return domain.Color(colorHex), nil
}

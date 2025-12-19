package github

import (
	_ "embed"
	"encoding/json"
	"sync"
)

//go:embed colors.json
var colorsJSON []byte

var (
	languageColors map[string]string
	colorOnce      sync.Once
)

type languageInfo struct {
	Color string `json:"color"`
	URL   string `json:"url"`
}

// GetLanguageColor は言語名から色を取得する
func GetLanguageColor(language string) string {
	colorOnce.Do(func() {
		languageColors = loadLanguageColors()
	})

	if color, ok := languageColors[language]; ok {
		return color
	}
	return "#586069" // デフォルトカラー（グレー）
}

// loadLanguageColors はcolors.jsonから言語の色情報を読み込む
func loadLanguageColors() map[string]string {
	var languageData map[string]languageInfo
	if err := json.Unmarshal(colorsJSON, &languageData); err != nil {
		return make(map[string]string)
	}

	colors := make(map[string]string)
	for lang, info := range languageData {
		if info.Color != "" {
			colors[lang] = info.Color
		}
	}

	return colors
}

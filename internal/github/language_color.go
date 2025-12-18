package github

import (
	"encoding/json"
	"os"
	"sync"
)

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
	file, err := os.Open("docs/colors.json")
	if err != nil {
		// ファイルが見つからない場合は空のマップを返す
		return make(map[string]string)
	}
	defer file.Close()

	var languageData map[string]languageInfo
	if err := json.NewDecoder(file).Decode(&languageData); err != nil {
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

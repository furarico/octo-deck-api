package domain

type Language struct {
	LanguageName string
	Color        string
}

func NewLanguage(languageName string, color string) *Language {
	return &Language{
		LanguageName: languageName,
		Color:        color,
	}
}

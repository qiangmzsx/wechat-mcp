package theme

import "strings"

func ConvertMarkdown(markdown string, themeID string) string {
	return Convert(markdown, themeID)
}

func ConvertMarkdownWithOptions(markdown string, themeID string, enableGrids bool) string {
	c := NewConverter(
		WithThemeID(themeID),
		EnableImageGrids(enableGrids),
	)
	return c.Convert(markdown)
}

func GetThemeByName(name string) *Theme {
	name = strings.ToLower(name)
	for _, theme := range AllThemes() {
		if strings.ToLower(theme.Name) == name || strings.ToLower(theme.ID) == name {
			return &theme
		}
	}
	return nil
}

func ListThemeIDs() []string {
	return ThemeIDs()
}

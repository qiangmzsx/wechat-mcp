package theme

import "strings"

var themeCache = make(map[string]Theme)

func init() {
	for _, t := range AllThemes() {
		themeCache[t.ID] = t
	}
}

func AllThemes() []Theme {
	themes := make([]Theme, 0)
	themes = append(themes, ClassicThemes()...)
	themes = append(themes, ModernThemes()...)
	themes = append(themes, ExtraThemes()...)
	return themes
}

func GetTheme(id string) Theme {
	if t, ok := themeCache[id]; ok {
		return t
	}
	return ClassicThemes()[0]
}

func ThemeGroups() []ThemeGroup {
	return []ThemeGroup{
		{Label: "经典", Themes: ClassicThemes()},
		{Label: "潮流", Themes: ModernThemes()},
		{Label: "更多风格", Themes: ExtraThemes()},
	}
}

func GetThemeByID(id string) *Theme {
	if t, ok := themeCache[id]; ok {
		return &t
	}
	return nil
}

func ThemeExists(id string) bool {
	_, ok := themeCache[id]
	return ok
}

func ThemeIDs() []string {
	ids := make([]string, 0, len(themeCache))
	for id := range themeCache {
		ids = append(ids, id)
	}
	return ids
}

func SearchThemes(keyword string) []Theme {
	keyword = strings.ToLower(keyword)
	var results []Theme
	for _, theme := range AllThemes() {
		name := strings.ToLower(theme.Name)
		desc := strings.ToLower(theme.Description)
		id := strings.ToLower(theme.ID)
		if strings.Contains(name, keyword) || strings.Contains(desc, keyword) || strings.Contains(id, keyword) {
			results = append(results, theme)
		}
	}
	return results
}

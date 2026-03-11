package theme

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// Theme represents a visual theme for WeChat article formatting
type Theme struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Styles      map[string]string `json:"styles"`
}

// ThemeGroup represents a group of themes
type ThemeGroup struct {
	Label  string  `json:"label"`
	Themes []Theme `json:"themes"`
}

var themeCache = make(map[string]Theme)
var themesLoaded = false

func init() {
	loadThemesFromDir()
	for _, t := range AllThemes() {
		themeCache[t.ID] = t
	}
}

func loadThemesFromDir() {
	if themesLoaded {
		return
	}
	themesLoaded = true

	// Get the themes directory
	themeDir := getThemesDir()

	// Check if directory exists
	if _, err := os.Stat(themeDir); os.IsNotExist(err) {
		return
	}

	// Read all TOML files
	entries, err := os.ReadDir(themeDir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".toml") {
			continue
		}

		themePath := filepath.Join(themeDir, entry.Name())
		theme, err := loadThemeFromFile(themePath)
		if err != nil {
			continue
		}

		themeCache[theme.ID] = theme
	}
}

func getThemesDir() string {
	// Try multiple possible locations
	dirs := []string{
		"../themes",
		"../themes/",
		"themes",
		"./themes",
		filepath.Join(".", "themes"),
		filepath.Join("..", "themes"),
	}

	for _, dir := range dirs {
		if _, err := os.Stat(dir); err == nil {
			return dir
		}
	}

	// Default to "../themes"
	return "../themes"
}

func loadThemeFromFile(path string) (Theme, error) {
	var theme Theme
	_, err := toml.DecodeFile(path, &theme)
	if err != nil {
		return Theme{}, fmt.Errorf("failed to parse theme file %s: %w", path, err)
	}

	// If ID is not set, use filename without extension
	if theme.ID == "" {
		theme.ID = strings.TrimSuffix(filepath.Base(path), ".toml")
	}

	// If Name is not set, use ID
	if theme.Name == "" {
		theme.Name = theme.ID
	}

	return theme, nil
}

func AllThemes() []Theme {
	themes := make([]Theme, 0)

	// Return themes loaded from TOML files
	for _, t := range themeCache {
		themes = append(themes, t)
	}

	return themes
}

func GetTheme(id string) Theme {
	if t, ok := themeCache[id]; ok {
		return t
	}

	// Try case-insensitive match
	idLower := strings.ToLower(id)
	for key, theme := range themeCache {
		if strings.ToLower(key) == idLower {
			return theme
		}
	}

	// Return default theme
	if len(themeCache) > 0 {
		for _, t := range themeCache {
			return t
		}
	}

	// Return first available theme as default
	for _, t := range themeCache {
		return t
	}

	// Return empty theme if no themes available
	return Theme{ID: "default", Name: "Default"}
}

func ThemeIDs() []string {
	ids := make([]string, 0, len(themeCache))
	for id := range themeCache {
		ids = append(ids, id)
	}
	return ids
}

func ThemeExists(id string) bool {
	_, ok := themeCache[id]
	if ok {
		return true
	}
	idLower := strings.ToLower(id)
	for key := range themeCache {
		if strings.ToLower(key) == idLower {
			return true
		}
	}
	return false
}

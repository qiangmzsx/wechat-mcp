package converter

type Theme struct {
	Name        string
	Type        string
	Description string
	Version     string
	Author      string
	Tags        []string
	Styles      map[string]string
	Prompt      string
	AIProvider  string
}

type ThemeManager interface {
	LoadThemes(dir string) error
	GetTheme(name string) (*Theme, error)
	ListThemes() []string
	GetAIPrompt(name string) (string, error)
	GetStyle(name string) (map[string]string, error)
}

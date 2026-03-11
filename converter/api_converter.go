package converter

import (
	"github.com/qiangmzsx/wechat-mcp/theme"
)

type apiConverter struct {
	themeMgr ThemeManager
}

func NewAPIConverter() Converter {
	return &apiConverter{
		themeMgr: NewThemeManager(),
	}
}

func (c *apiConverter) Convert(req *ConvertRequest) *ConvertResult {
	result := &ConvertResult{
		Theme:   req.Theme,
		Success: false,
	}

	if req.Markdown == "" {
		result.Error = "markdown content cannot be empty"
		return result
	}

	themeID := req.Theme
	if themeID == "" {
		themeID = "default"
	}

	html := theme.ConvertMarkdown(req.Markdown, themeID)
	if html == "" {
		result.Error = "theme conversion failed"
		return result
	}

	result.Images = c.ExtractImages(req.Markdown)
	result.HTML = html
	result.Success = true

	return result
}

func (c *apiConverter) ExtractImages(markdown string) []ImageRef {
	return ExtractImages(markdown)
}

func (c *apiConverter) GetThemeManager() ThemeManager {
	return c.themeMgr
}

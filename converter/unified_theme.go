package converter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/qiangmzsx/wechat-mcp/theme"
)

type themeManager struct {
	themes   map[string]*Theme
	themeDir string
}

func NewThemeManager() ThemeManager {
	return &themeManager{
		themes: make(map[string]*Theme),
	}
}

func (tm *themeManager) LoadThemes(dir string) error {
	if dir == "" {
		dir = tm.getDefaultThemeDir()
	}

	tm.themeDir = dir

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read theme directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if !strings.HasSuffix(entry.Name(), ".toml") {
			continue
		}

		themePath := filepath.Join(dir, entry.Name())
		if err := tm.loadThemeFromFile(themePath); err != nil {
			return fmt.Errorf("load theme from %s: %w", themePath, err)
		}
	}

	return nil
}

func (tm *themeManager) loadThemeFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	var themeObj Theme
	if err := toml.Unmarshal(data, &themeObj); err != nil {
		return fmt.Errorf("parse toml: %w", err)
	}

	if themeObj.Name == "" {
		return fmt.Errorf("theme name is required")
	}

	if themeObj.Type == "" {
		themeObj.Type = "ai"
	}

	if themeObj.Prompt == "" {
		themeObj.Prompt = getBuiltinPrompt(themeObj.Name)
	}

	tm.themes[themeObj.Name] = &themeObj
	return nil
}

func (tm *themeManager) getDefaultThemeDir() string {
	if _, err := os.Stat("themes"); err == nil {
		return "themes"
	}

	homeDir, _ := os.UserHomeDir()
	userThemeDir := filepath.Join(homeDir, ".config", "wechat-mcp", "themes")
	if _, err := os.Stat(userThemeDir); err == nil {
		return userThemeDir
	}

	return "themes"
}

func (tm *themeManager) GetTheme(name string) (*Theme, error) {
	if theme, ok := tm.themes[name]; ok {
		return theme, nil
	}

	t := theme.GetThemeByName(name)
	if t == nil {
		if name != "default" && name != "elegant" && name != "tech" && name != "minimalist" {
			return nil, fmt.Errorf("theme not found: %s", name)
		}
		t = theme.GetThemeByName("apple")
		if t == nil {
			return nil, fmt.Errorf("theme not found: %s", name)
		}
		name = "apple"
	}

	converterTheme := &Theme{
		Name:        t.Name,
		Type:        "api",
		Description: t.Description,
		Prompt:      getBuiltinPrompt(name),
		Styles:      t.Styles,
	}

	tm.themes[name] = converterTheme
	return converterTheme, nil
}

func (tm *themeManager) ListThemes() []string {
	ids := theme.ListThemeIDs()
	if len(ids) == 0 {
		ids = []string{"apple", "claude", "wechat", "default", "elegant", "tech", "minimalist"}
	}
	return ids
}

func (tm *themeManager) GetAIPrompt(name string) (string, error) {
	_, err := tm.GetTheme(name)
	if err != nil {
		return "", err
	}
	return getBuiltinPrompt(name), nil
}

func (tm *themeManager) GetStyle(name string) (map[string]string, error) {
	t, err := tm.GetTheme(name)
	if err != nil {
		return nil, err
	}
	return t.Styles, nil
}

func getBuiltinPrompt(name string) string {
	prompts := map[string]string{
		"default": `你是一个专业的微信公众号排版助手。请将以下 Markdown 内容转换为微信公众号兼容的 HTML。

## 样式要求
- 使用简洁大方的中文排版
- 字号适中（16px），行高舒适（1.75）
- 段落之间有适当间距
- 标题加粗醒目
- 引用使用左侧边框
- 代码块使用浅色背景

## 重要规则
1. 所有 CSS 必须使用内联 style 属性
2. 不使用外部样式表或 <style> 标签
3. 只使用安全的 HTML 标签（section, p, span, strong, em, a, h1-h6, ul, ol, li, blockquote, pre, code, table, img, br, hr）
4. 图片使用占位符格式：<!-- IMG:index -->，从0开始计数
5. 返回完整的 HTML，不需要其他说明文字
6. 确保 HTML 在微信中能正常显示`,

		"elegant": `你是一个优雅的微信公众号排版助手。请将以下 Markdown 内容转换为精美的微信公众号 HTML。

## 风格要求
- 优雅精致的排版风格
- 使用柔和的颜色搭配
- 适当的留白和间距
- 标题使用优雅的字体
- 引用使用精致的边框样式

## 重要规则
1. 所有 CSS 必须使用内联 style 属性
2. 不使用外部样式表或 <style> 标签
3. 只使用安全的 HTML 标签
4. 图片使用占位符格式：<!-- IMG:index -->，从0开始计数
5. 返回完整的 HTML，不需要其他说明文字`,

		"tech": `你是一个技术风格的微信公众号排版助手。请将以下 Markdown 内容转换为技术风格的微信公众号 HTML。

## 风格要求
- 简洁专业的技术风格
- 代码块使用深色或高对比度背景
- 适合技术文章阅读
- 清晰的层次结构
- 适当的行高

## 重要规则
1. 所有 CSS 必须使用内联 style 属性
2. 不使用外部样式表或 <style> 标签
3. 只使用安全的 HTML 标签
4. 图片使用占位符格式：<!-- IMG:index -->，从0开始计数
5. 返回完整的 HTML，不需要其他说明文字`,

		"minimalist": `你是一个极简风格的微信公众号排版助手。请将以下 Markdown 内容转换为极简风格的微信公众号 HTML。

## 风格要求
- 极简主义设计
- 大量留白
- 简洁的排版
- 去除不必要的装饰
- 专注于内容本身

## 重要规则
1. 所有 CSS 必须使用内联 style 属性
2. 不使用外部样式表或 <style> 标签
3. 只使用安全的 HTML 标签
4. 图片使用占位符格式：<!-- IMG:index -->，从0开始计数
5. 返回完整的 HTML，不需要其他说明文字`,

		"apple": `你是一个专业的微信公众号排版助手。请将以下 Markdown 内容转换为微信公众号兼容的 HTML。

## 样式要求
- 使用简洁大方的中文排版
- 字号适中（16px），行高舒适（1.75）
- 段落之间有适当间距
- 标题加粗醒目
- 引用使用左侧边框
- 代码块使用浅色背景

## 重要规则
1. 所有 CSS 必须使用内联 style 属性
2. 不使用外部样式表或 <style> 标签
3. 只使用安全的 HTML 标签（section, p, span, strong, em, a, h1-h6, ul, ol, li, blockquote, pre, code, table, img, br, hr）
4. 图片使用占位符格式：<!-- IMG:index -->，从0开始计数
5. 返回完整的 HTML，不需要其他说明文字
6. 确保 HTML 在微信中能正常显示`,
	}

	if prompt, ok := prompts[name]; ok {
		return prompt
	}

	return prompts["default"]
}

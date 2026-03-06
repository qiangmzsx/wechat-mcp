// Package converter
// Theme 模块 - 使用 TOML 格式管理主题配置
package converter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// Theme 主题定义 - TOML 结构
type Theme struct {
	Name        string   `toml:"name"`        // 主题名称
	Type        string   `toml:"type"`        // 类型: ai
	Description string   `toml:"description"` // 主题描述
	Version     string   `toml:"version"`     // 版本
	Author      string   `toml:"author"`      // 作者
	Tags        []string `toml:"tags"`        // 标签

	// 风格配置
	Style StyleConfig `toml:"style"`

	// AI 配置
	Prompt     string `toml:"prompt"`      // AI 提示词
	AIProvider string `toml:"ai_provider"` // 使用的 AI 提供商
}

// StyleConfig 风格配置
type StyleConfig struct {
	Mood       string `toml:"mood"`        // 风格情绪: professional, casual, elegant, tech, etc.
	Colors     string `toml:"colors"`      // 颜色描述
	BestFor    string `toml:"best_for"`    // 适用场景
	FontSize   string `toml:"font_size"`   // 字体大小: small, medium, large
	LineHeight string `toml:"line_height"` // 行高
	Background string `toml:"background"`  // 背景: white, light, grid, none

	// 颜色变量
	PrimaryColor     string `toml:"primary_color"`
	SecondaryColor   string `toml:"secondary_color"`
	BackgroundColor  string `toml:"background_color"`
	TextColor        string `toml:"text_color"`
	AccentColor      string `toml:"accent_color"`
	CodeBlockBg      string `toml:"code_block_bg"`
	QuoteBorderColor string `toml:"quote_border_color"`
}

// ThemeManager 主题管理器接口
type ThemeManager interface {
	// LoadThemes 从目录加载所有主题
	LoadThemes(dir string) error

	// GetTheme 获取指定主题
	GetTheme(name string) (*Theme, error)

	// ListThemes 列出所有主题
	ListThemes() []string

	// GetAIPrompt 获取主题的 AI 提示词
	GetAIPrompt(name string) (string, error)

	// GetStyle 获取主题的风格配置
	GetStyle(name string) (*StyleConfig, error)
}

// themeManager 主题管理器实现
type themeManager struct {
	themes   map[string]*Theme
	themeDir string
	builtin  map[string]string // 内置主题名 -> 内置提示词
}

// NewThemeManager 创建主题管理器
func NewThemeManager() ThemeManager {
	return &themeManager{
		themes:  make(map[string]*Theme),
		builtin: getBuiltinThemes(),
	}
}

// getBuiltinThemes 获取内置主题
func getBuiltinThemes() map[string]string {
	return map[string]string{
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
	}
}

// LoadThemes 从目录加载所有主题
func (tm *themeManager) LoadThemes(dir string) error {
	if dir == "" {
		dir = tm.getDefaultThemeDir()
	}

	tm.themeDir = dir

	// 检查目录是否存在
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// 目录不存在，创建并加载内置主题
		return nil
	}

	// 遍历目录加载主题文件
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read theme directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// 只处理 .toml 文件
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

// loadThemeFromFile 从文件加载主题
func (tm *themeManager) loadThemeFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	var theme Theme
	if err := toml.Unmarshal(data, &theme); err != nil {
		return fmt.Errorf("parse toml: %w", err)
	}

	// 验证主题
	if theme.Name == "" {
		return fmt.Errorf("theme name is required")
	}

	// 设置默认值
	if theme.Type == "" {
		theme.Type = "ai"
	}

	// 如果没有定义 prompt，使用内置的
	if theme.Prompt == "" {
		if builtin, ok := tm.builtin[theme.Name]; ok {
			theme.Prompt = builtin
		}
	}

	tm.themes[theme.Name] = &theme
	return nil
}

// getDefaultThemeDir 获取默认主题目录
func (tm *themeManager) getDefaultThemeDir() string {
	// 优先使用项目根目录的 themes/ 文件夹
	if _, err := os.Stat("themes"); err == nil {
		return "themes"
	}

	// 其次使用用户配置目录
	homeDir, _ := os.UserHomeDir()
	userThemeDir := filepath.Join(homeDir, ".config", "wechat-mcp", "themes")
	if _, err := os.Stat(userThemeDir); err == nil {
		return userThemeDir
	}

	// 返回默认路径
	return "themes"
}

// GetTheme 获取指定主题
func (tm *themeManager) GetTheme(name string) (*Theme, error) {
	// 优先从已加载的主题中获取
	if theme, ok := tm.themes[name]; ok {
		return theme, nil
	}

	// 尝试从内置主题获取
	if _, ok := tm.builtin[name]; ok {
		// 加载内置主题
		theme := &Theme{
			Name:        name,
			Type:        "ai",
			Description: getThemeDescription(name),
			Prompt:      tm.builtin[name],
		}
		tm.themes[name] = theme
		return theme, nil
	}

	return nil, fmt.Errorf("theme not found: %s", name)
}

// getThemeDescription 获取主题描述
func getThemeDescription(name string) string {
	descriptions := map[string]string{
		"default":    "默认主题 - 简洁大方的通用风格",
		"elegant":    "优雅主题 - 精致柔和的排版风格",
		"tech":       "技术主题 - 适合技术文章的简洁风格",
		"minimalist": "极简主题 - 简约清爽的阅读体验",
	}
	if desc, ok := descriptions[name]; ok {
		return desc
	}
	return name + " 主题"
}

// ListThemes 列出所有主题
func (tm *themeManager) ListThemes() []string {
	// 合并已加载的主题和内置主题
	themeMap := make(map[string]bool)

	// 添加已加载的主题
	for name := range tm.themes {
		themeMap[name] = true
	}

	// 添加内置主题
	for name := range tm.builtin {
		themeMap[name] = true
	}

	// 转换为切片
	names := make([]string, 0, len(themeMap))
	for name := range themeMap {
		names = append(names, name)
	}

	return names
}

// GetAIPrompt 获取主题的 AI 提示词
func (tm *themeManager) GetAIPrompt(name string) (string, error) {
	theme, err := tm.GetTheme(name)
	if err != nil {
		return "", err
	}

	if theme.Prompt == "" {
		return "", fmt.Errorf("theme '%s' has no prompt defined", name)
	}

	return theme.Prompt, nil
}

// GetStyle 获取主题的风格配置
func (tm *themeManager) GetStyle(name string) (*StyleConfig, error) {
	theme, err := tm.GetTheme(name)
	if err != nil {
		return nil, err
	}

	// 返回默认风格配置
	if theme.Style.FontSize == "" {
		defaultStyle := getDefaultStyle(name)
		return &defaultStyle, nil
	}

	return &theme.Style, nil
}

// getDefaultStyle 获取默认风格配置
func getDefaultStyle(name string) StyleConfig {
	defaultStyles := map[string]StyleConfig{
		"default": {
			FontSize:     "medium",
			LineHeight:   "1.75",
			Background:   "white",
			PrimaryColor: "#333333",
			TextColor:    "#333333",
			AccentColor:  "#4a90d9",
		},
		"elegant": {
			FontSize:     "medium",
			LineHeight:   "1.8",
			Background:   "light",
			PrimaryColor: "#2c3e50",
			TextColor:    "#34495e",
			AccentColor:  "#c0392b",
		},
		"tech": {
			FontSize:     "medium",
			LineHeight:   "1.7",
			Background:   "white",
			PrimaryColor: "#282c34",
			TextColor:    "#333333",
			AccentColor:  "#61afef",
		},
		"minimalist": {
			FontSize:     "medium",
			LineHeight:   "1.8",
			Background:   "white",
			PrimaryColor: "#000000",
			TextColor:    "#333333",
			AccentColor:  "#999999",
		},
	}

	if style, ok := defaultStyles[name]; ok {
		return style
	}

	return defaultStyles["default"]
}

// LoadThemeFromBytes 从字节数组加载主题
func (tm *themeManager) LoadThemeFromBytes(data []byte, name string) error {
	var theme Theme
	if err := toml.Unmarshal(data, &theme); err != nil {
		return fmt.Errorf("parse toml: %w", err)
	}

	if theme.Name == "" {
		theme.Name = name
	}

	tm.themes[theme.Name] = &theme
	return nil
}

// AddTheme 添加主题
func (tm *themeManager) AddTheme(theme *Theme) error {
	if theme.Name == "" {
		return fmt.Errorf("theme name is required")
	}
	tm.themes[theme.Name] = theme
	return nil
}

package converter

import (
	"fmt"

	"github.com/qiangmzsx/wechat-mcp/theme"
)

type unifiedThemeManager struct{}

func NewUnifiedThemeManager() ThemeManager {
	return &unifiedThemeManager{}
}

func (u *unifiedThemeManager) LoadThemes(dir string) error {
	return nil
}

func (u *unifiedThemeManager) GetTheme(name string) (*Theme, error) {
	t := theme.GetThemeByName(name)
	if t == nil {
		return nil, fmt.Errorf("theme not found: %s", name)
	}

	converterTheme := &Theme{
		Name:        t.Name,
		Type:        "api",
		Description: t.Description,
	}

	if t.Styles != nil {
		converterTheme.Style = StyleConfig{
			PrimaryColor:    t.Styles["primary_color"],
			SecondaryColor:  t.Styles["secondary_color"],
			BackgroundColor: t.Styles["background_color"],
			TextColor:       t.Styles["text_color"],
			AccentColor:     t.Styles["accent_color"],
		}
	}

	return converterTheme, nil
}

func (u *unifiedThemeManager) ListThemes() []string {
	return theme.ListThemeIDs()
}

func (u *unifiedThemeManager) GetAIPrompt(name string) (string, error) {
	themeObj := theme.GetThemeByName(name)
	if themeObj == nil {
		return "", fmt.Errorf("theme not found: %s", name)
	}
	return getBuiltinPrompt(name), nil
}

func (u *unifiedThemeManager) GetStyle(name string) (*StyleConfig, error) {
	t := theme.GetThemeByName(name)
	if t == nil {
		return nil, fmt.Errorf("theme not found: %s", name)
	}

	if t.Styles == nil {
		return &StyleConfig{}, nil
	}

	return &StyleConfig{
		PrimaryColor:    t.Styles["primary_color"],
		SecondaryColor:  t.Styles["secondary_color"],
		BackgroundColor: t.Styles["background_color"],
		TextColor:       t.Styles["text_color"],
		AccentColor:     t.Styles["accent_color"],
	}, nil
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
	}

	if prompt, ok := prompts[name]; ok {
		return prompt
	}

	return prompts["default"]
}

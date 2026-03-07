// Package converter 提供 Markdown 到微信公众号 HTML 的 AI 转换功能
// 使用策略模式和建造者模式，支持多主题配置
package converter

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/qiangmzsx/wechat-mcp/config"
	"github.com/qiangmzsx/wechat-mcp/provider"
	"github.com/qiangmzsx/wechat-mcp/provider/factory"
	"go.uber.org/zap"
)


// aiConverter AI 模式转换器
type aiConverter struct {
	log    *zap.Logger
	config config.ConverterConfig
	client provider.Provider
	theme  ThemeManager
	prompt *PromptBuilder
}


// NewAIConverter 创建 AI 转换器
func NewAIConverter(cfg *config.Config, log *zap.Logger) (Converter, error) {
	if !cfg.Converter.Enabled {
		return nil, fmt.Errorf("converter is disabled")
	}

	// 使用工厂创建 AI Provider
	client, err := factory.NewProvider(cfg)
	if err != nil {
		return nil, fmt.Errorf("create provider failed: %w", err)
	}


	// 创建主题管理器并加载主题
	themeMgr := NewThemeManager()
	if cfg.Converter.ThemeDir != "" {
		if err := themeMgr.LoadThemes(cfg.Converter.ThemeDir); err != nil {
			log.Warn("failed to load themes from directory",
				zap.String("dir", cfg.Converter.ThemeDir),
				zap.Error(err))
		}
	}

	return &aiConverter{
		log:    log,
		config: cfg.Converter,
		client: client,
		theme:  themeMgr,
		prompt: NewPromptBuilder(),
	}, nil
}

// Convert 执行转换
func (c *aiConverter) Convert(req *ConvertRequest) *ConvertResult {
	result := &ConvertResult{
		Theme:   req.Theme,
		Success: false,
	}

	// 验证请求
	if err := c.validateRequest(req); err != nil {
		result.Error = err.Error()
		return result
	}

	// 构建 Prompt
	prompt, err := c.buildPrompt(req)
	if err != nil {
		result.Error = fmt.Sprintf("build prompt failed: %s", err.Error())
		return result
	}

	// 提取图片引用
	images := c.ExtractImages(req.Markdown)
	result.Images = images

	// 调用 AI Provider 生成 HTML
	resp, err := c.client.Chat(
		context.Background(),
		provider.ChatRequest{
			Messages: []provider.Message{
				{Role: "system", Content: c.getSystemPrompt()},
				{Role: "user", Content: prompt},
			},
			Model: c.config.Model,
			Options: map[string]interface{}{
				provider.OptMaxTokens: c.config.MaxTokens,
			},
		},
	)
	if err != nil {
		result.Error = fmt.Sprintf("AI generation failed: %s", err.Error())
		c.log.Error("AI conversion failed",
			zap.String("theme", req.Theme),
			zap.Error(err))
		return result
	}

	// 提取 HTML 内容
	if resp.Content == "" {
		result.Error = "AI returned empty response"
		return result
	}

	html := resp.Content


	// 处理图片占位符
	html = c.processImagePlaceholders(html, images)

	result.HTML = html
	result.Success = true

	c.log.Info("AI conversion succeeded",
		zap.String("theme", req.Theme),
		zap.Int("image_count", len(images)),
		zap.Int("html_length", len(html)))

	return result
}

// validateRequest 验证请求
func (c *aiConverter) validateRequest(req *ConvertRequest) error {
	if req.Markdown == "" {
		return ErrEmptyMarkdownErr
	}

	// 如果没有指定主题，使用默认
	if req.Theme == "" {
		req.Theme = c.config.DefaultTheme
		if req.Theme == "" {
			req.Theme = "default"
		}
	}

	return nil
}

// buildPrompt 构建 AI 提示词
func (c *aiConverter) buildPrompt(req *ConvertRequest) (string, error) {
	// 优先使用自定义提示词
	if req.CustomPrompt != "" {
		return c.prompt.BuildCustomPrompt(req.CustomPrompt, req.Markdown), nil
	}

	// 获取主题提示词
	themePrompt, err := c.theme.GetAIPrompt(req.Theme)
	if err != nil {
		// 降级到默认主题
		c.log.Warn("theme not found, using default", zap.String("theme", req.Theme))
		themePrompt, _ = c.theme.GetAIPrompt("default")
	}

	// 获取主题风格配置
	style, err := c.theme.GetStyle(req.Theme)
	if err != nil {
		c.log.Warn("failed to get theme style, using defaults", zap.String("theme", req.Theme), zap.Error(err))
		style = &StyleConfig{
			PrimaryColor:    "#333333",
			SecondaryColor:  "#666666",
			BackgroundColor: "#ffffff",
			TextColor:       "#333333",
			AccentColor:     "#4a90d9",
		}
	}

	// 构建颜色变量
	vars := map[string]string{
		"PRIMARY_COLOR":    style.PrimaryColor,
		"SECONDARY_COLOR":  style.SecondaryColor,
		"BACKGROUND_COLOR": style.BackgroundColor,
		"TEXT_COLOR":       style.TextColor,
		"ACCENT_COLOR":     style.AccentColor,
	}

	return c.prompt.BuildPrompt(themePrompt, req.Markdown, vars), nil
}

// getSystemPrompt 获取系统提示词
func (c *aiConverter) getSystemPrompt() string {
	return `你是一个专业的微信公众号排版助手。请将 Markdown 内容转换为微信公众号兼容的 HTML。

## 重要规则
1. 所有 CSS 必须使用内联 style 属性
2. 不使用外部样式表或 <style> 标签
3. 只使用安全的 HTML 标签（section, p, span, strong, em, a, h1-h6, ul, ol, li, blockquote, pre, code, table, img, br, hr）
4. 图片使用占位符格式：<!-- IMG:index -->，从0开始计数
5. 返回完整的 HTML，不需要其他说明文字
6. 确保 HTML 在微信中能正常显示`
}

// extractHTML 从 AI 响应中提取 HTML
func (c *aiConverter) extractHTML(content string) string {
	// 尝试提取代码块中的 HTML
	// 格式: ```html ... ``` 或 ``` ... ```
	re := regexp.MustCompile("(?s)```(?:html)?\\s*(.*?)```")
	matches := re.FindStringSubmatch(content)
	if len(matches) >= 2 {
		return strings.TrimSpace(matches[1])
	}

	// 如果没有代码块，尝试提取 <html> 或 <section> 标签
	if strings.Contains(content, "<html") || strings.Contains(content, "<section") {
		return strings.TrimSpace(content)
	}

	// 尝试找到 HTML 开始和结束
	htmlStart := strings.Index(content, "<")
	if htmlStart >= 0 {
		return strings.TrimSpace(content[htmlStart:])
	}

	// 返回原始内容
	return strings.TrimSpace(content)
}

// processImagePlaceholders 处理图片占位符
// 优先级: WechatURL > Original
func (c *aiConverter) processImagePlaceholders(html string, images []ImageRef) string {
	result := html
	for _, img := range images {
		placeholder := fmt.Sprintf("<!-- IMG:%d -->", img.Index)
		if !strings.Contains(result, placeholder) {
			// 尝试其他格式
			placeholder = fmt.Sprintf("<!--IMG:%d-->", img.Index)
		}

		// 确定使用的图片 URL
		imgURL := img.WechatURL
		if imgURL == "" {
			// 如果没有微信 URL，使用原始路径
			imgURL = img.Original
		}

		// 替换占位符为实际图片标签
		if imgURL != "" {
			imgTag := "<img src=\"" + imgURL + "\" style=\"max-width:100%;height:auto;display:block;margin:20px auto;\" />"
			result = strings.ReplaceAll(result, placeholder, imgTag)
		}
	}
	return result
}

// ExtractImages 从 Markdown 中提取图片引用
func (c *aiConverter) ExtractImages(markdown string) []ImageRef {
	return ExtractImages(markdown)
}

// GetThemeManager 获取主题管理器
func (c *aiConverter) GetThemeManager() ThemeManager {
	return c.theme
}

// SimpleConverter 简化版转换器 - 用于不需要实际 API 调用的场景
type SimpleConverter struct {
	theme  ThemeManager
	prompt *PromptBuilder
}

// NewSimpleConverter 创建简化版转换器
func NewSimpleConverter() Converter {
	return &SimpleConverter{
		theme:  NewThemeManager(),
		prompt: NewPromptBuilder(),
	}
}

// Convert 执行转换（简化版，返回带提示的结果）
func (c *SimpleConverter) Convert(req *ConvertRequest) *ConvertResult {
	result := &ConvertResult{
		Theme:   req.Theme,
		Success: false,
	}

	if req.Markdown == "" {
		result.Error = "markdown content cannot be empty"
		return result
	}

	// 如果没有指定主题，使用默认
	if req.Theme == "" {
		req.Theme = "default"
	}

	// 构建 Prompt
	prompt, err := c.theme.GetAIPrompt(req.Theme)
	if err != nil {
		prompt, _ = c.theme.GetAIPrompt("default")
	}
	result.HTML = c.prompt.BuildPrompt(prompt, req.Markdown, nil)

	// 提取图片
	result.Images = c.ExtractImages(req.Markdown)
	result.Success = true

	return result
}

// ExtractImages 从 Markdown 中提取图片引用
func (c *SimpleConverter) ExtractImages(markdown string) []ImageRef {
	return ExtractImages(markdown)
}

// GetThemeManager 获取主题管理器
func (c *SimpleConverter) GetThemeManager() ThemeManager {
	return c.theme
}

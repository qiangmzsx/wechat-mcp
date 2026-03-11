// Package converter 提供 Markdown 到微信公众号 HTML 的 AI 转换功能
// 使用策略模式和建造者模式，支持多主题配置
package converter

import (
	"fmt"
	"regexp"

	"github.com/qiangmzsx/wechat-mcp/config"
)

// ImageType 图片类型
type ImageType string

const (
	ImageTypeLocal  ImageType = "local"  // 本地图片
	ImageTypeOnline ImageType = "online" // 在线图片
	ImageTypeAI     ImageType = "ai"     // AI 生成图片
)

// ImageRef 图片引用
type ImageRef struct {
	Index       int       // 位置索引
	Original    string    // 原始路径或提示词
	Placeholder string    // HTML 中的占位符 <!-- IMG:0 -->
	WechatURL   string    // 上传后的 URL (处理完成后)
	Type        ImageType // 图片类型
	AIPrompt    string    // AI 图片的生成提示词
}

// ConvertResult 转换结果
type ConvertResult struct {
	HTML    string     // 生成的 HTML（含占位符）
	Theme   string     // 使用的主题
	Images  []ImageRef // 图片引用列表
	Success bool       // 是否成功
	Error   string     // 错误信息
}

// ConvertRequest 转换请求
type ConvertRequest struct {
	Markdown      string               // Markdown 内容
	Theme         string               // 主题名称 (可选，默认使用配置的默认主题)
	CustomPrompt  string               // 自定义提示词（可选）
	ConverterType config.ConverterType // 转换器类型: api, ai (可选)
}

// Converter 转换器接口 - 策略模式
type Converter interface {
	// Convert 执行转换
	Convert(req *ConvertRequest) *ConvertResult

	// ExtractImages 从 Markdown 中提取图片引用
	ExtractImages(markdown string) []ImageRef

	// GetThemeManager 获取主题管理器
	GetThemeManager() ThemeManager
}

// AILLM AI LLM 接口 - 用于接入不同的 AI 提供商
type AILLM interface {
	// Generate 生成内容
	Generate(prompt string) (string, error)

	// GenerateWithSystem 使用系统提示词
	GenerateWithSystem(systemPrompt, userPrompt string) (string, error)
}

// ErrorCode 错误码
type ErrorCode string

const (
	ErrEmptyMarkdown  ErrorCode = "EMPTY_MARKDOWN"
	ErrMissingAPIKey  ErrorCode = "MISSING_API_KEY"
	ErrInvalidTheme   ErrorCode = "INVALID_THEME"
	ErrAIFailure      ErrorCode = "AI_FAILURE"
	ErrInvalidRequest ErrorCode = "INVALID_REQUEST"
)

// ConvertError 转换错误
type ConvertError struct {
	Code    ErrorCode
	Message string
	Err     error
}

func (e *ConvertError) Error() string {
	if e.Err != nil {
		return string(e.Code) + ": " + e.Message + ": " + e.Err.Error()
	}
	return string(e.Code) + ": " + e.Message
}

func (e *ConvertError) Unwrap() error {
	return e.Err
}

// 定义错误变量
var (
	ErrEmptyMarkdownErr = &ConvertError{Code: ErrEmptyMarkdown, Message: "markdown content cannot be empty"}
	ErrInvalidThemeErr  = &ConvertError{Code: ErrInvalidTheme, Message: "invalid theme name"}
	ErrAIFailureErr     = &ConvertError{Code: ErrAIFailure, Message: "AI generation failed"}
)

// imagePatterns 图片匹配正则表达式
var (
	localImagePattern  = regexp.MustCompile(`!\[([^\]]*)\]\((\.\/[^)]+)\)`)
	onlineImagePattern = regexp.MustCompile(`!\[([^\]]*)\]\((https?://[^)]+)\)`)
	aiImagePattern     = regexp.MustCompile(`!\[([^\]]*)\]\(__generate:([^)]+)__\)`)
)

// ExtractImages 从 Markdown 中提取图片引用
func ExtractImages(markdown string) []ImageRef {
	var images []ImageRef

	// 匹配本地图片: ![alt](./path/to/image.png)
	for i, match := range localImagePattern.FindAllStringSubmatch(markdown, -1) {
		if len(match) >= 3 {
			images = append(images, ImageRef{
				Index:       i,
				Original:    match[2],
				Placeholder: fmt.Sprintf("<!-- IMG:%d -->", i),
				Type:        ImageTypeLocal,
			})
		}
	}

	// 匹配在线图片: ![alt](https://...)
	offset := len(images)
	for i, match := range onlineImagePattern.FindAllStringSubmatch(markdown, -1) {
		if len(match) >= 3 {
			images = append(images, ImageRef{
				Index:       offset + i,
				Original:    match[2],
				Placeholder: fmt.Sprintf("<!-- IMG:%d -->", offset+i),
				Type:        ImageTypeOnline,
			})
		}
	}

	// 匹配 AI 生成图片: ![alt](__generate:prompt__)
	offset = len(images)
	for i, match := range aiImagePattern.FindAllStringSubmatch(markdown, -1) {
		if len(match) >= 3 {
			images = append(images, ImageRef{
				Index:       offset + i,
				Original:    match[2],
				Placeholder: fmt.Sprintf("<!-- IMG:%d -->", offset+i),
				Type:        ImageTypeAI,
				AIPrompt:    match[2],
			})
		}
	}

	return images
}

// ReplaceImagePlaceholders 在 HTML 中替换图片占位符
func ReplaceImagePlaceholders(html string, images []ImageRef) string {
	result := html
	for _, img := range images {
		if img.WechatURL != "" {
			imgTag := `<img src="` + img.WechatURL + `" style="max-width:100%;height:auto;display:block;margin:20px auto;" />`
			result = replacePlaceholder(result, img.Placeholder, imgTag)
		}
	}
	return result
}

// replacePlaceholder 替换占位符的辅助函数
func replacePlaceholder(html, placeholder, replacement string) string {
	// 如果占位符为空，使用默认格式
	if placeholder == "" {
		placeholder = `<!-- IMG:\d+ -->`
	}
	return regexp.MustCompile(placeholder).ReplaceAllString(html, replacement)
}

// GeneratePlaceholder 生成图片占位符
func GeneratePlaceholder(index int) string {
	return "<!-- IMG:" + string(rune('0'+index)) + " -->"
}

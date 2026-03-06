// Package converter
// Prompt 模块 - 构建 AI 转换提示词
package converter

import (
	"fmt"
	"regexp"
	"strings"
	"text/template"
)

// PromptBuilder Prompt 构建器 - 建造者模式
type PromptBuilder struct {
	systemPrompts map[string]string // 系统提示词模板
	variables     map[string]*PromptVariable
}

// PromptVariable Prompt 变量
type PromptVariable struct {
	Name         string
	Description  string
	DefaultValue string
	Required     bool
}

// NewPromptBuilder 创建 Prompt 构建器
func NewPromptBuilder() *PromptBuilder {
	pb := &PromptBuilder{
		systemPrompts: make(map[string]string),
		variables:     make(map[string]*PromptVariable),
	}
	pb.initBuiltInVariables()
	pb.initBuiltInSystemPrompts()
	return pb
}

// initBuiltInVariables 初始化内置变量
func (pb *PromptBuilder) initBuiltInVariables() {
	pb.variables = map[string]*PromptVariable{
		"MARKDOWN": {
			Name:         "MARKDOWN",
			Description:  "Markdown 内容",
			DefaultValue: "",
			Required:     true,
		},
		"THEME_NAME": {
			Name:         "THEME_NAME",
			Description:  "主题名称",
			DefaultValue: "default",
			Required:     false,
		},
		"TITLE": {
			Name:         "TITLE",
			Description:  "文章标题",
			DefaultValue: "未命名文章",
			Required:     false,
		},
		"FONT_SIZE": {
			Name:         "FONT_SIZE",
			Description:  "字体大小",
			DefaultValue: "16px",
			Required:     false,
		},
		"LINE_HEIGHT": {
			Name:         "LINE_HEIGHT",
			Description:  "行高",
			DefaultValue: "1.75",
			Required:     false,
		},
		"PRIMARY_COLOR": {
			Name:         "PRIMARY_COLOR",
			Description:  "主色调",
			DefaultValue: "#333333",
			Required:     false,
		},
	}
}

// initBuiltInSystemPrompts 初始化内置系统提示词
func (pb *PromptBuilder) initBuiltInSystemPrompts() {
	pb.systemPrompts = map[string]string{
		"default": `你是一个专业的微信公众号排版助手。请将以下 Markdown 内容转换为微信公众号兼容的 HTML。

## 样式要求
- 使用简洁大方的中文排版
- 字号适中，行高舒适
- 段落之间有适当间距
- 标题加粗醒目

## 重要规则
1. 所有 CSS 必须使用内联 style 属性
2. 不使用外部样式表或 <style> 标签
3. 只使用安全的 HTML 标签
4. 图片使用占位符格式：<!-- IMG:index -->，从0开始计数
5. 返回完整的 HTML，不需要其他说明文字`,

		"wechat_compatible": `你是一个微信公众号内容转换专家。请将 Markdown 转换为微信公众号兼容的 HTML。

## 微信公众号限制
- 不支持复杂的 CSS
- 图片需要使用微信素材库 URL
- 建议使用内联样式

## 重要规则
1. 所有 CSS 必须使用内联 style 属性
2. 不使用外部样式表
3. 只使用安全的 HTML 标签
4. 图片使用占位符格式：<!-- IMG:index -->，从0开始计数
5. 返回完整的 HTML`,

		"html_strict": `请将以下 Markdown 内容严格转换为微信公众号兼容的 HTML。

## 严格规则
1. 所有 CSS 必须使用内联 style 属性
2. 禁止使用 <style> 标签
3. 禁止使用外部样式表
4. 只允许使用安全的 HTML 标签
5. 图片使用占位符：<!-- IMG:index -->，从0开始计数
6. 返回纯 HTML，不要包含任何说明文字`,
	}
}

// BuildPrompt 构建完整的 Prompt
func (pb *PromptBuilder) BuildPrompt(themePrompt, markdown string, vars map[string]string) string {
	// 如果没有提供主题 prompt，使用默认
	if themePrompt == "" {
		themePrompt = pb.systemPrompts["default"]
	}

	// 合并变量
	if vars == nil {
		vars = make(map[string]string)
	}

	// 确保 MARKDOWN 变量存在
	if _, ok := vars["MARKDOWN"]; !ok {
		vars["MARKDOWN"] = markdown
	}

	// 构建完整 prompt
	result := themePrompt + "\n\n" + "```\n" + vars["MARKDOWN"] + "\n```"

	// 替换其他变量
	for key, value := range vars {
		if key != "MARKDOWN" {
			placeholder := "{{" + key + "}}"
			result = strings.ReplaceAll(result, placeholder, value)
		}
	}

	return result
}

// BuildSystemPrompt 构建系统提示词
func (pb *PromptBuilder) BuildSystemPrompt(style string) string {
	if prompt, ok := pb.systemPrompts[style]; ok {
		return prompt
	}
	return pb.systemPrompts["default"]
}

// AddSystemPrompt 添加自定义系统提示词
func (pb *PromptBuilder) AddSystemPrompt(name, prompt string) {
	pb.systemPrompts[name] = prompt
}

// GetSystemPrompt 获取系统提示词
func (pb *PromptBuilder) GetSystemPrompt(name string) (string, error) {
	if prompt, ok := pb.systemPrompts[name]; ok {
		return prompt, nil
	}
	return "", fmt.Errorf("system prompt not found: %s", name)
}

// ListSystemPrompts 列出所有系统提示词
func (pb *PromptBuilder) ListSystemPrompts() []string {
	names := make([]string, 0, len(pb.systemPrompts))
	for name := range pb.systemPrompts {
		names = append(names, name)
	}
	return names
}

// BuildPromptWithTemplate 使用 Go template 构建 Prompt
func (pb *PromptBuilder) BuildPromptWithTemplate(templateContent string, data map[string]interface{}) (string, error) {
	// 清理变量名
	cleanData := make(map[string]interface{})
	for k, v := range data {
		cleanKey := strings.ToLower(strings.ReplaceAll(k, "_", ""))
		cleanData[cleanKey] = v
	}

	tmpl, err := template.New("prompt").Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, cleanData); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}

	return buf.String(), nil
}

// ValidatePrompt 验证 Prompt 内容
func (pb *PromptBuilder) ValidatePrompt(prompt string) *ValidationResult {
	result := &ValidationResult{
		Valid:    true,
		Errors:   []string{},
		Warnings: []string{},
	}

	// 检查关键规则
	requiredPatterns := []struct {
		pattern string
		name    string
	}{
		{`style\s*=`, "内联样式"},
		{`IMG:\d+|<!-- IMG:`, "图片占位符"},
		{`<[a-z]+[^>]*>`, "HTML 标签"},
	}

	for _, rule := range requiredPatterns {
		matched, _ := regexp.MatchString(rule.pattern, prompt)
		if !matched {
			result.Warnings = append(result.Warnings, fmt.Sprintf("建议包含 %s 相关内容", rule.name))
		}
	}

	// 检查不安全内容
	dangerousPatterns := []string{
		`<script`,
		`javascript:`,
		`onload=`,
		`onerror=`,
	}

	for _, pattern := range dangerousPatterns {
		matched, _ := regexp.MatchString(pattern, prompt)
		if matched {
			result.Errors = append(result.Errors, fmt.Sprintf("包含不安全内容: %s", pattern))
			result.Valid = false
		}
	}

	return result
}

// ValidationResult 验证结果
type ValidationResult struct {
	Valid    bool
	Errors   []string
	Warnings []string
}

// BuildCustomPrompt 构建自定义提示词
func (pb *PromptBuilder) BuildCustomPrompt(customPrompt, markdown string) string {
	if customPrompt == "" {
		return pb.BuildPrompt("", markdown, nil)
	}

	// 确保包含基本规则
	baseRules := `
## 重要规则
1. 所有 CSS 必须使用内联 style 属性
2. 不使用外部样式表或 <style> 标签
3. 只使用安全的 HTML 标签
4. 图片使用占位符格式：<!-- IMG:index -->，从0开始计数
5. 返回完整的 HTML，不需要其他说明文字`

	if !strings.Contains(customPrompt, "重要规则") && !strings.Contains(customPrompt, "内联 style") {
		customPrompt += baseRules
	}

	return pb.BuildPrompt(customPrompt, markdown, nil)
}

// ExtractMarkdownTitle 提取 Markdown 标题
func (pb *PromptBuilder) ExtractMarkdownTitle(markdown string) string {
	lines := strings.Split(markdown, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			title := strings.TrimLeft(line, "#")
			title = strings.TrimSpace(title)
			if title != "" {
				return title
			}
		}
		// 第一行非空非图片作为标题
		if line != "" && !strings.HasPrefix(line, "!") && !strings.HasPrefix(line, ">") {
			return line
		}
	}
	return "未命名文章"
}

// EstimateTokens 估算 token 数量
func (pb *PromptBuilder) EstimateTokens(text string) int {
	chineseChars := 0
	otherChars := 0

	for _, r := range text {
		if r >= 0x4e00 && r <= 0x9fff {
			chineseChars++
		} else {
			otherChars++
		}
	}

	// 粗略估算: 中文约 1 字符/token，英文约 4 字符/token
	return chineseChars + (otherChars / 4)
}

// GetVariable 获取变量定义
func (pb *PromptBuilder) GetVariable(name string) (*PromptVariable, error) {
	v, ok := pb.variables[name]
	if !ok {
		return nil, fmt.Errorf("variable not found: %s", name)
	}
	return v, nil
}

// ListVariables 列出所有变量
func (pb *PromptBuilder) ListVariables() []string {
	names := make([]string, 0, len(pb.variables))
	for name := range pb.variables {
		names = append(names, name)
	}
	return names
}

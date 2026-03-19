package converter

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/qiangmzsx/wechat-mcp/config"
	"github.com/qiangmzsx/wechat-mcp/provider"
	"go.uber.org/zap"
)

func TestExtractImages(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		want     int
	}{
		{
			name:     "empty markdown",
			markdown: "",
			want:     0,
		},
		{
			name:     "no images",
			markdown: "# Hello\n\nThis is a test.",
			want:     0,
		},
		{
			name:     "local image",
			markdown: "![alt](./image.png)",
			want:     1,
		},
		{
			name:     "online image",
			markdown: "![alt](https://example.com/image.png)",
			want:     1,
		},
		{
			name:     "ai generated image",
			markdown: "![alt](__generate:a beautiful sunset__)",
			want:     1,
		},
		{
			name: "multiple images",
			markdown: `![alt1](./image1.png)

Some text

![alt2](https://example.com/image2.png)

![alt3](__generate:AI prompt__)`,
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			images := ExtractImages(tt.markdown)
			if len(images) != tt.want {
				t.Errorf("ExtractImages() returned %d images, want %d", len(images), tt.want)
			}
		})
	}
}

func TestThemeManager(t *testing.T) {
	tm := NewThemeManager()

	// 测试获取内置主题
	theme, err := tm.GetTheme("default")
	if err != nil {
		t.Errorf("GetTheme(default) failed: %v", err)
	}
	if theme == nil {
		t.Error("GetTheme(default) returned nil")
	}

	// 测试获取不存在的主题
	_, err = tm.GetTheme("nonexistent")
	if err == nil {
		t.Error("GetTheme(nonexistent) should return error")
	}

	// 测试列出主题
	themes := tm.ListThemes()
	if len(themes) == 0 {
		t.Error("ListThemes() returned empty list")
	}

	// 测试获取 AI 提示词
	prompt, err := tm.GetAIPrompt("default")
	if err != nil {
		t.Errorf("GetAIPrompt(default) failed: %v", err)
	}
	if prompt == "" {
		t.Error("GetAIPrompt(default) returned empty prompt")
	}
}

func TestPromptBuilder(t *testing.T) {
	pb := NewPromptBuilder()

	// 测试构建 Prompt
	markdown := "# Test\n\nHello world"
	prompt := pb.BuildPrompt("", markdown, nil)
	if prompt == "" {
		t.Error("BuildPrompt returned empty")
	}

	// 测试提取标题
	title := pb.ExtractMarkdownTitle(markdown)
	if title != "Test" {
		t.Errorf("ExtractMarkdownTitle() = %s, want Test", title)
	}

	// 测试估算 token
	tokens := pb.EstimateTokens("你好世界")
	if tokens == 0 {
		t.Error("EstimateTokens() returned 0")
	}

	// 测试验证 Prompt
	validation := pb.ValidatePrompt(prompt)
	if !validation.Valid {
		t.Errorf("ValidatePrompt() failed: %v", validation.Errors)
	}
}

func TestSimpleConverter(t *testing.T) {
	converter := NewSimpleConverter()

	// 测试空 Markdown
	result := converter.Convert(&ConvertRequest{
		Markdown: "",
		Theme:    "default",
	})
	if result.Success {
		t.Error("Convert(empty) should fail")
	}

	// 测试有效请求
	result = converter.Convert(&ConvertRequest{
		Markdown: "# Hello\n\nThis is a test.",
		Theme:    "default",
	})
	if !result.Success {
		t.Errorf("Convert() failed: %s", result.Error)
	}

	// 测试图片提取
	images := converter.ExtractImages("![alt](./image.png)")
	if len(images) != 1 {
		t.Errorf("ExtractImages() returned %d images, want 1", len(images))
	}
}

func TestAIConverter_NewAIConverter(t *testing.T) {
	// 测试 converter 禁用时
	disabledCfg := &config.Config{
		Converter: config.ConverterConfig{
			Enabled: false,
		},
	}
	_, err := NewAIConverter(disabledCfg, zap.NewNop())
	if err == nil {
		t.Error("NewAIConverter with disabled converter should return error")
	}
	if err.Error() != "converter is disabled" {
		t.Errorf("Expected 'converter is disabled' error, got: %s", err.Error())
	}
}

// TestAIConverter_ProviderAnthropic 测试 Anthropic Provider
func TestAIConverter_ProviderAnthropic(t *testing.T) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping: ANTHROPIC_API_KEY not set")
	}

	logger := zap.NewNop()

	cfg := &config.Config{
		Converter: config.ConverterConfig{
			Enabled:      true,
			Provider:     string(provider.ProviderAnthropic),
			APIKey:       apiKey,
			Model:        "claude-sonnet-4-20250514",
			MaxTokens:    1024,
			Timeout:      60 * time.Second,
			DefaultTheme: "default",
		},
	}

	conv, err := NewAIConverter(cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create AI converter: %v", err)
	}

	// 验证 Provider 类型
	if conv.GetThemeManager() == nil {
		t.Error("GetThemeManager() should not be nil")
	}

	// 测试转换
	req := &ConvertRequest{
		Markdown: "# 测试\n\n你好世界",
		Theme:    "default",
	}

	result := conv.Convert(req)
	if !result.Success {
		t.Errorf("Convert failed: %s", result.Error)
	}

	if result.HTML == "" {
		t.Error("Expected non-empty HTML")
	}

	t.Logf("Anthropic Provider test passed, HTML length: %d", len(result.HTML))
}

// TestAIConverter_ProviderOpenAI 测试 OpenAI Provider
func TestAIConverter_ProviderOpenAI(t *testing.T) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping: OPENAI_API_KEY not set")
	}

	logger := zap.NewNop()

	cfg := &config.Config{
		Converter: config.ConverterConfig{
			Enabled:      true,
			Provider:     string(provider.ProviderOpenAI),
			APIKey:       apiKey,
			Model:        "gpt-4o-mini",
			MaxTokens:    1024,
			Timeout:      60 * time.Second,
			DefaultTheme: "default",
		},
	}

	conv, err := NewAIConverter(cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create AI converter: %v", err)
	}

	// 验证 Provider 类型
	if conv.GetThemeManager() == nil {
		t.Error("GetThemeManager() should not be nil")
	}

	// 测试转换
	req := &ConvertRequest{
		Markdown: "# 测试\n\n你好世界",
		Theme:    "default",
	}

	result := conv.Convert(req)
	if !result.Success {
		t.Errorf("Convert failed: %s", result.Error)
	}

	if result.HTML == "" {
		t.Error("Expected non-empty HTML")
	}

	t.Logf("OpenAI Provider test passed, HTML length: %d", len(result.HTML))
}

// TestAIConverter_ProviderDefault 测试默认 Provider (Anthropic)
var cfg = &config.Config{
	Converter: config.ConverterConfig{
		Enabled:      true,
		Provider:     "anthropic",
		Timeout:      300 * time.Second,
		APIKey:       os.Getenv("ANTHROPIC_API_KEY"),
		BaseURL:      "https://api.minimaxi.com/anthropic",
		Model:        "MiniMax-M2.5",
		DefaultTheme: "nord",
	},
}

func TestAIConverter_Convert(t *testing.T) {
	logger := zap.NewNop()

	conv, err := NewAIConverter(cfg, logger)
	if err != nil {
		t.Skipf("Skipping: cannot create AI converter: %v", err)
	}

	// 测试空 Markdown
	req := &ConvertRequest{
		Markdown: `
# Raphael Publish - 公众号排版大师

> 欢迎使用 Raphael Publish，一款专为**微信公众号**与**内容创作者**设计的现代 Markdown 排版引擎！

## 核心功能

### 1. 魔法粘贴

- **跨平台粘贴**：直接从**飞书、Notion、Word**甚至任意网页复制富文本，粘贴瞬间自动转换为纯净 Markdown
- **智能清洗**：自动剥离冗余样式和乱码，只保留段落、粗体、列表、代码块等核心结构
- **零学习成本**：不需要会写 Markdown，粘贴进来就能用
- **图片直贴**：支持直接粘贴截图或剪贴板图片（Ctrl/Cmd + V），自动插入 Markdown 图片

### 2. 多图排版

支持朋友圈式的多列网格布局，比如下面自然形成的两图并排。通过 wechatCompat 引擎这些多图也能在微信中完美呈现：

![](https://images.unsplash.com/photo-1550745165-9bc0b252726f?w=600&h=400&fit=crop)
![](https://images.unsplash.com/photo-1555066931-4365d14bab8c?w=600&h=400&fit=crop)

### 3. 30 套高定样式

告别同质化的白底模板，30 套精心打磨的视觉主题任你切换（下方仅展示部分代表风格）：

1. **极简与经典**：Mac 纯净白、微信公众号原生、Medium 博客风
2. **深度阅读**：Claude 燕麦色、NYT 纽约时报、Retro 复古羊皮纸
3. **极客与商务**：Stripe 硅谷风、飞书效率蓝、Linear 暗夜模式、Bloomberg 终端机

`,
		Theme: "sspai",
	}
	result := conv.Convert(req)
	if !result.Success {
		t.Errorf("Empty markdown should fail,%v", result.Error)
	}

	t.Log(result.HTML)
	t.Log(result.Images)
}

func TestAIConverter_Convert_OpenAI(t *testing.T) {
	logger := zap.NewNop()
	var cfg = &config.Config{
		Converter: config.ConverterConfig{
			Enabled:      true,
			Provider:     "openai",
			Timeout:      300 * time.Second,
			APIKey:       os.Getenv("OPENAI_API_KEY"),
			BaseURL:      "https://api.minimaxi.com/v1",
			Model:        "MiniMax-M2.5",
			DefaultTheme: "default",
		},
	}

	conv, err := NewAIConverter(cfg, logger)
	if err != nil {
		t.Skipf("Skipping: cannot create AI converter: %v", err)
	}

	// 测试空 Markdown
	req := &ConvertRequest{
		Markdown: `
# Raphael Publish - 公众号排版大师

> 欢迎使用 Raphael Publish，一款专为**微信公众号**与**内容创作者**设计的现代 Markdown 排版引擎！

## 核心功能

### 1. 魔法粘贴

- **跨平台粘贴**：直接从**飞书、Notion、Word**甚至任意网页复制富文本，粘贴瞬间自动转换为纯净 Markdown
- **智能清洗**：自动剥离冗余样式和乱码，只保留段落、粗体、列表、代码块等核心结构
- **零学习成本**：不需要会写 Markdown，粘贴进来就能用
- **图片直贴**：支持直接粘贴截图或剪贴板图片（Ctrl/Cmd + V），自动插入 Markdown 图片

### 2. 多图排版

支持朋友圈式的多列网格布局，比如下面自然形成的两图并排。通过 wechatCompat 引擎这些多图也能在微信中完美呈现：

![](https://images.unsplash.com/photo-1550745165-9bc0b252726f?w=600&h=400&fit=crop)
![](https://images.unsplash.com/photo-1555066931-4365d14bab8c?w=600&h=400&fit=crop)

### 3. 30 套高定样式

告别同质化的白底模板，30 套精心打磨的视觉主题任你切换（下方仅展示部分代表风格）：

1. **极简与经典**：Mac 纯净白、微信公众号原生、Medium 博客风
2. **深度阅读**：Claude 燕麦色、NYT 纽约时报、Retro 复古羊皮纸
3. **极客与商务**：Stripe 硅谷风、飞书效率蓝、Linear 暗夜模式、Bloomberg 终端机

`,
		Theme: "elegant",
	}
	result := conv.Convert(req)
	if !result.Success {
		t.Errorf("Empty markdown should fail,%v", result.Error)
	}

	t.Log(result.HTML)
	t.Log(result.Images)
}

func TestAIConverter_GetSystemPrompt(t *testing.T) {
	logger := zap.NewNop()

	conv, err := NewAIConverter(cfg, logger)
	if err != nil {
		t.Skipf("Skipping: cannot create AI converter: %v", err)
	}

	prompt, _ := conv.GetThemeManager().GetAIPrompt("tech")
	if prompt == "" {
		t.Error("GetAIPrompt returned empty")
	}
	t.Log(prompt)
}

func TestAIConverter_ExtractHTML(t *testing.T) {
	// 直接创建 aiConverter 实例
	themeMgr := NewThemeManager()
	conv := &aiConverter{
		config: cfg.Converter,
		theme:  themeMgr,
		prompt: NewPromptBuilder(),
	}

	tests := []struct {
		name      string
		content   string
		wantEmpty bool
	}{
		{
			name:      "html in code block",
			content:   "```html\n<p>Hello</p>\n```",
			wantEmpty: false,
		},
		{
			name:      "html in plain code block",
			content:   "```\n<html><body>Test</body></html>\n```",
			wantEmpty: false,
		},
		{
			name:      "plain html",
			content:   "<section><p>Hello</p></section>",
			wantEmpty: false,
		},
		{
			name:      "text content",
			content:   "Just some text without HTML",
			wantEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html := conv.extractHTML(tt.content)
			if tt.wantEmpty && html != "" {
				t.Errorf("extractHTML() = %s, want empty", html)
			}
			if !tt.wantEmpty && html == "" {
				t.Error("extractHTML() returned empty")
			}
		})
	}
}

func TestAIConverter_ProcessImagePlaceholders(t *testing.T) {
	themeMgr := NewThemeManager()
	conv := &aiConverter{
		config: cfg.Converter,
		theme:  themeMgr,
		prompt: NewPromptBuilder(),
	}

	html := `<section>
<p>Image 1:</p>
<!-- IMG:0 -->
<p>Image 2:</p>
<!-- IMG:1 -->
</section>`

	images := []ImageRef{
		{Index: 0, Original: "./test1.png"},
		{Index: 1, Original: "https://example.com/test2.png"},
	}

	result := conv.processImagePlaceholders(html, images)
	if result == "" {
		t.Error("processImagePlaceholders returned empty")
	}

	if !strings.Contains(result, "<img") {
		t.Error("result should contain img tag")
	}

	t.Log(result)
}

func TestAIConverter_ProcessImagePlaceholders_WithRealImage(t *testing.T) {
	themeMgr := NewThemeManager()
	conv := &aiConverter{
		config: cfg.Converter,
		theme:  themeMgr,
		prompt: NewPromptBuilder(),
	}

	html := `<section><p>Image:</p><!-- IMG:0 --></section>`

	images := []ImageRef{
		{Index: 0, Original: "https://picsum.photos/200/300", Type: ImageTypeOnline},
	}

	result := conv.processImagePlaceholders(html, images)
	if result == "" {
		t.Error("processImagePlaceholders returned empty")
	}

	if !strings.Contains(result, "data:image") {
		t.Error("result should contain base64 image data")
	}

	if strings.Contains(result, "picsum.photos") {
		t.Error("result should not contain original URL")
	}

	t.Logf("Result length: %d", len(result))
}

func TestAIConverter_Convert_WithRealAPI(t *testing.T) {
	// 这个测试需要真实的 API Key，如果环境变量中有 ANTHROPIC_API_KEY 则执行
	apiKey := os.Getenv("ANTHROPIC_API_KEY")

	if apiKey == "" {
		t.Skip("Skipping: ANTHROPIC_API_KEY not set")
	}

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	cfg := &config.Config{
		Converter: config.ConverterConfig{
			Enabled:      true,
			APIKey:       apiKey,
			DefaultTheme: "default",
		},
	}

	conv, err := NewAIConverter(cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create AI converter: %v", err)
	}

	// 测试真实转换
	req := &ConvertRequest{
		Markdown: "# 测试标题\n\n这是一个测试段落。\n\n- 列表项1\n- 列表项2",
		Theme:    "default",
	}

	result := conv.Convert(req)
	if !result.Success {
		t.Errorf("Convert failed: %s", result.Error)
	}

	if result.HTML == "" {
		t.Error("Convert returned empty HTML")
	}

	// 验证图片提取
	if len(result.Images) != 0 {
		t.Logf("Found %d images", len(result.Images))
	}

	t.Logf("HTML length: %d", len(result.HTML))
	t.Logf("Theme used: %s", result.Theme)
}

func TestReplaceImagesWithBase64_EscapedURL(t *testing.T) {
	htmlContent := `<img src="https://picsum.photos/200/300?w=600&amp;h=400&amp;fit=crop" alt="test" />`
	images := []ImageRef{
		{
			Original: "https://picsum.photos/200/300?w=600&h=400&fit=crop",
			Type:     ImageTypeOnline,
		},
	}

	result := ReplaceImagesWithBase64(htmlContent, images)

	t.Logf("Original HTML: %s", htmlContent)
	t.Logf("Result HTML length: %d", len(result))

	if strings.Contains(result, "https://picsum.photos") {
		t.Error("URL should be replaced with base64, but still contains original URL")
	}

	if !strings.Contains(result, "data:image") {
		t.Error("Result should contain base64 data URI")
	}
}

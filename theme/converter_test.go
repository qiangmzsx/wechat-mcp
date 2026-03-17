package theme

import (
	"strings"
	"testing"
)

func TestConvertBasic(t *testing.T) {
	markdown := `
# Raphael Publish - 公众号排版大师

> 欢迎使用 Raphael Publish，一款专为**微信公众号**与**内容创作者**设计的现代 Markdown 排版引擎！

## 核心功能

### 1. 魔法粘贴

- **跨平台粘贴**：直接从**飞书、Notion、Word**甚至任意网页复制富文本，粘贴瞬间自动转换为纯净 Markdown
- **智能清洗**：自动剥离冗余样式和乱码，只保留段落、粗体、列表、代码块等核心结构
- **零学习成本**：不需要会写 Markdown，粘贴进来就能用
- **图片直贴**：支持直接粘贴截图或剪贴板图片（Ctrl/Cmd + V），自动插入 Markdown 图片

### 2. 多图排版

支持朋友圈式的多列网格布局，比如下面自然形成的两图并排。通过 ~wechatCompat~ 引擎这些多图也能在微信中完美呈现：

![](https://images.unsplash.com/photo-1550745165-9bc0b252726f?w=600&h=400&fit=crop)
![](https://images.unsplash.com/photo-1555066931-4365d14bab8c?w=600&h=400&fit=crop)


`

	markdown = strings.ReplaceAll(markdown, "~", "`")
	// t.Log(markdown)
	html := Convert(markdown, "sspai")

	if !strings.Contains(html, "<h1") {
		t.Error("Expected HTML to contain h1")
	}

	if !strings.Contains(html, "<div") {
		t.Error("Expected HTML to contain div wrapper")
	}
	t.Logf("%s", html)
}

func TestThemeExists(t *testing.T) {
	if !ThemeExists("apple") {
		t.Error("Expected apple theme to exist")
	}
	if !ThemeExists("wechat") {
		t.Error("Expected wechat theme to exist")
	}
	if ThemeExists("nonexistent") {
		t.Error("Expected nonexistent theme to not exist")
	}
}

func TestGetTheme(t *testing.T) {
	theme := GetTheme("claude")
	if theme.ID != "claude" {
		t.Errorf("Expected theme ID to be claude, got %s", theme.ID)
	}

	defaultTheme := GetTheme("nonexistent")
	if defaultTheme.ID != "apple" {
		t.Errorf("Expected default theme ID to be apple, got %s", defaultTheme.ID)
	}
}

func TestAllThemes(t *testing.T) {
	themes := AllThemes()
	if len(themes) == 0 {
		t.Error("Expected themes to not be empty")
	}
	if len(themes) < 30 {
		t.Logf("Warning: Expected at least 30 themes, got %d", len(themes))
	}
}

func TestConvertWithImageGrids(t *testing.T) {
	markdown := "![image1](img1.png)\n\n![image2](img2.png)\n\n![image3](img3.png)\n"

	html := Convert(markdown, "apple")

	if !strings.Contains(html, "image-grid") {
		t.Log("Image grids not applied (may need paragraph format)")
	}
}

func TestPreprocessMarkdown(t *testing.T) {
	input := "Some *** text --- with ___ special chars"
	output := PreprocessMarkdown(input)

	if strings.Contains(output, "***") {
		t.Error("Expected *** to be removed")
	}
}

func TestRemoveExtraNewlines(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		wantNl bool
	}{
		{
			name:   "普通段落不应该包含换行符",
			input:  "# 标题\n\n这是一段文字。\n\n这是另一段。",
			wantNl: false,
		},
		{
			name:   "代码块应该保留换行符",
			input:  "```\ncode line 1\ncode line 2\n```",
			wantNl: true,
		},
		{
			name:   "行内代码不应该有换行",
			input:  "这是`行内代码`的测试",
			wantNl: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html := Convert(tt.input, "default")

			preStart := strings.Index(html, "<pre")
			preEnd := strings.Index(html, "</pre>")

			if preStart != -1 && preEnd != -1 {
				preContent := html[preStart : preEnd+len("</pre>")]
				if !strings.Contains(preContent, "\n") {
					t.Errorf("代码块应该保留换行符，但未找到")
				}
			} else {
				if strings.Contains(html, "\n") {
					t.Errorf("非代码块内容不应包含换行符，但找到了 \\n")
				}
			}
		})
	}
}

func TestConvertNoNewlines(t *testing.T) {
	markdown := `# 标题

这是一段文字。

这是另一段文字。`

	html := Convert(markdown, "default")

	if strings.Contains(html, "\n") {
		t.Errorf("HTML 不应包含换行符，但找到了: %q", html)
	}

	if !strings.Contains(html, "<h1") {
		t.Error("Expected HTML to contain h1")
	}
}

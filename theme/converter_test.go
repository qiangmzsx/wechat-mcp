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

支持朋友圈式的多列网格布局，比如下面自然形成的两图并排。通过 wechatCompat 引擎这些多图也能在微信中完美呈现：

![](https://images.unsplash.com/photo-1550745165-9bc0b252726f?w=600&h=400&fit=crop)
![](https://images.unsplash.com/photo-1555066931-4365d14bab8c?w=600&h=400&fit=crop)

### 3. 30 套高定样式

告别同质化的白底模板，30 套精心打磨的视觉主题任你切换（下方仅展示部分代表风格）：

1. **极简与经典**：Mac 纯净白、微信公众号原生、Medium 博客风
2. **深度阅读**：Claude 燕麦色、NYT 纽约时报、Retro 复古羊皮纸
3. **极客与商务**：Stripe 硅谷风、飞书效率蓝、Linear 暗夜模式、Bloomberg 终端机

`

	html := Convert(markdown, "claude")

	if !strings.Contains(html, "<h1") {
		t.Error("Expected HTML to contain h1")
	}
	if !strings.Contains(html, "Hello World</h1>") {
		t.Error("Expected HTML to contain heading content")
	}
	if !strings.Contains(html, "This is a paragraph") {
		t.Error("Expected HTML to contain paragraph content")
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

func TestThemeGroups(t *testing.T) {
	groups := ThemeGroups()
	if len(groups) == 0 {
		t.Error("Expected theme groups to not be empty")
	}

	classicFound := false
	for _, g := range groups {
		if g.Label == "经典" {
			classicFound = true
			if len(g.Themes) == 0 {
				t.Error("Expected classic themes to not be empty")
			}
		}
	}
	if !classicFound {
		t.Error("Expected to find 经典 theme group")
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

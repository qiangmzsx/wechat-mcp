package converter

import (
	"testing"
)

func TestAPIConverter_Convert(t *testing.T) {
	conv := NewAPIConverter()

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

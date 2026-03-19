package theme

import (
	"fmt"
	"strings"

	"github.com/qiangmzsx/wechat-mcp/logger"
	"github.com/yuin/goldmark"
	"go.uber.org/zap"
)

type Converter struct {
	theme       Theme
	enableGrids bool
	imageBase64 bool
}

type Option func(*Converter)

func WithTheme(theme Theme) Option {
	return func(c *Converter) {
		c.theme = theme
	}
}

func WithThemeID(themeID string) Option {
	return func(c *Converter) {
		c.theme = GetTheme(themeID)
	}
}

func EnableImageGrids(enable bool) Option {
	return func(c *Converter) {
		c.enableGrids = enable
	}
}

func ConvertImageToBase64(enable bool) Option {
	return func(c *Converter) {
		c.imageBase64 = enable
	}
}

func NewConverter(opts ...Option) *Converter {
	c := &Converter{
		theme:       GetTheme("apple"),
		enableGrids: true,
		imageBase64: false,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func PreprocessMarkdown(content string) string {
	replacer := strings.NewReplacer(
		"***", "",
		"---", "",
		"___", "",
		"****", "",
	)
	return replacer.Replace(content)
}

func (c *Converter) Convert(markdown string) string {
	logger.Debug("converting markdown", zap.String("theme", c.theme.ID), zap.Bool("grids", c.enableGrids))
	markdown = PreprocessMarkdown(markdown)

	md := goldmark.New()

	var buf strings.Builder
	if err := md.Convert([]byte(markdown), &buf); err != nil {
		logger.Error("goldmark conversion failed", zap.Error(err))
		return ""
	}

	html := buf.String()

	html = c.processSimpleStyles(html)

	if c.enableGrids {
		html = c.processImageGridsSimple(html)
	}

	html = c.processListStylesSimple(html)

	html = c.removeExtraNewlines(html)

	logger.Debug("markdown conversion completed", zap.String("theme", c.theme.ID))
	return html
}

func (c *Converter) processSimpleStyles(html string) string {
	style := c.theme.Styles

	var result strings.Builder
	result.WriteString("<div style=\"")
	result.WriteString(style["container"])
	result.WriteString("\">")
	result.WriteString(html)
	result.WriteString("</div>")

	return result.String()
}

func (c *Converter) processImageGridsSimple(html string) string {
	lines := strings.Split(html, "\n")
	var result []string
	var imageBuffer []string

	flushImages := func() {
		if len(imageBuffer) >= 2 {
			result = append(result, `<p class="image-grid" style="display: flex; justify-content: center; gap: 8px; margin: 24px 0; align-items: flex-start;">`)
			for _, img := range imageBuffer {
				width := 100.0 / float64(len(imageBuffer))
				spacing := 8.0 * float64(len(imageBuffer)-1) / float64(len(imageBuffer))
				w := fmt.Sprintf("%.2f", width)
				s := fmt.Sprintf("%.2f", spacing)
				imgLine := strings.Replace(img, "<img ", "<img style=\"width: calc("+w+"% - "+s+"px); margin: 0; border-radius: 8px; height: auto;\" ", 1)
				result = append(result, imgLine)
			}
			result = append(result, "</p>")
		} else {
			result = append(result, imageBuffer...)
		}
		imageBuffer = nil
	}

	inParagraph := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "<p>") && strings.Contains(trimmed, "<img") {
			if !inParagraph {
				flushImages()
				inParagraph = true
			}
			imageBuffer = append(imageBuffer, line)
		} else if inParagraph && trimmed == "" {
			continue
		} else {
			if inParagraph {
				flushImages()
				inParagraph = false
			}
			result = append(result, line)
		}
	}

	if len(imageBuffer) > 0 {
		flushImages()
	}

	return strings.Join(result, "\n")
}

func (c *Converter) processListStylesSimple(html string) string {
	style := c.theme.Styles

	html = strings.Replace(html, "<ul>", "<ul style=\""+style["ul"]+" list-style-type: disc !important;\">", -1)
	html = strings.Replace(html, "<ol>", "<ol style=\""+style["ol"]+" list-style-type: decimal !important;\">", -1)

	selectors := []string{"h1", "h2", "h3", "h4", "h5", "h6", "p", "strong", "em", "a", "li", "blockquote", "code", "pre", "hr", "img", "table", "th", "td", "tr"}

	for _, sel := range selectors {
		if s, ok := style[sel]; ok {
			html = addStyleToAllElements(html, sel, s)
		}
	}

	return html
}

func addStyleToAllElements(html, tag, styleValue string) string {
	searchOpen := "<" + tag + ">"
	searchOpenWithStyle := "<" + tag + " style=\""

	if !strings.Contains(html, searchOpenWithStyle) {
		replaceOpen := "<" + tag + " style=\"" + styleValue + "\">"
		html = strings.ReplaceAll(html, searchOpen, replaceOpen)
	}

	html = addStyleToOpenTags(html, tag, styleValue)

	return html
}

func addStyleToOpenTags(html, tag, styleValue string) string {
	tagWithSpace := "<" + tag + " "

	var result strings.Builder
	searchStart := 0
	for {
		idx := strings.Index(html[searchStart:], tagWithSpace)
		if idx == -1 {
			result.WriteString(html[searchStart:])
			break
		}
		idx += searchStart

		result.WriteString(html[searchStart:idx])

		tagEnd := strings.Index(html[idx:], ">")
		if tagEnd == -1 {
			break
		}
		tagEnd += idx

		tagContent := html[idx+len(tagWithSpace) : tagEnd]
		if strings.Contains(tagContent, "style=\"") || strings.Contains(tagContent, " style=\"") {
			result.WriteString(html[idx : tagEnd+1])
			searchStart = tagEnd + 1
			continue
		}

		result.WriteString("<" + tag + " style=\"" + styleValue + "\" " + tagContent + ">")

		searchStart = tagEnd + 1
	}

	return result.String()
}

func (c *Converter) removeExtraNewlines(html string) string {
	preTag := "<pre"
	preEndTag := "</pre>"

	var result strings.Builder
	pos := 0

	for {
		preStart := strings.Index(html[pos:], preTag)
		if preStart == -1 {
			result.WriteString(c.stripNewlines(html[pos:]))
			break
		}
		preStart += pos

		result.WriteString(c.stripNewlines(html[pos:preStart]))

		preEnd := strings.Index(html[preStart:], preEndTag)
		if preEnd == -1 {
			result.WriteString(html[preStart:])
			break
		}
		preEnd += preStart + len(preEndTag)

		result.WriteString(html[preStart:preEnd])
		pos = preEnd
	}

	return result.String()
}

func (c *Converter) stripNewlines(s string) string {
	return strings.ReplaceAll(s, "\n", "")
}

func Convert(markdown string, themeID string) string {
	c := NewConverter(WithThemeID(themeID))
	return c.Convert(markdown)
}

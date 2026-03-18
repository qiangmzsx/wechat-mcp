package converter

import (
	"github.com/qiangmzsx/wechat-mcp/theme"
	"go.uber.org/zap"
)

type apiConverter struct {
	log      *zap.Logger
	themeMgr ThemeManager
}

func NewAPIConverter(log *zap.Logger) Converter {
	return &apiConverter{
		log:      log,
		themeMgr: NewThemeManager(),
	}
}

func (c *apiConverter) Convert(req *ConvertRequest) *ConvertResult {
	result := &ConvertResult{
		Theme:   req.Theme,
		Success: false,
	}

	if req.Markdown == "" {
		result.Error = "markdown content cannot be empty"
		return result
	}

	themeID := req.Theme
	if themeID == "" {
		themeID = "default"
	}

	if c.log != nil {
		c.log.Debug("API converter: converting markdown",
			zap.String("theme", themeID),
			zap.Int("markdown_length", len(req.Markdown)))
	}

	html := theme.ConvertMarkdown(req.Markdown, themeID)
	if html == "" {
		result.Error = "theme conversion failed"
		return result
	}

	images := c.ExtractImages(req.Markdown)

	if c.log != nil {
		c.log.Debug("API converter: extracted images",
			zap.Int("image_count", len(images)))
	}

	// 将图片转换为 base64 嵌入 HTML
	html = ReplaceImagesWithBase64WithLogger(html, images, c.log)

	result.Images = images
	result.HTML = FormatHTML(html)
	result.Success = true

	if c.log != nil {
		c.log.Info("API converter: conversion completed",
			zap.String("theme", themeID),
			zap.Int("image_count", len(images)),
			zap.Int("html_length", len(html)))
	}

	return result
}

func (c *apiConverter) ExtractImages(markdown string) []ImageRef {
	return ExtractImages(markdown)
}

func (c *apiConverter) GetThemeManager() ThemeManager {
	return c.themeMgr
}

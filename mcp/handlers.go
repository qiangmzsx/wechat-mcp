package mcp

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/qiangmzsx/wechat-mcp/config"
	"github.com/qiangmzsx/wechat-mcp/converter"
	"github.com/qiangmzsx/wechat-mcp/internal/util"
	"github.com/qiangmzsx/wechat-mcp/logger"
	"github.com/qiangmzsx/wechat-mcp/wechat"
	"github.com/silenceper/wechat/v2/officialaccount/draft"
	"go.uber.org/zap"
)

// uploadMaterialHandler 上传素材处理器
func (s *Server) uploadMaterialHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	startTime := time.Now()
	toolName := "upload_material"

	logger.Info("Tool called", zap.String("tool", toolName))

	args := request.GetArguments()
	filePath, ok := args["file_path"].(string)
	if !ok || filePath == "" {
		logger.Warn("Tool argument missing", zap.String("tool", toolName), zap.String("argument", "file_path"))
		return mcp.NewToolResultError("file_path is required"), nil
	}

	retry, _ := args["retry"].(bool)
	logger.Debug("Tool arguments",
		zap.String("tool", toolName),
		zap.String("file_path", filePath),
		zap.Bool("retry", retry),
	)

	var result *wechat.UploadMaterialResult
	var err error

	if retry {
		logger.Debug("Using retry mechanism", zap.String("tool", toolName), zap.Int("max_retries", 3))
		result, err = s.svc.UploadMaterialWithRetry(filePath, 3)
	} else {
		result, err = s.svc.UploadMaterial(filePath)
	}

	duration := time.Since(startTime)

	if err != nil {
		logger.Error("Tool execution failed",
			zap.String("tool", toolName),
			zap.Error(err),
			zap.Duration("duration", duration),
		)
		return mcp.NewToolResultError(fmt.Sprintf("upload failed: %v", err)), nil
	}

	logger.Info("Tool executed successfully",
		zap.String("tool", toolName),
		zap.String("media_id", util.MaskID(result.MediaID)),
		zap.Duration("duration", duration),
	)

	return mcp.NewToolResultText(fmt.Sprintf("素材上传成功!\nMediaID: %s\nURL: %s", result.MediaID, result.WechatURL)), nil
}

// createDraftHandler 创建草稿处理器
func (s *Server) createDraftHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	startTime := time.Now()
	toolName := "create_draft"

	logger.Info("Tool called", zap.String("tool", toolName))

	args := request.GetArguments()

	title, ok := args["title"].(string)
	if !ok || title == "" {
		logger.Warn("Tool argument missing", zap.String("tool", toolName), zap.String("argument", "title"))
		return mcp.NewToolResultError("title is required"), nil
	}

	content, ok := args["content"].(string)
	if !ok || content == "" {
		logger.Warn("Tool argument missing", zap.String("tool", toolName), zap.String("argument", "content"))
		return mcp.NewToolResultError("content is required"), nil
	}

	author, _ := args["author"].(string)
	digest, _ := args["digest"].(string)
	contentSourceURL, _ := args["content_source_url"].(string)
	thumbMediaID, _ := args["thumb_media_id"].(string)

	var needOpenComment uint = 0
	var onlyFansCanComment uint = 0
	if needOpen, ok := args["need_open_comment"].(bool); ok && needOpen {
		needOpenComment = 1
	}
	if onlyFans, ok := args["only_fans_can_comment"].(bool); ok && onlyFans {
		onlyFansCanComment = 1
	}

	logger.Debug("Tool arguments",
		zap.String("tool", toolName),
		zap.String("title", title),
		zap.String("author", author),
		zap.String("thumb_media_id", thumbMediaID),
		zap.Uint("need_open_comment", needOpenComment),
		zap.Uint("only_fans_can_comment", onlyFansCanComment),
	)

	article := &draft.Article{
		Title:              title,
		Content:            content,
		Author:             author,
		Digest:             digest,
		ContentSourceURL:   contentSourceURL,
		ThumbMediaID:       thumbMediaID,
		NeedOpenComment:    needOpenComment,
		OnlyFansCanComment: onlyFansCanComment,
	}

	result, err := s.svc.CreateDraft([]*draft.Article{article})
	duration := time.Since(startTime)

	if err != nil {
		logger.Error("Tool execution failed",
			zap.String("tool", toolName),
			zap.Error(err),
			zap.Duration("duration", duration),
		)
		return mcp.NewToolResultError(fmt.Sprintf("create draft failed: %v", err)), nil
	}

	logger.Info("Tool executed successfully",
		zap.String("tool", toolName),
		zap.String("media_id", util.MaskID(result.MediaID)),
		zap.Duration("duration", duration),
	)

	return mcp.NewToolResultText(fmt.Sprintf("草稿创建成功!\nMediaID: %s\n查看链接: %s", result.MediaID, result.DraftURL)), nil
}

// createNewspicDraftHandler 创建小绿书草稿处理器
func (s *Server) createNewspicDraftHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	startTime := time.Now()
	toolName := "create_newspic_draft"

	logger.Info("Tool called", zap.String("tool", toolName))

	args := request.GetArguments()

	title, ok := args["title"].(string)
	if !ok || title == "" {
		logger.Warn("Tool argument missing", zap.String("tool", toolName), zap.String("argument", "title"))
		return mcp.NewToolResultError("title is required"), nil
	}

	content, ok := args["content"].(string)
	if !ok || content == "" {
		logger.Warn("Tool argument missing", zap.String("tool", toolName), zap.String("argument", "content"))
		return mcp.NewToolResultError("content is required"), nil
	}

	imagePaths, ok := args["image_paths"].([]any)
	if !ok || len(imagePaths) == 0 {
		logger.Warn("Tool argument missing", zap.String("tool", toolName), zap.String("argument", "image_paths"))
		return mcp.NewToolResultError("image_paths is required and must not be empty"), nil
	}

	logger.Debug("Uploading images for newspic",
		zap.String("tool", toolName),
		zap.Int("image_count", len(imagePaths)),
	)

	imageList := make([]wechat.NewspicImageItem, 0, len(imagePaths))
	for i, path := range imagePaths {
		pathStr, ok := path.(string)
		if !ok {
			logger.Warn("Invalid image path type", zap.String("tool", toolName), zap.Int("index", i))
			continue
		}

		logger.Debug("Uploading image", zap.String("tool", toolName), zap.Int("index", i), zap.String("path", pathStr))

		result, err := s.svc.UploadMaterial(pathStr)
		if err != nil {
			logger.Error("Image upload failed",
				zap.String("tool", toolName),
				zap.Int("index", i),
				zap.Error(err),
			)
			return mcp.NewToolResultError(fmt.Sprintf("upload image failed: %v", err)), nil
		}

		imageList = append(imageList, wechat.NewspicImageItem{
			ImageMediaID: result.MediaID,
		})
	}

	newspicArticle := wechat.NewspicArticle{
		Title:       title,
		Content:     content,
		ArticleType: "newspic",
		ImageInfo: wechat.NewspicImageInfo{
			ImageList: imageList,
		},
	}

	logger.Debug("Creating newspic draft", zap.String("tool", toolName))

	result, err := s.svc.CreateNewspicDraft([]wechat.NewspicArticle{newspicArticle})
	duration := time.Since(startTime)

	if err != nil {
		logger.Error("Tool execution failed",
			zap.String("tool", toolName),
			zap.Error(err),
			zap.Duration("duration", duration),
		)
		return mcp.NewToolResultError(fmt.Sprintf("create newspic draft failed: %v", err)), nil
	}

	logger.Info("Tool executed successfully",
		zap.String("tool", toolName),
		zap.String("media_id", util.MaskID(result.MediaID)),
		zap.Duration("duration", duration),
	)

	return mcp.NewToolResultText(fmt.Sprintf("小绿书草稿创建成功!\nMediaID: %s\n查看链接: %s", result.MediaID, result.DraftURL)), nil
}

// getAccessTokenHandler 获取AccessToken处理器
func (s *Server) getAccessTokenHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	startTime := time.Now()
	toolName := "get_access_token"

	logger.Info("Tool called", zap.String("tool", toolName))

	result, err := s.svc.GetAccessToken()
	duration := time.Since(startTime)

	if err != nil {
		logger.Error("Tool execution failed",
			zap.String("tool", toolName),
			zap.Error(err),
			zap.Duration("duration", duration),
		)
		return mcp.NewToolResultError(fmt.Sprintf("get access token failed: %v", err)), nil
	}

	logger.Info("Tool executed successfully",
		zap.String("tool", toolName),
		zap.Duration("duration", duration),
	)

	return mcp.NewToolResultText(fmt.Sprintf("AccessToken: %s\nExpiresIn: %d秒", result.AccessToken, result.ExpiresIn)), nil
}

// downloadFileHandler 下载文件处理器
func (s *Server) downloadFileHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	startTime := time.Now()
	toolName := "download_file"

	logger.Info("Tool called", zap.String("tool", toolName))

	args := request.GetArguments()

	urlOrPath, ok := args["url_or_path"].(string)
	if !ok || urlOrPath == "" {
		logger.Warn("Tool argument missing", zap.String("tool", toolName), zap.String("argument", "url_or_path"))
		return mcp.NewToolResultError("url_or_path is required"), nil
	}

	logger.Debug("Downloading file",
		zap.String("tool", toolName),
		zap.String("url_or_path", urlOrPath),
	)

	path, err := wechat.DownloadFile(urlOrPath)
	duration := time.Since(startTime)

	if err != nil {
		logger.Error("Tool execution failed",
			zap.String("tool", toolName),
			zap.Error(err),
			zap.Duration("duration", duration),
		)
		return mcp.NewToolResultError(fmt.Sprintf("download file failed: %v", err)), nil
	}

	logger.Info("Tool executed successfully",
		zap.String("tool", toolName),
		zap.String("local_path", path),
		zap.Duration("duration", duration),
	)

	return mcp.NewToolResultText(fmt.Sprintf("文件路径: %s", path)), nil
}

// convertMarkdownHandler Markdown转换处理器
func (s *Server) convertMarkdownHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	startTime := time.Now()
	toolName := "convert_markdown"

	logger.Info("Tool called", zap.String("tool", toolName))

	args := request.GetArguments()

	markdown, ok := args["markdown"].(string)
	if !ok || markdown == "" {
		logger.Warn("Tool argument missing", zap.String("tool", toolName), zap.String("argument", "markdown"))
		return mcp.NewToolResultError("markdown is required"), nil
	}

	themeArg, _ := args["theme"].(string)
	customPrompt, _ := args["custom_prompt"].(string)
	converterTypeStr, _ := args["converter_type"].(string)

	logger.Debug("Tool arguments",
		zap.String("tool", toolName),
		zap.String("theme", themeArg),
		zap.Bool("has_custom_prompt", customPrompt != ""),
		zap.String("converter_type", converterTypeStr),
		zap.Int("markdown_length", len(markdown)),
	)

	req := &converter.ConvertRequest{
		Markdown:     markdown,
		Theme:        themeArg,
		CustomPrompt: customPrompt,
	}

	var result *converter.ConvertResult
	convType := req.ConverterType
	if convType == "" {
		convType = s.config.Converter.Type
	}

	if convType == config.ConverterTypeAI {
		logger.Debug("Using AI converter", zap.String("tool", toolName))
		result = s.aiConverter.Convert(req)
	} else {
		logger.Debug("Using API converter", zap.String("tool", toolName))
		result = s.converter.Convert(req)
	}
	duration := time.Since(startTime)

	if !result.Success {
		logger.Error("Tool execution failed",
			zap.String("tool", toolName),
			zap.String("error", result.Error),
			zap.Duration("duration", duration),
		)
		return mcp.NewToolResultError(fmt.Sprintf("convert failed: %s", result.Error)), nil
	}

	logger.Info("Tool executed successfully",
		zap.String("tool", toolName),
		zap.String("theme", result.Theme),
		zap.Int("image_count", len(result.Images)),
		zap.Int("html_length", len(result.HTML)),
		zap.Duration("duration", duration),
	)

	return mcp.NewToolResultText(result.HTML), nil
}

// listThemesHandler 列出主题处理器
func (s *Server) listThemesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	toolName := "list_themes"

	logger.Info("Tool called", zap.String("tool", toolName))

	themeMgr := s.converter.GetThemeManager()
	themes := themeMgr.ListThemes()

	logger.Info("Tool executed successfully",
		zap.String("tool", toolName),
		zap.Int("theme_count", len(themes)),
	)

	var result strings.Builder
	result.WriteString("可用主题:\n")
	for _, theme := range themes {
		result.WriteString(fmt.Sprintf("- %s\n", theme))
	}

	return mcp.NewToolResultText(result.String()), nil
}

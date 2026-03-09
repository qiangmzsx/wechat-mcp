package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/qiangmzsx/wechat-mcp/config"
	"github.com/qiangmzsx/wechat-mcp/converter"
	"github.com/qiangmzsx/wechat-mcp/wechat"
	"github.com/silenceper/wechat/v2/officialaccount/draft"
	"go.uber.org/zap"
)

// Server MCP服务器
type Server struct {
	svc       *wechat.Service
	converter converter.Converter
	config    *config.Config
	log       *zap.Logger
}

// New 创建MCP服务器
func New(cfg *config.Config, logger *zap.Logger) *Server {
	svc := wechat.NewService(cfg, logger)

	// 初始化 converter
	var conv converter.Converter
	if cfg.Converter.Enabled {
		var err error
		conv, err = converter.NewAIConverter(cfg, logger)
		if err != nil {
			logger.Warn("converter initialization failed, using simple converter", zap.Error(err))
			conv = converter.NewSimpleConverter()
		}
	} else {
		conv = converter.NewSimpleConverter()
	}

	return &Server{
		svc:       svc,
		converter: conv,
		config:    cfg,
		log:       logger,
	}
}

// Run 启动服务器
func (s *Server) Run() error {
	s.log.Info("Initializing MCP Server",
		zap.String("name", "WeChat MCP Server"),
		zap.String("version", "1.0.0"),
		zap.String("protocol", s.config.MCP.Protocol),
	)

	mcpServer := server.NewMCPServer(
		"WeChat MCP Server",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
		server.WithRecovery(),
	)

	// 注册工具
	s.registerTools(mcpServer)
	s.log.Info("MCP tools registered successfully")

	// 根据配置选择协议
	switch s.config.MCP.Protocol {
	case "http":
		return s.runHTTP(mcpServer)
	case "sse":
		return s.runSSE(mcpServer)
	default:
		s.log.Info("Starting STDIO server (waiting for client connection...)")
		return server.ServeStdio(mcpServer)
	}
}

// runHTTP 启动HTTP服务器
func (s *Server) runHTTP(mcpServer *server.MCPServer) error {
	addr := fmt.Sprintf("%s:%d", s.config.MCP.Host, s.config.MCP.Port)
	s.log.Info("Starting StreamableHTTP server",
		zap.String("addr", addr),
	)

	// 创建自定义 HTTP mux，添加健康检查端点
	mux := http.NewServeMux()
	mux.Handle("/mcp", server.NewStreamableHTTPServer(mcpServer))
	mux.HandleFunc("/health", s.healthHandler)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.log.Error("HTTP server failed to start", zap.Error(err))
		return err
	}
	return nil
}

// runSSE 启动SSE服务器
func (s *Server) runSSE(mcpServer *server.MCPServer) error {
	addr := fmt.Sprintf("%s:%d", s.config.MCP.Host, s.config.MCP.Port)
	s.log.Info("Starting SSE server",
		zap.String("addr", addr),
	)

	// 创建自定义 HTTP mux，添加健康检查端点
	mux := http.NewServeMux()
	mux.Handle("/sse", server.NewSSEServer(mcpServer))
	mux.HandleFunc("/health", s.healthHandler)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.log.Error("SSE server failed to start", zap.Error(err))
		return err
	}
	return nil
}

// healthHandler 健康检查处理器
func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"service":   "WeChat MCP Server",
		"version":   "1.0.0",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		s.log.Warn("Failed to encode health response", zap.Error(err))
	}
}

// registerTools 注册所有工具
func (s *Server) registerTools(mcpServer *server.MCPServer) {
	s.log.Debug("Registering MCP tools...")

	// 1. 上传素材工具
	mcpServer.AddTool(
		mcp.NewTool("upload_material",
			mcp.WithDescription("上传图片素材到微信公众号"),
			mcp.WithString("file_path",
				mcp.Required(),
				mcp.Description("本地文件路径或HTTP URL"),
			),
			mcp.WithBoolean("retry",
				mcp.Description("是否启用重试，默认false"),
			),
		),
		s.uploadMaterialHandler,
	)
	s.log.Debug("Tool registered: upload_material")

	// 2. 创建草稿工具
	mcpServer.AddTool(
		mcp.NewTool("create_draft",
			mcp.WithDescription("创建微信公众号文章草稿"),
			mcp.WithString("title",
				mcp.Required(),
				mcp.Description("文章标题"),
			),
			mcp.WithString("content",
				mcp.Required(),
				mcp.Description("文章正文内容(支持HTML)"),
			),
			mcp.WithString("author",
				mcp.Description("作者名称"),
			),
			mcp.WithString("digest",
				mcp.Description("文章摘要"),
			),
			mcp.WithString("content_source_url",
				mcp.Description("原文链接"),
			),
			mcp.WithString("thumb_media_id",
				mcp.Required(),
				mcp.Description("封面图片media_id"),
			),
			mcp.WithBoolean("need_open_comment",
				mcp.Description("是否开启评论，true开启，false关闭"),
			),
			mcp.WithBoolean("only_fans_can_comment",
				mcp.Description("是否仅粉丝可评论，true仅粉丝可评论，false所有人可评论"),
			),
		),
		s.createDraftHandler,
	)
	s.log.Debug("Tool registered: create_draft")

	// 3. 创建小绿书草稿工具
	mcpServer.AddTool(
		mcp.NewTool("create_newspic_draft",
			mcp.WithDescription("创建微信公众号小绿书草稿"),
			mcp.WithString("title",
				mcp.Required(),
				mcp.Description("文章标题"),
			),
			mcp.WithString("content",
				mcp.Required(),
				mcp.Description("文章正文内容"),
			),
			mcp.WithArray("image_paths",
				mcp.Required(),
				mcp.Description("图片路径列表，支持本地路径或HTTP URL，不要超过20张图片"),
				mcp.Items(map[string]any{"type": "string"}),
			),
		),
		s.createNewspicDraftHandler,
	)
	s.log.Debug("Tool registered: create_newspic_draft")

	// 4. 获取AccessToken工具
	mcpServer.AddTool(
		mcp.NewTool("get_access_token",
			mcp.WithDescription("获取微信公众号AccessToken（用于调试）"),
		),
		s.getAccessTokenHandler,
	)
	s.log.Debug("Tool registered: get_access_token")

	// 5. 下载文件工具
	mcpServer.AddTool(
		mcp.NewTool("download_file",
			mcp.WithDescription("下载文件到临时目录，或验证本地文件路径"),
			mcp.WithString("url_or_path",
				mcp.Required(),
				mcp.Description("文件URL或本地路径"),
			),
		),
		s.downloadFileHandler,
	)
	s.log.Debug("Tool registered: download_file")

	// 6. Markdown转HTML工具
	mcpServer.AddTool(
		mcp.NewTool("convert_markdown",
			mcp.WithDescription("将Markdown内容转换为微信公众号兼容的HTML"),
			mcp.WithString("markdown",
				mcp.Required(),
				mcp.Description("Markdown内容"),
			),
			mcp.WithString("theme",
				mcp.Description("主题名称: default, elegant, tech, minimalist"),
			),
			mcp.WithString("custom_prompt",
				mcp.Description("自定义提示词(可选)"),
			),
		),
		s.convertMarkdownHandler,
	)
	s.log.Debug("Tool registered: convert_markdown")

	// 7. 列出可用主题工具
	mcpServer.AddTool(
		mcp.NewTool("list_themes",
			mcp.WithDescription("列出所有可用的主题"),
		),
		s.listThemesHandler,
	)
	s.log.Debug("Tool registered: list_themes")
}

// ========== 工具处理器 ==========

// uploadMaterialHandler 上传素材处理器
func (s *Server) uploadMaterialHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	startTime := time.Now()
	toolName := "upload_material"

	s.log.Info("Tool called", zap.String("tool", toolName))

	args := request.GetArguments()
	filePath, ok := args["file_path"].(string)
	if !ok || filePath == "" {
		s.log.Warn("Tool argument missing", zap.String("tool", toolName), zap.String("argument", "file_path"))
		return mcp.NewToolResultError("file_path is required"), nil
	}

	retry, _ := args["retry"].(bool)
	s.log.Debug("Tool arguments",
		zap.String("tool", toolName),
		zap.String("file_path", filePath),
		zap.Bool("retry", retry),
	)

	var result *wechat.UploadMaterialResult
	var err error

	if retry {
		s.log.Debug("Using retry mechanism", zap.String("tool", toolName), zap.Int("max_retries", 3))
		result, err = s.svc.UploadMaterialWithRetry(filePath, 3)
	} else {
		result, err = s.svc.UploadMaterial(filePath)
	}

	duration := time.Since(startTime)

	if err != nil {
		s.log.Error("Tool execution failed",
			zap.String("tool", toolName),
			zap.Error(err),
			zap.Duration("duration", duration),
		)
		return mcp.NewToolResultError(fmt.Sprintf("upload failed: %v", err)), nil
	}

	s.log.Info("Tool executed successfully",
		zap.String("tool", toolName),
		zap.String("media_id", maskMediaID(result.MediaID)),
		zap.Duration("duration", duration),
	)

	return mcp.NewToolResultText(fmt.Sprintf("素材上传成功!\nMediaID: %s\nURL: %s", result.MediaID, result.WechatURL)), nil
}

// createDraftHandler 创建草稿处理器
func (s *Server) createDraftHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	startTime := time.Now()
	toolName := "create_draft"

	s.log.Info("Tool called", zap.String("tool", toolName))

	args := request.GetArguments()

	title, ok := args["title"].(string)
	if !ok || title == "" {
		s.log.Warn("Tool argument missing", zap.String("tool", toolName), zap.String("argument", "title"))
		return mcp.NewToolResultError("title is required"), nil
	}

	content, ok := args["content"].(string)
	if !ok || content == "" {
		s.log.Warn("Tool argument missing", zap.String("tool", toolName), zap.String("argument", "content"))
		return mcp.NewToolResultError("content is required"), nil
	}

	// 获取可选参数
	author, _ := args["author"].(string)
	digest, _ := args["digest"].(string)
	contentSourceURL, _ := args["content_source_url"].(string)
	thumbMediaID, _ := args["thumb_media_id"].(string)

	// 评论相关参数
	var needOpenComment uint = 0
	var onlyFansCanComment uint = 0
	if needOpen, ok := args["need_open_comment"].(bool); ok && needOpen {
		needOpenComment = 1
	}
	if onlyFans, ok := args["only_fans_can_comment"].(bool); ok && onlyFans {
		onlyFansCanComment = 1
	}

	s.log.Debug("Tool arguments",
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
		s.log.Error("Tool execution failed",
			zap.String("tool", toolName),
			zap.Error(err),
			zap.Duration("duration", duration),
		)
		return mcp.NewToolResultError(fmt.Sprintf("create draft failed: %v", err)), nil
	}

	s.log.Info("Tool executed successfully",
		zap.String("tool", toolName),
		zap.String("media_id", maskMediaID(result.MediaID)),
		zap.Duration("duration", duration),
	)

	return mcp.NewToolResultText(fmt.Sprintf("草稿创建成功!\nMediaID: %s\n查看链接: %s", result.MediaID, result.DraftURL)), nil
}

// createNewspicDraftHandler 创建小绿书草稿处理器
func (s *Server) createNewspicDraftHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	startTime := time.Now()
	toolName := "create_newspic_draft"

	s.log.Info("Tool called", zap.String("tool", toolName))

	args := request.GetArguments()

	title, ok := args["title"].(string)
	if !ok || title == "" {
		s.log.Warn("Tool argument missing", zap.String("tool", toolName), zap.String("argument", "title"))
		return mcp.NewToolResultError("title is required"), nil
	}

	content, ok := args["content"].(string)
	if !ok || content == "" {
		s.log.Warn("Tool argument missing", zap.String("tool", toolName), zap.String("argument", "content"))
		return mcp.NewToolResultError("content is required"), nil
	}

	// 处理图片
	imagePaths, ok := args["image_paths"].([]any)
	if !ok || len(imagePaths) == 0 {
		s.log.Warn("Tool argument missing", zap.String("tool", toolName), zap.String("argument", "image_paths"))
		return mcp.NewToolResultError("image_paths is required and must not be empty"), nil
	}

	s.log.Debug("Uploading images for newspic",
		zap.String("tool", toolName),
		zap.Int("image_count", len(imagePaths)),
	)

	// 上传图片素材
	imageList := make([]wechat.NewspicImageItem, 0, len(imagePaths))
	for i, path := range imagePaths {
		pathStr, ok := path.(string)
		if !ok {
			s.log.Warn("Invalid image path type", zap.String("tool", toolName), zap.Int("index", i))
			continue
		}

		s.log.Debug("Uploading image", zap.String("tool", toolName), zap.Int("index", i), zap.String("path", pathStr))

		result, err := s.svc.UploadMaterial(pathStr)
		if err != nil {
			s.log.Error("Image upload failed",
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

	s.log.Debug("Creating newspic draft", zap.String("tool", toolName))

	result, err := s.svc.CreateNewspicDraft([]wechat.NewspicArticle{newspicArticle})
	duration := time.Since(startTime)

	if err != nil {
		s.log.Error("Tool execution failed",
			zap.String("tool", toolName),
			zap.Error(err),
			zap.Duration("duration", duration),
		)
		return mcp.NewToolResultError(fmt.Sprintf("create newspic draft failed: %v", err)), nil
	}

	s.log.Info("Tool executed successfully",
		zap.String("tool", toolName),
		zap.String("media_id", maskMediaID(result.MediaID)),
		zap.Duration("duration", duration),
	)

	return mcp.NewToolResultText(fmt.Sprintf("小绿书草稿创建成功!\nMediaID: %s\n查看链接: %s", result.MediaID, result.DraftURL)), nil
}

// getAccessTokenHandler 获取AccessToken处理器
func (s *Server) getAccessTokenHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	startTime := time.Now()
	toolName := "get_access_token"

	s.log.Info("Tool called", zap.String("tool", toolName))

	result, err := s.svc.GetAccessToken()
	duration := time.Since(startTime)

	if err != nil {
		s.log.Error("Tool execution failed",
			zap.String("tool", toolName),
			zap.Error(err),
			zap.Duration("duration", duration),
		)
		return mcp.NewToolResultError(fmt.Sprintf("get access token failed: %v", err)), nil
	}

	s.log.Info("Tool executed successfully",
		zap.String("tool", toolName),
		zap.Duration("duration", duration),
	)

	return mcp.NewToolResultText(fmt.Sprintf("AccessToken: %s\nExpiresIn: %d秒", result.AccessToken, result.ExpiresIn)), nil
}

// downloadFileHandler 下载文件处理器
func (s *Server) downloadFileHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	startTime := time.Now()
	toolName := "download_file"

	s.log.Info("Tool called", zap.String("tool", toolName))

	args := request.GetArguments()

	urlOrPath, ok := args["url_or_path"].(string)
	if !ok || urlOrPath == "" {
		s.log.Warn("Tool argument missing", zap.String("tool", toolName), zap.String("argument", "url_or_path"))
		return mcp.NewToolResultError("url_or_path is required"), nil
	}

	s.log.Debug("Downloading file",
		zap.String("tool", toolName),
		zap.String("url_or_path", urlOrPath),
	)

	path, err := wechat.DownloadFile(urlOrPath)
	duration := time.Since(startTime)

	if err != nil {
		s.log.Error("Tool execution failed",
			zap.String("tool", toolName),
			zap.Error(err),
			zap.Duration("duration", duration),
		)
		return mcp.NewToolResultError(fmt.Sprintf("download file failed: %v", err)), nil
	}

	s.log.Info("Tool executed successfully",
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

	s.log.Info("Tool called", zap.String("tool", toolName))

	args := request.GetArguments()

	markdown, ok := args["markdown"].(string)
	if !ok || markdown == "" {
		s.log.Warn("Tool argument missing", zap.String("tool", toolName), zap.String("argument", "markdown"))
		return mcp.NewToolResultError("markdown is required"), nil
	}

	// 获取可选参数
	theme, _ := args["theme"].(string)
	customPrompt, _ := args["custom_prompt"].(string)

	s.log.Debug("Tool arguments",
		zap.String("tool", toolName),
		zap.String("theme", theme),
		zap.Bool("has_custom_prompt", customPrompt != ""),
		zap.Int("markdown_length", len(markdown)),
	)

	// 执行转换
	req := &converter.ConvertRequest{
		Markdown:     markdown,
		Theme:        theme,
		CustomPrompt: customPrompt,
	}

	result := s.converter.Convert(req)
	duration := time.Since(startTime)

	if !result.Success {
		s.log.Error("Tool execution failed",
			zap.String("tool", toolName),
			zap.Error(fmt.Errorf(result.Error)),
			zap.Duration("duration", duration),
		)
		return mcp.NewToolResultError(fmt.Sprintf("convert failed: %s", result.Error)), nil
	}

	s.log.Info("Tool executed successfully",
		zap.String("tool", toolName),
		zap.String("theme", result.Theme),
		zap.Int("image_count", len(result.Images)),
		zap.Int("html_length", len(result.HTML)),
		zap.Duration("duration", duration),
	)

	// 返回HTML内容
	return mcp.NewToolResultText(result.HTML), nil
}

// listThemesHandler 列出主题处理器
func (s *Server) listThemesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	toolName := "list_themes"

	s.log.Info("Tool called", zap.String("tool", toolName))

	themeMgr := s.converter.GetThemeManager()
	themes := themeMgr.ListThemes()

	s.log.Info("Tool executed successfully",
		zap.String("tool", toolName),
		zap.Int("theme_count", len(themes)),
	)

	// 构建主题列表
	var result strings.Builder
	result.WriteString("可用主题:\n")
	for _, theme := range themes {
		result.WriteString(fmt.Sprintf("- %s\n", theme))
	}

	return mcp.NewToolResultText(result.String()), nil
}

// maskMediaID 遮蔽 media_id 用于日志
func maskMediaID(id string) string {
	if id == "" || len(id) < 8 {
		return "***"
	}
	return id[:4] + "***" + id[len(id)-4:]
}

package mcp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/qiangmzsx/wechat-mcp/config"
	"github.com/qiangmzsx/wechat-mcp/converter"
	"github.com/qiangmzsx/wechat-mcp/logger"
	"github.com/qiangmzsx/wechat-mcp/wechat"
	"go.uber.org/zap"
)

// Server MCP服务器
type Server struct {
	svc         *wechat.Service
	converter   converter.Converter
	aiConverter converter.Converter
	config      *config.Config
}

// New 创建MCP服务器
func New(cfg *config.Config) *Server {
	svc := wechat.NewService(cfg)

	var conv converter.Converter
	var aiConv converter.Converter
	var err error

	// 创建默认转换器
	conv, err = converter.NewConverter(cfg)
	if err != nil {
		logger.Warn("converter initialization failed, using simple converter", zap.Error(err))
		conv = converter.NewSimpleConverter()
	}

	// 创建 AI 转换器（如果启用）
	if cfg.Converter.Enabled {
		aiConv, err = converter.NewAIConverter(cfg)
		if err != nil {
			logger.Warn("AI converter initialization failed", zap.Error(err))
			aiConv = conv
		}
	} else {
		aiConv = conv
	}

	return &Server{
		svc:         svc,
		converter:   conv,
		aiConverter: aiConv,
		config:      cfg,
	}
}

// Run 启动服务器
func (s *Server) Run() error {
	logger.Info("Initializing MCP Server",
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
	logger.Info("MCP tools registered successfully")

	// 根据配置选择协议
	switch s.config.MCP.Protocol {
	case "http":
		return s.runHTTP(mcpServer)
	case "sse":
		return s.runSSE(mcpServer)
	default:
		logger.Info("Starting STDIO server (waiting for client connection...)")
		return server.ServeStdio(mcpServer)
	}
}

// runHTTP 启动HTTP服务器
func (s *Server) runHTTP(mcpServer *server.MCPServer) error {
	addr := fmt.Sprintf("%s:%d", s.config.MCP.Host, s.config.MCP.Port)
	logger.Info("Starting StreamableHTTP server",
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
		logger.Error("HTTP server failed to start", zap.Error(err))
		return err
	}
	return nil
}

// runSSE 启动SSE服务器
func (s *Server) runSSE(mcpServer *server.MCPServer) error {
	addr := fmt.Sprintf("%s:%d", s.config.MCP.Host, s.config.MCP.Port)
	logger.Info("Starting SSE server",
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
		logger.Error("SSE server failed to start", zap.Error(err))
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
		logger.Warn("Failed to encode health response", zap.Error(err))
	}
}

// registerTools 注册所有工具
func (s *Server) registerTools(mcpServer *server.MCPServer) {
	logger.Debug("Registering MCP tools...")

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
	logger.Debug("Tool registered: upload_material")

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
	logger.Debug("Tool registered: create_draft")

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
	logger.Debug("Tool registered: create_newspic_draft")

	// 4. 获取AccessToken工具
	mcpServer.AddTool(
		mcp.NewTool("get_access_token",
			mcp.WithDescription("获取微信公众号AccessToken（用于调试）"),
		),
		s.getAccessTokenHandler,
	)
	logger.Debug("Tool registered: get_access_token")

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
	logger.Debug("Tool registered: download_file")

	// 6. Markdown转HTML工具
	mcpServer.AddTool(
		mcp.NewTool("convert_markdown",
			mcp.WithDescription("将Markdown内容转换为微信公众号兼容的HTML"),
			mcp.WithString("markdown",
				mcp.Required(),
				mcp.Description("Markdown内容"),
			),
			mcp.WithString("theme",
				mcp.Description("主题名称: default, elegant, tech, minimalist 等"),
			),
			mcp.WithString("custom_prompt",
				mcp.Description("自定义提示词(可选)"),
			),
			mcp.WithString("converter_type",
				mcp.Description("转换器类型: api (默认，基于goldmark), ai (基于LLM)"),
			),
		),
		s.convertMarkdownHandler,
	)
	logger.Debug("Tool registered: convert_markdown")

	// 7. 列出可用主题工具
	mcpServer.AddTool(
		mcp.NewTool("list_themes",
			mcp.WithDescription("列出所有可用的主题"),
		),
		s.listThemesHandler,
	)
	logger.Debug("Tool registered: list_themes")
}

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/qiangmzsx/wechat-mcp/config"
	"github.com/qiangmzsx/wechat-mcp/internal/util"
	"github.com/qiangmzsx/wechat-mcp/logger"
	"github.com/qiangmzsx/wechat-mcp/mcp"
	"go.uber.org/zap"
)

func main() {
	// 解析命令行参数
	configPath := flag.String("c", "config.toml", "配置文件路径")
	flag.Parse()

	// 加载配置
	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 创建日志
	zapLogger, err := cfg.NewLogger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "创建日志失败: %v\n", err)
		os.Exit(1)
	}
	defer zapLogger.Sync()

	// 初始化全局 logger
	logger.InitWithLogger(zapLogger)

	logger.Info("Starting WeChat MCP Server",
		zap.String("config_path", *configPath),
	)

	logger.Info("Config loaded successfully",
		zap.String("wechat_app_id", util.MaskAppID(cfg.WechatAppID)),
		zap.String("log_level", cfg.Log.Level),
		zap.String("log_format", cfg.Log.Format),
		zap.String("protocol", cfg.MCP.Protocol),
		zap.String("host", cfg.MCP.Host),
		zap.Int("port", cfg.MCP.Port),
	)

	logger.Info("Converter config",
		zap.Bool("enabled", cfg.Converter.Enabled),
		zap.String("model", cfg.Converter.Model),
		zap.Int("max_tokens", cfg.Converter.MaxTokens),
		zap.String("default_theme", cfg.Converter.DefaultTheme),
		zap.String("theme_dir", cfg.Converter.ThemeDir),
		zap.Duration("timeout", cfg.Converter.Timeout),
		zap.String("base_url", cfg.Converter.BaseURL),
	)

	// 创建并运行MCP服务器
	server := mcp.New(cfg)
	if err := server.Run(); err != nil {
		logger.Error("Server stopped with error", zap.Error(err))
		fmt.Fprintf(os.Stderr, "服务器运行失败: %v\n", err)
		os.Exit(1)
	}

	logger.Info("Server stopped gracefully")
}

package config

import (
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/qiangmzsx/wechat-mcp/provider"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// 环境变量名
	envWechatAppID     = "WECHAT_APP_ID"
	envWechatAppSecret = "WECHAT_APP_SECRET"

	// AI Provider 环境变量
	envAIAPIKey  = "AI_API_KEY"
	envAIBaseURL = "AI_BASE_URL"
)

// Config 应用配置
type Config struct {
	WechatAppID     string `toml:"wechat_app_id"`
	WechatAppSecret string `toml:"wechat_app_secret"`
	Log             LogConfig
	MCP             MCPConfig
	Converter       ConverterConfig
}

// LogConfig 日志配置
type LogConfig struct {
	Level  string `toml:"level"`
	Format string `toml:"format"`
}

// MCPConfig MCP服务配置
type MCPConfig struct {
	Protocol string `toml:"protocol"`
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
}

// ConverterType 转换器类型
type ConverterType string

const (
	// ConverterTypeAPI 使用 API 转换器（基于 goldmark，默认）
	ConverterTypeAPI ConverterType = "api"
	// ConverterTypeAI 使用 AI 转换器（基于 LLM）
	ConverterTypeAI ConverterType = "ai"
)

// ConverterConfig AI 转换器配置
type ConverterConfig struct {
	Type         ConverterType `toml:"type"`     // 转换器类型: api, ai
	Enabled      bool          `toml:"enabled"`  // 是否启用
	Provider     string        `toml:"provider"` // 供应商: anthropic, openai
	APIKey       string        `toml:"api_key"`
	BaseURL      string        `toml:"base_url"`
	Model        string        `toml:"model"`
	MaxTokens    int           `toml:"max_tokens"`
	DefaultTheme string        `toml:"default_theme"`
	ThemeDir     string        `toml:"theme_dir"`
	Timeout      time.Duration `toml:"timeout"`
}

// Load 加载配置文件
func Load(path string) (*Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", path)
	}

	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	// 设置默认值
	if cfg.Log.Level == "" {
		cfg.Log.Level = "info"
	}
	if cfg.Log.Format == "" {
		cfg.Log.Format = "json"
	}

	// Converter 默认配置
	// 设置默认转换器类型
	if cfg.Converter.Type == "" {
		cfg.Converter.Type = ConverterTypeAPI
	}

	// 如果使用 AI 转换器，自动启用
	if cfg.Converter.Type == ConverterTypeAI {
		cfg.Converter.Enabled = true
	}

	// 设置默认 Provider
	if cfg.Converter.Provider == "" {
		cfg.Converter.Provider = string(provider.ProviderAnthropic)
	}

	// 根据 Provider 类型设置默认模型
	if cfg.Converter.Model == "" {
		switch provider.ProviderType(cfg.Converter.Provider) {
		case provider.ProviderAnthropic:
			cfg.Converter.Model = "claude-sonnet-4-5-20250929"
		case provider.ProviderOpenAI:
			cfg.Converter.Model = "gpt-4.1"
		default:
			cfg.Converter.Model = "claude-sonnet-4-5-20250929"
		}
	}

	if cfg.Converter.MaxTokens == 0 {
		cfg.Converter.MaxTokens = 4096
	}
	if cfg.Converter.DefaultTheme == "" {
		cfg.Converter.DefaultTheme = "default"
	}
	if cfg.Converter.Timeout == 0 {
		cfg.Converter.Timeout = 60 * time.Second
	}

	// 优先从环境变量读取，环境变量优先级高于配置文件
	if appID := os.Getenv(envWechatAppID); appID != "" {
		cfg.WechatAppID = appID
	}
	if appSecret := os.Getenv(envWechatAppSecret); appSecret != "" {
		cfg.WechatAppSecret = appSecret
	}

	// AI Provider 环境变量（统一命名，支持所有供应商）
	if apiKey := os.Getenv(envAIAPIKey); apiKey != "" {
		cfg.Converter.APIKey = apiKey
	}
	if baseURL := os.Getenv(envAIBaseURL); baseURL != "" {
		cfg.Converter.BaseURL = baseURL
	}

	if cfg.WechatAppID == "" || cfg.WechatAppSecret == "" {
		return nil, fmt.Errorf("wechat_app_id and wechat_app_secret are required (set via config or env: %s, %s)", envWechatAppID, envWechatAppSecret)
	}

	return &cfg, nil
}

// NewLogger 创建日志实例
func (c *Config) NewLogger() (*zap.Logger, error) {
	level := zapcore.InfoLevel
	if err := level.UnmarshalText([]byte(c.Log.Level)); err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var encoder zapcore.Encoder
	if c.Log.Format == "console" {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stderr),
		level,
	)

	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)), nil
}

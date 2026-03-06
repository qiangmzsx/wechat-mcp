package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// 环境变量名
	envWechatAppID     = "WECHAT_APP_ID"
	envWechatAppSecret = "WECHAT_APP_SECRET"
)

// Config 应用配置
type Config struct {
	WechatAppID     string `toml:"wechat_app_id"`
	WechatAppSecret string `toml:"wechat_app_secret"`
	Log             LogConfig
	MCP             MCPConfig
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

	// 优先从环境变量读取，环境变量优先级高于配置文件
	if appID := os.Getenv(envWechatAppID); appID != "" {
		cfg.WechatAppID = appID
	}
	if appSecret := os.Getenv(envWechatAppSecret); appSecret != "" {
		cfg.WechatAppSecret = appSecret
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

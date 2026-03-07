// Package factory 提供 AI Provider 工厂
package factory

import (
	"fmt"

	"github.com/qiangmzsx/wechat-mcp/config"
	"github.com/qiangmzsx/wechat-mcp/provider"
	"github.com/qiangmzsx/wechat-mcp/provider/anthropic"
	"github.com/qiangmzsx/wechat-mcp/provider/openai"
)

// NewProvider 根据配置创建 AI Provider
func NewProvider(cfg *config.Config) (provider.Provider, error) {
	switch provider.ProviderType(cfg.Converter.Provider) {
	case provider.ProviderAnthropic:
		return anthropic.NewClient(
			cfg.Converter.APIKey,
			cfg.Converter.BaseURL,
			cfg.Converter.Model,
			cfg.Converter.MaxTokens,
		), nil

	case provider.ProviderOpenAI:
		return openai.NewClient(
			cfg.Converter.APIKey,
			cfg.Converter.BaseURL,
			cfg.Converter.Model,
			cfg.Converter.MaxTokens,
		), nil

	case "":
		// 默认使用 Anthropic
		return anthropic.NewClient(
			cfg.Converter.APIKey,
			cfg.Converter.BaseURL,
			cfg.Converter.Model,
			cfg.Converter.MaxTokens,
		), nil

	default:
		return nil, fmt.Errorf("unsupported provider: %s", cfg.Converter.Provider)
	}
}

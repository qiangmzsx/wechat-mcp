// Package provider 提供 AI 模型抽象接口
// 支持多种 AI 供应商：Anthropic、OpenAI 等
package provider

import (
	"context"
	"io"
)

// ProviderType AI 供应商类型
type ProviderType string

const (
	ProviderAnthropic ProviderType = "anthropic" // Anthropic (Claude)
	ProviderOpenAI    ProviderType = "openai"    // OpenAI 及兼容 API
)

// Options keys
const (
	OptMaxTokens   = "max_tokens"
	OptTemperature = "temperature"
)

// Provider 是所有 LLM 供应商需要实现的接口
type Provider interface {
	// Chat 发送消息并返回响应
	Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error)

	// ChatStream 流式返回响应
	ChatStream(ctx context.Context, req ChatRequest, onChunk func(StreamChunk)) (*ChatResponse, error)

	// Name 返回供应商名称
	Name() string

	// DefaultModel 返回默认模型名称
	DefaultModel() string
}

// ChatRequest 聊天请求
type ChatRequest struct {
	Messages []Message              `json:"messages"`
	Model    string                 `json:"model,omitempty"`
	Options  map[string]interface{} `json:"options,omitempty"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	Content      string `json:"content"`
	FinishReason string `json:"finish_reason"`
	Usage        *Usage `json:"usage,omitempty"`
}

// StreamChunk 流式响应片段
type StreamChunk struct {
	Content string `json:"content,omitempty"`
	Done    bool   `json:"done,omitempty"`
}

// Message 对话消息
type Message struct {
	Role    string `json:"role"` // system, user, assistant
	Content string `json:"content"`
}

// Usage 使用量统计
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// StreamReader 流式响应读取器
type StreamReader struct {
	reader io.ReadCloser
}

// NewStreamReader 创建流式读取器
func NewStreamReader(reader io.ReadCloser) *StreamReader {
	return &StreamReader{reader: reader}
}

// Close 关闭读取器
func (s *StreamReader) Close() error {
	return s.reader.Close()
}

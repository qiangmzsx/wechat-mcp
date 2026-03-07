// Package anthropic 提供 Anthropic (Claude) AI 供应商实现
package anthropic

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/qiangmzsx/wechat-mcp/provider"
)

// Client Anthropic API 客户端
type Client struct {
	apiKey    string
	baseURL   string
	model     string
	maxTokens int64
	client    *http.Client
}

// NewClient 创建 Anthropic 客户端
func NewClient(apiKey, baseURL, model string, maxTokens int) *Client {
	if baseURL == "" {
		baseURL = "https://api.anthropic.com"
	}
	baseURL = strings.TrimRight(baseURL, "/")

	return &Client{
		apiKey:    apiKey,
		baseURL:   baseURL,
		model:     model,
		maxTokens: int64(maxTokens),
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// Name 返回供应商名称
func (c *Client) Name() string {
	return "anthropic"
}

// DefaultModel 返回默认模型
func (c *Client) DefaultModel() string {
	return c.model
}

// Chat 发送聊天请求
func (c *Client) Chat(ctx context.Context, req provider.ChatRequest) (*provider.ChatResponse, error) {
	// 构建请求体
	body := map[string]interface{}{
		"model":      c.resolveModel(req.Model),
		"max_tokens": c.maxTokens,
		"messages":   convertMessages(req.Messages),
	}

	// 添加 system 消息
	for _, msg := range req.Messages {
		if msg.Role == "system" {
			body["system"] = []map[string]string{{"type": "text", "text": msg.Content}}
			break
		}
	}

	// 添加选项
	if v, ok := req.Options[provider.OptMaxTokens]; ok {
		if maxTokens, ok := v.(int); ok {
			body["max_tokens"] = maxTokens
		}
	}

	// 发送请求
	respBody, err := c.doRequest(ctx, "/v1/messages", body)
	if err != nil {
		return nil, err
	}
	defer respBody.Close()

	// 解析响应
	var resp messageResponse
	if err := json.NewDecoder(respBody).Decode(&resp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	// 提取内容
	var content string
	for _, block := range resp.Content {
		if block.Type == "text" {
			content = block.Text
			break
		}
	}

	return &provider.ChatResponse{
		Content:      content,
		FinishReason: resp.StopReason,
		Usage: &provider.Usage{
			PromptTokens:     resp.Usage.InputTokens,
			CompletionTokens: resp.Usage.OutputTokens,
			TotalTokens:      resp.Usage.InputTokens + resp.Usage.OutputTokens,
		},
	}, nil
}

// ChatStream 流式聊天
func (c *Client) ChatStream(ctx context.Context, req provider.ChatRequest, onChunk func(provider.StreamChunk)) (*provider.ChatResponse, error) {
	// 构建请求体
	body := map[string]interface{}{
		"model":      c.resolveModel(req.Model),
		"max_tokens": c.maxTokens,
		"messages":   convertMessages(req.Messages),
		"stream":     true,
	}

	// 添加 system 消息
	for _, msg := range req.Messages {
		if msg.Role == "system" {
			body["system"] = []map[string]string{{"type": "text", "text": msg.Content}}
			break
		}
	}

	// 发送请求
	respBody, err := c.doRequest(ctx, "/v1/messages", body)
	if err != nil {
		return nil, err
	}
	defer respBody.Close()

	// 读取流式响应
	scanner := bufio.NewScanner(respBody)
	var content string

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data:") {
			continue
		}

		data := strings.TrimPrefix(line, "data:")
		data = strings.TrimPrefix(data, " ")

		if data == "" || data == "[DONE]" {
			break
		}

		var chunk streamChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}

		if len(chunk.Content) > 0 && chunk.Content[0].Type == "text" {
			content += chunk.Content[0].Text
			if onChunk != nil {
				onChunk(provider.StreamChunk{Content: chunk.Content[0].Text})
			}
		}

		if chunk.Type == "message_delta" && chunk.Delta.StopReason != "" {
			if onChunk != nil {
				onChunk(provider.StreamChunk{Done: true})
			}
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("stream read error: %w", err)
	}

	return &provider.ChatResponse{
		Content:      content,
		FinishReason: "stop",
	}, nil
}

// resolveModel 解析模型名称
func (c *Client) resolveModel(model string) string {
	if model == "" {
		return c.model
	}
	return model
}

// doRequest 发送请求
func (c *Client) doRequest(ctx context.Context, path string, body interface{}) (io.ReadCloser, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+path, strings.NewReader(string(data)))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", c.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("API error: %s", string(respBody))
	}

	return resp.Body, nil
}

// convertMessages 转换消息格式
func convertMessages(messages []provider.Message) []map[string]string {
	result := make([]map[string]string, 0, len(messages))
	for _, msg := range messages {
		if msg.Role == "system" {
			continue // system 消息单独处理
		}
		result = append(result, map[string]string{
			"role":    msg.Role,
			"content": msg.Content,
		})
	}
	return result
}

// API 响应结构
type messageResponse struct {
	Type       string         `json:"type"`
	Content    []contentBlock `json:"content"`
	StopReason string         `json:"stop_reason"`
	Usage      usage          `json:"usage"`
}

type contentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

type streamChunk struct {
	Type    string         `json:"type"`
	Content []contentBlock `json:"content"`
	Delta   deltaBlock     `json:"delta"`
}

type deltaBlock struct {
	StopReason string `json:"stop_reason"`
}

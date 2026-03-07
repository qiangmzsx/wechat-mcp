// Package openai 提供 OpenAI 及兼容 API 供应商实现
package openai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/qiangmzsx/wechat-mcp/provider"
)

// Client OpenAI API 客户端
type Client struct {
	apiKey    string
	baseURL   string
	model     string
	maxTokens int
	client    *http.Client
}

// NewClient 创建 OpenAI 客户端
func NewClient(apiKey, baseURL, model string, maxTokens int) *Client {
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	baseURL = strings.TrimRight(baseURL, "/")

	return &Client{
		apiKey:    apiKey,
		baseURL:   baseURL,
		model:     model,
		maxTokens: maxTokens,
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// Name 返回供应商名称
func (c *Client) Name() string {
	return "openai"
}

// DefaultModel 返回默认模型
func (c *Client) DefaultModel() string {
	return c.model
}

// Chat 发送聊天请求
func (c *Client) Chat(ctx context.Context, req provider.ChatRequest) (*provider.ChatResponse, error) {
	// 构建请求体
	body := map[string]interface{}{
		"model":    c.resolveModel(req.Model),
		"messages": req.Messages,
		"stream":   false,
	}

	// 添加选项
	if v, ok := req.Options[provider.OptMaxTokens]; ok {
		if maxTokens, ok := v.(int); ok {
			body["max_tokens"] = maxTokens
		}
	} else {
		body["max_tokens"] = c.maxTokens
	}

	if v, ok := req.Options[provider.OptTemperature]; ok {
		body["temperature"] = v
	}

	// 发送请求
	respBody, err := c.doRequest(ctx, "/chat/completions", body)
	if err != nil {
		return nil, err
	}
	defer respBody.Close()

	// 解析响应
	var resp chatCompletionResponse
	if err := json.NewDecoder(respBody).Decode(&resp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response choices")
	}

	return &provider.ChatResponse{
		Content:      resp.Choices[0].Message.Content,
		FinishReason: resp.Choices[0].FinishReason,
		Usage: &provider.Usage{
			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		},
	}, nil
}

// ChatStream 流式聊天
func (c *Client) ChatStream(ctx context.Context, req provider.ChatRequest, onChunk func(provider.StreamChunk)) (*provider.ChatResponse, error) {
	// 构建请求体
	body := map[string]interface{}{
		"model":    c.resolveModel(req.Model),
		"messages": req.Messages,
		"stream":   true,
	}

	if v, ok := req.Options[provider.OptMaxTokens]; ok {
		if maxTokens, ok := v.(int); ok {
			body["max_tokens"] = maxTokens
		}
	} else {
		body["max_tokens"] = c.maxTokens
	}

	// 发送请求
	respBody, err := c.doRequest(ctx, "/chat/completions", body)
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

		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			content += chunk.Choices[0].Delta.Content
			if onChunk != nil {
				onChunk(provider.StreamChunk{Content: chunk.Choices[0].Delta.Content})
			}
		}

		if chunk.Choices[0].FinishReason != "" {
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

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+path, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

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

// API 响应结构
type chatCompletionResponse struct {
	Choices []choice `json:"choices"`
	Usage   usage    `json:"usage"`
}

type choice struct {
	Message      message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type streamChunk struct {
	Choices []streamChoice `json:"choices"`
}

type streamChoice struct {
	Delta        delta  `json:"delta"`
	FinishReason string `json:"finish_reason"`
}

type delta struct {
	Content string `json:"content"`
}

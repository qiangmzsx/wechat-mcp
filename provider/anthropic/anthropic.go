// Package anthropic provides Anthropic (Claude) AI provider implementation
package anthropic

import (
	"context"
	"fmt"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/qiangmzsx/wechat-mcp/logger"
	"github.com/qiangmzsx/wechat-mcp/provider"
	"go.uber.org/zap"
)

// Client Anthropic API client
type Client struct {
	client    anthropic.Client
	apiKey    string
	baseURL   string
	model     string
	maxTokens int64
}

// NewClient creates Anthropic client
func NewClient(apiKey, baseURL, model string, maxTokens int) *Client {
	if baseURL == "" {
		baseURL = "https://api.anthropic.com"
	}
	baseURL = strings.TrimRight(baseURL, "/")

	opts := []option.RequestOption{
		option.WithAPIKey(apiKey),
		option.WithBaseURL(baseURL),
	}

	return &Client{
		client:    anthropic.NewClient(opts...),
		apiKey:    apiKey,
		baseURL:   baseURL,
		model:     model,
		maxTokens: int64(maxTokens),
	}
}

// Name returns provider name
func (c *Client) Name() string {
	return "anthropic"
}

// DefaultModel returns default model
func (c *Client) DefaultModel() string {
	return c.model
}

// Chat sends chat request
func (c *Client) Chat(ctx context.Context, req provider.ChatRequest) (*provider.ChatResponse, error) {
	logger.Debug("Anthropic API: sending chat request", zap.String("model", string(c.resolveModel(req.Model))))

	// Build request params
	params := anthropic.MessageNewParams{
		Model:     c.resolveModel(req.Model),
		MaxTokens: c.maxTokens,
	}

	// Add system message
	for _, msg := range req.Messages {
		if msg.Role == "system" {
			params.System = append(params.System, anthropic.TextBlockParam{
				Text: msg.Content,
			})
			break
		}
	}

	// Add conversation messages
	params.Messages = convertMessages(req.Messages)

	// Add options
	if v, ok := req.Options[provider.OptMaxTokens]; ok {
		if maxTokens, ok := v.(int); ok {
			params.MaxTokens = int64(maxTokens)
		}
	}

	// Send request
	message, err := c.client.Messages.New(ctx, params)
	if err != nil {
		logger.Error("Anthropic API error", zap.Error(err))
		return nil, fmt.Errorf("anthropic API error: %w", err)
	}

	// Extract content
	var content string
	for _, block := range message.Content {
		if block.Type == "text" {
			content += block.Text
		}
	}

	logger.Debug("Anthropic API response received", zap.Int("tokens", int(message.Usage.InputTokens+message.Usage.OutputTokens)))
	return &provider.ChatResponse{
		Content:      content,
		FinishReason: string(message.StopReason),
		Usage: &provider.Usage{
			PromptTokens:     int(message.Usage.InputTokens),
			CompletionTokens: int(message.Usage.OutputTokens),
			TotalTokens:      int(message.Usage.InputTokens + message.Usage.OutputTokens),
		},
	}, nil
}

// ChatStream streaming chat
func (c *Client) ChatStream(ctx context.Context, req provider.ChatRequest, onChunk func(provider.StreamChunk)) (*provider.ChatResponse, error) {
	logger.Debug("Anthropic API: starting streaming chat", zap.String("model", string(c.resolveModel(req.Model))))

	// Build request params
	params := anthropic.MessageNewParams{
		Model:     c.resolveModel(req.Model),
		MaxTokens: c.maxTokens,
	}

	// Add system message
	for _, msg := range req.Messages {
		if msg.Role == "system" {
			params.System = append(params.System, anthropic.TextBlockParam{
				Text: msg.Content,
			})
			break
		}
	}

	// Add conversation messages
	params.Messages = convertMessages(req.Messages)

	// Send streaming request
	stream := c.client.Messages.NewStreaming(ctx, params)

	message := anthropic.Message{}
	var content string

	for stream.Next() {
		event := stream.Current()
		err := message.Accumulate(event)
		if err != nil {
			logger.Error("Anthropic stream accumulate error", zap.Error(err))
			return nil, fmt.Errorf("accumulate event: %w", err)
		}

		// Handle content block deltas
		switch eventVariant := event.AsAny().(type) {
		case anthropic.ContentBlockDeltaEvent:
			switch deltaVariant := eventVariant.Delta.AsAny().(type) {
			case anthropic.TextDelta:
				content += deltaVariant.Text
				if onChunk != nil {
					onChunk(provider.StreamChunk{Content: deltaVariant.Text})
				}
			}
		case anthropic.MessageDeltaEvent:
			if onChunk != nil {
				onChunk(provider.StreamChunk{Done: true})
			}
		}
	}

	if stream.Err() != nil {
		logger.Error("Anthropic stream error", zap.Error(stream.Err()))
		return nil, fmt.Errorf("stream error: %w", stream.Err())
	}

	logger.Debug("Anthropic streaming completed")
	return &provider.ChatResponse{
		Content:      content,
		FinishReason: string(message.StopReason),
	}, nil
}

// resolveModel resolves model name
func (c *Client) resolveModel(model string) anthropic.Model {
	if model == "" {
		return anthropic.Model(c.model)
	}
	return anthropic.Model(model)
}

// convertMessages converts message format
func convertMessages(messages []provider.Message) []anthropic.MessageParam {
	result := make([]anthropic.MessageParam, 0, len(messages))
	for _, msg := range messages {
		if msg.Role == "system" {
			continue // system message handled separately
		}
		role := anthropic.MessageParamRole(msg.Role)
		result = append(result, anthropic.MessageParam{
			Content: []anthropic.ContentBlockParamUnion{{
				OfText: &anthropic.TextBlockParam{
					Text: msg.Content,
				},
			}},
			Role: role,
		})
	}
	return result
}

# AGENTS.md - WeChat MCP Server Developer Guide

## Project Overview

This is a Go-based MCP (Model Context Protocol) server for interacting with WeChat Official Accounts. It provides tools for uploading materials, creating drafts, converting Markdown to HTML, and more.

- **Language**: Go 1.23+
- **Framework**: mcp-go (MCP SDK), Gin (implied from skill)
- **Dependencies**: See `go.mod`

---

## Build, Lint & Test Commands

### Build
```bash
# Build binary
go build -o wechat-mcp .

# Build with version info
go build -ldflags="-X main.version=1.0.0" -o wechat-mcp .
```

### Run
```bash
# Run with default config.toml
./wechat-mcp

# Run with custom config
./wechat-mcp -c /path/to/config.toml

# With environment variables
WECHAT_APP_ID=xxx WECHAT_APP_SECRET=xxx ./wechat-mcp
```

### Test
```bash
# Run all tests
go test ./...

# Run all tests with verbose output
go test -v ./...

# Run specific test
go test -v -run TestExtractImages ./converter/

# Run tests in specific package
go test -v ./converter/

# Run tests with coverage
go test -cover ./...

# Skip integration tests (that require API keys)
go test -v -short ./...
```

### Development
```bash
# Format code
go fmt ./...

# Run go vet
go vet ./...

# Tidy dependencies
go mod tidy
```

---

## Code Style Guidelines

### 1. Imports

Group imports with blank line between groups:
```go
import (
    "context"
    "fmt"
    "time"

    "github.com/qiangmzsx/wechat-mcp/config"
    "github.com/qiangmzsx/wechat-mcp/converter"
    "go.uber.org/zap"
)
```

**Order**: Standard library → External dependencies → Project packages.

### 2. Naming Conventions

- **Variables/functions**: `camelCase`
- **Exported types/functions**: `PascalCase`
- **Constants**: `PascalCase` or `camelCase` with grouping
- **Package names**: Short, lowercase, no underscores

### 3. Comments

- **Chinese comments** for exported functions (matches project style):
```go
// Load 加载配置文件
func Load(path string) (*Config, error) {}

// NewLogger 创建日志实例
func (c *Config) NewLogger() (*zap.Logger, error) {}
```

- **English comments** for internal/unexported functions are acceptable
- Use period at end of comments

### 4. Error Handling

- Return errors as last return value
- Use `fmt.Errorf` with `%w` for error wrapping:
```go
return nil, fmt.Errorf("parse config: %w", err)
```
- Handle errors explicitly, never ignore with `_` unless intentionally
- Log errors with zap before returning:
```go
if err != nil {
    logger.Error("operation failed", zap.Error(err))
    return nil, fmt.Errorf("operation: %w", err)
}
```

### 5. Logging

- Use `go.uber.org/zap` for structured logging
- Use key-value pairs:
```go
logger.Info("message",
    zap.String("key", value),
    zap.Int("count", n),
    zap.Error(err),
)
```
- Use appropriate log levels: Debug → Info → Warn → Error
- Add `defer logger.Sync()` for sync on exit

### 6. Types & Structs

- Use struct tags for configuration (TOML):
```go
type Config struct {
    WechatAppID     string `toml:"wechat_app_id"`
    WechatAppSecret string `toml:"wechat_app_secret"`
}
```
- Use meaningful type names (e.g., `ConverterConfig`, not `ConfigType`)

### 7. Context Usage

- Pass `context.Context` as first parameter for functions that may need timeout/cancellation:
```go
func (s *Server) Handler(ctx context.Context, request Request) error {}
```

### 8. Testing Patterns

- Use table-driven tests with `t.Run`:
```go
tests := []struct {
    name     string
    input    string
    expected string
}{
    {"case1", "input1", "expected1"},
    {"case2", "input2", "expected2"},
}
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        result := fn(tt.input)
        if result != tt.expected {
            t.Errorf("got %s, want %s", result, tt.expected)
        }
    })
}
```
- Skip integration tests when API keys not set:
```go
if apiKey := os.Getenv("ANTHROPIC_API_KEY"); apiKey == "" {
    t.Skip("Skipping: ANTHROPIC_API_KEY not set")
}
```
- Use `zap.NewNop()` for tests that don't need logging

### 9. Configuration

- Use `config.toml` for configuration
- Support environment variables with higher priority than config file
- Provide sensible defaults for all optional fields
- Validate required fields early

### 10. MCP Tool Registration

- Use `mcp-go` server pattern
- Register tools in `registerTools()` method
- Use structured handler functions with context and request parameters
- Return `mcp.NewToolResultError()` for errors, `mcp.NewToolResultText()` for success

---

## Project Structure

```
wechat-mcp/
├── main.go              # Entry point
├── config/              # Configuration loading
├── converter/           # Markdown to HTML conversion
│   ├── converter.go     # Main converter logic
│   ├── converter_test.go
│   ├── types.go         # Request/Response types
│   ├── theme.go         # Theme management
│   └── prompt.go        # Prompt building
├── mcp/                 # MCP server
│   └── server.go        # Server and tool handlers
├── provider/            # AI providers
│   ├── provider.go      # Provider interface
│   ├── anthropic/       # Anthropic implementation
│   ├── openai/         # OpenAI implementation
│   └── factory/        # Provider factory
├── theme/              # Theme definitions
├── wechat/             # WeChat API integration
├── config.example.toml  # Example configuration
└── themes/             # Theme files
```

---

## Key Patterns

### Factory Pattern (Provider Factory)
```go
client, err := factory.NewProvider(cfg)
```

### Strategy Pattern (Converter)
```go
switch convType {
case config.ConverterTypeAI:
    return NewAIConverter(cfg, log)
case config.ConverterTypeAPI:
    return NewAPIConverter(), nil
}
```

### Handler Pattern (MCP Tools)
```go
func (s *Server) toolNameHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    // Validate args
    // Execute logic
    // Return result
}
```

---

## Common Tasks

### Adding a New MCP Tool
1. Define tool in `registerTools()` with `mcp.NewTool()`
2. Create handler function in `mcp/server.go`
3. Register handler with tool

### Adding a New Theme
1. Add theme definition in `theme/themes.go` or `converter/theme.go`
2. Add AI prompt in theme's AI prompt section

### Adding a New AI Provider
1. Implement `provider.Provider` interface in new package under `provider/`
2. Add provider type constant in `provider/provider.go`
3. Update factory to create new provider

---

## Configuration Reference

See `config.example.toml` for all available options:

| Section | Options |
|---------|---------|
| wechat | wechat_app_id, wechat_app_secret |
| log | level (debug/info/warn/error), format (json/console) |
| mcp | protocol (stdio/http/sse), host, port |
| converter | enabled, provider (anthropic/openai), api_key, base_url, model, max_tokens, default_theme, theme_dir, timeout |

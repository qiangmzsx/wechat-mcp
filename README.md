# WeChat MCP Server

微信公众号 MCP (Model Context Protocol) 服务器，提供与微信公众号平台交互的能力。

## 功能特性

- 支持 MCP 协议：stdio、HTTP、SSE
- 上传图片素材到微信公众号
- 创建和管理文章草稿
- 创建小绿书（图片）草稿
- 获取 AccessToken
- 文件下载
- **AI Markdown 转 HTML**（支持 Anthropic 和 OpenAI）

## 环境要求

- Go 1.23+
- 微信公众号 AppID 和 AppSecret

## 快速开始

### 1. 克隆项目

```bash
git clone https://github.com/qiangmzsx/wechat-mcp.git
cd wechat-mcp
```

### 2. 配置

复制配置文件：

```bash
cp config.example.toml config.toml
```

编辑 `config.toml`：

```toml
# 微信公众号配置（必填）
wechat_app_id = "your_app_id"
wechat_app_secret = "your_app_secret"

# 日志配置
[log]
level = "debug"      # debug, info, warn, error
format = "json"      # json, console

# MCP服务配置
[mcp]
protocol = "http"    # stdio, http, sse
host = "0.0.0.0"
port = 7990

# AI 转换器配置（可选）
[converter]
enabled = false                      # 是否启用 AI 转换功能
provider = "anthropic"              # 供应商: anthropic, openai
api_key = ""                        # API Key
base_url = ""                       # 自定义 API 地址
# model 默认值: anthropic 为 claude-sonnet-4-20250514, openai 为 gpt-4o-mini
model = "claude-sonnet-4-20250514"
max_tokens = 4096                    # 最大 token 数
default_theme = "default"           # 默认主题
theme_dir = ""                      # 主题目录
timeout = "60s"                     # 超时时间
```

**配置说明：**

| 配置项 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| `wechat_app_id` | string | ✅ | - | 微信公众号 AppID |
| `wechat_app_secret` | string | ✅ | - | 微信公众号 AppSecret |
| `log.level` | string | - | `debug` | 日志级别：debug, info, warn, error |
| `log.format` | string | - | `json` | 日志格式：json, console |
| `mcp.protocol` | string | - | - | MCP 协议：stdio, http, sse |
| `mcp.host` | string | - | - | MCP 服务地址 |
| `mcp.port` | int | - | - | MCP 服务端口 |
| `converter.enabled` | bool | - | `false` | 是否启用 AI 转换功能 |
| `converter.provider` | string | - | `anthropic` | AI 供应商：anthropic, openai |
| `converter.api_key` | string | - | - | API Key |
| `converter.base_url` | string | - | - | 自定义 API 地址 |
| `converter.model` | string | - | `claude-sonnet-4-20250514` (anthropic) / `gpt-4o-mini` (openai) | 使用的模型 |
| `converter.max_tokens` | int | - | `4096` | 最大 token 数 |
| `converter.default_theme` | string | - | `default` | 默认主题 |
| `converter.theme_dir` | string | - | - | 主题目录 |
| `converter.timeout` | string | - | `60s` | 超时时间 |

**环境变量（优先级高于配置文件）：**

| 环境变量 | 说明 |
|---------|------|
| `WECHAT_APP_ID` | 微信公众号 AppID |
| `WECHAT_APP_SECRET` | 微信公众号 AppSecret |
| `AI_API_KEY` | AI API Key（支持 Anthropic 和 OpenAI） |
| `AI_BASE_URL` | 自定义 API 地址（支持所有供应商） |

### 3. 运行

```bash
./wechat-mcp
```

指定配置文件：

```bash
./wechat-mcp -c /path/to/config.toml
```

使用环境变量：

```bash
export WECHAT_APP_ID="your_app_id"
export WECHAT_APP_SECRET="your_app_secret"
./wechat-mcp
```

或一次性运行：

```bash
WECHAT_APP_ID="xxx" WECHAT_APP_SECRET="xxx" ./wechat-mcp -c config.toml
```

## MCP 工具

### upload_material

上传图片素材到微信公众号。

**参数：**
- `file_path` (必填): 本地文件路径或 HTTP URL
- `retry` (可选): 是否启用重试，默认 false

**返回：**
- MediaID
- 微信图片 URL

---

### create_draft

创建微信公众号文章草稿。

**参数：**
- `title` (必填): 文章标题
- `content` (必填): 文章正文内容（支持 HTML）
- `thumb_media_id` (必填): 封面图片 media_id
- `author` (可选): 作者名称
- `digest` (可选): 文章摘要
- `content_source_url` (可选): 原文链接

**返回：**
- MediaID
- 草稿查看链接

---

### create_newspic_draft

创建小绿书（图片）草稿。

**参数：**
- `title` (必填): 文章标题
- `content` (必填): 文章正文内容
- `image_paths` (必填): 图片路径列表，支持本地路径或 HTTP URL

**返回：**
- MediaID
- 草稿查看链接

---

### get_access_token

获取微信公众号 AccessToken（用于调试）。

**返回：**
- AccessToken
- 过期时间

---

### download_file

下载文件到临时目录，或验证本地文件路径。

**参数：**
- `url_or_path` (必填): 文件 URL 或本地路径

**返回：**
- 本地文件路径

---

### convert_markdown

将 Markdown 内容转换为微信公众号兼容的 HTML（需要启用 AI 转换器）。

**参数：**
- `markdown` (必填): Markdown 内容
- `theme` (可选): 主题名称 (default, elegant, tech, minimalist)
- `custom_prompt` (可选): 自定义提示词

**返回：**
- 转换后的 HTML 内容

---

### list_themes

列出所有可用的主题。

**参数：** 无

**返回：**
- 可用主题列表
---

## 主题说明

本项目内置了 4 个主题，可通过 `convert_markdown` 工具的 `theme` 参数指定。

### 可用主题

| 主题 | 说明 | 适用场景 |
|------|------|----------|
| `default` | 简洁大方的通用风格 | 通用场景，技术文章，生活分享 |
| `elegant` | 精致柔和的排版风格 | 情感文章，文艺分享，生活随笔 |
| `tech` | 简洁专业的技术风格 | 技术文章，编程教程，代码分享 |
| `minimalist` | 简约清爽的阅读体验 | 深度阅读，思考类文章，简洁主义 |

### 主题详细说明

#### default (默认主题)

简洁大方的通用风格，适合大多数公众号文章。

- **配色**: 蓝灰色调，专业简洁
- **字号**: 中等（16px）
- **行高**: 1.75
- **标签**: default, simple, general

#### elegant (优雅主题)

精致柔和的排版风格，给人温暖舒适的阅读体验。

- **配色**: 柔和暖色调
- **字号**: 中等
- **行高**: 1.8
- **标签**: elegant, soft, literary

#### tech (技术主题)

简洁专业的技术风格，适合程序员阅读。

- **配色**: 深色系，高对比度
- **代码块**: Monokai 风格深色背景
- **行高**: 1.7
- **标签**: tech, code, developer

#### minimalist (极简主题)

简约清爽的阅读体验，去除一切不必要的装饰。

- **配色**: 黑白灰，纯粹阅读
- **特点**: 大量留白，专注于内容
- **行高**: 1.8
- **标签**: minimalist, simple, clean

### 使用示例

```bash
# 使用默认主题转换
convert_markdown(markdown="# Hello", theme="default")

# 使用技术主题转换
convert_markdown(markdown="# 代码示例", theme="tech")

# 使用优雅主题转换
convert_markdown(markdown="# 随笔", theme="elegant")

# 使用极简主题转换
convert_markdown(markdown="# 深度思考", theme="minimalist")
```

### 自定义主题

你可以在 `themes/` 目录下创建自定义主题配置文件（参考现有主题格式），或通过 `converter.theme_dir` 配置指定自定义主题目录。

## 协议说明

| 协议 | 说明 | 适用场景 |
|------|------|----------|
| stdio | 标准输入输出 | 本地调试、CLI 工具集成 |
| HTTP | Streamable HTTP | 生产环境部署 |
| SSE | Server-Sent Events | 需要实时推送的场景 |

## 技术栈

- [mcp-go](https://github.com/mark3labs/mcp-go) - MCP Go SDK
- [silenceper/wechat](https://github.com/silenceper/wechat) - 微信 SDK
- [zap](https://github.com/uber-go/zap) - 日志库

## 许可证

本项目基于 GNU General Public License v3 (GPLv3) 开源，详见 [LICENSE](LICENSE) 文件。

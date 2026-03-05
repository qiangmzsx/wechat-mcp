# WeChat MCP Server

微信公众号 MCP (Model Context Protocol) 服务器，提供与微信公众号平台交互的能力。

## 功能特性

- 支持 MCP 协议：stdio、HTTP、SSE
- 上传图片素材到微信公众号
- 创建和管理文章草稿
- 创建小绿书（图片）草稿
- 获取 AccessToken
- 文件下载

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

复制配置文件并填入微信公众平台凭证：

```bash
cp config.example.toml config.toml
```

编辑 `config.toml`：

```toml
# 微信公众号配置
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
```

### 3. 运行

```bash
./wechat-mcp
```

指定配置文件：

```bash
./wechat-mcp -c /path/to/config.toml
```

```bash
go run .
```

指定配置文件：

```bash
go run . /path/to/config.toml
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
- `author` (可选): 作者名称
- `digest` (可选): 文章摘要
- `content_source_url` (可选): 原文链接
- `thumb_media_id` (必填): 封面图片 media_id

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

MIT

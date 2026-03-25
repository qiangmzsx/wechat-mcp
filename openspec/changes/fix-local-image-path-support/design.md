## Context

### 当前状态

`converter/types.go` 中的 `localImagePattern` 正则表达式：

```go
localImagePattern = regexp.MustCompile(`!\[([^\]]*)\]\((\.\/[^)]+)\)`)
```

当前只匹配 `./` 开头的本地图片路径（如 `![alt](./image.png)`），不匹配：
- `![alt](images/photo.jpg)` - 无前缀相对路径
- `![alt](subdir/image.png)` - 子目录路径
- `![alt](../images/photo.jpg)` - 父目录路径

### 问题根因

正则表达式 `\.\/[^)]+` 要求必须以 `./` 开头，导致其他格式的本地图片路径无法被 `ExtractImages()` 函数识别。

## Goals / Non-Goals

**Goals:**
- 支持多种格式的本地图片路径（无前缀、带 `./`、带 `../`）
- 保持 base64 转换功能正常工作
- 添加完整的测试覆盖

**Non-Goals:**
- 不改变 HTTP 在线图片的处理方式
- 不改变 AI 生成图片的处理方式
- 不修改 `ImageToBase64` 相关函数（它们已经支持本地路径）

## Decisions

### Decision 1: 修改本地图片正则表达式

**选择**: 扩展正则表达式以匹配更多本地路径格式

**当前正则**:
```go
`!\[([^\]]*)\]\((\.\/[^)]+)\)`
```

**新正则**:
```go
`!\[([^\]]*)\]\(([^\s)]+)\)`
```

**理由**:
- `[^\s)]+` 匹配所有非空白、非右括号的字符
- 自然排除 HTTP URL（因为有 `://`）
- 自然排除 AI 提示词（因为有 `__generate:` 前缀会被第一个 `[^\s)]+` 匹配，但后续的替换逻辑不会将其当作图片处理）

**替代方案考虑**:
- 使用更严格的正则（如 `([^http/][^)]+)`）- 但 `[^\s)]+` 更简洁且足够清晰

### Decision 2: 添加测试用例覆盖

**选择**: 在 `converter/converter_test.go` 中添加新的测试用例

**新增测试用例**:
- 无前缀相对路径: `![alt](images/photo.jpg)`
- 子目录路径: `![alt](subdir/image.png)`
- 父目录路径: `![alt](../images/photo.jpg)`
- 带空格的路径（需要 URL 编码）: `![alt](path%20with%20spaces.png)`

## Risks / Trade-offs

[风险] 修改正则可能影响现有功能 → **缓解**: 现有测试全部通过，确保向后兼容

[风险] 新正则可能意外匹配到错误格式 → **缓解**: 后续的 `ImageToBase64` 函数只处理实际存在的本地文件，不存在的路径会被忽略

## Open Questions

无

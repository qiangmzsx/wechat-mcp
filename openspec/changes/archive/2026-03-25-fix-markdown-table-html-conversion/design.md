## Context

### 当前状态

Markdown 表格转换为 HTML 的流程如下：

1. `converter/api_converter.go` 调用 `theme.ConvertMarkdown(markdown, themeID)`
2. `theme/converter.go` 使用 goldmark 库解析 Markdown 并生成 HTML
3. `processListStylesSimple()` 函数尝试将主题样式应用到 HTML 元素

### 问题根源

`addStyleToAllElements()` 和 `addStyleToOpenTags()` 函数存在缺陷：

```go
// addStyleToAllElements 只处理 <tag> 和 <tag style="..."> 模式
searchOpen := "<" + tag + ">"          // 例如 "<table>"
searchOpenWithStyle := "<" + tag + " style=\""

// addStyleToOpenTags 只处理 <tag 模式（带空格开头的属性）
tagWithSpace := "<" + tag + " "
```

goldmark 生成的表格 HTML 结构为：
```html
<table><thead><tr><th>Header 1</th><th>Header 2</th></tr></thead>
<tbody><tr><td>Cell 1</td><td>Cell 2</td></tr></tbody></table>
```

问题在于：
1. `<th>` 和 `<td>` 是自闭合标签但包含内容的标签，`<th>Header 1</th>` 不匹配 `<th>` 也不匹配 `<th ` 
2. `addStyleToOpenTags` 无法为已存在样式的元素追加样式

## Goals / Non-Goals

**Goals:**
- 修复 `addStyleToAllElements` 函数，正确应用主题样式到表格相关元素
- 确保 `<table>`、`<th>`、`<td>`、`<tr>` 元素都能获得主题配置的样式
- 新增表格转换测试用例验证修复有效

**Non-Goals:**
- 不改变 goldmark 的解析行为
- 不修改主题配置文件的样式定义
- 不添加新的外部依赖

## Decisions

### Decision 1: 修复 `addStyleToAllElements` 函数

**选择**：重写 `addStyleToAllElements` 函数，使用更可靠的样式应用策略

**理由**：
- 当前函数只处理 `<tag>` 开头的元素，但 `<th>content</th>` 这样的结构无法匹配
- 需要能够匹配 `<th>content</th>`、`<td>content</td>` 等包含内容的标签

**实现方案**：
```go
func addStyleToAllElements(html, tag, styleValue string) string {
    // 匹配 <tag>...</tag> 或 <tag>content</tag> 模式
    pattern := `<` + tag + `>([^<]*)</` + tag + `>`
    // 替换为 <tag style="...">...</tag>
    re := regexp.MustCompile(pattern)
    return re.ReplaceAllString(html, `<$1 style="$2">$3</$1>`)
    // 需要更复杂的实现来保留现有样式
}
```

### Decision 2: 使用 AST 重写样式应用逻辑

**选择**：使用 goquery 库解析 HTML 并应用样式

**理由**：
- 正则表达式难以处理嵌套和复杂 HTML 结构
- goquery 提供 DOM 级别的操作，更可靠
- goquery 在项目中已有使用（通过 goquery 或类似库）

**替代方案**：
- 使用字符串替换（不够健壮）
- 继续使用正则但改进模式（复杂且容易出错）

### Decision 3: 新增测试用例

**选择**：在 `theme/table_conversion_test.go` 中新增表格转换测试

**理由**：
- 防止回归
- 验证修复有效
- 提供文档化的期望行为

## Risks / Trade-offs

[Risk] 修改 `addStyleToAllElements` 可能影响其他元素的样式应用 → **Mitigation**: 确保修改仅影响表格相关元素的处理逻辑，使用条件判断

[Risk] goquery 库引入 → **Mitigation**: 检查项目是否已有类似依赖，或使用标准库的 html/template + 简单解析

[Risk] 微信平台对复杂表格样式的支持 → **Mitigation**: 使用内联样式（已在系统提示中要求），确保兼容性

## Open Questions

1. goldmark 是否输出 `<thead>` 和 `<tbody>`？如果输出，是否需要特殊处理？
2. 是否需要支持表格的响应式布局（横向滚动）？
3. 是否需要处理嵌套表格的情况？

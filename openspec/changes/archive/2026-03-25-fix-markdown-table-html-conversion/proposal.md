## Why

当前 Markdown 表格转换为 HTML 时，表格的样式（如边框、背景色、内边距等）未能正确应用。虽然主题（theme）配置文件中定义了 `table`、`th`、`td`、`tr` 的样式，但实际转换后的 HTML 表格缺乏预期的视觉效果。

这导致微信公众号文章中的表格显示为浏览器默认样式，与主题风格不统一，影响阅读体验。

## What Changes

1. **修复表格样式应用逻辑**：确保 goldmark 生成的 HTML 表格元素（`<table>`、`<th>`、`<td>`、`<tr>`）能正确应用主题配置的 CSS 样式

2. **添加表格转换测试**：新增针对 Markdown 表格到 HTML 转换的单元测试，验证样式正确应用

3. **确保表格响应式适配**：验证表格在小屏幕设备上的显示效果

## Capabilities

### New Capabilities
- `markdown-table-conversion`: 完善 Markdown 表格到 HTML 的转换功能，确保主题样式（边框、背景、内边距等）正确应用到生成的 HTML 表格元素

### Modified Capabilities
- 无（现有能力的行为需求未改变，仅修复实现缺陷）

## Impact

- **影响的代码模块**：`theme/converter.go` 中的 `addStyleToAllElements` 和 `addStyleToOpenTags` 函数
- **测试文件**：需要新增 `theme/table_conversion_test.go`
- **主题配置**：29 个主题的 `table`、`th`、`td`、`tr` 样式定义已存在，无需修改

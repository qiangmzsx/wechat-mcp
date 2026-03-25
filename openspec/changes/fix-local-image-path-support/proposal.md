## Why

当前 `localImagePattern` 正则表达式只支持 `./` 开头的本地图片路径（如 `![alt](./image.png)`），但用户经常使用不带 `./` 前缀的相对路径（如 `![alt](images/photo.jpg)`），这些路径目前无法被正确识别和转换为 base64。

## What Changes

1. **修改本地图片正则表达式**：扩展 `localImagePattern` 以支持更多本地图片路径格式
2. **支持无前缀相对路径**：如 `images/photo.jpg`、`./subdir/image.png`、`../images/photo.jpg`
3. **保持向后兼容**：确保现有使用 `./` 前缀的图片继续正常工作

## Capabilities

### New Capabilities

- `local-image-path-support`: 支持多种格式的本地图片路径（无前缀相对路径、带 `./` 前缀、带 `../` 前缀）

## Impact

- **影响代码**：`converter/types.go` 中的 `localImagePattern` 正则表达式
- **测试文件**：`converter/converter_test.go` 需要新增测试用例验证新路径格式

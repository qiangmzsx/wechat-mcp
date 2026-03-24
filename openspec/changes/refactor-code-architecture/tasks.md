## 1. Code Quality Fixes (Prerequisites)

- [x] 1.1 Fix `NewAIConverter` signature mismatch in `converter/converter_test.go` - tests pass 2 args but function accepts 1
- [x] 1.2 Fix `zap.Error(fmt.Errorf(result.Error))` in `mcp/server.go:634` - use `zap.String("error", result.Error)` instead
- [x] 1.3 Fix `GeneratePlaceholder` bug in `converter/types.go:191` - use `fmt.Sprintf` instead of `string(rune('0'+index))`

## 2. Extract Shared Utilities

- [x] 2.1 Extract `maskMediaID` function from `wechat/service.go` and `mcp/server.go` to shared utility package (created `internal/util/util.go`)
- [ ] 2.2 Review and consolidate HTTP download logic duplication between `DownloadFile` and `imageToBase64FromURL`

## 3. MCP Server Refactoring

- [x] 3.1 Split `mcp/server.go` - move tool handlers to new `mcp/handlers.go` file
- [x] 3.2 Keep only bootstrap/registration logic in `mcp/server.go`
- [ ] 3.3 Verify all 7 tools still work after split

## 4. Theme Management Consolidation

- [ ] 4.1 Rename `converter/unified_theme.go` to `converter/theme_registry.go` for clarity
- [ ] 4.2 Move `ThemeManager` interface definition into `converter/theme_registry.go` alongside implementation
- [ ] 4.3 Remove unused `AILLM` interface from `converter/types.go` (or document if it's intentional for future use)

## 5. Error Handling Standardization

- [ ] 5.1 Audit all `fmt.Errorf` calls to ensure consistent `%w` wrapping pattern
- [ ] 5.2 Replace raw string errors with proper error types where appropriate

## 6. Context Propagation

- [ ] 6.1 Audit functions that create `context.Background()` instead of using passed `ctx`
- [ ] 6.2 Update `aiConverter.Convert()` to propagate context to AI provider calls

## 7. Verification

- [ ] 7.1 Run `go build ./...` to verify no build errors
- [ ] 7.2 Run `go test ./...` to verify all tests pass
- [ ] 7.3 Run `go vet ./...` to verify no warnings
- [ ] 7.4 Manual test of all 7 MCP tools to verify backward compatibility

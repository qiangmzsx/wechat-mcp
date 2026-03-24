## Why

The current codebase has architectural issues that make it difficult to maintain and extend:

1. **Monolithic MCP Server** (`mcp/server.go` - 682 lines): All tool handlers are in a single file, violating single responsibility principle
2. **Mixed concerns in converter package**: Interface definitions, implementations, and theme management are interleaved across files
3. **Inconsistent error handling**: Some places use `fmt.Errorf` with `%w`, others don't wrap errors at all
4. **Theme architecture split**: `ThemeManager` interface is in `converter/theme.go` but implementation is in `converter/unified_theme.go`, creating confusion
5. **Duplicated patterns**: `SimpleConverter` and `aiConverter` share some logic but are not cleanly abstracted

These issues slow down development and increase bug risk. The refactoring will improve code quality while maintaining 100% backward compatibility for all MCP tools.

## What Changes

### Code Organization
- Split `mcp/server.go` into smaller, focused files by domain (handlers, tools registration, server bootstrap)
- Move `ThemeManager` interface and implementation to the same package (`converter/`)
- Consolidate converter-related types into `converter/types.go`

### Design Pattern Improvements
- Apply cleaner error wrapping convention throughout
- Extract common patterns into shared utilities
- Improve dependency injection in the converter package

### Backward Compatibility (MUST maintain)
- **All 7 MCP tools unchanged**: `upload_material`, `create_draft`, `create_newspic_draft`, `get_access_token`, `download_file`, `convert_markdown`, `list_themes`
- Configuration file format unchanged
- All environment variables unchanged
- API behavior identical

## Capabilities

### New Capabilities
<!-- No new capabilities - pure refactoring -->

### Modified Capabilities
<!-- No spec-level behavior changes - all capabilities preserved as-is -->

## Impact

### Affected Code
- `mcp/server.go` → split into multiple files
- `converter/theme.go` → move interface to proper location
- `converter/unified_theme.go` → rename for clarity
- Error handling utilities across packages

### Unaffected (Backward Compatible)
- All MCP tool signatures and behavior
- Configuration schema (`config.toml`)
- External dependencies
- Theme files in `themes/` directory

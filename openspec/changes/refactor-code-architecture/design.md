## Context

The wechat-mcp codebase has grown organically without consistent architectural governance. Key issues identified:

1. **Monolithic `mcp/server.go`** (682 lines): All 7 tool handlers plus protocol logic in single file
2. **Scattered theme management**: `ThemeManager` interface defined in `converter/theme.go`, implementation in `converter/unified_theme.go`
3. **Inconsistent error wrapping**: Some use `fmt.Errorf("...: %w", err)`, others plain `fmt.Errorf("...")`
4. **Context ignored in places**: Functions receive `context.Context` but create `context.Background()` internally
5. **Converter types spread across files**: Interface in `types.go`, AI impl in `converter.go`, API impl in `api_converter.go`

**Constraints**:
- Zero behavior changes to MCP tools
- Configuration schema unchanged
- All existing tests must pass

## Goals / Non-Goals

**Goals:**
- Improve code organization and maintainability
- Apply consistent error handling patterns
- Clean up theme management architecture
- Make the codebase easier to navigate and extend

**Non-Goals:**
- No new features or capabilities
- No API/interface changes
- No dependency updates
- No performance optimization (this is a code quality refactor only)

## Decisions

### Decision 1: Split `mcp/server.go` into package files

**Choice**: Create `mcp/handlers.go` for tool handlers, keep `mcp/server.go` for bootstrap

**Rationale**: Handlers are ~400 lines of repetitive switch/case argument parsing. Separating by concern makes the code easier to scan. Each handler in its own file would be over-engineering for 7 handlers.

**Alternative considered**: One file per handler (`upload_material.go`, `create_draft.go`, etc.) → rejected as over-engineering for 7 handlers that share similar structure.

### Decision 2: Consolidate theme management in `converter/`

**Choice**: Keep `ThemeManager` interface in `converter/theme.go`, move implementation from `unified_theme.go` to a file that better describes its purpose (e.g., `converter/theme_registry.go`)

**Rationale**: The `unified_theme.go` name is confusing - it doesn't describe what the file does. The theme system is primarily used by the converter, so keeping it in the converter package makes sense.

**Alternative considered**: Move all theme code to `theme/` package entirely → rejected because the converter package needs tight integration with themes.

### Decision 3: Consistent error wrapping

**Choice**: Use `fmt.Errorf("operation: %w", err)` pattern everywhere

**Rationale**: This is the Go standard for error wrapping. The `%w` verb allows errors to be unwrapped with `errors.Is()` / `errors.As()`.

**Enforcement**: Will add to code review checklist.

### Decision 4: Context propagation

**Choice**: Pass `ctx` through call chains instead of creating `context.Background()`

**Rationale**: Allows callers to cancel long-running operations. Standard Go practice.

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| Breaking MCP tool behavior | Extensive test suite exists; will run full test suite after refactor |
| Introducing bugs in error paths | Focus refactoring on structure only, not logic |
| Merge conflicts if multiple devs work simultaneously | Complete refactoring in single branch before merge |

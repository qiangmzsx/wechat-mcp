## 1. Investigation & Analysis

- [x] 1.1 Write a test to capture current goldmark HTML output for markdown tables
- [x] 1.2 Verify goldmark's exact HTML output structure for `<table>`, `<thead>`, `<tbody>`, `<tr>`, `<th>`, `<td>` elements
- [x] 1.3 Confirm which selectors (table, th, td, tr) are NOT receiving styles

## 2. Fix Style Application Function

- [x] 2.1 Enable goldmark table extension (fixes `<th>content</th>` and `<td>content</td>` patterns)
- [x] 2.2 Verify styles are applied to `<table>` element correctly
- [x] 2.3 Handle style merging when element already has existing inline styles

## 3. Testing

- [x] 3.1 Create `theme/table_conversion_test.go` with test cases for table styling
- [x] 3.2 Add test for basic table with header (single column)
- [x] 3.3 Add test for multi-column table header
- [x] 3.4 Add test for table with multiple data rows
- [x] 3.5 Add test for mixed content (bold/italic) in table cells
- [x] 3.6 Run all existing tests to ensure no regression

## 4. Verification

- [x] 4.1 Verify table styling works with "default" theme
- [x] 4.2 Verify table styling works with "wechat" theme
- [x] 4.3 Verify table styling works with 3-4 other themes (tech, minimalist, dracula)
- [x] 4.4 Run `go vet` and `go fmt` on modified files

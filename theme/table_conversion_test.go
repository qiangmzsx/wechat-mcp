package theme

import (
	"strings"
	"testing"
)

// TestTableConversion_BasicTable tests basic table with header (single column)
func TestTableConversion_BasicTable(t *testing.T) {
	markdown := `| Header |
|-------|
| Cell  |`

	html := Convert(markdown, "default")

	// Verify table element has style
	if !strings.Contains(html, `<table style="`) {
		t.Error("table element should have style attribute")
	}

	// Verify th element has style
	if !strings.Contains(html, `<th style="`) {
		t.Error("th element should have style attribute")
	}

	// Verify td element has style
	if !strings.Contains(html, `<td style="`) {
		t.Error("td element should have style attribute")
	}

	// Verify specific style properties
	if !strings.Contains(html, "border-collapse: collapse") {
		t.Error("table should have border-collapse style")
	}

	if !strings.Contains(html, "background-color:") {
		t.Error("th should have background-color style")
	}

	t.Logf("Basic table HTML:\n%s\n", html)
}

// TestTableConversion_MultiColumnHeader tests multi-column table header
func TestTableConversion_MultiColumnHeader(t *testing.T) {
	markdown := `| Header 1 | Header 2 | Header 3 |
|----------|----------|----------|
| Cell 1   | Cell 2   | Cell 3   |`

	html := Convert(markdown, "default")

	// Count occurrences of styled elements
	thCount := strings.Count(html, `<th style="`)
	tdCount := strings.Count(html, `<td style="`)

	if thCount != 3 {
		t.Errorf("expected 3 th elements with styles, got %d", thCount)
	}

	if tdCount != 3 {
		t.Errorf("expected 3 td elements with styles, got %d", tdCount)
	}

	t.Logf("Multi-column table HTML:\n%s\n", html)
}

// TestTableConversion_MultipleDataRows tests table with multiple data rows
func TestTableConversion_MultipleDataRows(t *testing.T) {
	markdown := `| Column |
|--------|
| Row 1  |
| Row 2  |
| Row 3  |
| Row 4  |`

	html := Convert(markdown, "default")

	// Should have 1 header and 4 data cells
	thCount := strings.Count(html, `<th style="`)
	tdCount := strings.Count(html, `<td style="`)

	if thCount != 1 {
		t.Errorf("expected 1 th element, got %d", thCount)
	}

	if tdCount != 4 {
		t.Errorf("expected 4 td elements, got %d", tdCount)
	}

	// Verify each td has style
	for i := 1; i <= 4; i++ {
		if !strings.Contains(html, `<td style="`) {
			t.Error("td elements should have style attribute")
		}
	}

	t.Logf("Multiple rows table HTML:\n%s\n", html)
}

// TestTableConversion_MixedContent tests bold/italic content in table cells
func TestTableConversion_MixedContent(t *testing.T) {
	markdown := `| Header |
|--------|
| **Bold** |
| *Italic* |
| ***Bold Italic*** |`

	html := Convert(markdown, "default")

	// Verify strong tags exist (for bold)
	if !strings.Contains(html, "<strong") {
		t.Error("HTML should contain strong tags for bold text")
	}

	// Verify em tags exist (for italic)
	if !strings.Contains(html, "<em") {
		t.Error("HTML should contain em tags for italic text")
	}

	// Verify table structure is preserved
	if !strings.Contains(html, "<th") || !strings.Contains(html, "<td") {
		t.Error("table structure should be preserved with th and td elements")
	}

	t.Logf("Mixed content table HTML:\n%s\n", html)
}

// TestTableConversion_DefaultTheme tests table styling with default theme
func TestTableConversion_DefaultTheme(t *testing.T) {
	markdown := `| H1 | H2 |
|----|----|
| A  | B  |`

	html := Convert(markdown, "default")

	// Verify styles are applied
	if !strings.Contains(html, `<table style="`) {
		t.Error("default theme: table should have style")
	}

	if !strings.Contains(html, `width: 100%`) {
		t.Error("default theme: table should have width: 100%")
	}

	t.Logf("Default theme table HTML:\n%s\n", html)
}

// TestTableConversion_WechatTheme tests table styling with wechat theme
func TestTableConversion_WechatTheme(t *testing.T) {
	markdown := `| H1 | H2 |
|----|----|
| A  | B  |`

	html := Convert(markdown, "wechat")

	// Verify styles are applied
	if !strings.Contains(html, `<table style="`) {
		t.Error("wechat theme: table should have style")
	}

	// Wechat theme uses green accent
	if !strings.Contains(html, "background-color:") {
		t.Error("wechat theme: should have background-color style")
	}

	t.Logf("Wechat theme table HTML:\n%s\n", html)
}

// TestTableConversion_MultipleThemes tests table styling with multiple themes
func TestTableConversion_MultipleThemes(t *testing.T) {
	themes := []string{"default", "wechat", "tech", "minimalist", "dracula"}

	markdown := `| H1 | H2 |
|----|----|
| A  | B  |`

	for _, themeID := range themes {
		html := Convert(markdown, themeID)

		// All themes should apply table styles
		if !strings.Contains(html, `<table style="`) {
			t.Errorf("theme %s: table should have style", themeID)
		}

		if !strings.Contains(html, `<th style="`) {
			t.Errorf("theme %s: th should have style", themeID)
		}

		if !strings.Contains(html, `<td style="`) {
			t.Errorf("theme %s: td should have style", themeID)
		}

		t.Logf("Theme %s table HTML:\n%s\n", themeID, html)
	}
}

// TestTableConversion_ComplexTable tests complex table with all elements
func TestTableConversion_ComplexTable(t *testing.T) {
	markdown := `| Header 1 | Header 2 | Header 3 |
|----------|----------|----------|
| Cell 1-1 | Cell 1-2 | Cell 1-3 |
| Cell 2-1 | Cell 2-2 | Cell 2-3 |
| Cell 3-1 | Cell 3-2 | Cell 3-3 |`

	html := Convert(markdown, "default")

	// Verify structure
	if !strings.Contains(html, "<table") {
		t.Error("HTML should contain table element")
	}

	if !strings.Contains(html, "<thead>") {
		t.Error("HTML should contain thead element")
	}

	if !strings.Contains(html, "<tbody>") {
		t.Error("HTML should contain tbody element")
	}

	// Verify th elements (3 headers) - styled th elements
	thCount := strings.Count(html, `<th style="`)
	if thCount != 3 {
		t.Errorf("expected 3 styled th elements, got %d", thCount)
	}

	// Verify td elements (9 data cells) - styled td elements
	tdCount := strings.Count(html, `<td style="`)
	if tdCount != 9 {
		t.Errorf("expected 9 styled td elements, got %d", tdCount)
	}

	// Verify all th elements have styles
	styledThCount := strings.Count(html, `<th style="`)
	if styledThCount != 3 {
		t.Errorf("expected 3 styled th elements, got %d", styledThCount)
	}

	// Verify all td elements have styles
	styledTdCount := strings.Count(html, `<td style="`)
	if styledTdCount != 9 {
		t.Errorf("expected 9 styled td elements, got %d", styledTdCount)
	}

	t.Logf("Complex table HTML:\n%s\n", html)
}

## ADDED Requirements

### Requirement: Table Header Styling
The system SHALL apply the theme's `th` style to all table header cells (`<th>`) in the converted HTML output.

#### Scenario: Basic table with header
- **WHEN** Markdown containing a table with header row is converted
- **THEN** each `<th>` element SHALL contain the style attributes defined in the theme's `th` configuration, including background-color, padding, text-align, font-weight, color, and border

#### Scenario: Multi-column header
- **WHEN** Markdown containing a table with multiple header columns is converted
- **THEN** each `<th>` element in the header row SHALL all receive the same theme-defined styles

### Requirement: Table Data Cell Styling
The system SHALL apply the theme's `td` style to all table data cells (`<td>`) in the converted HTML output.

#### Scenario: Basic table with data cells
- **WHEN** Markdown containing a table with data rows is converted
- **THEN** each `<td>` element SHALL contain the style attributes defined in the theme's `td` configuration, including padding, border, and color

#### Scenario: Mixed content in data cells
- **WHEN** Markdown table cell contains bold or italic text
- **THEN** the `<td>` element SHALL have the theme styles applied AND the inner formatting SHALL be preserved

### Requirement: Table Row Styling
The system SHALL apply the theme's `tr` style to all table rows (`<tr>`) in the converted HTML output.

#### Scenario: Table row background
- **WHEN** Markdown table is converted
- **THEN** each `<tr>` element SHALL contain the style attributes defined in the theme's `tr` configuration

### Requirement: Table Container Styling
The system SHALL apply the theme's `table` style to the table container (`<table>`) element in the converted HTML output.

#### Scenario: Table width and border
- **WHEN** Markdown table is converted
- **THEN** the `<table>` element SHALL contain the style attributes defined in the theme's `table` configuration, including width, margin, border-collapse, and font-size

### Requirement: Style Merging with Existing Styles
The system SHALL merge theme styles with any existing inline styles on table elements without losing existing style information.

#### Scenario: Element has existing inline style
- **WHEN** goldmark outputs a table element that already has inline styles
- **THEN** the theme styles SHALL be merged with (not replace) the existing styles

### Requirement: All Themes Support Table Styling
The system SHALL apply table styles consistently across all available themes.

#### Scenario: Convert table with wechat theme
- **WHEN** Markdown table is converted using the "wechat" theme
- **THEN** the table elements SHALL have the styles defined in `themes/wechat.toml` for table, th, td, tr

#### Scenario: Convert table with default theme
- **WHEN** Markdown table is converted using the "default" theme
- **THEN** the table elements SHALL have the styles defined in `themes/default.toml` for table, th, td, tr

#### Scenario: Convert table with all themes
- **WHEN** Markdown table is converted using any available theme
- **THEN** the table elements SHALL have styles appropriate to that theme (as defined in each theme's .toml file)

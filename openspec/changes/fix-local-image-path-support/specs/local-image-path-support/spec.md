## ADDED Requirements

### Requirement: Local Image Path Pattern Matching
The system SHALL extract local image references from Markdown using an expanded regex pattern that matches multiple path formats.

#### Scenario: Extract local image with ./ prefix
- **WHEN** Markdown contains `![alt](./image.png)`
- **THEN** ExtractImages SHALL return an ImageRef with Type=ImageTypeLocal and Original="./image.png"

#### Scenario: Extract local image without ./ prefix
- **WHEN** Markdown contains `![alt](images/photo.jpg)`
- **THEN** ExtractImages SHALL return an ImageRef with Type=ImageTypeLocal and Original="images/photo.jpg"

#### Scenario: Extract local image from subdirectory
- **WHEN** Markdown contains `![alt](subdir/image.png)`
- **THEN** ExtractImages SHALL return an ImageRef with Type=ImageTypeLocal and Original="subdir/image.png"

#### Scenario: Extract local image from parent directory
- **WHEN** Markdown contains `![alt](../images/photo.jpg)`
- **THEN** ExtractImages SHALL return an ImageRef with Type=ImageTypeLocal and Original="../images/photo.jpg"

#### Scenario: Extract local image with spaces (URL encoded)
- **WHEN** Markdown contains `![alt](path%20with%20spaces.png)`
- **THEN** ExtractImages SHALL return an ImageRef with Type=ImageTypeLocal and Original="path%20with%20spaces.png"

### Requirement: Local Image Base64 Conversion
The system SHALL convert local image files to base64 data URIs for embedding in HTML.

#### Scenario: Convert existing local image to base64
- **WHEN** ImageToBase64 is called with an existing local file path
- **THEN** it SHALL return a data URI (data:image/xxx;base64,...) with the file's binary content base64 encoded

#### Scenario: Skip non-existent local image files
- **WHEN** ImageToBase64 is called with a non-existent file path
- **THEN** it SHALL return an error and SHALL NOT crash

### Requirement: HTML Replacement with Base64
The system SHALL replace local image src attributes in HTML with their base64 data URIs.

#### Scenario: Replace local image src with base64
- **WHEN** HTML contains `<img src="./image.png" />`
- **AND** ExtractImages has extracted this as ImageTypeLocal
- **THEN** ReplaceImagesWithBase64 SHALL replace the src with the base64 data URI

#### Scenario: Preserve HTTP image URLs
- **WHEN** HTML contains `<img src="https://example.com/image.png" />`
- **THEN** ReplaceImagesWithBase64 SHALL NOT modify HTTP/HTTPS URLs

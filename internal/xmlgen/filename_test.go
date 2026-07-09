package xmlgen

import (
	"strings"
	"testing"

	"github.com/leokr/blogger-xml-exporter/internal/config"
)

// TestFilenameValidTemplate tests filename generation with valid template.
func TestFilenameValidTemplate(t *testing.T) {
	xmlCfg := config.XMLConfig{
		Filename: "{{ .title }}",
	}

	values := map[string]interface{}{
		"title": "Test Report",
	}

	result := Filename(xmlCfg, map[string]interface{}{}, values)

	if !strings.Contains(result, "Test") || !strings.Contains(result, "Report") {
		t.Errorf("expect filename to contain 'Test' and 'Report', got: %s", result)
	}
	if !strings.HasSuffix(result, ".xml") {
		t.Errorf("expect filename to end with .xml, got: %s", result)
	}
}

// TestFilenameEmptyTemplate tests fallback behavior with empty template.
func TestFilenameEmptyTemplate(t *testing.T) {
	xmlCfg := config.XMLConfig{
		Filename: "",
	}

	result := Filename(xmlCfg, map[string]interface{}{}, map[string]interface{}{})

	if result != "post.xml" {
		t.Errorf("expect 'post.xml' for empty template, got: %s", result)
	}
}

// TestFilenameMissingFormField tests when referenced form field doesn't exist.
func TestFilenameMissingFormField(t *testing.T) {
	xmlCfg := config.XMLConfig{
		Filename: "{{ .missing_field | default \"fallback\" }}",
	}

	result := Filename(xmlCfg, map[string]interface{}{}, map[string]interface{}{})

	if !strings.Contains(result, "fallback") {
		t.Errorf("expect default fallback in filename, got: %s", result)
	}
	if !strings.HasSuffix(result, ".xml") {
		t.Errorf("expect filename to end with .xml, got: %s", result)
	}
}

// TestFilenameSpecialCharactersEscaped tests that unsafe chars are replaced.
func TestFilenameSpecialCharactersEscaped(t *testing.T) {
	xmlCfg := config.XMLConfig{
		Filename: "{{ .name }}",
	}

	values := map[string]interface{}{
		"name": "Test<Report>:file*?name",
	}

	result := Filename(xmlCfg, map[string]interface{}{}, values)

	// Special chars should be replaced with underscores
	if strings.ContainsAny(result, "<>:*?\"\\|") {
		t.Errorf("expect no unsafe characters in filename, got: %s", result)
	}
	if !strings.HasSuffix(result, ".xml") {
		t.Errorf("expect filename to end with .xml, got: %s", result)
	}
}

// TestFilenameLongNameTruncated tests that very long filenames are truncated.
func TestFilenameLongNameTruncated(t *testing.T) {
	longName := strings.Repeat("a", 200)
	xmlCfg := config.XMLConfig{
		Filename: "{{ .name }}",
	}

	values := map[string]interface{}{
		"name": longName,
	}

	result := Filename(xmlCfg, map[string]interface{}{}, values)

	// Result should be max 150 chars + .xml extension
	if len(result) > 154 {
		t.Errorf("expect truncated filename, got length: %d", len(result))
	}
	if !strings.HasSuffix(result, ".xml") {
		t.Errorf("expect filename to end with .xml, got: %s", result)
	}
}

// TestFilenameOnlyUnsafeChars tests with only unsafe characters.
func TestFilenameOnlyUnsafeChars(t *testing.T) {
	xmlCfg := config.XMLConfig{
		Filename: "{{ .name }}",
	}

	values := map[string]interface{}{
		"name": "<>:*?",
	}

	result := Filename(xmlCfg, map[string]interface{}{}, values)

	if result != "post.xml" {
		t.Errorf("expect fallback 'post.xml' for all-unsafe filename, got: %s", result)
	}
}

// TestFilenameWithSource tests template with source function for post data.
func TestFilenameWithSource(t *testing.T) {
	xmlCfg := config.XMLConfig{
		Filename: "{{ source \"title\" }}",
	}

	post := map[string]interface{}{
		"title": "Post Title From Source",
	}

	result := Filename(xmlCfg, post, map[string]interface{}{})

	// Result should contain the title parts and end with .xml
	if !strings.Contains(result, "Post") || !strings.Contains(result, "Title") {
		t.Errorf("expect source-derived title in filename, got: %s", result)
	}
	if !strings.HasSuffix(result, ".xml") {
		t.Errorf("expect filename to end with .xml, got: %s", result)
	}
}

// TestFilenameTrimsWhitespace tests that leading/trailing whitespace is trimmed.
func TestFilenameTrimsWhitespace(t *testing.T) {
	xmlCfg := config.XMLConfig{
		Filename: "{{ .name }}",
	}

	values := map[string]interface{}{
		"name": "  Report Name  ",
	}

	result := Filename(xmlCfg, map[string]interface{}{}, values)

	if strings.HasPrefix(result, " ") || strings.HasPrefix(result, "-") {
		t.Errorf("expect trimmed filename, got: %s", result)
	}
	if !strings.HasSuffix(result, ".xml") {
		t.Errorf("expect filename to end with .xml, got: %s", result)
	}
}

// TestFilenameInvalidTemplate tests with syntactically invalid template.
func TestFilenameInvalidTemplate(t *testing.T) {
	xmlCfg := config.XMLConfig{
		Filename: "{{ .name | invalid_function }}",
	}

	values := map[string]interface{}{
		"name": "Test",
	}

	result := Filename(xmlCfg, map[string]interface{}{}, values)

	// Should fallback to "post.xml" on template error
	if result != "post.xml" {
		t.Errorf("expect 'post.xml' fallback for invalid template, got: %s", result)
	}
}

// TestFilenameDotOnlyString tests with string containing only dots.
func TestFilenameDotOnlyString(t *testing.T) {
	xmlCfg := config.XMLConfig{
		Filename: "{{ .name }}",
	}

	values := map[string]interface{}{
		"name": "...",
	}

	result := Filename(xmlCfg, map[string]interface{}{}, values)

	if result != "post.xml" {
		t.Errorf("expect 'post.xml' for dot-only filename, got: %s", result)
	}
}

// TestFilenameMultipleExtensions tests that extra dots are handled.
func TestFilenameMultipleExtensions(t *testing.T) {
	xmlCfg := config.XMLConfig{
		Filename: "{{ .name }}",
	}

	values := map[string]interface{}{
		"name": "report.data.backup",
	}

	result := Filename(xmlCfg, map[string]interface{}{}, values)

	if !strings.HasSuffix(result, ".xml") {
		t.Errorf("expect filename to end with .xml, got: %s", result)
	}
	if !strings.Contains(result, "report_data_backup") && !strings.Contains(result, "report.data.backup") {
		t.Errorf("expect report name preserved in filename, got: %s", result)
	}
}

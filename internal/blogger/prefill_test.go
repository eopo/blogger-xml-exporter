package blogger

import (
	"testing"

	"github.com/leokr/blogger-xml-exporter/internal/config"
)

// TestStripHTML tests HTML tag removal.
func TestStripHTML(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"<p>Hello</p>", "Hello"},
		{"<div><span>Nested</span></div>", "Nested"},
		{"Text &amp; more", "Text & more"},
		{"<br/><hr/>", ""},
		{"No tags here", "No tags here"},
	}

	for _, tt := range tests {
		result := stripHTML(tt.input)
		if result != tt.expected {
			t.Errorf("stripHTML(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

// TestUnescapeUnicode tests Unicode escape sequence resolution.
func TestUnescapeUnicode(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"\\u0041", "A"},
		{"\\u00E9", "é"},
		{"\\u4E2D", "中"},
		{"Hello \\u0041", "Hello A"},
		{"Invalid \\uXXXX remains as-is", "Invalid \\uXXXX remains as-is"},
	}

	for _, tt := range tests {
		result := unescapeUnicode(tt.input)
		if result != tt.expected {
			t.Errorf("unescapeUnicode(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

// TestResolveFieldSimplePath tests accessing a direct field.
func TestResolveFieldSimplePath(t *testing.T) {
	post := map[string]interface{}{
		"title": "Test Post",
		"id":    "123",
	}

	if ResolveField(post, "title") != "Test Post" {
		t.Error("expect 'Test Post' for title")
	}
	if ResolveField(post, "id") != "123" {
		t.Error("expect '123' for id")
	}
}

// TestResolveFieldNestedPath tests accessing nested fields via dot path.
func TestResolveFieldNestedPath(t *testing.T) {
	post := map[string]interface{}{
		"author": map[string]interface{}{
			"displayName": "John Doe",
			"email":       "john@example.com",
		},
	}

	if ResolveField(post, "author.displayName") != "John Doe" {
		t.Error("expect nested author displayName")
	}
	if ResolveField(post, "author.email") != "john@example.com" {
		t.Error("expect nested author email")
	}
}

// TestResolveFieldNilPost tests with nil post.
func TestResolveFieldNilPost(t *testing.T) {
	result := ResolveField(nil, "any.path")
	if result != nil {
		t.Error("expect nil for nil post")
	}
}

// TestResolveFieldEmptyPath tests with empty path.
func TestResolveFieldEmptyPath(t *testing.T) {
	post := map[string]interface{}{"field": "value"}
	result := ResolveField(post, "")
	if result != nil {
		t.Error("expect nil for empty path")
	}
}

// TestResolveFieldMissingField tests accessing non-existent field.
func TestResolveFieldMissingField(t *testing.T) {
	post := map[string]interface{}{"title": "Test"}
	result := ResolveField(post, "missing")
	if result != nil {
		t.Error("expect nil for missing field")
	}
}

// TestResolveFieldBrokenPath tests accessing non-existent nested path.
func TestResolveFieldBrokenPath(t *testing.T) {
	post := map[string]interface{}{
		"author": map[string]interface{}{
			"name": "John",
		},
	}

	result := ResolveField(post, "author.missing.deep")
	if result != nil {
		t.Error("expect nil for broken nested path")
	}
}

// TestResolveFieldsEmptyFields tests with no form fields.
func TestResolveFieldsEmptyFields(t *testing.T) {
	post := map[string]interface{}{"title": "Test"}
	fields := []config.FormItem{}

	result := ResolveFields(post, fields)

	if len(result) != 0 {
		t.Error("expect empty values map for no fields")
	}
}

// TestResolveFieldsWithSource tests fields with Source pointing to post data.
func TestResolveFieldsWithSource(t *testing.T) {
	post := map[string]interface{}{
		"title": "Test Post",
	}

	fields := []config.FormItem{
		{
			Name:   "postTitle",
			Type:   "text",
			Source: "title",
		},
	}

	result := ResolveFields(post, fields)

	if result["postTitle"] != "Test Post" {
		t.Errorf("expect 'Test Post', got: %v", result["postTitle"])
	}
}

// TestResolveFieldsWithTemplate tests fields with Template for computed values.
func TestResolveFieldsWithTemplate(t *testing.T) {
	post := map[string]interface{}{}
	fields := []config.FormItem{
		{
			Name:     "timestamp",
			Type:     "text",
			Template: "{{ now | date \"2006-01-02\" }}",
		},
	}

	result := ResolveFields(post, fields)

	// Just check that the template was evaluated (Sprig now() should work)
	if result["timestamp"] == "" {
		t.Error("expect non-empty timestamp from template")
	}
}

// TestToDisplayValueString tests string conversion.
func TestToDisplayValueString(t *testing.T) {
	if ToDisplayValue("hello") != "hello" {
		t.Error("expect string to pass through")
	}
}

// TestToDisplayValueNil tests nil value.
func TestToDisplayValueNil(t *testing.T) {
	if ToDisplayValue(nil) != "" {
		t.Error("expect empty string for nil")
	}
}

// TestToDisplayValueList tests list-to-csv conversion.
func TestToDisplayValueList(t *testing.T) {
	items := []interface{}{"apple", "banana", "cherry"}
	result := ToDisplayValue(items)
	expected := "apple, banana, cherry"

	if result != expected {
		t.Errorf("expect %q, got %q", expected, result)
	}
}

// TestToDisplayValueEmptyList tests empty list.
func TestToDisplayValueEmptyList(t *testing.T) {
	items := []interface{}{}
	result := ToDisplayValue(items)

	if result != "" {
		t.Errorf("expect empty string for empty list, got %q", result)
	}
}

// TestToDisplayValueNumber tests number conversion.
func TestToDisplayValueNumber(t *testing.T) {
	result := ToDisplayValue(42)
	expected := "42"

	if result != expected {
		t.Errorf("expect %q, got %q", expected, result)
	}
}

// TestBuildTemplateFuncsHasSource tests that source function is available.
func TestBuildTemplateFuncsHasSource(t *testing.T) {
	post := map[string]interface{}{
		"title": "Test",
	}

	funcs := BuildTemplateFuncs(post)

	if _, ok := funcs["source"]; !ok {
		t.Error("expect 'source' function in template funcs")
	}
}

// TestBuildTemplateFuncsHasSprig tests that Sprig functions are available.
func TestBuildTemplateFuncsHasSprig(t *testing.T) {
	funcs := BuildTemplateFuncs(map[string]interface{}{})

	// Check for some common Sprig functions
	if _, ok := funcs["upper"]; !ok {
		t.Error("expect 'upper' Sprig function")
	}
	if _, ok := funcs["lower"]; !ok {
		t.Error("expect 'lower' Sprig function")
	}
}

// TestBuildTemplateFuncsSourceFunctionWorks tests source function evaluation.
func TestBuildTemplateFuncsSourceFunctionWorks(t *testing.T) {
	post := map[string]interface{}{
		"author": map[string]interface{}{
			"name": "Alice",
		},
	}

	funcs := BuildTemplateFuncs(post)
	sourceFunc := funcs["source"].(func(string) string)

	result := sourceFunc("author.name")
	if result != "Alice" {
		t.Errorf("expect 'Alice', got %q", result)
	}
}

// TestRenderTemplateSimple tests basic template rendering.
func TestRenderTemplateSimple(t *testing.T) {
	post := map[string]interface{}{}
	result, err := RenderTemplate(post, "test", "Hello {{ .name }}", map[string]interface{}{"name": "World"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "Hello World" {
		t.Errorf("expect 'Hello World', got %q", result)
	}
}

// TestRenderTemplateInvalidSyntax tests template with syntax error.
func TestRenderTemplateInvalidSyntax(t *testing.T) {
	post := map[string]interface{}{}
	_, err := RenderTemplate(post, "test", "{{ invalid", map[string]interface{}{})

	if err == nil {
		t.Error("expect error for invalid syntax")
	}
}

// TestRenderTemplateWithSource tests template using source function.
func TestRenderTemplateWithSource(t *testing.T) {
	post := map[string]interface{}{
		"title": "Test Post",
	}

	result, err := RenderTemplate(post, "test", "{{ source \"title\" }}", map[string]interface{}{})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "Test Post" {
		t.Errorf("expect 'Test Post', got %q", result)
	}
}

// TestBaseFieldFuncMapExists tests that base function map is initialized.
func TestBaseFieldFuncMapExists(t *testing.T) {
	if BaseFieldFuncMap == nil {
		t.Error("expect BaseFieldFuncMap to be initialized")
	}

	// Check for expected functions
	if _, ok := BaseFieldFuncMap["stripHTML"]; !ok {
		t.Error("expect stripHTML in BaseFieldFuncMap")
	}
	if _, ok := BaseFieldFuncMap["unescapeUnicode"]; !ok {
		t.Error("expect unescapeUnicode in BaseFieldFuncMap")
	}
}

// TestResolvePresetsEmptyItems tests presets with no form items.
func TestResolvePresetsEmptyItems(t *testing.T) {
	post := map[string]interface{}{}
	items := []config.FormItem{}
	values := map[string]interface{}{}

	result := ResolvePresets(post, items, values)

	if len(result) != 0 {
		t.Error("expect no presets for empty items")
	}
}

// TestResolvePresetsNoGroups tests presets when no groups exist.
func TestResolvePresetsNoGroups(t *testing.T) {
	post := map[string]interface{}{}
	items := []config.FormItem{
		{
			Type: "text",
			Name: "field1",
		},
	}
	values := map[string]interface{}{}

	result := ResolvePresets(post, items, values)

	if len(result) != 0 {
		t.Error("expect no presets when items have no groups")
	}
}

// TestResolvePresetsNestedWalk tests that presets walks nested groups.
func TestResolvePresetsNestedWalk(t *testing.T) {
	post := map[string]interface{}{}
	items := []config.FormItem{
		{
			Type:  "group",
			Title: "Outer",
			Items: []config.FormItem{
				{
					Type:  "group",
					Title: "Inner",
					Presets: []config.Preset{
						{
							Label: "Inner Preset",
							Values: map[string]string{
								"field": "value",
							},
						},
					},
				},
			},
		},
	}
	values := map[string]interface{}{}

	result := ResolvePresets(post, items, values)

	// Should find the inner group's presets
	if len(result) == 0 {
		t.Error("expect presets to be found in nested walk")
	}
}

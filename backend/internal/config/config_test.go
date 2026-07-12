package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// testConfigPath returns the absolute path to config.test.yaml.
func testConfigPath() string {
	_, thisFile, _, _ := runtime.Caller(0)
	testDir := filepath.Dir(thisFile)
	// Navigate: internal/config/config_test.go -> internal -> backend -> root
	return filepath.Join(testDir, "..", "..", "..", "config.test.yaml")
}

// TestLoadExampleConfig verifies that config.test.yaml is valid at all times.
func TestLoadExampleConfig(t *testing.T) {
	cfg, err := Load(testConfigPath())
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if len(cfg.Form.Items) == 0 {
		t.Error("expect at least one form item")
	}
	if len(cfg.Form.Fields()) == 0 {
		t.Error("expect at least one form field")
	}
	if len(cfg.XML.Fields) == 0 {
		t.Error("expect at least one XML field")
	}
}

// TestLoadMissingFile tests loading a non-existent file.
func TestLoadMissingFile(t *testing.T) {
	_, err := Load("/nonexistent/config.yaml")
	if err == nil {
		t.Error("expect error for missing file")
	}
}

// TestLoadValidConfig tests loading valid configuration.
func TestLoadValidConfig(t *testing.T) {
	cfg, err := Load(testConfigPath())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg == nil {
		t.Fatal("expect non-nil config")
	}

	// Verify required fields are loaded
	if cfg.Blogger.BlogID == "" {
		t.Error("expect blogger.blogId to be set")
	}
	if cfg.Site.Title == "" {
		t.Error("expect site.title to be set")
	}
}

// TestLoadSetsDefaults tests that Load applies default values.
func TestLoadSetsDefaults(t *testing.T) {
	// Create temp config with minimal required fields
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.yaml")

	content := `
blogger:
  blogId: "test-blog"
server:
  port: 0
site:
  title: ""
form:
  items:
    - type: "text"
      name: "test"
xml:
  root: "root"
  fields:
    - xmlPath: "root/test"
      type: "scalar"
      formField: "test"
`

	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}

	cfg, err := Load(tmpFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check defaults
	if cfg.Blogger.MaxResults != 20 {
		t.Errorf("expect MaxResults default 20, got: %d", cfg.Blogger.MaxResults)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("expect Port default 8080, got: %d", cfg.Server.Port)
	}
	if cfg.Site.Title != "Blogpost → XML Exporter" {
		t.Errorf("expect default title, got: %s", cfg.Site.Title)
	}
}

// TestBloggerConfigValidation tests blogger config validation.
func TestBloggerConfigValidation(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.yaml")

	content := `
blogger:
  blogId: ""
server:
  port: 8080
site:
  title: "Test"
form:
  items: []
xml:
  root: "root"
`

	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}

	_, err := Load(tmpFile)
	if err == nil {
		t.Error("expect error for missing blogId")
	}
}

// TestFormFieldsExtraction tests Form.Fields() returns leaf items only.
func TestFormFieldsExtraction(t *testing.T) {
	form := FormConfig{
		Items: []FormItem{
			{
				Type:  "group",
				Title: "Group 1",
				Items: []FormItem{
					{
						Type: "text",
						Name: "field1",
					},
					{
						Type: "text",
						Name: "field2",
					},
				},
			},
			{
				Type: "text",
				Name: "field3",
			},
		},
	}

	fields := form.Fields()

	if len(fields) != 3 {
		t.Errorf("expect 3 fields, got: %d", len(fields))
	}

	// Verify only leaf fields are returned
	for _, field := range fields {
		if field.Type == "group" {
			t.Error("expect no groups in Fields()")
		}
	}
}

// TestFormFieldsNestedGroups tests Fields() with nested groups.
func TestFormFieldsNestedGroups(t *testing.T) {
	form := FormConfig{
		Items: []FormItem{
			{
				Type:  "group",
				Title: "Outer",
				Items: []FormItem{
					{
						Type:  "group",
						Title: "Inner",
						Items: []FormItem{
							{
								Type: "text",
								Name: "nested_field",
							},
						},
					},
				},
			},
		},
	}

	fields := form.Fields()

	if len(fields) != 1 {
		t.Errorf("expect 1 nested field, got: %d", len(fields))
	}
	if fields[0].Name != "nested_field" {
		t.Error("expect nested_field to be found")
	}
}

// TestFormFieldsEmpty tests Fields() with no items.
func TestFormFieldsEmpty(t *testing.T) {
	form := FormConfig{
		Items: []FormItem{},
	}

	fields := form.Fields()

	if len(fields) != 0 {
		t.Errorf("expect 0 fields, got: %d", len(fields))
	}
}

// TestFormItemStructure tests FormItem field assignments.
func TestFormItemStructure(t *testing.T) {
	item := FormItem{
		Type:     "text",
		Name:     "username",
		Label:    "User Name",
		Source:   "author.name",
		Hidden:   false,
		Template: "{{ .name | upper }}",
	}

	if item.Type != "text" {
		t.Error("expect type to be set")
	}
	if item.Name != "username" {
		t.Error("expect name to be set")
	}
	if item.Source != "author.name" {
		t.Error("expect source to be set")
	}
}

// TestSelectOptionStructure tests SelectOption fields.
func TestSelectOptionStructure(t *testing.T) {
	opt := SelectOption{
		Value: "opt1",
		Label: "Option 1",
	}

	if opt.Value != "opt1" {
		t.Error("expect value to be set")
	}
	if opt.Label != "Option 1" {
		t.Error("expect label to be set")
	}
}

// TestXMLFieldStructure tests XMLField nesting.
func TestXMLFieldStructure(t *testing.T) {
	field := XMLField{
		XMLPath:   "root/child",
		Type:      "scalar",
		FormField: "field1",
		Template:  "{{ .value }}",
		Fields: []XMLField{
			{
				XMLPath: "root/child/grandchild",
				Type:    "scalar",
			},
		},
	}

	if field.XMLPath != "root/child" {
		t.Error("expect XMLPath to be set")
	}
	if len(field.Fields) != 1 {
		t.Error("expect nested fields")
	}
}

// TestXMLConfigStructure tests XMLConfig fields.
func TestXMLConfigStructure(t *testing.T) {
	cfg := XMLConfig{
		Root:     "document",
		Filename: "output_{{ .id }}",
		Namespaces: []XMLAttr{
			{Name: "xmlns", Value: "http://example.com"},
		},
	}

	if cfg.Root != "document" {
		t.Error("expect root to be set")
	}
	if len(cfg.Namespaces) != 1 {
		t.Error("expect namespaces")
	}
}

// TestServerConfigDefaults tests server config defaults.
func TestServerConfigDefaults(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.yaml")

	content := `
blogger:
  blogId: "test"
server:
  port: 0
site:
  title: "Test"
form:
  items:
    - type: "text"
      name: "field"
xml:
  root: "root"
  fields:
    - xmlPath: "root/field"
      type: "scalar"
      formField: "field"
`

	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}

	cfg, err := Load(tmpFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Server.Port != 8080 {
		t.Errorf("expect port 8080 (default), got: %d", cfg.Server.Port)
	}
}

// TestSiteConfigTitle tests site title defaults.
func TestSiteConfigTitle(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.yaml")

	content := `
blogger:
  blogId: "test"
site:
  title: ""
  heading: ""
form:
  items:
    - type: "text"
      name: "field"
xml:
  root: "root"
  fields:
    - xmlPath: "root/field"
      type: "scalar"
      formField: "field"
`

	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}

	cfg, err := Load(tmpFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Site.Title == "" {
		t.Error("expect default title to be set")
	}
	if cfg.Site.Heading != cfg.Site.Title {
		t.Error("expect heading to match title when empty")
	}
}

// TestBloggerConfigMaxResults tests MaxResults defaults.
func TestBloggerConfigMaxResults(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.yaml")

	content := `
blogger:
  blogId: "test"
  maxResults: 0
site:
  title: "Test"
form:
  items:
    - type: "text"
      name: "field"
xml:
  root: "root"
  fields:
    - xmlPath: "root/field"
      type: "scalar"
      formField: "field"
`

	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}

	cfg, err := Load(tmpFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Blogger.MaxResults != 20 {
		t.Errorf("expect MaxResults 20 (default), got: %d", cfg.Blogger.MaxResults)
	}
}

// TestThemeConfigStructure tests ThemeConfig fields.
func TestThemeConfigStructure(t *testing.T) {
	theme := ThemeConfig{
		PrimaryColor: "#FF0000",
		DarkColor:    "#000000",
		LightColor:   "#FFFFFF",
	}

	if theme.PrimaryColor != "#FF0000" {
		t.Error("expect primary color")
	}
	if theme.DarkColor != "#000000" {
		t.Error("expect dark color")
	}
	if theme.LightColor != "#FFFFFF" {
		t.Error("expect light color")
	}
}

// TestAssetsConfigStructure tests AssetsConfig fields.
func TestAssetsConfigStructure(t *testing.T) {
	assets := AssetsConfig{
		Dir:     "/assets",
		Favicon: "favicon.ico",
		Logo:    "logo.png",
	}

	if assets.Dir != "/assets" {
		t.Error("expect assets dir")
	}
	if assets.Favicon != "favicon.ico" {
		t.Error("expect favicon")
	}
	if assets.Logo != "logo.png" {
		t.Error("expect logo")
	}
}

// TestFormItemArrayFields tests array field configuration.
func TestFormItemArrayFields(t *testing.T) {
	item := FormItem{
		Type: "array",
		Name: "items",
		Fields: []FormItem{
			{
				Type:   "text",
				Name:   "item_name",
				Source: "name",
			},
		},
	}

	if item.Type != "array" {
		t.Error("expect array type")
	}
	if len(item.Fields) != 1 {
		t.Error("expect array fields")
	}
}

// TestFormItemSelectOptions tests select field configuration.
func TestFormItemSelectOptions(t *testing.T) {
	item := FormItem{
		Type: "select",
		Name: "status",
		Options: []SelectOption{
			{Value: "active", Label: "Active"},
			{Value: "inactive", Label: "Inactive"},
		},
		AllowCustom: true,
	}

	if item.Type != "select" {
		t.Error("expect select type")
	}
	if len(item.Options) != 2 {
		t.Error("expect select options")
	}
	if !item.AllowCustom {
		t.Error("expect allowCustom to be true")
	}
}

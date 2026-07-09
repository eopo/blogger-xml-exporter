package xmlgen

import (
	"strings"
	"testing"

	"github.com/leokr/blogger-xml-exporter/internal/config"
)

// TestRenderArrayFieldPositive tests successful rendering of array fields with multiple rows.
func TestRenderArrayFieldPositive(t *testing.T) {
	xmlCfg := config.XMLConfig{
		Root: "doc:Document",
		Namespaces: []config.XMLAttr{
			{Name: "doc", Value: "http://example.com/doc"},
		},
		Fields: []config.XMLField{
			{
				XMLPath: "doc:items/doc:item", Type: "array", FormField: "items",
				Fields: []config.XMLField{
					{XMLPath: "doc:name", FormField: "name"},
					{XMLPath: "doc:value", FormField: "value"},
				},
			},
		},
	}

	values := map[string]interface{}{
		"items": []interface{}{
			map[string]interface{}{"name": "Item1", "value": "100"},
			map[string]interface{}{"name": "Item2", "value": "200"},
		},
	}

	out, err := Render(xmlCfg, map[string]interface{}{}, values)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	xmlStr := string(out)

	// Verify array generates two separate elements with same name.
	if got := strings.Count(xmlStr, "<doc:item>"); got != 2 {
		t.Errorf("expect exactly two <doc:item>, got: %d\n%s", got, xmlStr)
	}
	if !strings.Contains(xmlStr, "<doc:name>Item1</doc:name>") {
		t.Errorf("expect Item1 in output:\n%s", xmlStr)
	}
	if !strings.Contains(xmlStr, "<doc:value>200</doc:value>") {
		t.Errorf("expect value 200 in output:\n%s", xmlStr)
	}
}

// TestRenderArrayFieldEmpty tests array field with empty list.
func TestRenderArrayFieldEmpty(t *testing.T) {
	xmlCfg := config.XMLConfig{
		Root: "doc:Document",
		Namespaces: []config.XMLAttr{
			{Name: "doc", Value: "http://example.com/doc"},
		},
		Fields: []config.XMLField{
			{
				XMLPath: "doc:items/doc:item", Type: "array", FormField: "items",
				Fields: []config.XMLField{
					{XMLPath: "doc:name", FormField: "name"},
				},
			},
		},
	}

	values := map[string]interface{}{
		"items": []interface{}{}, // Empty array
	}

	out, err := Render(xmlCfg, map[string]interface{}{}, values)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	xmlStr := string(out)

	// Empty array should not generate any item elements.
	if strings.Contains(xmlStr, "<doc:item>") {
		t.Errorf("expect no <doc:item> for empty array:\n%s", xmlStr)
	}
}

// TestRenderListFieldPositive tests list type which splits comma-separated values.
func TestRenderListFieldPositive(t *testing.T) {
	xmlCfg := config.XMLConfig{
		Root: "doc:Document",
		Namespaces: []config.XMLAttr{
			{Name: "doc", Value: "http://example.com/doc"},
		},
		Fields: []config.XMLField{
			{XMLPath: "doc:tags/doc:tag", Type: "list", FormField: "tags"},
		},
	}

	values := map[string]interface{}{
		"tags": "alpha, beta, gamma",
	}

	out, err := Render(xmlCfg, map[string]interface{}{}, values)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	xmlStr := string(out)

	// List field should generate three separate tag elements.
	if got := strings.Count(xmlStr, "<doc:tag>"); got != 3 {
		t.Errorf("expect exactly three <doc:tag>, got: %d\n%s", got, xmlStr)
	}
	if !strings.Contains(xmlStr, "<doc:tag>alpha</doc:tag>") {
		t.Errorf("expect alpha tag in output:\n%s", xmlStr)
	}
	if !strings.Contains(xmlStr, "<doc:tag>beta</doc:tag>") {
		t.Errorf("expect beta tag in output:\n%s", xmlStr)
	}
}

// TestRenderListFieldEmpty tests list field with empty string.
func TestRenderListFieldEmpty(t *testing.T) {
	xmlCfg := config.XMLConfig{
		Root: "doc:Document",
		Namespaces: []config.XMLAttr{
			{Name: "doc", Value: "http://example.com/doc"},
		},
		Fields: []config.XMLField{
			{XMLPath: "doc:tags/doc:tag", Type: "list", FormField: "tags"},
		},
	}

	values := map[string]interface{}{
		"tags": "",
	}

	out, err := Render(xmlCfg, map[string]interface{}{}, values)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	xmlStr := string(out)

	// Empty list should not generate tag elements.
	if strings.Contains(xmlStr, "<doc:tag>") {
		t.Errorf("expect no <doc:tag> for empty list:\n%s", xmlStr)
	}
}

// TestRenderFormFieldPositive tests simple form field mapping.
func TestRenderFormFieldPositive(t *testing.T) {
	xmlCfg := config.XMLConfig{
		Root: "doc:Document",
		Namespaces: []config.XMLAttr{
			{Name: "doc", Value: "http://example.com/doc"},
		},
		Fields: []config.XMLField{
			{XMLPath: "doc:title", FormField: "title"},
			{XMLPath: "doc:author", FormField: "author"},
		},
	}

	values := map[string]interface{}{
		"title":  "Test Title",
		"author": "John Doe",
	}

	out, err := Render(xmlCfg, map[string]interface{}{}, values)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	xmlStr := string(out)

	if !strings.Contains(xmlStr, "<doc:title>Test Title</doc:title>") {
		t.Errorf("expect title in output:\n%s", xmlStr)
	}
	if !strings.Contains(xmlStr, "<doc:author>John Doe</doc:author>") {
		t.Errorf("expect author in output:\n%s", xmlStr)
	}
}

// TestRenderFormFieldMissing tests form field that is not provided in values.
func TestRenderFormFieldMissing(t *testing.T) {
	xmlCfg := config.XMLConfig{
		Root: "doc:Document",
		Namespaces: []config.XMLAttr{
			{Name: "doc", Value: "http://example.com/doc"},
		},
		Fields: []config.XMLField{
			{XMLPath: "doc:title", FormField: "title"},
		},
	}

	values := map[string]interface{}{} // title not provided

	out, err := Render(xmlCfg, map[string]interface{}{}, values)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	xmlStr := string(out)

	// Missing form field should produce empty element.
	if !strings.Contains(xmlStr, "<doc:title></doc:title>") && !strings.Contains(xmlStr, "<doc:title/>") {
		t.Errorf("expect empty title element:\n%s", xmlStr)
	}
}

// TestRenderTemplateWithSourcePositive tests template evaluation with source function.
func TestRenderTemplateWithSourcePositive(t *testing.T) {
	xmlCfg := config.XMLConfig{
		Root: "doc:Document",
		Namespaces: []config.XMLAttr{
			{Name: "doc", Value: "http://example.com/doc"},
		},
		Fields: []config.XMLField{
			{XMLPath: "doc:publisher", Template: `{{ source "publisher.name" }}`},
		},
	}

	post := map[string]interface{}{
		"publisher": map[string]interface{}{
			"name": "Example Publisher Ltd",
		},
	}

	out, err := Render(xmlCfg, post, map[string]interface{}{})
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	xmlStr := string(out)

	if !strings.Contains(xmlStr, "<doc:publisher>Example Publisher Ltd</doc:publisher>") {
		t.Errorf("expect publisher from source in output:\n%s", xmlStr)
	}
}

// TestRenderTemplateWithMissingSource tests template when source path doesn't exist.
func TestRenderTemplateWithMissingSource(t *testing.T) {
	xmlCfg := config.XMLConfig{
		Root: "doc:Document",
		Namespaces: []config.XMLAttr{
			{Name: "doc", Value: "http://example.com/doc"},
		},
		Fields: []config.XMLField{
			{XMLPath: "doc:publisher", Template: `{{ source "missing.field" }}`},
		},
	}

	post := map[string]interface{}{} // No publisher data

	out, err := Render(xmlCfg, post, map[string]interface{}{})
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	xmlStr := string(out)

	// Missing source should produce empty value.
	if !strings.Contains(xmlStr, "<doc:publisher></doc:publisher>") && !strings.Contains(xmlStr, "<doc:publisher/>") {
		t.Errorf("expect empty publisher element:\n%s", xmlStr)
	}
}

// TestRenderStaticTemplatePositive tests static template values.
func TestRenderStaticTemplatePositive(t *testing.T) {
	xmlCfg := config.XMLConfig{
		Root: "doc:Document",
		Namespaces: []config.XMLAttr{
			{Name: "doc", Value: "http://example.com/doc"},
		},
		Fields: []config.XMLField{
			{XMLPath: "doc:license", Template: "CC-BY-4.0"},
			{XMLPath: "doc:version", Template: "1.0"},
		},
	}

	out, err := Render(xmlCfg, map[string]interface{}{}, map[string]interface{}{})
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	xmlStr := string(out)

	if !strings.Contains(xmlStr, "<doc:license>CC-BY-4.0</doc:license>") {
		t.Errorf("expect static license in output:\n%s", xmlStr)
	}
	if !strings.Contains(xmlStr, "<doc:version>1.0</doc:version>") {
		t.Errorf("expect static version in output:\n%s", xmlStr)
	}
}

// TestRenderNestedElementsPositive tests nested XML elements with shared intermediates.
func TestRenderNestedElementsPositive(t *testing.T) {
	xmlCfg := config.XMLConfig{
		Root: "doc:Document",
		Namespaces: []config.XMLAttr{
			{Name: "doc", Value: "http://example.com/doc"},
		},
		Fields: []config.XMLField{
			{XMLPath: "doc:contact/doc:name", FormField: "name"},
			{XMLPath: "doc:contact/doc:email", FormField: "email"},
		},
	}

	values := map[string]interface{}{
		"name":  "Alice",
		"email": "alice@example.com",
	}

	out, err := Render(xmlCfg, map[string]interface{}{}, values)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	xmlStr := string(out)

	// Shared intermediate element should appear exactly once.
	if got := strings.Count(xmlStr, "<doc:contact>"); got != 1 {
		t.Errorf("expect exactly one <doc:contact>, got: %d\n%s", got, xmlStr)
	}
	if !strings.Contains(xmlStr, "<doc:name>Alice</doc:name>") {
		t.Errorf("expect name in output:\n%s", xmlStr)
	}
	if !strings.Contains(xmlStr, "<doc:email>alice@example.com</doc:email>") {
		t.Errorf("expect email in output:\n%s", xmlStr)
	}
}

// TestRenderNamespaceHandling tests correct namespace prefix application.
func TestRenderNamespaceHandling(t *testing.T) {
	xmlCfg := config.XMLConfig{
		Root: "app:Root",
		Namespaces: []config.XMLAttr{
			{Name: "app", Value: "http://example.com/app"},
			{Name: "meta", Value: "http://example.com/meta"},
		},
		Fields: []config.XMLField{
			{XMLPath: "app:data", FormField: "data"},
			{XMLPath: "meta:info", FormField: "info"},
		},
	}

	values := map[string]interface{}{
		"data": "TestData",
		"info": "TestInfo",
	}

	out, err := Render(xmlCfg, map[string]interface{}{}, values)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	xmlStr := string(out)

	if !strings.Contains(xmlStr, `<app:Root xmlns:app="http://example.com/app" xmlns:meta="http://example.com/meta">`) {
		t.Errorf("expect root element with both namespaces:\n%s", xmlStr)
	}
	if !strings.Contains(xmlStr, "<app:data>TestData</app:data>") {
		t.Errorf("expect app-prefixed data:\n%s", xmlStr)
	}
	if !strings.Contains(xmlStr, "<meta:info>TestInfo</meta:info>") {
		t.Errorf("expect meta-prefixed info:\n%s", xmlStr)
	}
}

// TestRenderXMLDeclaration tests XML declaration at start of document.
func TestRenderXMLDeclaration(t *testing.T) {
	xmlCfg := config.XMLConfig{
		Root: "doc:Document",
		Namespaces: []config.XMLAttr{
			{Name: "doc", Value: "http://example.com/doc"},
		},
		Fields: []config.XMLField{},
	}

	out, err := Render(xmlCfg, map[string]interface{}{}, map[string]interface{}{})
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	xmlStr := string(out)

	if !strings.HasPrefix(xmlStr, `<?xml version="1.0" encoding="UTF-8"?>`) {
		t.Errorf("expect XML declaration at start:\n%s", xmlStr)
	}
}

// TestRenderSpecialCharactersEscaped tests that special XML characters are properly escaped.
func TestRenderSpecialCharactersEscaped(t *testing.T) {
	xmlCfg := config.XMLConfig{
		Root: "doc:Document",
		Namespaces: []config.XMLAttr{
			{Name: "doc", Value: "http://example.com/doc"},
		},
		Fields: []config.XMLField{
			{XMLPath: "doc:text", FormField: "text"},
		},
	}

	values := map[string]interface{}{
		"text": `This & That <test> "quotes"`,
	}

	out, err := Render(xmlCfg, map[string]interface{}{}, values)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	xmlStr := string(out)

	// Special characters should be escaped.
	if !strings.Contains(xmlStr, "&amp;") {
		t.Errorf("expect & to be escaped as &amp;:\n%s", xmlStr)
	}
	if !strings.Contains(xmlStr, "&lt;") {
		t.Errorf("expect < to be escaped as &lt;:\n%s", xmlStr)
	}
	// Quotes can be escaped as &quot; or &#34; - just verify one is present
	hasQuoteEscape := strings.Contains(xmlStr, "&quot;") || strings.Contains(xmlStr, "&#34;")
	if !hasQuoteEscape {
		t.Errorf("expect \" to be escaped:\n%s", xmlStr)
	}
}

// TestRenderComplexNested tests complex nested structure with arrays and templates.
func TestRenderComplexNested(t *testing.T) {
	xmlCfg := config.XMLConfig{
		Root: "doc:Document",
		Namespaces: []config.XMLAttr{
			{Name: "doc", Value: "http://example.com/doc"},
		},
		Fields: []config.XMLField{
			{XMLPath: "doc:metadata/doc:title", FormField: "title"},
			{
				XMLPath: "doc:metadata/doc:tags/doc:tag", Type: "list", FormField: "tags",
			},
			{
				XMLPath: "doc:content/doc:section", Type: "array", FormField: "sections",
				Fields: []config.XMLField{
					{XMLPath: "doc:heading", FormField: "heading"},
					{XMLPath: "doc:body", FormField: "body"},
				},
			},
			{XMLPath: "doc:footer", Template: "Generated 2025"},
		},
	}

	values := map[string]interface{}{
		"title": "Main Title",
		"tags":  "tech, news",
		"sections": []interface{}{
			map[string]interface{}{"heading": "Section 1", "body": "Content 1"},
			map[string]interface{}{"heading": "Section 2", "body": "Content 2"},
		},
	}

	out, err := Render(xmlCfg, map[string]interface{}{}, values)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	xmlStr := string(out)

	// Verify title is in metadata.
	if !strings.Contains(xmlStr, "<doc:metadata>") || !strings.Contains(xmlStr, "<doc:title>Main Title</doc:title>") {
		t.Errorf("expect title in metadata:\n%s", xmlStr)
	}

	// Verify list creates two tag elements.
	if got := strings.Count(xmlStr, "<doc:tag>"); got != 2 {
		t.Errorf("expect exactly two <doc:tag>, got: %d\n%s", got, xmlStr)
	}

	// Verify array creates two section elements.
	if got := strings.Count(xmlStr, "<doc:section>"); got != 2 {
		t.Errorf("expect exactly two <doc:section>, got: %d\n%s", got, xmlStr)
	}

	// Verify footer is static.
	if !strings.Contains(xmlStr, "<doc:footer>Generated 2025</doc:footer>") {
		t.Errorf("expect static footer:\n%s", xmlStr)
	}
}

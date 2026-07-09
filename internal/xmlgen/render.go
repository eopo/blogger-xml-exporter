// Package xmlgen generates XML output from config-defined element trees.
package xmlgen

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/leokr/blogger-xml-exporter/internal/blogger"
	"github.com/leokr/blogger-xml-exporter/internal/config"
)

// node represents an XML tree node.
type node struct {
	name     string
	text     string
	hasText  bool
	children []*node
}

// getOrCreateChild returns an existing child or creates one.
// Shared paths reuse intermediate elements.
func (n *node) getOrCreateChild(name string) *node {
	for _, c := range n.children {
		if c.name == name {
			return c
		}
	}
	return n.addChild(name)
}

// addChild always creates a new child.
// Used for repeated array/list entries.
func (n *node) addChild(name string) *node {
	c := &node{name: name}
	n.children = append(n.children, c)
	return c
}

// Render generates the output XML from config and form values.
func Render(xmlCfg config.XMLConfig, post map[string]interface{}, values map[string]interface{}) ([]byte, error) {
	root := &node{name: xmlCfg.Root}
	for _, field := range xmlCfg.Fields {
		if err := applyXMLField(root, field, values, post); err != nil {
			return nil, fmt.Errorf("xml field %q: %w", field.XMLPath, err)
		}
	}

	var buf bytes.Buffer
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	if err := writeNode(&buf, root, 0, xmlCfg); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// applyXMLField navigates to the target element via XMLPath and applies the field value.
func applyXMLField(parent *node, field config.XMLField, data map[string]interface{}, post map[string]interface{}) error {
	segments := strings.Split(field.XMLPath, "/")
	target := parent
	for _, seg := range segments[:len(segments)-1] {
		target = target.getOrCreateChild(seg)
	}
	leafName := segments[len(segments)-1]

	switch field.Type {
	case "array":
		rows, err := toRows(data[field.FormField])
		if err != nil {
			return fmt.Errorf("formField %q: %w", field.FormField, err)
		}
		for i, row := range rows {
			item := target.addChild(leafName)
			for _, sub := range field.Fields {
				if err := applyXMLField(item, sub, row, post); err != nil {
					return fmt.Errorf("row %d: %w", i, err)
				}
			}
		}
	case "list":
		value, err := resolveXMLValue(field, data, post)
		if err != nil {
			return err
		}
		for _, item := range splitList(value) {
			leaf := target.addChild(leafName)
			leaf.text = item
			leaf.hasText = true
		}
	default:
		value, err := resolveXMLValue(field, data, post)
		if err != nil {
			return err
		}
		leaf := target.getOrCreateChild(leafName)
		leaf.text = value
		leaf.hasText = true
	}
	return nil
}

// resolveXMLValue computes field value: from FormField (direct), from Template (evaluated), or empty.
func resolveXMLValue(field config.XMLField, data map[string]interface{}, post map[string]interface{}) (string, error) {
	if field.Template != "" {
		value, err := blogger.RenderTemplate(post, field.XMLPath, field.Template, data)
		if err != nil {
			return "", fmt.Errorf("template failed: %w", err)
		}
		return value, nil
	}
	if field.FormField != "" {
		return toText(data[field.FormField]), nil
	}
	return "", nil
}

// toRows normalizes array field values to []map[string]interface{}.
func toRows(value interface{}) ([]map[string]interface{}, error) {
	switch v := value.(type) {
	case nil:
		return nil, nil
	case []map[string]interface{}:
		return v, nil
	case []interface{}:
		rows := make([]map[string]interface{}, 0, len(v))
		for _, item := range v {
			row, ok := item.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("expected object in array")
			}
			rows = append(rows, row)
		}
		return rows, nil
	default:
		return nil, fmt.Errorf("expected list of objects, got %T", value)
	}
}

// toText converts field values to string representation.
func toText(value interface{}) string {
	switch v := value.(type) {
	case nil:
		return ""
	case string:
		return v
	default:
		return fmt.Sprint(v)
	}
}

// splitList splits comma-separated input into trimmed, non-empty parts.
func splitList(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

// writeNode writes node and children recursively with indentation.
// Namespaces and attributes are written at root level only.
func writeNode(buf *bytes.Buffer, n *node, depth int, xmlCfg config.XMLConfig) error {
	indent := strings.Repeat("\t", depth)
	buf.WriteString(indent)
	fmt.Fprintf(buf, "<%s", n.name)

	if depth == 0 {
		for _, ns := range xmlCfg.Namespaces {
			escaped, err := escapeXML(ns.Value)
			if err != nil {
				return err
			}
			fmt.Fprintf(buf, ` xmlns:%s="%s"`, ns.Name, escaped)
		}
		for _, attr := range xmlCfg.Attributes {
			escaped, err := escapeXML(attr.Value)
			if err != nil {
				return err
			}
			fmt.Fprintf(buf, ` %s="%s"`, attr.Name, escaped)
		}
	}
	buf.WriteString(">")

	if len(n.children) > 0 {
		buf.WriteString("\n")
		for _, c := range n.children {
			if err := writeNode(buf, c, depth+1, xmlCfg); err != nil {
				return err
			}
		}
		buf.WriteString(indent)
	} else if n.hasText {
		escaped, err := escapeXML(n.text)
		if err != nil {
			return err
		}
		buf.WriteString(escaped)
	}

	fmt.Fprintf(buf, "</%s>\n", n.name)
	return nil
}

// escapeXML escapes special characters for safe XML output.
func escapeXML(s string) (string, error) {
	var buf bytes.Buffer
	if err := xml.EscapeText(&buf, []byte(s)); err != nil {
		return "", err
	}
	return buf.String(), nil
}

package blogger

import (
	"fmt"
	"html"
	"log"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"

	"github.com/eopo/blogger-xml-exporter/internal/config"
)

var tagPattern = regexp.MustCompile(`<[^>]*>`)
var unicodeEscapePattern = regexp.MustCompile(`\\u([0-9a-fA-F]{4})`)

// BaseFieldFuncMap contains template helper functions used in form and XML templates.
// Includes Sprig v3 plus stripHTML and unescapeUnicode.
// Note: "source" is added separately by BuildTemplateFuncs.
var BaseFieldFuncMap = buildBaseFieldFuncMap()

func buildBaseFieldFuncMap() template.FuncMap {
	funcs := sprig.TxtFuncMap()
	funcs["stripHTML"] = stripHTML
	funcs["unescapeUnicode"] = unescapeUnicode
	return funcs
}

// stripHTML removes HTML tags and unescapes entities.
func stripHTML(s string) string {
	return html.UnescapeString(tagPattern.ReplaceAllString(s, ""))
}

// unescapeUnicode resolves \uXXXX escape sequences to Unicode characters.
func unescapeUnicode(s string) string {
	return unicodeEscapePattern.ReplaceAllStringFunc(s, func(match string) string {
		codepoint, err := strconv.ParseInt(match[2:], 16, 32)
		if err != nil {
			return match
		}
		return string(rune(codepoint))
	})
}

// ResolveField navigates a dot-path (e.g., "author.displayName") in post and returns the value.
// Returns nil if post is nil or path does not exist.
func ResolveField(post map[string]interface{}, dotPath string) interface{} {
	if post == nil || dotPath == "" {
		return nil
	}
	parts := strings.Split(dotPath, ".")
	var current interface{} = post
	for _, part := range parts {
		if part == "" {
			continue
		}
		m, ok := current.(map[string]interface{})
		if !ok {
			return nil
		}
		current = m[part]
	}
	return current
}

// ResolveFields computes pre-filled values for all form fields.
// Uses Source (direct post access) or Template (computed).
// Field errors are logged but do not halt processing.
func ResolveFields(post map[string]interface{}, fields []config.FormItem) map[string]interface{} {
	values := make(map[string]interface{}, len(fields))
	for _, field := range fields {
		values[field.Name] = resolveField(post, field, values)
		if field.Type == "select" && field.OptionsSource != "" {
			values[field.Name+"__options"] = resolveSelectOptions(post, field.OptionsSource)
		}
	}
	return values
}

// resolveSelectOptions reads dynamic options for select fields from post via dot path.
func resolveSelectOptions(post map[string]interface{}, dotPath string) []string {
	items, ok := ResolveField(post, dotPath).([]interface{})
	if !ok {
		return nil
	}
	opts := make([]string, 0, len(items))
	for _, item := range items {
		opts = append(opts, fmt.Sprint(item))
	}
	return opts
}

// resolveField computes the value of a single field.
func resolveField(post map[string]interface{}, field config.FormItem, priorValues map[string]interface{}) interface{} {
	if field.Type == "array" {
		return resolveArrayField(post, field)
	}
	if field.Template != "" {
		v, err := renderFieldTemplate(post, field.Name, field.Template, priorValues)
		if err != nil {
			log.Printf("field %q: %v (using empty value)", field.Name, err)
			return ""
		}
		return v
	}
	if field.Source != "" {
		return toDisplayValue(ResolveField(post, field.Source))
	}
	return ""
}

// resolveArrayField resolves array field rows.
// If Source is set, reads items from post. Otherwise, one default row from post if fields have defaults.
func resolveArrayField(post map[string]interface{}, field config.FormItem) []map[string]interface{} {
	if field.Source != "" {
		items, ok := ResolveField(post, field.Source).([]interface{})
		if !ok {
			return []map[string]interface{}{}
		}

		rows := make([]map[string]interface{}, 0, len(items))
		for _, item := range items {
			itemPost, _ := item.(map[string]interface{})
			rows = append(rows, resolveArrayRow(itemPost, field.Fields))
		}
		return rows
	}

	if !fieldsHaveDefault(field.Fields) {
		return []map[string]interface{}{}
	}
	return []map[string]interface{}{resolveArrayRow(post, field.Fields)}
}

// fieldsHaveDefault reports whether at least one subfield has Source or Template.
func fieldsHaveDefault(fields []config.FormItem) bool {
	for _, sub := range fields {
		if sub.Source != "" || sub.Template != "" {
			return true
		}
	}
	return false
}

// resolveArrayRow computes subfield values for a single array row.
func resolveArrayRow(itemPost map[string]interface{}, fields []config.FormItem) map[string]interface{} {
	row := make(map[string]interface{}, len(fields))
	for _, sub := range fields {
		row[sub.Name] = resolveField(itemPost, sub, row)
	}
	return row
}

// renderFieldTemplate evaluates a field's template.
func renderFieldTemplate(post map[string]interface{}, fieldName, tmplText string, values map[string]interface{}) (string, error) {
	return RenderTemplate(post, fieldName, tmplText, values)
}

// BuildTemplateFuncs returns the complete template function map including Sprig, stripHTML, unescapeUnicode, and "source".
// The "source" function accesses arbitrary dot-paths in post; returns empty string if post is nil.
func BuildTemplateFuncs(post map[string]interface{}) template.FuncMap {
	funcs := make(template.FuncMap, len(BaseFieldFuncMap)+1)
	for name, fn := range BaseFieldFuncMap {
		funcs[name] = fn
	}
	funcs["source"] = func(dotPath string) string {
		return ToDisplayValue(ResolveField(post, dotPath))
	}
	return funcs
}

// RenderTemplate parses and executes a Go template using BuildTemplateFuncs.
func RenderTemplate(post map[string]interface{}, templateName, tmplText string, data interface{}) (string, error) {
	tmpl, err := template.New(templateName).Funcs(BuildTemplateFuncs(post)).Parse(tmplText)
	if err != nil {
		return "", fmt.Errorf("template is invalid: %w", err)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("template execution failed: %w", err)
	}
	return buf.String(), nil
}

func toDisplayValue(v interface{}) string {
	return ToDisplayValue(v)
}

// ToDisplayValue converts a post value to display text: strings direct, lists comma-separated.
func ToDisplayValue(v interface{}) string {
	switch val := v.(type) {
	case nil:
		return ""
	case string:
		return val
	case []interface{}:
		parts := make([]string, 0, len(val))
		for _, item := range val {
			parts = append(parts, fmt.Sprint(item))
		}
		return strings.Join(parts, ", ")
	default:
		return fmt.Sprint(val)
	}
}

// ResolvedPreset is a computed preset with label and resolved field values.
type ResolvedPreset struct {
	Label  string                 `json:"label"`
	Values map[string]interface{} `json:"values"`
}

// ResolvePresets computes applicable presets for all form groups.
// Returns a map of group title -> presets.
func ResolvePresets(post map[string]interface{}, items []config.FormItem, values map[string]interface{}) map[string][]ResolvedPreset {
	out := make(map[string][]ResolvedPreset)
	var walk func(items []config.FormItem)
	walk = func(items []config.FormItem) {
		for _, item := range items {
			if item.Type != "group" {
				continue
			}
			if len(item.Presets) > 0 {
				out[item.Title] = resolveGroupPresets(post, item, values)
			}
			walk(item.Items)
		}
	}
	walk(items)
	return out
}

// resolveGroupPresets computes presets for one group.
// Without PresetsSource, evaluates each preset template once.
// With PresetsSource, uses the single preset as template for each item.
func resolveGroupPresets(post map[string]interface{}, group config.FormItem, values map[string]interface{}) []ResolvedPreset {
	if group.PresetsSource == "" {
		out := make([]ResolvedPreset, 0, len(group.Presets))
		for _, preset := range group.Presets {
			out = append(out, ResolvedPreset{
				Label:  preset.Label,
				Values: renderPresetValues(post, preset.Values, values),
			})
		}
		return out
	}

	items, ok := ResolveField(post, group.PresetsSource).([]interface{})
	if !ok || len(group.Presets) == 0 {
		return nil
	}
	presetTemplate := group.Presets[0]
	out := make([]ResolvedPreset, 0, len(items))
	for _, item := range items {
		label, err := RenderTemplate(post, "preset-label", presetTemplate.Label, item)
		if err != nil {
			log.Printf("preset label: %v (skip item)", err)
			continue
		}
		out = append(out, ResolvedPreset{
			Label:  label,
			Values: renderPresetValues(post, presetTemplate.Values, item),
		})
	}
	return out
}

// renderPresetValues evaluates field templates of a preset.
func renderPresetValues(post map[string]interface{}, tmpl map[string]string, data interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(tmpl))
	for name, tmplText := range tmpl {
		value, err := RenderTemplate(post, "preset-"+name, tmplText, data)
		if err != nil {
			log.Printf("preset field %q: %v (using empty value)", name, err)
			value = ""
		}
		out[name] = value
	}
	return out
}

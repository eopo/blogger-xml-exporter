// Package config loads and validates application configuration from a YAML file.
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// FormItem represents a node in the form tree.
// Type discriminates the node kind: "group", "array", or scalar field types.
type FormItem struct {
	Type string `yaml:"type" json:"type"`

	// Group-only fields
	Title         string     `yaml:"title,omitempty" json:"title,omitempty"`
	Collapsible   bool       `yaml:"collapsible,omitempty" json:"collapsible,omitempty"`
	Collapsed     bool       `yaml:"collapsed,omitempty" json:"collapsed,omitempty"`
	Items         []FormItem `yaml:"items,omitempty" json:"items,omitempty"`
	Presets       []Preset   `yaml:"presets,omitempty" json:"presets,omitempty"`
	PresetsSource string     `yaml:"presetsSource,omitempty" json:"presetsSource,omitempty"`

	// Form data fields
	Name     string `yaml:"name,omitempty" json:"name,omitempty"`
	Label    string `yaml:"label,omitempty" json:"label,omitempty"`
	Hidden   bool   `yaml:"hidden,omitempty" json:"hidden,omitempty"`
	Source   string `yaml:"source,omitempty" json:"source,omitempty"`
	Template string `yaml:"template,omitempty" json:"template,omitempty"`

	// Array fields
	Fields []FormItem `yaml:"fields,omitempty" json:"fields,omitempty"`

	// Select fields
	Options       []SelectOption `yaml:"options,omitempty" json:"options,omitempty"`
	OptionsSource string         `yaml:"optionsSource,omitempty" json:"optionsSource,omitempty"`
	AllowCustom   bool           `yaml:"allowCustom,omitempty" json:"allowCustom,omitempty"`

	// Date field
	IncludeTime bool `yaml:"includeTime,omitempty" json:"includeTime,omitempty"`

	// Layout (within parent group)
	Row   int `yaml:"row,omitempty" json:"row,omitempty"`
	Width int `yaml:"width,omitempty" json:"width,omitempty"`
}

// Preset is a quick-fill template for a form group.
type Preset struct {
	Label  string            `yaml:"label" json:"label"`
	Values map[string]string `yaml:"values" json:"values"`
}

// SelectOption is a selectable value in a "select" field.
type SelectOption struct {
	Value string `yaml:"value" json:"value"`
	Label string `yaml:"label" json:"label"`
}

// XMLField describes an output XML element, independent from FormItem.
type XMLField struct {
	XMLPath   string     `yaml:"xmlPath"`
	Type      string     `yaml:"type"` // "" (scalar), "list", or "array"
	FormField string     `yaml:"formField"`
	Template  string     `yaml:"template"`
	Fields    []XMLField `yaml:"fields"`
}

// XMLAttr is a key-value pair for namespaces or root attributes.
type XMLAttr struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// XMLConfig describes the output XML structure.
type XMLConfig struct {
	Root       string     `yaml:"root"`
	Namespaces []XMLAttr  `yaml:"namespaces"`
	Attributes []XMLAttr  `yaml:"attributes"`
	Fields     []XMLField `yaml:"fields"`
	Filename   string     `yaml:"filename"`
}

// BloggerConfig contains Blogger API settings.
type BloggerConfig struct {
	BlogID     string `yaml:"blogId"`
	MaxResults int    `yaml:"maxResults"`
}

// ServerConfig contains HTTP server settings.
type ServerConfig struct {
	Port int `yaml:"port"`
}

// SiteConfig controls UI labels.
type SiteConfig struct {
	Title   string `yaml:"title"`
	Heading string `yaml:"heading"`
}

// AssetsConfig describes optional assets directory (favicon/logo).
type AssetsConfig struct {
	Dir     string `yaml:"dir"`
	Favicon string `yaml:"favicon"`
	Logo    string `yaml:"logo"`
}

// ThemeConfig defines application theme colors.
type ThemeConfig struct {
	PrimaryColor string `yaml:"primaryColor,omitempty"`
	DarkColor    string `yaml:"darkColor,omitempty"`
	LightColor   string `yaml:"lightColor,omitempty"`
}

// FormConfig contains the form tree.
type FormConfig struct {
	Items []FormItem `yaml:"items"`
}

// Fields returns all leaf items (non-groups) in document order.
func (fc FormConfig) Fields() []FormItem {
	var out []FormItem
	var walk func(items []FormItem)
	walk = func(items []FormItem) {
		for _, item := range items {
			if item.Type == "group" {
				walk(item.Items)
			} else {
				out = append(out, item)
			}
		}
	}
	walk(fc.Items)
	return out
}

// Config is the root configuration structure.
type Config struct {
	Blogger BloggerConfig `yaml:"blogger"`
	Server  ServerConfig  `yaml:"server"`
	Site    SiteConfig    `yaml:"site"`
	Theme   ThemeConfig   `yaml:"theme,omitempty"`
	Assets  AssetsConfig  `yaml:"assets"`
	Form    FormConfig    `yaml:"form"`
	XML     XMLConfig     `yaml:"xml"`
}

// Load reads and validates the configuration file.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if cfg.Blogger.BlogID == "" {
		return nil, fmt.Errorf("config: blogger.blogId is required")
	}
	if cfg.Blogger.MaxResults <= 0 {
		cfg.Blogger.MaxResults = 20
	}
	if cfg.Server.Port <= 0 {
		cfg.Server.Port = 8080
	}
	if cfg.Site.Title == "" {
		cfg.Site.Title = "Blogpost → XML Exporter"
	}
	if cfg.Site.Heading == "" {
		cfg.Site.Heading = cfg.Site.Title
	}
	if cfg.Theme.PrimaryColor == "" {
		cfg.Theme.PrimaryColor = "#2563eb"
	}
	if cfg.Theme.DarkColor == "" {
		cfg.Theme.DarkColor = "#1e40af"
	}
	if cfg.Theme.LightColor == "" {
		cfg.Theme.LightColor = "#3b82f6"
	}
	if cfg.XML.Root == "" {
		return nil, fmt.Errorf("config: xml.root is required")
	}
	if cfg.XML.Filename == "" {
		cfg.XML.Filename = "post"
	}
	if len(cfg.Form.Items) == 0 {
		return nil, fmt.Errorf("config: form.items must not be empty")
	}
	if err := validateFormItems(cfg.Form.Items); err != nil {
		return nil, err
	}
	if len(cfg.XML.Fields) == 0 {
		return nil, fmt.Errorf("config: xml.fields must not be empty")
	}
	formFields := cfg.Form.Fields()
	formFieldNames := make(map[string]bool, len(formFields))
	for _, f := range formFields {
		formFieldNames[f.Name] = true
	}
	if err := validateXMLFields(cfg.XML.Fields, formFieldNames); err != nil {
		return nil, err
	}

	return &cfg, nil
}

var validItemTypes = map[string]bool{
	"group": true, "array": true,
	"text": true, "textarea": true, "date": true, "list": true, "select": true,
}

func validateFormItems(items []FormItem) error {
	for i := range items {
		item := items[i]
		if !validItemTypes[item.Type] {
			return fmt.Errorf("config: form item has invalid type %q", item.Type)
		}
		if item.Row < 0 {
			return fmt.Errorf("config: item %q has invalid row %d", itemLabel(item), item.Row)
		}
		if item.Width < 0 || item.Width > 12 {
			return fmt.Errorf("config: item %q has invalid width %d (0-12)", itemLabel(item), item.Width)
		}

		if item.Type == "group" {
			if len(item.Items) == 0 {
				return fmt.Errorf("config: group %q requires items", item.Title)
			}
			if err := validateFormItems(item.Items); err != nil {
				return err
			}
			if err := validatePresets(item); err != nil {
				return err
			}
			continue
		}

		if err := validateFormField(item); err != nil {
			return err
		}
	}
	return nil
}

func itemLabel(item FormItem) string {
	if item.Type == "group" {
		return item.Title
	}
	return item.Name
}

func validateFormField(field FormItem) error {
	if field.Name == "" {
		return fmt.Errorf("config: form item of type %q requires name", field.Type)
	}
	if field.Source != "" && field.Template != "" {
		return fmt.Errorf("config: field %q cannot have both source and template", field.Name)
	}
	if field.Type == "array" {
		if len(field.Fields) == 0 {
			return fmt.Errorf("config: array field %q requires fields", field.Name)
		}
		if err := validateFormItems(field.Fields); err != nil {
			return err
		}
	} else if len(field.Fields) > 0 {
		return fmt.Errorf("config: non-array field %q cannot define fields", field.Name)
	}
	return nil
}

func validatePresets(group FormItem) error {
	if len(group.Presets) == 0 {
		return nil
	}
	directNames := make(map[string]bool)
	for _, child := range group.Items {
		if child.Type != "group" && child.Name != "" {
			directNames[child.Name] = true
		}
	}
	for _, preset := range group.Presets {
		if preset.Label == "" {
			return fmt.Errorf("config: group %q has preset without label", group.Title)
		}
		if len(preset.Values) == 0 {
			return fmt.Errorf("config: preset %q in group %q requires values", preset.Label, group.Title)
		}
		if group.PresetsSource == "" {
			for name := range preset.Values {
				if !directNames[name] {
					return fmt.Errorf("config: preset %q references unknown field %q", preset.Label, name)
				}
			}
		}
	}
	if group.PresetsSource != "" && len(group.Presets) != 1 {
		return fmt.Errorf("config: group %q: presetsSource requires exactly one preset template", group.Title)
	}
	return nil
}

func validateXMLFields(fields []XMLField, formFieldNames map[string]bool) error {
	for _, field := range fields {
		if field.XMLPath == "" {
			return fmt.Errorf("config: xml field requires xmlPath")
		}
		if field.FormField != "" && field.Template != "" {
			return fmt.Errorf("config: xml field %q cannot have both formField and template", field.XMLPath)
		}
		if field.Type == "array" {
			if field.FormField == "" {
				return fmt.Errorf("config: array xml field %q requires formField", field.XMLPath)
			}
			if !formFieldNames[field.FormField] {
				return fmt.Errorf("config: xml field %q references unknown form field %q", field.XMLPath, field.FormField)
			}
			if len(field.Fields) == 0 {
				return fmt.Errorf("config: array xml field %q requires fields", field.XMLPath)
			}
			if err := validateXMLFields(field.Fields, nil); err != nil {
				return err
			}
			continue
		}
		if len(field.Fields) > 0 {
			return fmt.Errorf("config: non-array xml field %q cannot define fields", field.XMLPath)
		}
		if field.FormField != "" && formFieldNames != nil && !formFieldNames[field.FormField] {
			return fmt.Errorf("config: xml field %q references unknown form field %q", field.XMLPath, field.FormField)
		}
	}
	return nil
}

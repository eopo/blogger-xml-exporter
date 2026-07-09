package config

import "testing"

// TestLoadExampleConfig verifies that config.yaml is valid at all times.
func TestLoadExampleConfig(t *testing.T) {
	cfg, err := Load("../../config.yaml")
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

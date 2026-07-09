package xmlgen

import (
	"strings"
	"testing"

	"github.com/leokr/blogger-xml-exporter/internal/config"
)

func TestRenderArrayHiddenAndList(t *testing.T) {
	xmlCfg := config.XMLConfig{
		Root: "mat:Material",
		Namespaces: []config.XMLAttr{
			{Name: "mat", Value: "http://www.mdr.de/material"},
		},
		Fields: []config.XMLField{
			{XMLPath: "mat:titel", FormField: "titel"},
			{XMLPath: "mat:lieferant/mat:ort", Template: "Essen"},
			{XMLPath: "mat:lieferant/mat:plz", Template: "45359"},
			{
				XMLPath: "mat:zusatzinformationen/mat:zusatzinformation", Type: "array", FormField: "zusatzinformationen",
				Fields: []config.XMLField{
					{XMLPath: "mat:Typ", FormField: "Typ"},
					{XMLPath: "mat:Wert", FormField: "Wert"},
				},
			},
			{XMLPath: "mat:labels/mat:label", Type: "list", FormField: "labels"},
			{XMLPath: "mat:autor", Template: `{{ source "author.displayName" }}`},
		},
	}

	post := map[string]interface{}{
		"author": map[string]interface{}{"displayName": "Max Mustermann"},
	}

	values := map[string]interface{}{
		"titel": "Beispiel-Titel",
		"zusatzinformationen": []interface{}{
			map[string]interface{}{"Typ": "StoryID", "Wert": "12345"},
			map[string]interface{}{"Typ": "MediaID", "Wert": "567"},
		},
		"labels": "Politik, Wirtschaft",
	}

	out, err := Render(xmlCfg, post, values)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	xmlStr := string(out)

	// Shared intermediate element: mat:lieferant should appear exactly once.
	if got := strings.Count(xmlStr, "<mat:lieferant>"); got != 1 {
		t.Errorf("expect exactly one <mat:lieferant>, got: %d\n%s", got, xmlStr)
	}
	if !strings.Contains(xmlStr, "<mat:ort>Essen</mat:ort>") {
		t.Errorf("expect <mat:ort>Essen</mat:ort> in output:\n%s", xmlStr)
	}

	// Array: two separate elements with same name.
	if got := strings.Count(xmlStr, "<mat:zusatzinformation>"); got != 2 {
		t.Errorf("expect exactly two <mat:zusatzinformation>, got: %d\n%s", got, xmlStr)
	}
	if !strings.Contains(xmlStr, "<mat:Typ>StoryID</mat:Typ>") || !strings.Contains(xmlStr, "<mat:Wert>12345</mat:Wert>") {
		t.Errorf("expect StoryID/12345 in output:\n%s", xmlStr)
	}

	// list: comma-separated values as repeated elements.
	if got := strings.Count(xmlStr, "<mat:label>"); got != 2 {
		t.Errorf("expect exactly two <mat:label>, got: %d\n%s", got, xmlStr)
	}
	if !strings.Contains(xmlStr, "<mat:label>Politik</mat:label>") || !strings.Contains(xmlStr, "<mat:label>Wirtschaft</mat:label>") {
		t.Errorf("expect Politik/Wirtschaft in output:\n%s", xmlStr)
	}

	// Template with direct post access via "source".
	if !strings.Contains(xmlStr, "<mat:autor>Max Mustermann</mat:autor>") {
		t.Errorf("expect <mat:autor>Max Mustermann</mat:autor> in output:\n%s", xmlStr)
	}

	if !strings.HasPrefix(xmlStr, `<?xml version="1.0" encoding="UTF-8"?>`) {
		t.Errorf("expect XML declaration at start:\n%s", xmlStr)
	}
	if !strings.Contains(xmlStr, `<mat:Material xmlns:mat="http://www.mdr.de/material">`) {
		t.Errorf("expect root element with namespace:\n%s", xmlStr)
	}
}

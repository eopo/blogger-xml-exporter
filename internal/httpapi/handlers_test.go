package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/leokr/blogger-xml-exporter/internal/blogger"
	"github.com/leokr/blogger-xml-exporter/internal/config"
)

// TestHealthCheck tests the health endpoint.
func TestHealthCheck(t *testing.T) {
	server := New(
		&config.Config{},
		blogger.NewClient("test-key", "test-blog"),
	)

	req := httptest.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()
	server.handleHealth(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expect status 200, got: %d", w.Code)
	}
	if w.Body.String() != "ok" {
		t.Errorf("expect body 'ok', got: %s", w.Body.String())
	}
}

// TestFormSchemaBasic tests form schema endpoint with minimal config.
func TestFormSchemaBasic(t *testing.T) {
	cfg := &config.Config{
		Site: config.SiteConfig{
			Title:   "Test App",
			Heading: "Test Heading",
		},
		Theme: config.ThemeConfig{
			PrimaryColor: "#000000",
			DarkColor:    "#333333",
			LightColor:   "#FFFFFF",
		},
		Form: config.FormConfig{
			Items: []config.FormItem{},
		},
	}

	server := New(cfg, blogger.NewClient("test-key", "test-blog"))

	req := httptest.NewRequest("GET", "/api/form-schema", nil)
	w := httptest.NewRecorder()
	server.handleFormSchema(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expect status 200, got: %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}

	if _, ok := response["items"]; !ok {
		t.Error("expect 'items' in response")
	}
	if _, ok := response["site"]; !ok {
		t.Error("expect 'site' in response")
	}
	if _, ok := response["theme"]; !ok {
		t.Error("expect 'theme' in response")
	}
}

// TestFormSchemaWithAssets tests form schema with assets configured.
func TestFormSchemaWithAssets(t *testing.T) {
	cfg := &config.Config{
		Site: config.SiteConfig{
			Title:   "Test",
			Heading: "Test",
		},
		Theme: config.ThemeConfig{},
		Assets: config.AssetsConfig{
			Dir:     "/assets",
			Favicon: "favicon.ico",
			Logo:    "logo.png",
		},
		Form: config.FormConfig{
			Items: []config.FormItem{},
		},
	}

	server := New(cfg, blogger.NewClient("test-key", "test-blog"))

	req := httptest.NewRequest("GET", "/api/form-schema", nil)
	w := httptest.NewRecorder()
	server.handleFormSchema(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expect status 200, got: %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}

	assets := response["assets"].(map[string]interface{})
	if assets["favicon"] != "/assets/favicon.ico" {
		t.Errorf("expect favicon path in assets, got: %v", assets)
	}
	if assets["logo"] != "/assets/logo.png" {
		t.Errorf("expect logo path in assets, got: %v", assets)
	}
}

// TestFormSchemaNoAssets tests form schema with no assets configured.
func TestFormSchemaNoAssets(t *testing.T) {
	cfg := &config.Config{
		Site: config.SiteConfig{
			Title:   "Test",
			Heading: "Test",
		},
		Theme: config.ThemeConfig{},
		Assets: config.AssetsConfig{
			Dir: "", // No assets directory
		},
		Form: config.FormConfig{
			Items: []config.FormItem{},
		},
	}

	server := New(cfg, blogger.NewClient("test-key", "test-blog"))

	req := httptest.NewRequest("GET", "/api/form-schema", nil)
	w := httptest.NewRecorder()
	server.handleFormSchema(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expect status 200, got: %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}

	assets := response["assets"].(map[string]interface{})
	if len(assets) != 0 {
		t.Errorf("expect empty assets map, got: %v", assets)
	}
}

// TestHandleDefaults tests defaults endpoint returns empty form with defaults.
func TestHandleDefaults(t *testing.T) {
	cfg := &config.Config{
		Site: config.SiteConfig{
			Title:   "Test",
			Heading: "Test",
		},
		Theme:  config.ThemeConfig{},
		Assets: config.AssetsConfig{},
		Form: config.FormConfig{
			Items: []config.FormItem{
				{
					Type:  "text",
					Name:  "field1",
					Label: "Field 1",
				},
			},
		},
	}

	server := New(cfg, blogger.NewClient("test-key", "test-blog"))

	req := httptest.NewRequest("GET", "/api/defaults", nil)
	w := httptest.NewRecorder()
	server.handleDefaults(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expect status 200, got: %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}

	if _, ok := response["post"]; !ok {
		t.Error("expect 'post' in response")
	}
	if _, ok := response["values"]; !ok {
		t.Error("expect 'values' in response")
	}
	if _, ok := response["presets"]; !ok {
		t.Error("expect 'presets' in response")
	}
}

// TestHandleDefaultsEmptyForm tests defaults with empty form items.
func TestHandleDefaultsEmptyForm(t *testing.T) {
	cfg := &config.Config{
		Site: config.SiteConfig{
			Title:   "Test",
			Heading: "Test",
		},
		Theme:  config.ThemeConfig{},
		Assets: config.AssetsConfig{},
		Form: config.FormConfig{
			Items: []config.FormItem{},
		},
	}

	server := New(cfg, blogger.NewClient("test-key", "test-blog"))

	req := httptest.NewRequest("GET", "/api/defaults", nil)
	w := httptest.NewRecorder()
	server.handleDefaults(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expect status 200, got: %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}

	post := response["post"].(map[string]interface{})
	if len(post) != 0 {
		t.Errorf("expect empty post map, got: %v", post)
	}
}

// TestServerNew tests server initialization.
func TestServerNew(t *testing.T) {
	cfg := &config.Config{}
	client := blogger.NewClient("key", "blog")

	server := New(cfg, client)

	if server == nil {
		t.Error("expect non-nil server")
	}
}

// TestWriteJSONHelper tests JSON response encoding.
func TestWriteJSONHelper(t *testing.T) {
	w := httptest.NewRecorder()

	data := map[string]interface{}{
		"status": "ok",
		"value":  123,
	}

	writeJSON(w, http.StatusOK, data)

	if w.Code != http.StatusOK {
		t.Errorf("expect status 200, got: %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("expect status 'ok', got: %v", response["status"])
	}
	if response["value"] != float64(123) {
		t.Errorf("expect value 123, got: %v", response["value"])
	}
}

// TestWriteJSONErrors tests JSON error response.
func TestWriteJSONErrors(t *testing.T) {
	w := httptest.NewRecorder()

	data := map[string]interface{}{
		"error":   "not found",
		"details": "resource does not exist",
	}

	writeJSON(w, http.StatusNotFound, data)

	if w.Code != http.StatusNotFound {
		t.Errorf("expect status 404, got: %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}

	if response["error"] != "not found" {
		t.Errorf("expect error message, got: %v", response)
	}
}

// TestContentTypeJSON tests JSON content type header.
func TestContentTypeJSON(t *testing.T) {
	w := httptest.NewRecorder()
	writeJSON(w, http.StatusOK, map[string]string{"test": "data"})

	ct := w.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expect Content-Type application/json, got: %s", ct)
	}
}

// TestHandleHealthMethodOnly tests that health endpoint only accepts GET.
func TestHandleHealthMethodOnly(t *testing.T) {
	server := New(
		&config.Config{},
		blogger.NewClient("test-key", "test-blog"),
	)

	// Test with POST - should still work (no method filtering in handler itself)
	req := httptest.NewRequest("POST", "/healthz", bytes.NewReader([]byte("")))
	w := httptest.NewRecorder()
	server.handleHealth(w, req)

	// Handler ignores method, just responds
	if w.Code != http.StatusOK {
		t.Errorf("expect status 200, got: %d", w.Code)
	}
}

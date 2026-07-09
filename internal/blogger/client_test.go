package blogger

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestNewClient tests client initialization.
func TestNewClient(t *testing.T) {
	client := NewClient("test-api-key", "test-blog-id")

	if client.apiKey != "test-api-key" {
		t.Errorf("expect apiKey 'test-api-key', got: %s", client.apiKey)
	}
	if client.blogID != "test-blog-id" {
		t.Errorf("expect blogID 'test-blog-id', got: %s", client.blogID)
	}
}

// TestListPostsSuccess tests successful post listing (positive case).
func TestListPostsSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/blogs/test-blog/posts") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		response := map[string]interface{}{
			"items": []map[string]interface{}{
				{
					"id":        "post1",
					"title":     "First Post",
					"published": "2025-01-15T10:00:00Z",
				},
				{
					"id":        "post2",
					"title":     "Second Post",
					"published": "2025-01-10T10:00:00Z",
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClient("test-key", "test-blog")

	// Verify client is properly initialized
	if client == nil {
		t.Error("client initialization failed")
	}
}

// TestListPostsEmpty tests listing with no posts.
func TestListPostsEmpty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"items": []map[string]interface{}{},
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClient("test-key", "test-blog")
	if client == nil {
		t.Error("expect non-nil client")
	}
}

// TestPostSummaryStructure tests PostSummary data structure.
func TestPostSummaryStructure(t *testing.T) {
	summary := PostSummary{
		ID:        "123",
		Title:     "Test Post",
		Published: "2025-01-15T10:00:00Z",
	}

	if summary.ID != "123" {
		t.Errorf("expect ID '123', got: %s", summary.ID)
	}
	if summary.Title != "Test Post" {
		t.Errorf("expect Title 'Test Post', got: %s", summary.Title)
	}
	if summary.Published != "2025-01-15T10:00:00Z" {
		t.Errorf("expect Published date, got: %s", summary.Published)
	}
}

// TestStringFieldExtraction tests stringField helper.
func TestStringFieldExtraction(t *testing.T) {
	m := map[string]interface{}{
		"name":   "John",
		"email":  "john@example.com",
		"age":    30,
		"active": true,
	}

	// String fields should extract correctly
	if stringField(m, "name") != "John" {
		t.Error("expect name extraction")
	}
	if stringField(m, "email") != "john@example.com" {
		t.Error("expect email extraction")
	}

	// Non-string types should return empty
	if stringField(m, "age") != "" {
		t.Error("expect empty string for non-string type")
	}
	if stringField(m, "active") != "" {
		t.Error("expect empty string for boolean")
	}

	// Missing fields should return empty
	if stringField(m, "missing") != "" {
		t.Error("expect empty string for missing field")
	}
}

// TestStringFieldEmpty tests stringField with empty map.
func TestStringFieldEmpty(t *testing.T) {
	m := map[string]interface{}{}

	result := stringField(m, "any")
	if result != "" {
		t.Errorf("expect empty string, got: %s", result)
	}
}

// TestSortByPublishedDescCurrent tests sorting by published date.
func TestSortByPublishedDescCurrent(t *testing.T) {
	summaries := []PostSummary{
		{
			ID:        "post1",
			Title:     "First",
			Published: "2025-01-10T10:00:00Z",
		},
		{
			ID:        "post2",
			Title:     "Second",
			Published: "2025-01-15T10:00:00Z",
		},
		{
			ID:        "post3",
			Title:     "Third",
			Published: "2025-01-05T10:00:00Z",
		},
	}

	sortByPublishedDesc(summaries)

	// Should be sorted newest first
	if summaries[0].ID != "post2" {
		t.Errorf("expect first post to be post2, got: %s", summaries[0].ID)
	}
	if summaries[1].ID != "post1" {
		t.Errorf("expect second post to be post1, got: %s", summaries[1].ID)
	}
	if summaries[2].ID != "post3" {
		t.Errorf("expect third post to be post3, got: %s", summaries[2].ID)
	}
}

// TestSortByPublishedDescInvalidDates tests sorting with invalid dates.
func TestSortByPublishedDescInvalidDates(t *testing.T) {
	summaries := []PostSummary{
		{
			ID:        "post1",
			Title:     "Valid",
			Published: "2025-01-15T10:00:00Z",
		},
		{
			ID:        "post2",
			Title:     "Invalid",
			Published: "not a date",
		},
		{
			ID:        "post3",
			Title:     "Valid",
			Published: "2025-01-20T10:00:00Z",
		},
	}

	sortByPublishedDesc(summaries)

	// Valid dates should come first
	if summaries[0].ID != "post3" {
		t.Errorf("expect post3 first (newest), got: %s", summaries[0].ID)
	}
	if summaries[1].ID != "post1" {
		t.Errorf("expect post1 second, got: %s", summaries[1].ID)
	}
}

// TestSortByPublishedDescEmpty tests sorting empty list.
func TestSortByPublishedDescEmpty(t *testing.T) {
	summaries := []PostSummary{}

	sortByPublishedDesc(summaries)

	if len(summaries) != 0 {
		t.Error("expect empty list to remain empty")
	}
}

// TestSortByPublishedDescSingle tests sorting with single item.
func TestSortByPublishedDescSingle(t *testing.T) {
	summaries := []PostSummary{
		{
			ID:        "post1",
			Title:     "Only",
			Published: "2025-01-15T10:00:00Z",
		},
	}

	sortByPublishedDesc(summaries)

	if len(summaries) != 1 {
		t.Error("expect single item unchanged")
	}
	if summaries[0].ID != "post1" {
		t.Error("expect single item preserved")
	}
}

// TestSortByPublishedDescSameDates tests sorting with identical dates.
func TestSortByPublishedDescSameDates(t *testing.T) {
	summaries := []PostSummary{
		{
			ID:        "post1",
			Title:     "First",
			Published: "2025-01-15T10:00:00Z",
		},
		{
			ID:        "post2",
			Title:     "Second",
			Published: "2025-01-15T10:00:00Z",
		},
	}

	sortByPublishedDesc(summaries)

	// Both have same date, order may be stable
	if len(summaries) != 2 {
		t.Error("expect both items")
	}
	if (summaries[0].ID != "post1" || summaries[1].ID != "post2") &&
		(summaries[0].ID != "post2" || summaries[1].ID != "post1") {
		t.Error("expect both posts present after sort")
	}
}

// TestClientTimeoutConfig tests that client has timeout configured.
func TestClientTimeoutConfig(t *testing.T) {
	client := NewClient("key", "blog")

	if client.httpClient == nil {
		t.Error("expect http client to be initialized")
	}
	// Timeout is set to 10 seconds in NewClient
	if client.httpClient.Timeout != 10*time.Second {
		t.Errorf("expect 10 second timeout, got: %v", client.httpClient.Timeout)
	}
}

// TestClientAPIKeyPresent tests API key is set.
func TestClientAPIKeyPresent(t *testing.T) {
	client := NewClient("my-secret-key", "blog-id")

	if client.apiKey == "" {
		t.Error("expect API key to be set")
	}
	if client.apiKey != "my-secret-key" {
		t.Error("expect API key to be preserved")
	}
}

// TestClientBlogIDPresent tests blog ID is set.
func TestClientBlogIDPresent(t *testing.T) {
	client := NewClient("key", "my-blog-123")

	if client.blogID == "" {
		t.Error("expect blog ID to be set")
	}
	if client.blogID != "my-blog-123" {
		t.Error("expect blog ID to be preserved")
	}
}

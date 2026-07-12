// Package blogger provides read-only access to Google Blogger API v3.
package blogger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const baseURL = "https://www.googleapis.com/blogger/v3"

// Client queries Blogger API for posts.
type Client struct {
	httpClient *http.Client
	apiKey     string
	blogID     string
}

// NewClient creates a new Blogger API client.
func NewClient(apiKey, blogID string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		apiKey:     apiKey,
		blogID:     blogID,
	}
}

// PostSummary is a post listing for dropdown selection.
type PostSummary struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Published string `json:"published"`
}

type postsListResponse struct {
	Items []map[string]interface{} `json:"items"`
}

// ListPosts retrieves recent posts from the blog.
func (c *Client) ListPosts(ctx context.Context, maxResults int) ([]PostSummary, error) {
	q := url.Values{}
	q.Set("key", c.apiKey)
	q.Set("maxResults", fmt.Sprintf("%d", maxResults))
	q.Set("fetchBodies", "false")

	endpoint := fmt.Sprintf("%s/blogs/%s/posts?%s", baseURL, c.blogID, q.Encode())
	return c.fetchPostSummaries(ctx, endpoint)
}

// SearchPosts performs a full-text search across all blog posts.
func (c *Client) SearchPosts(ctx context.Context, query string) ([]PostSummary, error) {
	q := url.Values{}
	q.Set("key", c.apiKey)
	q.Set("q", query)
	q.Set("fetchBodies", "false")

	endpoint := fmt.Sprintf("%s/blogs/%s/posts/search?%s", baseURL, c.blogID, q.Encode())
	return c.fetchPostSummaries(ctx, endpoint)
}

func (c *Client) fetchPostSummaries(ctx context.Context, endpoint string) ([]PostSummary, error) {
	var body postsListResponse
	if err := c.getJSON(ctx, endpoint, &body); err != nil {
		return nil, err
	}

	summaries := make([]PostSummary, 0, len(body.Items))
	for _, item := range body.Items {
		summaries = append(summaries, PostSummary{
			ID:        stringField(item, "id"),
			Title:     stringField(item, "title"),
			Published: stringField(item, "published"),
		})
	}
	sortByPublishedDesc(summaries)
	return summaries, nil
}

// sortByPublishedDesc sorts posts by published date, newest first.
func sortByPublishedDesc(summaries []PostSummary) {
	sort.SliceStable(summaries, func(i, j int) bool {
		ti, iErr := time.Parse(time.RFC3339, summaries[i].Published)
		tj, jErr := time.Parse(time.RFC3339, summaries[j].Published)
		if iErr != nil && jErr != nil {
			return false
		}
		if iErr != nil {
			return false
		}
		if jErr != nil {
			return true
		}
		return ti.After(tj)
	})
}

// GetPost retrieves a post by ID with full content as JSON.
func (c *Client) GetPost(ctx context.Context, postID string) (map[string]interface{}, error) {
	q := url.Values{}
	q.Set("key", c.apiKey)

	endpoint := fmt.Sprintf("%s/blogs/%s/posts/%s?%s", baseURL, c.blogID, url.PathEscape(postID), q.Encode())
	var post map[string]interface{}
	if err := c.getJSON(ctx, endpoint, &post); err != nil {
		return nil, err
	}
	return post, nil
}

func (c *Client) getJSON(ctx context.Context, endpoint string, out interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("blogger API request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		// Include the (bounded) response body so container logs reveal the actual
		// reason (e.g. "API key not valid") instead of just the status code.
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("blogger API returned status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("failed to decode blogger API response: %w", err)
	}
	return nil
}

func stringField(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

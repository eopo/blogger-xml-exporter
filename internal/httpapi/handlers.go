// Package httpapi provides HTTP request handlers for the application.
package httpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/eopo/blogger-xml-exporter/internal/blogger"
	"github.com/eopo/blogger-xml-exporter/internal/config"
	"github.com/eopo/blogger-xml-exporter/internal/xmlgen"
)

// Server holds HTTP handler dependencies.
type Server struct {
	cfg    *config.Config
	client *blogger.Client
}

// New creates a new API server.
func New(cfg *config.Config, client *blogger.Client) *Server {
	return &Server{cfg: cfg, client: client}
}

// Routes registers all API endpoints and serves the static frontend.
func (s *Server) Routes(mux *http.ServeMux, staticDir string) {
	mux.HandleFunc("GET /healthz", s.handleHealth)
	mux.HandleFunc("GET /api/form-schema", s.handleFormSchema)
	mux.HandleFunc("GET /api/defaults", s.handleDefaults)
	mux.HandleFunc("GET /api/posts", s.handleListPosts)
	mux.HandleFunc("GET /api/posts/{id}", s.handleGetPost)
	mux.HandleFunc("POST /api/generate", s.handleGenerate)
	if s.cfg.Assets.Dir != "" {
		mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(s.cfg.Assets.Dir))))
	}
	mux.Handle("/", http.FileServer(http.Dir(staticDir)))
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (s *Server) handleFormSchema(w http.ResponseWriter, _ *http.Request) {
	assets := map[string]string{}
	if s.cfg.Assets.Dir != "" {
		if s.cfg.Assets.Favicon != "" {
			assets["favicon"] = "/assets/" + s.cfg.Assets.Favicon
		}
		if s.cfg.Assets.Logo != "" {
			assets["logo"] = "/assets/" + s.cfg.Assets.Logo
		}
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"items": s.cfg.Form.Items,
		"site": map[string]string{
			"title":   s.cfg.Site.Title,
			"heading": s.cfg.Site.Heading,
		},
		"theme": map[string]string{
			"primaryColor": s.cfg.Theme.PrimaryColor,
			"darkColor":    s.cfg.Theme.DarkColor,
			"lightColor":   s.cfg.Theme.LightColor,
		},
		"assets": assets,
	})
}

// handleDefaults returns field defaults without a selected post.
func (s *Server) handleDefaults(w http.ResponseWriter, _ *http.Request) {
	post := map[string]interface{}{}
	values := blogger.ResolveFields(post, s.cfg.Form.Fields())
	presets := blogger.ResolvePresets(post, s.cfg.Form.Items, values)
	writeJSON(w, http.StatusOK, map[string]interface{}{"post": post, "values": values, "presets": presets})
}

// handleListPosts returns recent posts or search results if query is set.
func (s *Server) handleListPosts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	query := strings.TrimSpace(r.URL.Query().Get("q"))

	var posts []blogger.PostSummary
	var err error
	if query != "" {
		posts, err = s.client.SearchPosts(ctx, query)
	} else {
		posts, err = s.client.ListPosts(ctx, s.cfg.Blogger.MaxResults)
	}
	if err != nil {
		log.Printf("failed to load posts: %v", err)
		writeError(w, http.StatusBadGateway, "failed to load posts")
		return
	}
	writeJSON(w, http.StatusOK, posts)
}

func (s *Server) handleGetPost(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	if postID == "" {
		writeError(w, http.StatusBadRequest, "post ID is required")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	post, err := s.client.GetPost(ctx, postID)
	if err != nil {
		log.Printf("failed to load post %q: %v", postID, err)
		writeError(w, http.StatusBadGateway, "failed to load post")
		return
	}

	values := blogger.ResolveFields(post, s.cfg.Form.Fields())
	presets := blogger.ResolvePresets(post, s.cfg.Form.Items, values)
	writeJSON(w, http.StatusOK, map[string]interface{}{"post": post, "values": values, "presets": presets})
}

type generateRequest struct {
	Post   map[string]interface{} `json:"post"`
	Values map[string]interface{} `json:"values"`
}

func (s *Server) handleGenerate(w http.ResponseWriter, r *http.Request) {
	var req generateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	xmlBytes, err := xmlgen.Render(s.cfg.XML, req.Post, req.Values)
	if err != nil {
		log.Printf("failed to render XML: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to generate XML")
		return
	}
	filename := xmlgen.Filename(s.cfg.XML, req.Post, req.Values)

	w.Header().Set("Content-Type", "application/xml")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename=%q`, filename))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(xmlBytes)
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

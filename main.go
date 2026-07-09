// Command blogger-xml-exporter starts the HTTP server.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/leokr/blogger-xml-exporter/internal/blogger"
	"github.com/leokr/blogger-xml-exporter/internal/config"
	"github.com/leokr/blogger-xml-exporter/internal/httpapi"
)

func main() {
	configPath := envOrDefault("CONFIG_PATH", "config.yaml")
	staticDir := envOrDefault("STATIC_DIR", "web/static")

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	apiKey := os.Getenv("BLOGGER_API_KEY")
	if apiKey == "" {
		log.Fatal("BLOGGER_API_KEY environment variable is not set")
	}

	client := blogger.NewClient(apiKey, cfg.Blogger.BlogID)
	server := httpapi.New(cfg, client)

	mux := http.NewServeMux()
	server.Routes(mux, staticDir)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

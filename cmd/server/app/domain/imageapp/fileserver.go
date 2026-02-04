package imageapp

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/kawai-network/veridium/cmd/server/foundation/web"
	"github.com/kawai-network/veridium/internal/paths"
)

// serveFile serves generated image files
func (a *app) serveFile(w http.ResponseWriter, r *http.Request) {
	if a.engine == nil {
		http.Error(w, "image service not available", http.StatusNotImplemented)
		return
	}

	// Extract filename from path (/files/filename.png)
	path := strings.TrimPrefix(r.URL.Path, "/files/")
	filename := filepath.Base(path)

	if filename == "." || filename == "/" || filename == "" {
		http.Error(w, "invalid filename", http.StatusBadRequest)
		return
	}

	// Build full path
	outputDir := paths.StableDiffusionOutputs()
	fullPath := filepath.Join(outputDir, filename)

	// Security: ensure path is within output directory
	// Clean both paths and ensure directory boundary matching
	cleanOutputDir := filepath.Clean(outputDir)
	cleanPath := filepath.Clean(fullPath)

	// Ensure cleanPath starts with cleanOutputDir followed by separator
	if !strings.HasPrefix(cleanPath, cleanOutputDir+string(filepath.Separator)) &&
		cleanPath != cleanOutputDir {
		http.Error(w, "invalid path", http.StatusForbidden)
		return
	}

	// Check if file exists
	if _, err := os.Stat(cleanPath); os.IsNotExist(err) {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}

	// Serve the file
	http.ServeFile(w, r, cleanPath)
}

// AddFileServerRoute adds the file server route to the app
func AddFileServerRoute(app *web.App, cfg Config) {
	api := newApp(cfg)

	// Serve files without authentication (public access to generated images)
	app.RawHandlerFunc(http.MethodGet, "", "/files/*", api.serveFile)
}

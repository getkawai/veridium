package gateway

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ServerConfig holds configuration for the gateway server.
type ServerConfig struct {
	Host      string
	Port      int
	ModelName string
	StaticDir string // Path to serve static files from
}

// Server represents the OpenAI-compatible gateway HTTP server.
type Server struct {
	engine  *gin.Engine
	server  *http.Server
	handler *Handler
	config  ServerConfig
}

// NewServer creates a new gateway server with the given configuration.
func NewServer(cfg ServerConfig, executor LLMExecutor, whisperExecutor *WhisperExecutor, imageExecutor ImageExecutor) *Server {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())

	handler := NewHandler(executor, whisperExecutor, imageExecutor, cfg.ModelName)

	s := &Server{
		engine:  engine,
		handler: handler,
		config:  cfg,
	}

	s.setupRoutes()
	return s
}

// setupRoutes configures the API routes.
func (s *Server) setupRoutes() {
	// Health check
	s.engine.GET("/health", s.handler.HealthCheck)

	// Serve static files if directory is provided
	if s.config.StaticDir != "" {
		s.engine.Static("/files", s.config.StaticDir)
		fmt.Printf("File server: serving %s at /files\n", s.config.StaticDir)
	}

	// OpenAI-compatible endpoints
	v1 := s.engine.Group("/v1")
	{
		v1.POST("/chat/completions", s.handler.ChatCompletions)
		v1.POST("/audio/transcriptions", s.handler.AudioTranscriptions)
		v1.POST("/images/generations", s.handler.ImageGenerations)
	}
}

// Start starts the HTTP server and blocks until it stops.
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	s.server = &http.Server{
		Addr:         addr,
		Handler:      s.engine,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 5 * time.Minute, // Long timeout for streaming
		IdleTimeout:  120 * time.Second,
	}

	fmt.Printf("Gateway server starting on %s\n", addr)
	fmt.Printf("Model: %s\n", s.config.ModelName)
	fmt.Printf("Endpoint: POST %s/v1/chat/completions\n", addr)

	return s.server.ListenAndServe()
}

// Stop gracefully stops the HTTP server.
func (s *Server) Stop(ctx context.Context) error {
	if s.server == nil {
		return nil
	}
	return s.server.Shutdown(ctx)
}

// Engine returns the underlying Gin engine for testing.
func (s *Server) Engine() *gin.Engine {
	return s.engine
}

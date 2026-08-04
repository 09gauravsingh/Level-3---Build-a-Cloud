// cmd/api/
// ├── main.go                  # Configuration, routing, and server startup
// ├── api.go                   # API structure and service interface
// ├── handlers.go              # REST endpoint logic
// ├── middleware.go            # Authentication and error handling
// ├── models.go                # JSON request and response structures
// ├── kubernetes_client.go     # Kubernetes connection
// └── kubernetes_instances.go  # Kubernetes product operations

package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create structured JSON logging.
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, nil),
	)

	// The API refuses to start without authentication configured.
	token := os.Getenv("API_TOKEN")
	if token == "" {
		logger.Error("API_TOKEN is required")
		os.Exit(1)
	}

	namespace := envOrDefault(
		"PAAS_NAMESPACE",
		"database-services",
	)

	port := envOrDefault("PORT", "8080")

	// Connect the application to Kubernetes.
	service, err := NewKubeService(namespace)
	if err != nil {
		logger.Error(
			"failed to connect to Kubernetes",
			"error",
			err,
		)
		os.Exit(1)
	}

	api := NewAPI(service, logger, token)

	// Create the Gin HTTP router.
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// Public endpoint for Kubernetes health checks.
	router.GET("/healthz", api.Health)

	// Protected product-instance API routes.
	v1 := router.Group("/api/v1")
	v1.Use(api.Authenticate)
	{
		v1.POST("/instances", api.CreateInstance)
		v1.GET("/instances", api.ListInstances)
		v1.DELETE("/instances/:name", api.DeleteInstance)
		v1.GET(
			"/instances/:name/connection",
			api.GetConnection,
		)
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Listen for Control+C or Kubernetes Pod termination.
	stopContext, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	// Start the HTTP server without blocking shutdown handling.
	go func() {
		logger.Info(
			"API server started",
			"port",
			port,
			"namespace",
			namespace,
		)

		err := server.ListenAndServe()

		if err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			logger.Error(
				"API server failed",
				"error",
				err,
			)
		}
	}()

	// Wait until a shutdown signal is received.
	<-stopContext.Done()

	shutdownContext, cancel :=
		context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
	defer cancel()

	// Allow active requests to finish before stopping.
	if err := server.Shutdown(shutdownContext); err != nil {
		logger.Error(
			"server shutdown failed",
			"error",
			err,
		)
		os.Exit(1)
	}

	logger.Info("API server stopped")
}

// envOrDefault reads configuration from an environment variable.
func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}

	return fallback
}

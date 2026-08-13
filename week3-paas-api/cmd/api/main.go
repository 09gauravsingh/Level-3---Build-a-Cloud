// cmd/api/ // ├── main.go # Configuration, routing, and server startup
// // ├── api.go # API structure and service interface // ├── handlers.go # REST endpoint logic
// // ├── middleware.go # Authentication and error handling // ├── models.go # JSON request and response structures
//  // ├── kubernetes_client.go # Kubernetes connection // └── kubernetes_instances.go # Kubernetes product operations
// // The program starts in main.go. // main.go creates the API object from api.go, // registers middleware and handlers, and starts the server.
// // When a request comes, it goes through middleware first,
// then the handler, then the Kubernetes service. // Program entry point, configuration, logging and // server lifecycle

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

	apihttp "codeberg.org/gauravsingh78945/build-a-cloud/week3-paas-api/internal/api"
	"codeberg.org/gauravsingh78945/build-a-cloud/week3-paas-api/internal/platform"
)

func main() {
	// Create structured JSON logging.
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, nil),
	)

	// Week 4: JWT authentication requires these values.
	if os.Getenv("ADMIN_USERNAME") == "" ||
		os.Getenv("ADMIN_PASSWORD") == "" ||
		os.Getenv("JWT_SECRET") == "" {

		logger.Error("JWT authentication configuration is required")
		os.Exit(1)
	}

	// Namespace where PostgreSQL instances are managed.
	namespace := envOrDefault(
		"PAAS_NAMESPACE",
		"database-services",
	)

	port := envOrDefault("PORT", "8080")

	// Connect the REST API to Kubernetes.
	service, err := platform.NewKubeService(namespace)
	if err != nil {
		logger.Error(
			"failed to connect to Kubernetes",
			"error",
			err,
		)
		os.Exit(1)
	}

	// Configure HTTP server.
	server := &http.Server{
		Addr: ":" + port,

		Handler: apihttp.NewRouter(
			service,
			logger,
		),

		ReadHeaderTimeout: 5 * time.Second,
	}

	// Handle Ctrl+C and Kubernetes termination.
	stopContext, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

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

	<-stopContext.Done()

	// Give active requests time to finish.
	shutdownContext, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownContext); err != nil {
		logger.Error("shutdown failed", "error", err)
		os.Exit(1)
	}

	logger.Info("API server stopped")
}

// envOrDefault returns an environment value or its fallback.
func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}

	return fallback
}

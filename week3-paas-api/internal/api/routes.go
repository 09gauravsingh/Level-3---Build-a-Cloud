// Creates the Gin router and registers
// public/protected routes.
// This file defines which routes exist and which handler runs.
package api

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"codeberg.org/gauravsingh78945/build-a-cloud/week3-paas-api/internal/models"
)

// PlatformService defines the Kubernetes operations needed by the REST API.
type PlatformService interface {
	Create(
		context.Context,
		models.CreateInstanceRequest,
	) (models.Instance, error)

	List(context.Context) ([]models.Instance, error)

	Delete(context.Context, string) error

	Connection(
		context.Context,
		string,
	) (models.ConnectionData, error)
}

// API contains dependencies used by all handlers.
type API struct {
	service PlatformService
	logger  *slog.Logger
	token   string
}

// NewRouter creates the Gin router and registers all REST API routes.
func NewRouter(
	service PlatformService,
	logger *slog.Logger,
	token string,
) *gin.Engine {
	api := &API{
		service: service,
		logger:  logger,
		token:   token,
	}

	// Create a new Gin router.
	router := gin.New()

	// Log requests and prevent panics from crashing the API.
	router.Use(gin.Logger(), gin.Recovery())

	// Allow Swagger Editor to call the local API from the browser.
	// This middleware must run before authentication and route handlers.
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"https://editor.swagger.io",
		},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
		ExposeHeaders: []string{
			"Location",
		},
		AllowPrivateNetwork: true,
		MaxAge:              12 * time.Hour,
	}))

	// Public health endpoint.
	router.GET("/healthz", api.Health)
	router.POST("/api/v1/login", api.login)

	// All product endpoints require bearer-token authentication.
	v1 := router.Group("/api/v1")
	v1.Use(api.Authenticate)

	v1.POST("/instances", api.CreateInstance)
	v1.GET("/instances", api.ListInstances)
	v1.DELETE("/instances/:name", api.DeleteInstance)
	v1.GET("/instances/:name/connection", api.GetConnection)

	return router
}

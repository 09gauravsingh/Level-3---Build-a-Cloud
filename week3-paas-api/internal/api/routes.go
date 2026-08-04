//This file only defines which routes exist and which handler runs.

package api

import (
	"context"
	"log/slog"

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

// NewRouter registers all REST API routes.
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

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// Public health endpoint.
	router.GET("/healthz", api.Health)

	// All product endpoints require authentication.
	v1 := router.Group("/api/v1")
	v1.Use(api.Authenticate)

	v1.POST("/instances", api.CreateInstance)
	v1.GET("/instances", api.ListInstances)
	v1.DELETE("/instances/:name", api.DeleteInstance)
	v1.GET("/instances/:name/connection", api.GetConnection)

	return router
}

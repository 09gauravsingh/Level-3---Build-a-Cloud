package main

import (
	"context"
	"log/slog"
)

// PlatformService defines the Kubernetes operations used by the REST API.
// Unit tests can later replace KubeService with a fake service.
type PlatformService interface {
	Create(context.Context, CreateInstanceRequest) (Instance, error)
	List(context.Context) ([]Instance, error)
	Delete(context.Context, string) error
	Connection(context.Context, string) (ConnectionData, error)
}

// API stores the dependencies required by the handlers.
type API struct {
	service PlatformService
	logger  *slog.Logger
	token   string
}

// NewAPI creates the API handler object.
func NewAPI(
	service PlatformService,
	logger *slog.Logger,
	token string,
) *API {
	return &API{
		service: service,
		logger:  logger,
		token:   token,
	}
}

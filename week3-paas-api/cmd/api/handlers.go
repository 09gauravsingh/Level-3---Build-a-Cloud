package main

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

// Health confirms that the API process is running.
func (api *API) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "week3-paas-api",
	})
}

// CreateInstance creates a CloudNativePG Cluster resource.
func (api *API) CreateInstance(c *gin.Context) {
	var request CreateInstanceRequest

	// Convert the incoming JSON body into a Go struct.
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid JSON request",
		})
		return
	}

	setDefaults(&request)

	if message := validateCreateRequest(request); message != "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Message: message,
		})
		return
	}

	// Ask Kubernetes to create the PostgreSQL product instance.
	instance, err := api.service.Create(
		c.Request.Context(),
		request,
	)
	if err != nil {
		api.ServiceError(c, err)
		return
	}

	// The Operator continues provisioning asynchronously.
	c.Header(
		"Location",
		"/api/v1/instances/"+instance.Name,
	)

	c.JSON(http.StatusAccepted, instance)
}

// ListInstances returns all databases created through this API.
func (api *API) ListInstances(c *gin.Context) {
	instances, err := api.service.List(
		c.Request.Context(),
	)
	if err != nil {
		api.ServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, InstanceList{
		Items: instances,
		Count: len(instances),
	})
}

// DeleteInstance starts deletion of one PostgreSQL instance.
func (api *API) DeleteInstance(c *gin.Context) {
	name := c.Param("name")

	if err := api.service.Delete(
		c.Request.Context(),
		name,
	); err != nil {
		api.ServiceError(c, err)
		return
	}

	// Kubernetes and CloudNativePG perform the actual cleanup.
	c.JSON(http.StatusAccepted, OperationResponse{
		Name:   name,
		Status: "Deleting",
	})
}

// GetConnection returns database connection data from a Kubernetes Secret.
func (api *API) GetConnection(c *gin.Context) {
	connection, err := api.service.Connection(
		c.Request.Context(),
		c.Param("name"),
	)
	if err != nil {
		api.ServiceError(c, err)
		return
	}

	// Prevent browsers and proxies from caching credentials.
	c.Header("Cache-Control", "no-store")
	c.JSON(http.StatusOK, connection)
}

// setDefaults supplies simple values for optional fields.
func setDefaults(request *CreateInstanceRequest) {
	if request.Instances == 0 {
		request.Instances = 1
	}

	if request.StorageSize == "" {
		request.StorageSize = "1Gi"
	}

	if request.Database == "" {
		request.Database = "app"
	}

	if request.Owner == "" {
		request.Owner = request.Database
	}
}

var validInstanceName = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)

// validateCreateRequest checks basic product limits.
func validateCreateRequest(request CreateInstanceRequest) string {
	if request.Name == "" ||
		!validInstanceName.MatchString(request.Name) {
		return "name must use lowercase letters, numbers or hyphens"
	}

	if request.Instances < 1 || request.Instances > 3 {
		return "instances must be between 1 and 3"
	}

	return ""
}

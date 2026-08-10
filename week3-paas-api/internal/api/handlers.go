// HTTP request logic: bind JSON, defaults, validation,
// call service, return JSON.

package api

import (
	"net/http"
	"os"
	"regexp"
	"time"

	"codeberg.org/gauravsingh78945/build-a-cloud/week3-paas-api/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWT Authentication Code:
func (api *API) login(c *gin.Context) {
	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	//Read login JSON.
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	//Check Credentials:
	if request.Username != os.Getenv("ADMIN_USERNAME") || request.Password != os.Getenv("ADMIN_PASSWORD") {
		c.JSON(401, gin.H{"error": "Invalid Credentials"})
		return
	}

	//JWT valid for one hour
	claims := jwt.MapClaims{
		"sub": request.Username,
		"exp": time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//Sign JWT using our secret.
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(500, gin.H{"error": "could not create token"})
		return
	}

	c.JSON(200, gin.H{"token": signedToken})
}

// Health confirms that the API process is running.
func (api *API) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "week3-paas-api",
	})
}

// CreateInstance creates a CloudNativePG Cluster resource.
func (api *API) CreateInstance(c *gin.Context) {
	var request models.CreateInstanceRequest

	// Convert the incoming JSON into a Go struct.
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid JSON request",
		})
		return
	}

	setDefaults(&request)

	if message := validateCreateRequest(request); message != "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: message,
		})
		return
	}

	// Ask Kubernetes to create the PostgreSQL instance.
	instance, err := api.service.Create(
		c.Request.Context(),
		request,
	)
	if err != nil {
		api.ServiceError(c, err)
		return
	}

	// CloudNativePG continues provisioning asynchronously.
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

	c.JSON(http.StatusOK, models.InstanceList{
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

	// Kubernetes performs the actual cleanup asynchronously.
	c.JSON(http.StatusAccepted, models.OperationResponse{
		Name:   name,
		Status: "Deleting",
	})
}

// GetConnection returns data from the CloudNativePG Secret.
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

// setDefaults supplies values for optional request fields.
func setDefaults(request *models.CreateInstanceRequest) {
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

var validInstanceName = regexp.MustCompile(
	`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`,
)

// validateCreateRequest checks basic product limits.
func validateCreateRequest(
	request models.CreateInstanceRequest,
) string {
	if request.Name == "" ||
		!validInstanceName.MatchString(request.Name) {
		return "name must use lowercase letters, numbers or hyphens"
	}

	if request.Instances < 1 || request.Instances > 3 {
		return "instances must be between 1 and 3"
	}

	return ""
}

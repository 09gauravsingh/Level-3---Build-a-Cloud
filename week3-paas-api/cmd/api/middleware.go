package main

import (
	"crypto/subtle"
	"net/http"

	"github.com/gin-gonic/gin"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

// Authenticate checks the bearer token before protected handlers run.
func (api *API) Authenticate(c *gin.Context) {
	expected := "Bearer " + api.token
	provided := c.GetHeader("Authorization")

	// Constant-time comparison reduces timing-based token attacks.
	if subtle.ConstantTimeCompare(
		[]byte(provided),
		[]byte(expected),
	) != 1 {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			ErrorResponse{
				Error:   "Unauthorized",
				Message: "invalid bearer token",
			},
		)
		return
	}

	// Continue to the requested endpoint.
	c.Next()
}

// ServiceError converts Kubernetes errors into REST responses.
func (api *API) ServiceError(
	c *gin.Context,
	err error,
) {
	switch {
	case apierrors.IsNotFound(err):
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Not Found",
			Message: "instance not found",
		})

	case apierrors.IsAlreadyExists(err):
		c.JSON(http.StatusConflict, ErrorResponse{
			Error:   "Conflict",
			Message: "instance already exists",
		})

	default:
		// Log internal details but return a safe message to the client.
		api.logger.Error(
			"Kubernetes operation failed",
			"error",
			err,
		)

		c.JSON(
			http.StatusInternalServerError,
			ErrorResponse{
				Error:   "Internal Server Error",
				Message: "operation failed",
			},
		)
	}
}

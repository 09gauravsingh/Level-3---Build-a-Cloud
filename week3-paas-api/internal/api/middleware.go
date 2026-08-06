//Bearer-token authentication and shared error
//translation. This file handles authentication and error handling.

package api

import (
	"crypto/subtle"
	"net/http"

	"github.com/gin-gonic/gin"
	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"codeberg.org/gauravsingh78945/build-a-cloud/week3-paas-api/internal/models"
)

// Authenticate protects product endpoints with a bearer token.
func (api *API) Authenticate(c *gin.Context) {
	expected := "Bearer " + api.token
	provided := c.GetHeader("Authorization")

	if subtle.ConstantTimeCompare(
		[]byte(provided),
		[]byte(expected),
	) != 1 {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			models.ErrorResponse{
				Error:   "Unauthorized",
				Message: "invalid bearer token",
			},
		)
		return
	}

	c.Next()
}

// ServiceError converts Kubernetes errors into REST responses.
func (api *API) ServiceError(
	c *gin.Context,
	err error,
) {
	switch {
	case apierrors.IsNotFound(err):
		c.JSON(
			http.StatusNotFound,
			models.ErrorResponse{
				Error:   "Not Found",
				Message: "instance not found",
			},
		)

	case apierrors.IsAlreadyExists(err):
		c.JSON(
			http.StatusConflict,
			models.ErrorResponse{
				Error:   "Conflict",
				Message: "instance already exists",
			},
		)

	default:
		api.logger.Error(
			"Kubernetes operation failed",
			"error",
			err,
		)

		c.JSON(
			http.StatusInternalServerError,
			models.ErrorResponse{
				Error:   "Internal Server Error",
				Message: "operation failed",
			},
		)
	}
}

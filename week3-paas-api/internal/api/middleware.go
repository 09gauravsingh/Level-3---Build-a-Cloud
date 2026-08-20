//Bearer-token authentication(implemented using JWT) and shared error
//translation. This file handles authentication and error handling.

package api

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"codeberg.org/gauravsingh78945/build-a-cloud/week3-paas-api/internal/models"
)

// Authenticate checks the JWT token before allowing access.
func (api *API) Authenticate(c *gin.Context) {
	auth := c.GetHeader("Authorization")

	// Header must look like:
	// Authorization: Bearer <JWT>
	if !strings.HasPrefix(auth, "Bearer ") {
		c.JSON(401, gin.H{"error": "missing bearer token"})
		c.Abort()
		return
	}

	// Remove "Bearer " and keep only the JWT.
	tokenString := strings.TrimPrefix(auth, "Bearer ")

	//verify signature and expiration
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (any, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
	)

	if err != nil || !token.Valid {
		c.JSON(401, gin.H{"error": "Invalid or expired token"})
		c.Abort()
		return
	}

	//Remember who is calling so handlers can scope data to one user.
	claims, _ := token.Claims.(jwt.MapClaims)

	subject, _ := claims["sub"].(string)
	isAdmin := claims["role"] == roleAdmin

	//A non-admin without a subject has no owner scope of its own and
	//would otherwise be able to read every instance.
	if !isAdmin && subject == "" {
		c.JSON(401, gin.H{"error": "Invalid or expired token"})
		c.Abort()
		return
	}

	c.Set(contextUsername, subject)
	c.Set(contextIsAdmin, isAdmin)

	//JWT is valid --> continue to Handler(or allow access).
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

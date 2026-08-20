// Self-service registration and the per-request owner scope used to keep
// one user's instances invisible to another.

package api

import (
	"errors"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"

	"codeberg.org/gauravsingh78945/build-a-cloud/week3-paas-api/internal/models"
)

// Claim values for the "role" field of the JWT.
const (
	roleAdmin = "admin"
	roleUser  = "user"
)

// Gin context keys filled in by the Authenticate middleware.
const (
	contextUsername = "username"
	contextIsAdmin  = "isAdmin"
)

// The username ends up as a Kubernetes label value, so it must stay within
// what Kubernetes accepts: lowercase letters, digits and hyphens only.
var validUsername = regexp.MustCompile(
	`^[a-z0-9]([a-z0-9-]{1,30}[a-z0-9])?$`,
)

const minPasswordLength = 8

// register creates a new account from a username and password.
func (api *API) register(c *gin.Context) {
	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid JSON request",
		})
		return
	}

	username := strings.ToLower(strings.TrimSpace(request.Username))

	if message := validateRegistration(
		username,
		request.Password,
	); message != "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: message,
		})
		return
	}

	err := api.users.CreateUser(username, request.Password)

	if errors.Is(err, ErrUserExists) {
		c.JSON(http.StatusConflict, models.ErrorResponse{
			Error:   "Conflict",
			Message: "username is already taken",
		})
		return
	}

	if err != nil {
		api.logger.Error(
			"failed to create user",
			"error",
			err,
		)

		c.JSON(
			http.StatusInternalServerError,
			models.ErrorResponse{
				Error:   "Internal Server Error",
				Message: "could not create account",
			},
		)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"username": username})
}

// validateRegistration returns a message when the credentials are unusable.
func validateRegistration(username, password string) string {
	if !validUsername.MatchString(username) {
		return "username must be 3-32 lowercase letters, numbers or hyphens"
	}

	// The admin signs in with environment credentials, not with a row.
	if username == strings.ToLower(os.Getenv("ADMIN_USERNAME")) {
		return "username is reserved"
	}

	if len(password) < minPasswordLength {
		return "password must be at least 8 characters"
	}

	return ""
}

// ownerScope returns the username whose instances the caller may see.
// Administrators get an empty scope, which means every instance.
func (api *API) ownerScope(c *gin.Context) string {
	if isAdmin, _ := c.Get(contextIsAdmin); isAdmin == true {
		return ""
	}

	username, _ := c.Get(contextUsername)

	name, _ := username.(string)

	return name
}

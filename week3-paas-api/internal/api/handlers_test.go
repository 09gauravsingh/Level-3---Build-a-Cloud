package api

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"codeberg.org/gauravsingh78945/build-a-cloud/week3-paas-api/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Credentials the API is configured with for the whole test binary.
const (
	testSecret   = "test-secret"
	testUsername = "test-admin"
	testPassword = "test-password"
)

// TestMain configures the environment the login handler and the
// authentication middleware read from, before any test runs.
func TestMain(m *testing.M) {
	os.Setenv("JWT_SECRET", testSecret)
	os.Setenv("ADMIN_USERNAME", testUsername)
	os.Setenv("ADMIN_PASSWORD", testPassword)

	os.Exit(m.Run())
}

// signedTestToken mints an administrator JWT that the middleware accepts:
// signed with the same secret, using the same algorithm, and not yet expired.
func signedTestToken(t *testing.T) string {
	t.Helper()

	return signedTokenFor(t, testUsername, roleAdmin)
}

// signedTokenFor mints a JWT for one username and role, so tests can act as
// an ordinary user whose view is limited to their own instances.
func signedTokenFor(t *testing.T, username, role string) string {
	t.Helper()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  username,
		"role": role,
		"exp":  time.Now().Add(time.Hour).Unix(),
	})

	signed, err := token.SignedString([]byte(testSecret))
	if err != nil {
		t.Fatal(err)
	}

	return signed
}

// Readable arguments for performRequest, so call sites read as
// performRequest(router, http.MethodGet, path, noBody, withAuth)
// instead of a trail of bare "" / true / false.
const (
	noBody = ""

	withAuth    = true
	withoutAuth = false
)

// fakeService is a stand-in for the real Kubernetes-backed PlatformService.
//
// Each field holds the behavior for one interface method. A test sets only the
// function it cares about; any method left nil returns a harmless zero value,
// so tests stay focused on the endpoint under test.
//
// The trailing owner argument is the scope the handler derived from the JWT.
type fakeService struct {
	createFn     func(context.Context, models.CreateInstanceRequest, string) (models.Instance, error)
	listFn       func(context.Context, string) ([]models.Instance, error)
	deleteFn     func(context.Context, string, string) error
	connectionFn func(context.Context, string, string) (models.ConnectionData, error)
}

func (f *fakeService) Create(ctx context.Context, request models.CreateInstanceRequest, owner string) (models.Instance, error) {
	if f.createFn == nil {
		return models.Instance{}, nil
	}

	return f.createFn(ctx, request, owner)
}

func (f *fakeService) List(ctx context.Context, owner string) ([]models.Instance, error) {
	if f.listFn == nil {
		return []models.Instance{}, nil
	}

	return f.listFn(ctx, owner)
}

func (f *fakeService) Delete(ctx context.Context, name, owner string) error {
	if f.deleteFn == nil {
		return nil
	}

	return f.deleteFn(ctx, name, owner)
}

func (f *fakeService) Connection(ctx context.Context, name, owner string) (models.ConnectionData, error) {
	if f.connectionFn == nil {
		return models.ConnectionData{}, nil
	}

	return f.connectionFn(ctx, name, owner)
}

// newTestRouter builds the real router (real routes, real middleware) wired to
// a fake service, so the tests exercise routing and auth but never touch a
// cluster. Logs go to io.Discard to keep test output readable. Registered
// users live in a throwaway SQLite file that disappears with the test.
func newTestRouter(t *testing.T, service PlatformService) http.Handler {
	t.Helper()

	gin.SetMode(gin.TestMode)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	users, err := NewUserStore(filepath.Join(t.TempDir(), "users.db"))
	if err != nil {
		t.Fatal(err)
	}

	return NewRouter(service, logger, users)
}

// performRequest runs one in-memory HTTP request through the router and
// returns the recorded response. No network or listening port is involved.
func performRequest(
	t *testing.T,
	router http.Handler,
	method, path, body string,
	authenticated bool,
) *httptest.ResponseRecorder {
	t.Helper()

	request := httptest.NewRequest(method, path, strings.NewReader(body))

	if body != noBody {
		request.Header.Set("Content-Type", "application/json")
	}

	if authenticated {
		request.Header.Set("Authorization", "Bearer "+signedTestToken(t))
	}

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	return response
}

// requireStatus fails the test unless the response carries the expected code.
func requireStatus(t *testing.T, response *httptest.ResponseRecorder, want int) {
	t.Helper()

	if response.Code != want {
		t.Fatalf("expected %d, got %d", want, response.Code)
	}
}

// decodeBody parses the JSON response body into target, failing on bad JSON.
func decodeBody(t *testing.T, response *httptest.ResponseRecorder, target any) {
	t.Helper()

	if err := json.NewDecoder(response.Body).Decode(target); err != nil {
		t.Fatal(err)
	}
}

// TestHealth: /healthz is public and answers without a token.
func TestHealth(t *testing.T) {
	router := newTestRouter(t, &fakeService{})

	response := performRequest(t, router, http.MethodGet, "/healthz", noBody, withoutAuth)

	requireStatus(t, response, http.StatusOK)
}

// TestAuthentication: /api/v1 routes are rejected when the token is missing.
func TestAuthentication(t *testing.T) {
	router := newTestRouter(t, &fakeService{})

	response := performRequest(t, router, http.MethodGet, "/api/v1/instances", noBody, withoutAuth)

	requireStatus(t, response, http.StatusUnauthorized)
}

// TestAuthenticationRejectsForeignToken: a well-formed JWT signed with the
// wrong secret must not grant access.
func TestAuthenticationRejectsForeignToken(t *testing.T) {
	router := newTestRouter(t, &fakeService{})

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": testUsername,
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	signed, err := token.SignedString([]byte("another-secret"))
	if err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest(http.MethodGet, "/api/v1/instances", nil)
	request.Header.Set("Authorization", "Bearer "+signed)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	requireStatus(t, response, http.StatusUnauthorized)
}

// TestLogin: valid credentials return a token the protected routes accept,
// and wrong credentials are refused.
func TestLogin(t *testing.T) {
	router := newTestRouter(t, &fakeService{})

	body := `{"username":"` + testUsername + `","password":"` + testPassword + `"}`

	response := performRequest(t, router, http.MethodPost, "/api/v1/login", body, withoutAuth)

	requireStatus(t, response, http.StatusOK)

	var result struct {
		Token string `json:"token"`
	}

	decodeBody(t, response, &result)

	if result.Token == "" {
		t.Fatal("expected a token in the login response")
	}

	wrong := `{"username":"` + testUsername + `","password":"wrong"}`

	response = performRequest(t, router, http.MethodPost, "/api/v1/login", wrong, withoutAuth)

	requireStatus(t, response, http.StatusUnauthorized)
}

// TestCreateInstance: POST accepts the request body, reports 202 Accepted,
// points Location at the new instance, and echoes the created instance back.
func TestCreateInstance(t *testing.T) {
	// The fake reflects the request back so we can prove the handler decoded
	// the JSON body and passed it through to the service unchanged.
	service := &fakeService{
		createFn: func(_ context.Context, request models.CreateInstanceRequest, _ string) (models.Instance, error) {
			return models.Instance{
				Name:             request.Name,
				Status:           "Provisioning",
				DesiredInstances: request.Instances,
			}, nil
		},
	}

	router := newTestRouter(t, service)

	body := `{
		"name": "demo-db",
		"instances": 1,
		"storageSize": "1Gi",
		"database": "demoapp",
		"owner": "demoowner"
	}`

	response := performRequest(t, router, http.MethodPost, "/api/v1/instances", body, withAuth)

	requireStatus(t, response, http.StatusAccepted)

	if got := response.Header().Get("Location"); got != "/api/v1/instances/demo-db" {
		t.Fatalf("expected Location /api/v1/instances/demo-db, got %q", got)
	}

	var instance models.Instance

	decodeBody(t, response, &instance)

	if instance.Name != "demo-db" {
		t.Fatalf("expected demo-db, got %s", instance.Name)
	}
}

// TestListInstances: GET returns the service's instances wrapped in a list
// whose Count matches the number of items.
func TestListInstances(t *testing.T) {
	service := &fakeService{
		listFn: func(_ context.Context, _ string) ([]models.Instance, error) {
			return []models.Instance{
				{Name: "demo-db", Status: "Ready"},
			}, nil
		},
	}

	router := newTestRouter(t, service)

	response := performRequest(t, router, http.MethodGet, "/api/v1/instances", noBody, withAuth)

	requireStatus(t, response, http.StatusOK)

	var result models.InstanceList

	decodeBody(t, response, &result)

	if result.Count != 1 {
		t.Fatalf("expected 1 instance, got %d", result.Count)
	}
}

// TestDeleteInstance: DELETE returns 202 and forwards the :name path
// parameter to the service.
func TestDeleteInstance(t *testing.T) {
	// Captured inside the fake so we can assert which name the handler parsed
	// out of the URL.
	var deletedName string

	service := &fakeService{
		deleteFn: func(_ context.Context, name, _ string) error {
			deletedName = name

			return nil
		},
	}

	router := newTestRouter(t, service)

	response := performRequest(t, router, http.MethodDelete, "/api/v1/instances/demo-db", noBody, withAuth)

	requireStatus(t, response, http.StatusAccepted)

	if deletedName != "demo-db" {
		t.Fatalf("expected demo-db, got %s", deletedName)
	}
}

// TestConnection: the connection endpoint returns the credentials and marks
// the response no-store, since the body contains a password.
func TestConnection(t *testing.T) {
	service := &fakeService{
		connectionFn: func(_ context.Context, _, _ string) (models.ConnectionData, error) {
			return models.ConnectionData{
				Host:     "demo-db-rw",
				Port:     5432,
				Database: "demoapp",
				Username: "demoowner",
				Password: "test-password",
			}, nil
		},
	}

	router := newTestRouter(t, service)

	response := performRequest(t, router, http.MethodGet, "/api/v1/instances/demo-db/connection", noBody, withAuth)

	requireStatus(t, response, http.StatusOK)

	if got := response.Header().Get("Cache-Control"); got != "no-store" {
		t.Fatalf("expected Cache-Control no-store, got %q", got)
	}

	var connection models.ConnectionData

	decodeBody(t, response, &connection)

	if connection.Username != "demoowner" {
		t.Fatalf("expected demoowner, got %s", connection.Username)
	}
}

// performRequestWithToken runs one request signed by a specific caller, so a
// test can act as an ordinary user instead of the administrator.
func performRequestWithToken(
	t *testing.T,
	router http.Handler,
	method, path, token string,
) *httptest.ResponseRecorder {
	t.Helper()

	request := httptest.NewRequest(method, path, nil)
	request.Header.Set("Authorization", "Bearer "+token)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	return response
}

// registerBody builds a registration payload.
func registerBody(username, password string) string {
	return `{"username":"` + username + `","password":"` + password + `"}`
}

// TestRegister: a new account is created and can immediately sign in.
func TestRegister(t *testing.T) {
	router := newTestRouter(t, &fakeService{})

	body := registerBody("alice", "alice-password")

	response := performRequest(t, router, http.MethodPost, "/api/v1/register", body, withoutAuth)

	requireStatus(t, response, http.StatusCreated)

	response = performRequest(t, router, http.MethodPost, "/api/v1/login", body, withoutAuth)

	requireStatus(t, response, http.StatusOK)

	var result struct {
		Token string `json:"token"`
	}

	decodeBody(t, response, &result)

	if result.Token == "" {
		t.Fatal("expected a token for the registered user")
	}
}

// TestRegisterRejectsDuplicate: the second attempt at a taken username fails.
func TestRegisterRejectsDuplicate(t *testing.T) {
	router := newTestRouter(t, &fakeService{})

	body := registerBody("bob", "bob-password")

	requireStatus(
		t,
		performRequest(t, router, http.MethodPost, "/api/v1/register", body, withoutAuth),
		http.StatusCreated,
	)

	requireStatus(
		t,
		performRequest(t, router, http.MethodPost, "/api/v1/register", body, withoutAuth),
		http.StatusConflict,
	)
}

// TestRegisterRejectsInvalidInput: usernames must be label-safe, passwords
// long enough, and the administrator name is reserved. Mixed case is not an
// error: the handler lowercases the username before storing it.
func TestRegisterRejectsInvalidInput(t *testing.T) {
	router := newTestRouter(t, &fakeService{})

	cases := map[string]string{
		"illegal character": registerBody("al_ice", "alice-password"),
		"short username":    registerBody("al", "alice-password"),
		"short password":    registerBody("alice", "short"),
		"reserved username": registerBody(testUsername, "alice-password"),
	}

	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			response := performRequest(t, router, http.MethodPost, "/api/v1/register", body, withoutAuth)

			requireStatus(t, response, http.StatusBadRequest)
		})
	}
}

// TestLoginRejectsUnknownUser: credentials that were never registered fail.
func TestLoginRejectsUnknownUser(t *testing.T) {
	router := newTestRouter(t, &fakeService{})

	body := registerBody("nobody", "nobody-password")

	response := performRequest(t, router, http.MethodPost, "/api/v1/login", body, withoutAuth)

	requireStatus(t, response, http.StatusUnauthorized)
}

// TestOwnerScope: an ordinary user's requests are limited to their own
// instances, while the administrator sees every instance.
func TestOwnerScope(t *testing.T) {
	var listedOwner, createdOwner, deletedOwner, connectionOwner string

	service := &fakeService{
		listFn: func(_ context.Context, owner string) ([]models.Instance, error) {
			listedOwner = owner

			return []models.Instance{}, nil
		},

		createFn: func(_ context.Context, request models.CreateInstanceRequest, owner string) (models.Instance, error) {
			createdOwner = owner

			return models.Instance{Name: request.Name}, nil
		},

		deleteFn: func(_ context.Context, _, owner string) error {
			deletedOwner = owner

			return nil
		},

		connectionFn: func(_ context.Context, _, owner string) (models.ConnectionData, error) {
			connectionOwner = owner

			return models.ConnectionData{}, nil
		},
	}

	router := newTestRouter(t, service)

	userToken := signedTokenFor(t, "alice", roleUser)

	requireStatus(
		t,
		performRequestWithToken(t, router, http.MethodGet, "/api/v1/instances", userToken),
		http.StatusOK,
	)

	if listedOwner != "alice" {
		t.Fatalf("expected list scoped to alice, got %q", listedOwner)
	}

	requireStatus(
		t,
		performRequestWithToken(t, router, http.MethodDelete, "/api/v1/instances/demo-db", userToken),
		http.StatusAccepted,
	)

	if deletedOwner != "alice" {
		t.Fatalf("expected delete scoped to alice, got %q", deletedOwner)
	}

	requireStatus(
		t,
		performRequestWithToken(t, router, http.MethodGet, "/api/v1/instances/demo-db/connection", userToken),
		http.StatusOK,
	)

	if connectionOwner != "alice" {
		t.Fatalf("expected connection scoped to alice, got %q", connectionOwner)
	}

	// A created instance is labelled with its owner.
	body := `{"name":"demo-db","instances":1}`

	requireStatus(
		t,
		performRequestWithToken(t, router, http.MethodPost, "/api/v1/instances", userToken),
		http.StatusBadRequest,
	)

	request := httptest.NewRequest(http.MethodPost, "/api/v1/instances", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+userToken)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	requireStatus(t, response, http.StatusAccepted)

	if createdOwner != "alice" {
		t.Fatalf("expected create scoped to alice, got %q", createdOwner)
	}

	// The administrator has no scope, which means every instance.
	requireStatus(
		t,
		performRequest(t, router, http.MethodGet, "/api/v1/instances", noBody, withAuth),
		http.StatusOK,
	)

	if listedOwner != "" {
		t.Fatalf("expected an unscoped admin list, got %q", listedOwner)
	}
}

// TestAuthenticationRejectsTokenWithoutSubject: a non-admin token with no
// subject has no owner scope and must not fall back to seeing everything.
func TestAuthenticationRejectsTokenWithoutSubject(t *testing.T) {
	router := newTestRouter(t, &fakeService{})

	response := performRequestWithToken(
		t,
		router,
		http.MethodGet,
		"/api/v1/instances",
		signedTokenFor(t, "", roleUser),
	)

	requireStatus(t, response, http.StatusUnauthorized)
}

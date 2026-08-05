package api

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"codeberg.org/gauravsingh78945/build-a-cloud/week3-paas-api/internal/models"
	"github.com/gin-gonic/gin"
)

const testToken = "test-token"

// fakeService replaces Kubernetes during unit tests.
type fakeService struct {
	createFn     func(context.Context, models.CreateInstanceRequest) (models.Instance, error)
	listFn       func(context.Context) ([]models.Instance, error)
	deleteFn     func(context.Context, string) error
	connectionFn func(context.Context, string) (models.ConnectionData, error)
}

func (f *fakeService) Create(
	ctx context.Context,
	request models.CreateInstanceRequest,
) (models.Instance, error) {
	if f.createFn != nil {
		return f.createFn(ctx, request)
	}

	return models.Instance{}, nil
}

func (f *fakeService) List(
	ctx context.Context,
) ([]models.Instance, error) {
	if f.listFn != nil {
		return f.listFn(ctx)
	}

	return []models.Instance{}, nil
}

func (f *fakeService) Delete(
	ctx context.Context,
	name string,
) error {
	if f.deleteFn != nil {
		return f.deleteFn(ctx, name)
	}

	return nil
}

func (f *fakeService) Connection(
	ctx context.Context,
	name string,
) (models.ConnectionData, error) {
	if f.connectionFn != nil {
		return f.connectionFn(ctx, name)
	}

	return models.ConnectionData{}, nil
}

// newTestRouter creates the real Gin router with a fake service.
func newTestRouter(service PlatformService) http.Handler {
	gin.SetMode(gin.TestMode)

	// Discard logs so test output stays clean.
	logger := slog.New(
		slog.NewTextHandler(io.Discard, nil),
	)

	return NewRouter(service, logger, testToken)
}

// performRequest sends a simulated HTTP request through Gin.
func performRequest(
	router http.Handler,
	method string,
	path string,
	body string,
	withToken bool,
) *httptest.ResponseRecorder {
	request := httptest.NewRequest(
		method,
		path,
		strings.NewReader(body),
	)

	if body != "" {
		request.Header.Set("Content-Type", "application/json")
	}

	if withToken {
		request.Header.Set(
			"Authorization",
			"Bearer "+testToken,
		)
	}

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	return response
}

// TestHealth checks the public health endpoint.
func TestHealth(t *testing.T) {
	router := newTestRouter(&fakeService{})

	response := performRequest(
		router,
		http.MethodGet,
		"/healthz",
		"",
		false,
	)

	if response.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", response.Code)
	}
}

// TestAuthentication checks that product routes require a token.
func TestAuthentication(t *testing.T) {
	router := newTestRouter(&fakeService{})

	response := performRequest(
		router,
		http.MethodGet,
		"/api/v1/instances",
		"",
		false,
	)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", response.Code)
	}
}

// TestCreateInstance checks the POST endpoint.
func TestCreateInstance(t *testing.T) {
	service := &fakeService{
		createFn: func(
			_ context.Context,
			request models.CreateInstanceRequest,
		) (models.Instance, error) {
			return models.Instance{
				Name:             request.Name,
				Status:           "Provisioning",
				DesiredInstances: request.Instances,
			}, nil
		},
	}

	router := newTestRouter(service)

	response := performRequest(
		router,
		http.MethodPost,
		"/api/v1/instances",
		`{
			"name": "demo-db",
			"instances": 1,
			"storageSize": "1Gi",
			"database": "demoapp",
			"owner": "demoowner"
		}`,
		true,
	)

	if response.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", response.Code)
	}

	if response.Header().Get("Location") != "/api/v1/instances/demo-db" {
		t.Fatal("expected Location header for demo-db")
	}

	var instance models.Instance

	if err := json.NewDecoder(response.Body).Decode(&instance); err != nil {
		t.Fatal(err)
	}

	if instance.Name != "demo-db" {
		t.Fatalf("expected demo-db, got %s", instance.Name)
	}
}

// TestListInstances checks the GET list endpoint.
func TestListInstances(t *testing.T) {
	service := &fakeService{
		listFn: func(
			_ context.Context,
		) ([]models.Instance, error) {
			return []models.Instance{
				{
					Name:   "demo-db",
					Status: "Ready",
				},
			}, nil
		},
	}

	router := newTestRouter(service)

	response := performRequest(
		router,
		http.MethodGet,
		"/api/v1/instances",
		"",
		true,
	)

	if response.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", response.Code)
	}

	var result models.InstanceList

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}

	if result.Count != 1 {
		t.Fatalf("expected 1 instance, got %d", result.Count)
	}
}

// TestDeleteInstance checks the DELETE endpoint and path parameter.
func TestDeleteInstance(t *testing.T) {
	var deletedName string

	service := &fakeService{
		deleteFn: func(
			_ context.Context,
			name string,
		) error {
			deletedName = name
			return nil
		},
	}

	router := newTestRouter(service)

	response := performRequest(
		router,
		http.MethodDelete,
		"/api/v1/instances/demo-db",
		"",
		true,
	)

	if response.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", response.Code)
	}

	if deletedName != "demo-db" {
		t.Fatalf("expected demo-db, got %s", deletedName)
	}
}

// TestConnection checks the connection-data endpoint.
func TestConnection(t *testing.T) {
	service := &fakeService{
		connectionFn: func(
			_ context.Context,
			_ string,
		) (models.ConnectionData, error) {
			return models.ConnectionData{
				Host:     "demo-db-rw",
				Port:     5432,
				Database: "demoapp",
				Username: "demoowner",
				Password: "test-password",
			}, nil
		},
	}

	router := newTestRouter(service)

	response := performRequest(
		router,
		http.MethodGet,
		"/api/v1/instances/demo-db/connection",
		"",
		true,
	)

	if response.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", response.Code)
	}

	if response.Header().Get("Cache-Control") != "no-store" {
		t.Fatal("expected Cache-Control: no-store")
	}

	var connection models.ConnectionData

	if err := json.NewDecoder(response.Body).Decode(&connection); err != nil {
		t.Fatal(err)
	}

	if connection.Username != "demoowner" {
		t.Fatalf(
			"expected demoowner, got %s",
			connection.Username,
		)
	}
}

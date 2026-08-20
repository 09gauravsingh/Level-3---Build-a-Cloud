package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"codeberg.org/gauravsingh78945/build-a-cloud/week3-paas-api/internal/models"
)

func TestInstanceAuditQueryDoesNotFilterStreamInstanceLabel(t *testing.T) {
	query := instanceAuditQuery("happydb", "")
	want := `{namespace="paas-system", app="week3-paas-api"} |= "\"instance\":\"happydb\"" | json | type="audit"`

	if query != want {
		t.Fatalf("unexpected query\n got: %s\nwant: %s", query, want)
	}

	if strings.Contains(query, `| instance="happydb"`) {
		t.Fatalf("query still filters the stream instance label: %s", query)
	}
}

func TestInstanceAuditQueryScopesOrdinaryUsersToTheirOwnActions(t *testing.T) {
	query := instanceAuditQuery("happydb", "alice")
	want := `{namespace="paas-system", app="week3-paas-api"} |= "\"instance\":\"happydb\"" |= "\"user\":\"alice\"" | json | type="audit"`

	if query != want {
		t.Fatalf("unexpected query\n got: %s\nwant: %s", query, want)
	}
}

func TestGetInstanceAuditReturnsActivity(t *testing.T) {
	loki := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		if query != instanceAuditQuery("happydb", "") {
			t.Errorf("unexpected Loki query %q", query)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"result": []map[string]any{
					{
						"values": [][]string{
							{
								"1787231852828550009",
								`{"time":"2026-08-20T13:17:32.828550009Z","level":"INFO","msg":"audit","type":"audit","user":"admin","action":"connection.read","instance":"happydb","result":"success"}`,
							},
						},
					},
				},
			},
		})
	}))
	defer loki.Close()

	t.Setenv("LOKI_URL", loki.URL)

	router := newTestRouter(t, &fakeService{
		listFn: func(_ context.Context, _ string) ([]models.Instance, error) {
			return []models.Instance{{Name: "happydb"}}, nil
		},
	})

	response := performRequest(
		t,
		router,
		http.MethodGet,
		"/api/v1/instances/happydb/audit",
		noBody,
		withAuth,
	)

	requireStatus(t, response, http.StatusOK)

	var body struct {
		Instance string          `json:"instance"`
		Activity []ActivityEntry `json:"activity"`
	}
	decodeBody(t, response, &body)

	if body.Instance != "happydb" {
		t.Fatalf("expected instance happydb, got %q", body.Instance)
	}

	if len(body.Activity) != 1 {
		t.Fatalf("expected 1 activity entry, got %d", len(body.Activity))
	}

	got := body.Activity[0]
	if got.User != "admin" || got.Action != "connection.read" || got.Result != "success" {
		t.Fatalf("unexpected activity entry %#v", got)
	}
}

func TestGetInstanceAuditHidesOtherUsersActivity(t *testing.T) {
	loki := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("query") != instanceAuditQuery("happydb", "alice") {
			t.Errorf("unexpected Loki query %q", r.URL.Query().Get("query"))
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"result": []map[string]any{
					{
						"values": [][]string{
							{
								"1787231852828550009",
								`{"type":"audit","user":"admin","action":"connection.read","instance":"happydb","result":"success"}`,
							},
							{
								"1787231852828550010",
								`{"type":"audit","user":"alice","action":"connection.read","instance":"happydb","result":"success"}`,
							},
							{
								"1787231852828550011",
								`{"type":"audit","user":"bob","action":"instance.create","instance":"happydb","result":"success"}`,
							},
						},
					},
				},
			},
		})
	}))
	defer loki.Close()

	t.Setenv("LOKI_URL", loki.URL)

	router := newTestRouter(t, &fakeService{
		listFn: func(_ context.Context, owner string) ([]models.Instance, error) {
			if owner != "alice" {
				t.Errorf("expected owner scope alice, got %q", owner)
			}

			return []models.Instance{{Name: "happydb"}}, nil
		},
	})

	response := performRequestWithToken(
		t,
		router,
		http.MethodGet,
		"/api/v1/instances/happydb/audit",
		signedTokenFor(t, "alice", roleUser),
	)

	requireStatus(t, response, http.StatusOK)

	var body struct {
		Activity []ActivityEntry `json:"activity"`
	}
	decodeBody(t, response, &body)

	if len(body.Activity) != 1 {
		t.Fatalf("expected only alice's activity, got %#v", body.Activity)
	}

	if body.Activity[0].User != "alice" {
		t.Fatalf("expected user alice, got %#v", body.Activity[0])
	}
}

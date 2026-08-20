//this file is basically a Loki client inside your Go API.
//the purpose of this file is Query Loki for logs from the last 24 hours,
//convert Loki’s JSON response into a simple Go format, and return those logs so your frontend/API can display them.

package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// LogEntry is one log line that we return to the frontend.
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Line      string `json:"line"`
}

// ActivityEntry is the safe audit information
// that normal users can see in the UI.
type ActivityEntry struct {
	Timestamp string `json:"timestamp"`
	User      string `json:"user"`
	Action    string `json:"action"`
	Result    string `json:"result"`
}

// Only the audit fields that we need from the raw JSON log.
type auditLogLine struct {
	User   string `json:"user"`
	Action string `json:"action"`
	Result string `json:"result"`
}

// This represents only the part of Loki's response that we need.
//
// Loki gives us something like:
//
// "values": [
//
//	["timestamp", "log line"],
//	["timestamp", "log line"]
//
// ]
type lokiResponse struct {
	Data struct {
		Result []struct {
			Values [][]string `json:"values"`
		} `json:"result"`
	} `json:"data"`
}

// Reuse one HTTP client instead of creating a new client
// for every request.
var lokiHTTPClient = &http.Client{
	Timeout: 10 * time.Second,
}

// -------------------------------------------------------
// Task 7
// queryLoki talks directly to Loki for some logs
// -------------------------------------------------------

func queryLoki(
	ctx context.Context,
	query string,
) ([]LogEntry, error) {

	// Example:
	// http://loki-gateway.logging.svc.cluster.local
	lokiURL := strings.TrimRight(
		os.Getenv("LOKI_URL"),
		"/",
	)

	if lokiURL == "" {
		return nil, errors.New("LOKI_URL is not configured")
	}

	// These become URL parameters for Loki.
	params := url.Values{}

	// Example query:
	// {namespace="database-services", instance="my-db"}
	params.Set("query", query)

	// Retrieve logs from the last 24 hours.
	start := time.Now().Add(-24 * time.Hour)

	params.Set(
		"start",
		strconv.FormatInt(
			start.UnixNano(),
			10,
		),
	)

	params.Set(
		"end",
		strconv.FormatInt(
			time.Now().UnixNano(),
			10,
		),
	)

	// Maximum 100 log lines.
	params.Set("limit", "100")

	// Newest logs first.
	params.Set("direction", "backward")

	// Final URL becomes something like:
	//
	// http://loki-gateway.../loki/api/v1/query_range?query=...
	requestURL :=
		lokiURL +
			"/loki/api/v1/query_range?" +
			params.Encode()

	// Create HTTP request to Loki.
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		requestURL,
		nil,
	)

	if err != nil {
		return nil, err
	}

	// Send request to Loki.
	resp, err := lokiHTTPClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Loki should return HTTP 200.
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"loki returned status %d",
			resp.StatusCode,
		)
	}

	// Convert Loki JSON response into Go struct.
	var result lokiResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// This will contain the simplified logs
	// that we return to Vue.
	logs := []LogEntry{}

	for _, stream := range result.Data.Result {

		for _, value := range stream.Values {

			// Loki gives:
			//
			// value[0] = timestamp
			// value[1] = actual log
			if len(value) != 2 {
				continue
			}

			timestamp := value[0]

			// Loki timestamp is nanoseconds.
			// Convert it into readable UTC time.
			if ns, err := strconv.ParseInt(
				timestamp,
				10,
				64,
			); err == nil {

				timestamp = time.Unix(
					0,
					ns,
				).UTC().Format(time.RFC3339)
			}

			logs = append(
				logs,
				LogEntry{
					Timestamp: timestamp,
					Line:      value[1],
				},
			)
		}
	}

	return logs, nil
}

// -------------------------------------------------------
// Task 8
// This function checks whether the logged-in user can access this DB.
// -------------------------------------------------------

func (api *API) canAccessInstance(
	c *gin.Context,
	name string,
) bool {

	// IMPORTANT:
	//
	// I already implemented user-specific ownership.
	//
	// ownerScope(c) tells your service which user is
	// currently logged in.
	//
	// Therefore List() gives only the instances this
	// user is allowed to see.
	instances, err := api.service.List(
		c.Request.Context(),
		api.ownerScope(c),
	)

	if err != nil {
		api.ServiceError(c, err)
		return false
	}

	// Search this user's instances.
	for _, instance := range instances {

		if instance.Name == name {
			// Found it.
			// User is allowed to access it.
			return true
		}
	}

	// We return 404 instead of exposing whether
	// another user's database exists.
	c.JSON(
		http.StatusNotFound,
		gin.H{
			"error": "instance not found",
		},
	)

	return false
}

// -------------------------------------------------------
// Task 9
// Get the instance logs.(Technical/system/database logs)(PostgreSQl started, connection accepted, etc)
// API endpoint:
// GET /api/v1/instances/:name/logs
// -------------------------------------------------------

func (api *API) GetInstanceLogs(
	c *gin.Context,
) {

	// Example:
	//
	// /instances/gaurav-db/logs
	//
	// name = "gaurav-db"
	name := c.Param("name")

	// FIRST:
	// Check that this logged-in user owns/can access the DB.
	if !api.canAccessInstance(c, name) {
		return
	}

	// Build Loki query.
	//
	// Example:
	//
	// {namespace="database-services",
	//  instance="gaurav-db"}
	//
	// strconv.Quote safely puts quotes around the name.
	query := fmt.Sprintf(
		`{namespace="database-services", instance=%s}`,
		strconv.Quote(name),
	)

	// TASK 7 is used here.
	//
	// Ask Loki for the actual logs.
	logs, err := queryLoki(
		c.Request.Context(),
		query,
	)

	if err != nil {
		c.JSON(
			http.StatusBadGateway,
			gin.H{
				"error": "could not retrieve instance logs",
			},
		)

		return
	}

	// Return logs to Vue.
	c.JSON(
		http.StatusOK,
		gin.H{
			"instance": name,
			"logs":     logs,
		},
	)
}

// instanceAuditQuery returns LogQL that finds this instance's audit events.
//
// Alloy already sets a stream label called instance (namespace/pod:container
// for the API, or the CloudNativePG cluster name for database pods). Filtering
// with `| instance="happydb"` therefore matches that label, not the JSON field
// in the audit log, and returns no rows. Search the raw line instead.
//
// actor is the signed-in username for ordinary users, so they only retrieve
// their own actions. Administrators pass an empty actor and see every event.
func instanceAuditQuery(name, actor string) string {
	query := fmt.Sprintf(
		`{namespace="paas-system", app="week3-paas-api"} |= %s`,
		strconv.Quote(`"instance":"`+name+`"`),
	)

	if actor != "" {
		query += fmt.Sprintf(
			` |= %s`,
			strconv.Quote(`"user":"`+actor+`"`),
		)
	}

	return query + ` | json | type="audit"`
}

// -------------------------------------------------------
// TASK 10
//the purpose of this audit is to get: User action history(create,delete,etc)
// API endpoint:
// GET /api/v1/instances/:name/audit
// -------------------------------------------------------

func (api *API) GetInstanceAudit(
	c *gin.Context,
) {
	name := c.Param("name")

	// Only allow access to an instance
	// that belongs to this logged-in user.
	if !api.canAccessInstance(c, name) {
		return
	}

	// Administrators see every action on this instance.
	// Ordinary users see only the actions they performed.
	actor := api.ownerScope(c)
	query := instanceAuditQuery(name, actor)

	logs, err := queryLoki(
		c.Request.Context(),
		query,
	)

	if err != nil {
		c.JSON(
			http.StatusBadGateway,
			gin.H{
				"error": "could not retrieve activity",
			},
		)
		return
	}

	// Convert the raw audit logs into safe
	// user-facing activity entries.
	activity := []ActivityEntry{}

	for _, logEntry := range logs {
		var auditLine auditLogLine

		if err := json.Unmarshal(
			[]byte(logEntry.Line),
			&auditLine,
		); err != nil {
			continue
		}

		if actor != "" && auditLine.User != actor {
			continue
		}

		activity = append(
			activity,
			ActivityEntry{
				Timestamp: logEntry.Timestamp,
				User:      auditLine.User,
				Action:    auditLine.Action,
				Result:    auditLine.Result,
			},
		)
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"instance": name,
			"activity": activity,
		},
	)
}

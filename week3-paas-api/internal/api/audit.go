// this file creates a separate audit logger for important user actions, and write those audit events
// as JSON logs to stdout so Alloy can collect and send them to loki.
package api

import (
	"log/slog"
	"os"
)

// creates a logger that output logs in JSON format.
var auditLogger = slog.New(
	slog.NewJSONHandler(os.Stdout, nil),
)

// Records an important user action.
func audit(user, action, instance, result string) {
	auditLogger.Info("audit",
		"type", "audit",
		"user", user,
		"action", action,
		"instance", instance,
		"result", result,
	)
}

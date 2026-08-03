package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const serviceName = "week3-paas-api"

type healthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Time    string `json:"time"`
}

func main() {
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, nil),
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	mux.HandleFunc(
		"GET /healthz",
		func(w http.ResponseWriter, _ *http.Request) {
			response := healthResponse{
				Status:  "ok",
				Service: serviceName,
				Time:    time.Now().UTC().Format(time.RFC3339),
			}

			if err := writeJSON(
				w,
				http.StatusOK,
				response,
			); err != nil {
				logger.Error(
					"failed to write health response",
					"error",
					err,
				)
			}
		},
	)

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	shutdownSignal, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	serverError := make(chan error, 1)

	go func() {
		logger.Info(
			"API server starting",
			"address",
			server.Addr,
		)

		serverError <- server.ListenAndServe()
	}()

	select {
	case err := <-serverError:
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Error(
				"API server failed",
				"error",
				err,
			)
			os.Exit(1)
		}

	case <-shutdownSignal.Done():
		logger.Info("shutdown signal received")
	}

	shutdownContext, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownContext); err != nil {
		logger.Error(
			"graceful shutdown failed",
			"error",
			err,
		)
		os.Exit(1)
	}

	logger.Info("API server stopped")
}

func writeJSON(
	w http.ResponseWriter,
	statusCode int,
	payload any,
) error {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(payload)
}

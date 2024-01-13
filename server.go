package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// ServerImpl is an implementation of the OpenAPI ServerInterface.
type ServerImpl struct{}

// Get handles the root endpoint request.
//
// It responds with a JSON payload describing the API surface.
func (s *ServerImpl) Get(ctx echo.Context) error {
	apiDescription := map[string]string{
		"/":                   "Root Endpoint",
		"/latency/{duration}": "Get a response within the provided latency duration",
		"/data/{size}":        "Get a response with a payload matching the provided size criteria",
	}

	if err := ctx.JSON(http.StatusOK, apiDescription); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed marshaling response: %s", err))
	}

	return nil
}

// GetDataSize is a handler returning a payload of the requested size.
//
// It parses the {size} parameter from the URL to determine the size bounds.
// It generates a random payload matching those bounds.
//
// The payload is returned as a Blob with Content-Type "application/octet-stream".
func (s *ServerImpl) GetDataSize(ctx echo.Context, size string) error {
	// Compute size bounds
	sizeBounds, err := ParseSize(size)
	if err != nil {
		slog.Error(
			"failed parsing size",
			"handler", "GetDataSize",
			"size", size,
			"error_message", err.Error(),
		)
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	// Validate the size bounds
	if err := sizeBounds.Validate(); err != nil {
		slog.Error(
			"failed validating size",
			"handler", "GetDataSize",
			"size", size,
			"error_message", err.Error(),
		)

		return ctx.String(http.StatusBadRequest, err.Error())
	}

	// Write the data to the response
	payload := sizeBounds.Payload()

	return ctx.Blob(http.StatusOK, "application/octet-stream", payload)
}

// GetLatencyDuration is a handler that waits for the specified duration before responding.
//
// It parses the {duration} parameter from the URL to determine the wait time.
//
// The response returns a JSON object indicating the duration that was waited.
func (s *ServerImpl) GetLatencyDuration(ctx echo.Context, duration string) error {
	// parse the latency duration from the value passed within the URL
	latency, err := ParseLatency(duration)
	if err != nil {
		slog.Error(
			"failed parsing latency duration",
			"handler", "GetLatencyDuration",
			"duration", duration,
			"error_message", err.Error(),
		)

		// return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed parsing latency duration: %s", err))
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	// Validate the latency duration
	if err := latency.Validate(); err != nil {
		slog.Error(
			"failed validating latency duration",
			"handler", "GetLatencyDuration",
			"duration", duration,
			"error_message", err.Error(),
		)

		return ctx.String(http.StatusBadRequest, err.Error())
	}

	// Wait for the specified duration
	waited := latency.Wait()

	// Return a 200 response with the time waited
	return ctx.JSON(http.StatusOK, latencyResponse{
		Waited: waited,
	})
}

// NewLatencyResponse represents the response for the GetLatency handler.
// It contains the time that was waited.
type latencyResponse struct {
	// Waited indicates the duration that was waited before responding.
	Waited time.Duration `json:"waited"`
}

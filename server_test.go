package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetDataSizeHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		size       string
		wantStatus int
	}{
		{"valid size", "10kb", http.StatusOK},
		{"valid size bounds", "10kb-20kb", http.StatusOK},
		{"upper size bound less than lower size bound", "20kb-10kb", http.StatusBadRequest},
		{"invalid size", "invalid", http.StatusBadRequest},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/data/"+tt.size, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			handler := &ServerImpl{}

			assert.NoError(t, handler.GetDataSize(c, tt.size))
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestGetLatencyDurationHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		duration   string
		wantStatus int
	}{
		{"valid duration", "1s", http.StatusOK},
		{"valid duration bounds", "1s-2s", http.StatusOK},
		{"upper duration bound less than lower duration bound", "2s-1s", http.StatusBadRequest},
		{"invalid duration", "invalid", http.StatusBadRequest},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/latency/"+tt.duration, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			handler := &ServerImpl{}

			assert.NoError(t, handler.GetLatencyDuration(c, tt.duration))
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestGetResponse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		status         int
		contentType    string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "default status and content type",
			expectedStatus: http.StatusOK,
			expectedBody:   "Custom response body based on parameters",
		},
		{
			name:           "custom status code",
			status:         http.StatusAccepted,
			expectedStatus: http.StatusAccepted,
			expectedBody:   "Custom response body based on parameters",
		},
		{
			name:           "json content type",
			contentType:    "application/json",
			expectedStatus: http.StatusOK,
			expectedBody:   "{\"content_type\":\"application/json\",\"status\":200}\n",
		},
		{
			name:           "no content status code",
			status:         http.StatusNoContent,
			expectedStatus: http.StatusNoContent,
			expectedBody:   "",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := echo.New()

			req := httptest.NewRequest(http.MethodGet, "/response", nil)
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}

			responseParams := GetResponseParams{}
			if tt.status != 0 {
				responseParams.Status = &tt.status
			}

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handler := &ServerImpl{}
			err := handler.GetResponse(c, responseParams)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.Equal(t, tt.expectedBody, rec.Body.String())
		})
	}
}

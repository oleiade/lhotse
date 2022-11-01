package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", root).Methods(http.MethodOptions)
	r.HandleFunc("/latency/{duration}", latency).Methods(http.MethodOptions, http.MethodGet)
	r.HandleFunc("/data/{size}", data).Methods(http.MethodOptions, http.MethodGet)

	log.WithFields(log.Fields{"port": "3434"}).Info("starting server")

	err := http.ListenAndServe(":3434", r)
	if errors.Is(err, http.ErrServerClosed) {
		log.Info("server closed")
	} else if err != nil {
		log.WithFields(log.Fields{
			"error_message": err.Error(),
		}).Error("failed starting the server")
		os.Exit(1)
	}

	log.WithFields(log.Fields{"port": "3434"}).Info("server stopped")
}

type Endpoint struct {
	Method      string
	Path        string
	Description string
}

func root(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("allow", "GET, OPTIONS")
	w.WriteHeader(http.StatusNoContent)
}

func latency(w http.ResponseWriter, r *http.Request) {
	// Indicate the supported methods on this endpoint
	w.Header().Set("allow", "GET, OPTIONS")

	// If the OPTIONS method is used, return immediately
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// get the duration from the url
	duration, ok := mux.Vars(r)["duration"]
	if !ok {
		_, err := io.WriteString(w, "latency duration value is missing from the request")
		if err != nil {
			log.WithFields(log.Fields{"error_message": err.Error()}).Error("failed writing response")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("x-missing-field", "duration")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// parse the latency from the value passed within the URL
	latency, err := ParseLatency(duration)
	if err != nil {
		log.WithFields(log.Fields{"error_message": err.Error()}).Error("failed parsing latency")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Wait for the specified duration
	waited := latency.Wait()

	// Return a 200 response indicating the time waited
	_, err = io.WriteString(w, fmt.Sprintf("waited %s", waited))
	if err != nil {
		log.WithFields(log.Fields{"error_message": err.Error()}).Error("failed writing response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func data(w http.ResponseWriter, r *http.Request) {
	// Indicate the supported methods on this endpoint
	w.Header().Set("allow", "GET, OPTIONS")

	// If the OPTIONS method is used, return immediately
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// get the duration from the url
	size, ok := mux.Vars(r)["size"]
	if !ok {
		_, err := io.WriteString(w, "data size value is missing from the request")
		if err != nil {
			log.WithFields(log.Fields{"error_message": err.Error()}).Error("failed writing response")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("x-missing-field", "size")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sizeBounds, err := ParseSize(size)
	if err != nil {
		log.WithFields(log.Fields{"error_message": err.Error()}).Error("failed parsing size")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Write the data to the response
	payload := sizeBounds.Payload()

	_, err = w.Write(payload)
	if err != nil {
		log.WithFields(log.Fields{"error_message": err.Error()}).Error("failed writing response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

package handlers

import (
	"errors"
	"log"
	"net/http"
)

// HTTPError represents an error with an HTTP status code.
type HTTPError struct {
	Code    int
	Message string
}

func (e *HTTPError) Error() string {
	return e.Message
}

// AppHandlerFunc is an HTTP handler that returns an error.
type AppHandlerFunc func(w http.ResponseWriter, r *http.Request) error

// ServeHTTP implements http.Handler. It calls the underlying function
// and handles any returned error: HTTPError is sent with its status code,
// any other error is logged and sent as a 500.
func (f AppHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := f(w, r); err != nil {
		var httpErr *HTTPError
		if errors.As(err, &httpErr) {
			http.Error(w, httpErr.Message, httpErr.Code)
		} else {
			log.Printf("unexpected error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

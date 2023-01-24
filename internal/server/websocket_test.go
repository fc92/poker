package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServeHome(t *testing.T) {
	// Test non-GET request
	req, _ := http.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()
	serveHome(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Incorrect status code for non-GET request, got %d, expected %d", w.Code, http.StatusMethodNotAllowed)
	}

	// Test request to non-root URL
	req, _ = http.NewRequest(http.MethodGet, "/test", nil)
	w = httptest.NewRecorder()
	serveHome(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("Incorrect status code for non-root URL request, got %d, expected %d", w.Code, http.StatusNotFound)
	}

}

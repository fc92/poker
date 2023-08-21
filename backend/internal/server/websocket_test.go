package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockedResponseWriter struct {
	mock.Mock
}

func (m *MockedResponseWriter) Header() http.Header {
	return http.Header{}
}

func (m *MockedResponseWriter) Write([]byte) (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *MockedResponseWriter) WriteHeader(statusCode int) {
	m.Called(statusCode)
}
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

func TestStartServer(t *testing.T) {
	t.Run("should start server and listen for connections", func(t *testing.T) {
		originalListenAndServe := httpListenAndServe
		defer func() { httpListenAndServe = originalListenAndServe }()

		var called bool
		httpListenAndServe = func(addr string, handler http.Handler) error {
			called = true
			return errors.New("test error")
		}

		StartServer("localhost:8080")
		assert.True(t, called)
	})
}

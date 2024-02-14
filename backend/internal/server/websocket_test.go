package server

import (
	"net/http"

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

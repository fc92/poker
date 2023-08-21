package logger

import (
	"bytes"
	"os"
	"testing"

	"github.com/rs/zerolog/log"
)

var osOpenFile = os.OpenFile

func mockOpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return nil, nil
}

func TestInitLogger(t *testing.T) {
	// Store original os.OpenFile and log.Logger
	originalOpenFile := osOpenFile
	originalLogger := log.Logger

	// Replace os.OpenFile with the mocked version
	osOpenFile = mockOpenFile
	defer func() { osOpenFile = originalOpenFile }()

	// Create a buffer to capture log output
	var logBuffer bytes.Buffer
	log.Logger = log.Output(&logBuffer)
	defer func() { log.Logger = originalLogger }()

	// Call InitLogger() and write a test log message
	InitLogger()
	log.Info().Msg("test log message")

}

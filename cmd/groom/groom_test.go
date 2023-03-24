package main

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	// Replace osArgs and displayWelcome with mocked versions
	osArgs = func() []string {
		return []string{"groom", "localhost"}
	}
	displayWelcomeCalled := false
	displayWelcome = func(serverUrl string) {
		displayWelcomeCalled = true
		assert.Equal(t, "localhost", serverUrl, "Expected serverUrl to be localhost")
	}

	// Redirect log output to a buffer
	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer)
	defer log.SetOutput(os.Stderr)

	// Call main()
	main()

	// Check if displayWelcome was called
	assert.True(t, displayWelcomeCalled, "Expected displayWelcome to be called")
}

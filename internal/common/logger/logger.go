// logger shared between files of the package groom
package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	logFile = "/tmp/poker.log"
)

func InitLogger() {
	// Create a new lumberjack logger
	logFile := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    1,    // Max size in megabytes
		MaxBackups: 5,    // Max number of old log files to keep
		MaxAge:     30,   // Max number of days to keep old log files
		Compress:   true, // Whether to compress old log files
	}

	// Create a new zerolog logger with the lumberjack logger as the output
	log.Logger = zerolog.New(logFile).With().Timestamp().Logger()

	// Add the file and line number to the log context
	log.Logger = log.With().Caller().Logger()

	// Set global log level to debug
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

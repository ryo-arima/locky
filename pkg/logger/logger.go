package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/ryo-arima/locky/pkg/code"
	"github.com/ryo-arima/locky/pkg/config"
)

// Logger represents the application logger
type Logger struct {
	config *config.LoggerConfig
	output io.Writer
}

var globalLogger *Logger

// Initialize initializes the global logger
func Initialize(loggerConfig config.LoggerConfig) {
	globalLogger = &Logger{
		config: &loggerConfig,
		output: os.Stdout,
	}

	// Set output
	switch loggerConfig.Output {
	case "stderr":
		globalLogger.output = os.Stderr
	case "stdout", "":
		globalLogger.output = os.Stdout
	default:
		globalLogger.output = os.Stdout
	}
}

// GetLogger returns the global logger instance
func GetLogger() *Logger {
	if globalLogger == nil {
		// Initialize with default config if not initialized
		Initialize(config.LoggerConfig{
			Level:  "INFO",
			Output: "stdout",
		})
	}
	return globalLogger
}

// writeLogEntry writes a log entry to the output
func (rcvr *Logger) writeLogEntry(timestamp, level, code, requestID, message string, fields map[string]interface{}) {
	if rcvr == nil {
		return
	}

	if requestID == "" {
		requestID = "xxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	}

	// Format: [timestamp] [level] [code] [request_id] message
	fmt.Fprintf(rcvr.output, "[%s] [%s] [%s] [%s] %s",
		timestamp, level, code, requestID, message)

	// Only show fields for DEBUG level
	if len(fields) > 0 && level == "DEBUG  " {
		if fieldsJSON, err := json.Marshal(fields); err == nil {
			fmt.Fprintf(rcvr.output, " %s", string(fieldsJSON))
		}
	}

	fmt.Fprintln(rcvr.output)
}

// log is the internal logging method
func (rcvr *Logger) log(levelStr string, mc code.MCode, requestID string, message string, fields map[string]interface{}) {
	if rcvr == nil {
		return
	}

	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.000000000") + " UTC"
	codeStr := mc.PaddedCode()

	rcvr.writeLogEntry(timestamp, levelStr, codeStr, requestID, message, fields)
}

// Public logging methods
func Debug(mc code.MCode, message string, fields map[string]interface{}) {
	requestID := ""
	if fields != nil {
		if reqID, ok := fields["request_id"].(string); ok {
			requestID = reqID
		}
	}
	GetLogger().log("DEBUG  ", mc, requestID, message, fields)
}

func Info(mc code.MCode, requestID string, message string) {
	GetLogger().log("INFO   ", mc, requestID, message, nil)
}

func Warn(mc code.MCode, requestID string, message string) {
	GetLogger().log("WARN   ", mc, requestID, message, nil)
}

func Error(mc code.MCode, requestID string, message string) {
	GetLogger().log("ERROR  ", mc, requestID, message, nil)
}

func Fatal(mc code.MCode, requestID string, message string) {
	GetLogger().log("FATAL  ", mc, requestID, message, nil)
	os.Exit(1)
}

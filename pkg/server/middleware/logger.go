package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/locky/pkg/config"
)

// LogLevel represents the log level
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

// String returns string representation of log level with padding to 7 characters
func (rcvr LogLevel) String() string {
	switch rcvr {
	case DEBUG:
		return "DEBUG  "
	case INFO:
		return "INFO   "
	case WARN:
		return "WARN   "
	case ERROR:
		return "ERROR  "
	case FATAL:
		return "FATAL  "
	default:
		return "UNKNOWN"
	}
}

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Code      string                 `json:"code"`
	Component string                 `json:"component"`
	Service   string                 `json:"service"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	File      string                 `json:"file,omitempty"`
	Function  string                 `json:"function,omitempty"`
	Line      int                    `json:"line,omitempty"`
	TraceID   string                 `json:"trace_id,omitempty"`
	RequestID string                 `json:"request_id,omitempty"`
	UserID    string                 `json:"user_id,omitempty"`
	Error     string                 `json:"error,omitempty"`
}

// Logger represents the application logger
type Logger struct {
	config     *config.LoggerConfig
	level      LogLevel
	output     io.Writer
	baseConfig *config.BaseConfig
}

// NewLogger creates a new logger instance
func NewLogger(loggerConfig config.LoggerConfig, baseConfig *config.BaseConfig) config.LoggerInterface {
	logger := &Logger{
		config:     &loggerConfig,
		baseConfig: baseConfig,
		output:     os.Stdout,
	}

	// Set log level
	switch strings.ToUpper(loggerConfig.Level) {
	case "DEBUG":
		logger.level = DEBUG
	case "INFO":
		logger.level = INFO
	case "WARN":
		logger.level = WARN
	case "ERROR":
		logger.level = ERROR
	case "FATAL":
		logger.level = FATAL
	default:
		logger.level = INFO
	}

	// Set output
	switch loggerConfig.Output {
	case "stderr":
		logger.output = os.Stderr
	case "stdout", "":
		logger.output = os.Stdout
	default:
		// File output
		if file, err := os.OpenFile(loggerConfig.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666); err == nil {
			logger.output = file
		} else {
			logger.output = os.Stdout
			fmt.Fprintf(os.Stderr, "Failed to open log file %s: %v\n", loggerConfig.Output, err)
		}
	}

	return logger
}

// ToConfigMCode converts middleware.MCode to config.MCode
func ToConfigMCode(m MCode) config.MCode {
	return config.MCode{
		Code:    m.Code,
		Message: m.Message,
	}
}

// FormatWithOptional formats the message with optional additional message
func formatWithOptional(mcode config.MCode, optionalMessage string) string {
	if optionalMessage == "" {
		return mcode.Message
	}
	return fmt.Sprintf("%s: %s", mcode.Message, optionalMessage)
}

// log writes a log entry using MCode
func (l *Logger) log(level LogLevel, mcode config.MCode, optionalMessage string, fields map[string]interface{}) {
	if level < l.level {
		return
	}

	finalMessage := formatWithOptional(mcode, optionalMessage)

	// Get current time in UTC and format with " UTC" string (with space before UTC)
	now := time.Now().UTC()
	// Format: 2025-11-08T05:01:15.791560000 UTC
	timestamp := now.Format("2006-01-02T15:04:05.000000000") + " UTC"

	entry := LogEntry{
		Timestamp: timestamp,
		Level:     level.String(),
		Code:      mcode.PaddedCode(GetMaxCodeLength()),
		Component: l.config.Component,
		Service:   l.config.Service,
		Message:   finalMessage,
		Fields:    fields,
	}

	l.writeLogEntry(entry)
}

// writeLogEntry writes the actual log entry to output
func (l *Logger) writeLogEntry(entry LogEntry) {
	// Add caller information if enabled or DEBUG level
	if l.config.EnableCaller || l.level == DEBUG {
		if pc, file, line, ok := runtime.Caller(4); ok {
			entry.File = file
			entry.Line = line
			if fn := runtime.FuncForPC(pc); fn != nil {
				entry.Function = fn.Name()
			}
		}
	}

	// Extract common fields from fields map
	if entry.Fields != nil {
		if traceID, ok := entry.Fields["trace_id"].(string); ok {
			entry.TraceID = traceID
			delete(entry.Fields, "trace_id")
		}
		if requestID, ok := entry.Fields["request_id"].(string); ok {
			entry.RequestID = requestID
			delete(entry.Fields, "request_id")
		}
		if userID, ok := entry.Fields["user_id"].(string); ok {
			entry.UserID = userID
			delete(entry.Fields, "user_id")
		}
		if err, ok := entry.Fields["error"].(string); ok {
			entry.Error = err
			delete(entry.Fields, "error")
		}
		if err, ok := entry.Fields["error"].(error); ok {
			entry.Error = err.Error()
			delete(entry.Fields, "error")
		}
	}

	if l.config.Structured {
		// JSON format
		if jsonBytes, err := json.Marshal(entry); err == nil {
			fmt.Fprintln(l.output, string(jsonBytes))
		} else {
			// Fallback to simple format
			fmt.Fprintf(l.output, "[%s] [%s] [%s] %s\n",
				entry.Timestamp, entry.Level, entry.Code, entry.Message)
		}
	} else {
		// Human-readable format
		// Format: [timestamp] [level] [code] [request_id] message
		requestIDStr := entry.RequestID
		if requestIDStr == "" {
			requestIDStr = "xxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
		}
		fmt.Fprintf(l.output, "[%s] [%s] [%s] [%s] %s",
			entry.Timestamp, entry.Level, entry.Code, requestIDStr, entry.Message)
		// Only output fields for DEBUG level
		if len(entry.Fields) > 0 && entry.Level == "DEBUG  " {
			if fieldsJSON, err := json.Marshal(entry.Fields); err == nil {
				fmt.Fprintf(l.output, " %s", string(fieldsJSON))
			}
		}
		fmt.Fprintln(l.output)
	}
}

// DEBUG logs a debug message using MCode
func (l *Logger) DEBUG(mcode config.MCode, optionalMessage string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(DEBUG, mcode, optionalMessage, f)
}

// INFO logs an info message using MCode
func (l *Logger) INFO(mcode config.MCode, optionalMessage string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(INFO, mcode, optionalMessage, f)
}

// WARN logs a warning message using MCode
func (l *Logger) WARN(mcode config.MCode, optionalMessage string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(WARN, mcode, optionalMessage, f)
}

// ERROR logs an error message using MCode
func (l *Logger) ERROR(mcode config.MCode, optionalMessage string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(ERROR, mcode, optionalMessage, f)
}

// FATAL logs a fatal message using MCode and exits
func (l *Logger) FATAL(mcode config.MCode, optionalMessage string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(FATAL, mcode, optionalMessage, f)
	os.Exit(1)
}

// GinLoggerWriter wraps our custom logger to implement io.Writer for Gin
type GinLoggerWriter struct {
	logger config.LoggerInterface
}

// Write implements io.Writer interface for Gin logging
func (w *GinLoggerWriter) Write(p []byte) (n int, err error) {
	msg := string(p)
	// Remove trailing newline if present
	if len(msg) > 0 && msg[len(msg)-1] == '\n' {
		msg = msg[:len(msg)-1]
	}

	// Skip empty messages
	if msg == "" {
		return len(p), nil
	}

	// Parse Gin HTTP request logs (format: "[GIN] 2025/11/08 - 14:04:27 | 400 |    2.531667ms |       127.0.0.1 | POST     "/v1/public/user"")
	// These are from gin.Default() logger, which we skip because LoggerWithConfig middleware handles request logging
	if strings.HasPrefix(msg, "[GIN] ") && strings.Contains(msg, " | ") {
		// Skip Gin's default HTTP request logs (they contain " | " separator)
		return len(p), nil
	}

	// Parse Gin debug/warning/error messages - use empty message to avoid redundancy
	mcode := ToConfigMCode(MCode{"GINLOG", ""})

	// Clean up the message - remove [GIN-debug], [GIN-warning], etc prefixes and redundant [WARNING], [ERROR] text
	cleanMsg := msg
	cleanMsg = strings.TrimPrefix(cleanMsg, "[GIN-debug] ")
	cleanMsg = strings.TrimPrefix(cleanMsg, "[GIN-warning] ")
	cleanMsg = strings.TrimPrefix(cleanMsg, "[GIN-error] ")
	cleanMsg = strings.TrimPrefix(cleanMsg, "[GIN] ")

	// Remove redundant [WARNING], [ERROR] brackets from the message
	cleanMsg = strings.ReplaceAll(cleanMsg, "[WARNING] ", "")
	cleanMsg = strings.ReplaceAll(cleanMsg, "[ERROR] ", "")
	cleanMsg = strings.ReplaceAll(cleanMsg, "[DEBUG] ", "")

	// Compact route registration logs: ": DELETE /v1/..." -> "DELETE /v1/..."
	cleanMsg = strings.TrimPrefix(cleanMsg, ": ")

	// Determine log level based on message content
	if strings.Contains(msg, "[WARNING]") || strings.Contains(msg, "[GIN-warning]") || strings.Contains(msg, "WARNING") {
		w.logger.WARN(mcode, cleanMsg, nil)
	} else if strings.Contains(msg, "[ERROR]") || strings.Contains(msg, "[GIN-error]") || strings.Contains(msg, "ERROR") {
		w.logger.ERROR(mcode, cleanMsg, nil)
	} else if strings.Contains(msg, "[GIN-debug]") || strings.Contains(msg, "[debug]") {
		w.logger.DEBUG(mcode, cleanMsg, nil)
	} else {
		w.logger.INFO(mcode, cleanMsg, nil)
	}

	return len(p), nil
}

// NewGinLoggerWriter creates a new GinLoggerWriter
func NewGinLoggerWriter(logger config.LoggerInterface) *GinLoggerWriter {
	return &GinLoggerWriter{logger: logger}
}

// LoggerWithConfig returns a Gin middleware for HTTP request logging
func LoggerWithConfig(conf config.BaseConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log request details
		end := time.Now()
		latency := end.Sub(start)

		if raw != "" {
			path = path + "?" + raw
		}

		requestID := GetRequestID(c)
		fields := map[string]interface{}{
			"request_id": requestID,
			"method":     c.Request.Method,
			"path":       path,
			"client_ip":  c.ClientIP(),
			"status":     c.Writer.Status(),
			"latency_ms": latency.Milliseconds(),
		}

		// Add error if exists
		if len(c.Errors) > 0 {
			fields["error"] = c.Errors.Last().Error()
		}

		// Log based on status code
		status := c.Writer.Status()
		// Format: "METHOD /path STATUS" (e.g., "POST /v1/public/user 200")
		requestInfo := fmt.Sprintf("%s %s %d", c.Request.Method, path, status)

		if status >= 500 {
			conf.Logger.ERROR(ToConfigMCode(MLWC5), requestInfo, fields)
		} else if status >= 400 {
			conf.Logger.WARN(ToConfigMCode(MLWC4), requestInfo, fields)
		} else {
			conf.Logger.INFO(ToConfigMCode(MLWC3), requestInfo, fields)
		}
	}
}

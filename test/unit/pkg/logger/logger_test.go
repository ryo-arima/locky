package logger_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/ryo-arima/locky/pkg/code"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitialize(t *testing.T) {
	tests := []struct {
		name   string
		config config.LoggerConfig
	}{
		{
			name: "Initialize with stdout",
			config: config.LoggerConfig{
				Component:    "test",
				Service:      "test-service",
				Level:        "DEBUG",
				Structured:   false,
				EnableCaller: true,
				Output:       "stdout",
			},
		},
		{
			name: "Initialize with stderr",
			config: config.LoggerConfig{
				Component:    "test",
				Service:      "test-service",
				Level:        "INFO",
				Structured:   false,
				EnableCaller: false,
				Output:       "stderr",
			},
		},
		{
			name: "Initialize with default output",
			config: config.LoggerConfig{
				Component:    "test",
				Service:      "test-service",
				Level:        "WARN",
				Structured:   false,
				EnableCaller: true,
				Output:       "",
			},
		},
		{
			name: "Initialize with unknown output defaults to stdout",
			config: config.LoggerConfig{
				Component:    "test",
				Service:      "test-service",
				Level:        "ERROR",
				Structured:   false,
				EnableCaller: true,
				Output:       "unknown",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger.Initialize(tt.config)
			l := logger.GetLogger()
			assert.NotNil(t, l, "Logger should not be nil after initialization")
		})
	}
}

func TestGetLogger(t *testing.T) {
	// Reset logger by initializing with a known config
	logger.Initialize(config.LoggerConfig{
		Level:  "INFO",
		Output: "stdout",
	})

	l := logger.GetLogger()
	assert.NotNil(t, l, "GetLogger should return non-nil logger")

	l2 := logger.GetLogger()
	assert.Equal(t, l, l2, "GetLogger should return the same instance")
}

func TestDebug(t *testing.T) {
	// Initialize logger
	logger.Initialize(config.LoggerConfig{
		Level:  "DEBUG",
		Output: "stdout",
	})

	// Test with fields
	fields := map[string]interface{}{
		"request_id": "test-request-123",
		"user_id":    "user-456",
		"action":     "test_action",
	}

	// This should not panic
	assert.NotPanics(t, func() {
		logger.Debug(code.SM1, "Test debug message", fields)
	})

	// Test with nil fields
	assert.NotPanics(t, func() {
		logger.Debug(code.SM2, "Test debug without fields", nil)
	})

	// Test with empty fields
	assert.NotPanics(t, func() {
		logger.Debug(code.SM3, "Test debug with empty fields", map[string]interface{}{})
	})
}

func TestInfo(t *testing.T) {
	logger.Initialize(config.LoggerConfig{
		Level:  "INFO",
		Output: "stdout",
	})

	assert.NotPanics(t, func() {
		logger.Info(code.MLWC1, "test-req-123", "Test info message")
	})

	assert.NotPanics(t, func() {
		logger.Info(code.MLWC2, "", "Test info without request ID")
	})
}

func TestWarn(t *testing.T) {
	logger.Initialize(config.LoggerConfig{
		Level:  "WARN",
		Output: "stdout",
	})

	assert.NotPanics(t, func() {
		logger.Warn(code.MLWC3, "test-req-123", "Test warning message")
	})

	assert.NotPanics(t, func() {
		logger.Warn(code.MLWC4, "", "Test warning without request ID")
	})
}

func TestError(t *testing.T) {
	logger.Initialize(config.LoggerConfig{
		Level:  "ERROR",
		Output: "stdout",
	})

	assert.NotPanics(t, func() {
		logger.Error(code.CNDBC3, "test-req-123", "Test error message")
	})

	assert.NotPanics(t, func() {
		logger.Error(code.MLWC5, "", "Test error without request ID")
	})
}

func TestLoggerOutputFormat(t *testing.T) {
	// Capture output by testing the logger behavior
	// Note: Since we can't directly capture the output in the current implementation,
	// we'll just verify the logger doesn't panic with various inputs

	logger.Initialize(config.LoggerConfig{
		Level:  "DEBUG",
		Output: "stdout",
	})

	testCases := []struct {
		name      string
		logFunc   func()
		shouldRun bool
	}{
		{
			name: "Debug with all fields",
			logFunc: func() {
				logger.Debug(code.RCHK1, "Debug test", map[string]interface{}{
					"request_id": "req-123",
					"key1":       "value1",
					"key2":       123,
				})
			},
			shouldRun: true,
		},
		{
			name: "Info with request ID",
			logFunc: func() {
				logger.Info(code.RURP1, "req-456", "Info test")
			},
			shouldRun: true,
		},
		{
			name: "Warn with empty request ID",
			logFunc: func() {
				logger.Warn(code.RUCR1, "", "Warn test")
			},
			shouldRun: true,
		},
		{
			name: "Error with special characters in message",
			logFunc: func() {
				logger.Error(code.RUUP1, "req-789", "Error test with special chars: @#$%^&*()")
			},
			shouldRun: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.shouldRun {
				assert.NotPanics(t, tc.logFunc)
			}
		})
	}
}

func TestLoggerWithNilConfig(t *testing.T) {
	// Test that GetLogger initializes with default config if not initialized
	// We can't easily reset the global logger, but we can verify it doesn't panic
	l := logger.GetLogger()
	assert.NotNil(t, l)
	
	// Should not panic even with nil fields
	assert.NotPanics(t, func() {
		logger.Debug(code.SM1, "test", nil)
		logger.Info(code.SM2, "req", "test")
		logger.Warn(code.SM3, "req", "test")
		logger.Error(code.MLWC1, "req", "test")
	})
}

func TestMCodePaddedCodeIntegration(t *testing.T) {
	// Test that MCode.PaddedCode() works correctly with logger
	logger.Initialize(config.LoggerConfig{
		Level:  "DEBUG",
		Output: "stdout",
	})

	testCodes := []code.MCode{
		code.SM1,
		code.CNDBC1,
		code.UCPCU0,
		code.GINLOG,
	}

	for _, mc := range testCodes {
		t.Run(mc.Code, func(t *testing.T) {
			paddedCode := mc.PaddedCode()
			assert.GreaterOrEqual(t, len(paddedCode), len(mc.Code))
			assert.True(t, strings.HasPrefix(paddedCode, mc.Code))

			// Test that logging with this code doesn't panic
			assert.NotPanics(t, func() {
				logger.Debug(mc, "Test message", map[string]interface{}{
					"request_id": "test-123",
				})
			})
		})
	}
}

func TestLoggerJSONFieldsFormat(t *testing.T) {
	logger.Initialize(config.LoggerConfig{
		Level:  "DEBUG",
		Output: "stdout",
	})

	// Test that complex nested structures in fields don't cause issues
	complexFields := map[string]interface{}{
		"request_id": "req-complex-123",
		"nested": map[string]interface{}{
			"level1": map[string]interface{}{
				"level2": "deep value",
			},
		},
		"array": []string{"item1", "item2", "item3"},
		"number": 12345,
		"bool":   true,
	}

	assert.NotPanics(t, func() {
		logger.Debug(code.UUGU1, "Complex fields test", complexFields)
	})

	// Verify the fields can be marshaled to JSON
	jsonData, err := json.Marshal(complexFields)
	require.NoError(t, err)
	assert.NotEmpty(t, jsonData)
}

// TestLoggerConcurrency tests that logger is safe for concurrent use
func TestLoggerConcurrency(t *testing.T) {
	logger.Initialize(config.LoggerConfig{
		Level:  "DEBUG",
		Output: "stdout",
	})

	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(id int) {
			logger.Debug(code.SM1, "Concurrent test", map[string]interface{}{
				"request_id": "concurrent-req",
				"goroutine":  id,
			})
			logger.Info(code.SM2, "concurrent-req", "Concurrent info")
			logger.Warn(code.SM3, "concurrent-req", "Concurrent warn")
			logger.Error(code.MLWC1, "concurrent-req", "Concurrent error")
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}

// Benchmark tests
func BenchmarkLogger_Debug(b *testing.B) {
	logger.Initialize(config.LoggerConfig{
		Level:  "DEBUG",
		Output: "stdout",
	})

	fields := map[string]interface{}{
		"request_id": "bench-req-123",
		"key1":       "value1",
		"key2":       123,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Debug(code.SM1, "Benchmark test", fields)
	}
}

func BenchmarkLogger_Info(b *testing.B) {
	logger.Initialize(config.LoggerConfig{
		Level:  "INFO",
		Output: "stdout",
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info(code.SM2, "bench-req-123", "Benchmark test")
	}
}

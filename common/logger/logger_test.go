package logger

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewLogger(t *testing.T) {
	// Test with stdout only
	logger, err := NewLogger(INFO, "")
	if err != nil {
		t.Errorf("NewLogger() error = %v", err)
	}
	if logger == nil {
		t.Error("NewLogger() returned nil logger")
	}

	// Test with file output
	tmpDir := os.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")
	logger, err = NewLogger(INFO, logFile)
	if err != nil {
		t.Errorf("NewLogger() error = %v", err)
	}
	if logger == nil {
		t.Error("NewLogger() returned nil logger")
	}
	if logger.file == nil {
		t.Error("NewLogger() did not create log file")
	}
	os.Remove(logFile)
}

func TestLogLevels(t *testing.T) {
	// Create a buffer to capture log output
	var buf bytes.Buffer
	logger := &Logger{
		out:    &buf,
		level:  DEBUG,
		logger: log.New(&buf, "", 0),
	}
	defaultLogger = logger

	tests := []struct {
		name    string
		level   Level
		logFunc func(string, ...interface{})
		want    string
	}{
		{
			name:    "Debug level",
			level:   DEBUG,
			logFunc: Debug,
			want:    "[DEBUG]",
		},
		{
			name:    "Info level",
			level:   INFO,
			logFunc: Info,
			want:    "[INFO]",
		},
		{
			name:    "Warn level",
			level:   WARN,
			logFunc: Warn,
			want:    "[WARN]",
		},
		{
			name:    "Error level",
			level:   ERROR,
			logFunc: Error,
			want:    "[ERROR]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.logFunc("test message")
			if !strings.Contains(buf.String(), tt.want) {
				t.Errorf("Log output does not contain %s: %s", tt.want, buf.String())
			}
		})
	}
}

func TestLoggerLevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		out:    &buf,
		level:  INFO,
		logger: log.New(&buf, "", 0),
	}
	defaultLogger = logger

	// Debug message should not be logged when level is INFO
	Debug("debug message")
	if buf.String() != "" {
		t.Error("Debug message was logged when level is INFO")
	}

	// Info message should be logged
	buf.Reset()
	Info("info message")
	if buf.String() == "" {
		t.Error("Info message was not logged when level is INFO")
	}
}

func TestLoggerClose(t *testing.T) {
	// Create a temporary log file
	tmpDir := os.TempDir()
	logFile := filepath.Join(tmpDir, "test_close.log")
	logger, _ := NewLogger(INFO, logFile)
	defaultLogger = logger

	// Write something to the log
	Info("test message")

	// Close the logger
	if err := Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}

	// Try to write to the closed file
	file := logger.file
	if _, err := file.Write([]byte("test")); err == nil {
		t.Error("Write to closed log file should fail")
	}

	os.Remove(logFile)
}

func TestLoggerConcurrency(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		out:    &buf,
		level:  DEBUG,
		logger: log.New(&buf, "", 0),
	}
	defaultLogger = logger

	// Test concurrent logging
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 100; j++ {
				Info("test message from goroutine %d", id)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to finish
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify that all messages were logged
	output := buf.String()
	count := strings.Count(output, "[INFO]")
	if count != 1000 {
		t.Errorf("Expected 1000 log messages, got %d", count)
	}
}

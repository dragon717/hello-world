package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// Level type
type Level int

const (
	// Log levels
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

var levelNames = map[Level]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
	FATAL: "FATAL",
}

// Logger represents the logger instance
type Logger struct {
	mu       sync.Mutex
	out      io.Writer
	file     *os.File
	level    Level
	logger   *log.Logger
	filename string
}

var defaultLogger *Logger

// InitLogger initializes the default logger
func InitLogger(level Level, filename string) error {
	logger, err := NewLogger(level, filename)
	if err != nil {
		return err
	}
	defaultLogger = logger
	return nil
}

// NewLogger creates a new logger instance
func NewLogger(level Level, filename string) (*Logger, error) {
	var writer io.Writer
	var file *os.File
	var err error

	if filename != "" {
		// Ensure directory exists
		dir := filepath.Dir(filename)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %v", err)
		}

		file, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %v", err)
		}
		writer = io.MultiWriter(os.Stdout, file)
	} else {
		writer = os.Stdout
	}

	return &Logger{
		out:      writer,
		file:     file,
		level:    level,
		logger:   log.New(writer, "", 0),
		filename: filename,
	}, nil
}

// log formats and writes the log message
func (l *Logger) log(level Level, format string, v ...interface{}) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// Get caller information
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}

	// Format the message
	msg := fmt.Sprintf(format, v...)
	timestamp := time.Now().Format("2006/01/02 15:04:05.000")
	logLine := fmt.Sprintf("[%s] [%s] %s:%d: %s\n",
		levelNames[level],
		timestamp,
		filepath.Base(file),
		line,
		msg)

	l.logger.Output(0, logLine)
}

// Debug logs a debug message
func Debug(format string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.log(DEBUG, format, v...)
	}
}

// Info logs an info message
func Info(format string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.log(INFO, format, v...)
	}
}

// Warn logs a warning message
func Warn(format string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.log(WARN, format, v...)
	}
}

// Error logs an error message
func Error(format string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.log(ERROR, format, v...)
	}
}

// Fatal logs a fatal message and exits
func Fatal(format string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.log(FATAL, format, v...)
	}
	os.Exit(1)
}

// Close closes the logger
func Close() error {
	if defaultLogger != nil && defaultLogger.file != nil {
		return defaultLogger.file.Close()
	}
	return nil
}

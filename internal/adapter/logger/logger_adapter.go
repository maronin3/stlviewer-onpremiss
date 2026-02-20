package logger

import (
	"fmt"
	"log"
	"os"
)

// FileLogger implements the Logger interface using file-based logging
type FileLogger struct {
	logger *log.Logger
	file   *os.File
}

// NewFileLogger creates a new file-based logger
func NewFileLogger(filePath string) (*FileLogger, error) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	logger := log.New(file, "", log.LstdFlags)

	return &FileLogger{
		logger: logger,
		file:   file,
	}, nil
}

// Info logs an informational message
func (l *FileLogger) Info(msg string, args ...interface{}) {
	l.logger.Printf("[INFO] "+msg, args...)
}

// Error logs an error message
func (l *FileLogger) Error(msg string, args ...interface{}) {
	l.logger.Printf("[ERROR] "+msg, args...)
}

// Debug logs a debug message
func (l *FileLogger) Debug(msg string, args ...interface{}) {
	l.logger.Printf("[DEBUG] "+msg, args...)
}

// Close closes the log file
func (l *FileLogger) Close() error {
	if l.file != nil {
		err := l.file.Close()
		l.file = nil
		return err
	}
	return nil
}

// ConsoleLogger implements the Logger interface using console output
type ConsoleLogger struct {
	logger *log.Logger
}

// NewConsoleLogger creates a new console-based logger
func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

// Info logs an informational message
func (l *ConsoleLogger) Info(msg string, args ...interface{}) {
	l.logger.Printf("[INFO] "+msg, args...)
}

// Error logs an error message
func (l *ConsoleLogger) Error(msg string, args ...interface{}) {
	l.logger.Printf("[ERROR] "+msg, args...)
}

// Debug logs a debug message
func (l *ConsoleLogger) Debug(msg string, args ...interface{}) {
	l.logger.Printf("[DEBUG] "+msg, args...)
}

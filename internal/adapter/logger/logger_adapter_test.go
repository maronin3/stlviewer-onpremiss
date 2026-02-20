package logger

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileLogger_Creation(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")

	logger, err := NewFileLogger(logPath)
	require.NoError(t, err)
	require.NotNil(t, logger)
	defer logger.Close()

	// Verify file was created
	_, err = os.Stat(logPath)
	assert.NoError(t, err)
}

func TestFileLogger_Info(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "info.log")

	logger, err := NewFileLogger(logPath)
	require.NoError(t, err)
	defer logger.Close()

	logger.Info("Test message: %s", "hello")

	// Read log file
	content, err := os.ReadFile(logPath)
	require.NoError(t, err)

	logContent := string(content)
	assert.Contains(t, logContent, "[INFO]")
	assert.Contains(t, logContent, "Test message: hello")
}

func TestFileLogger_Error(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "error.log")

	logger, err := NewFileLogger(logPath)
	require.NoError(t, err)
	defer logger.Close()

	logger.Error("Error occurred: %v", "something went wrong")

	// Read log file
	content, err := os.ReadFile(logPath)
	require.NoError(t, err)

	logContent := string(content)
	assert.Contains(t, logContent, "[ERROR]")
	assert.Contains(t, logContent, "Error occurred: something went wrong")
}

func TestFileLogger_Debug(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "debug.log")

	logger, err := NewFileLogger(logPath)
	require.NoError(t, err)
	defer logger.Close()

	logger.Debug("Debug info: %d", 42)

	// Read log file
	content, err := os.ReadFile(logPath)
	require.NoError(t, err)

	logContent := string(content)
	assert.Contains(t, logContent, "[DEBUG]")
	assert.Contains(t, logContent, "Debug info: 42")
}

func TestFileLogger_MultipleLogs(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "multiple.log")

	logger, err := NewFileLogger(logPath)
	require.NoError(t, err)
	defer logger.Close()

	logger.Info("First message")
	logger.Error("Second message")
	logger.Debug("Third message")

	// Read log file
	file, err := os.Open(logPath)
	require.NoError(t, err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	assert.Len(t, lines, 3)
	assert.True(t, strings.Contains(lines[0], "[INFO]") && strings.Contains(lines[0], "First message"))
	assert.True(t, strings.Contains(lines[1], "[ERROR]") && strings.Contains(lines[1], "Second message"))
	assert.True(t, strings.Contains(lines[2], "[DEBUG]") && strings.Contains(lines[2], "Third message"))
}

func TestFileLogger_Close(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "close.log")

	logger, err := NewFileLogger(logPath)
	require.NoError(t, err)

	err = logger.Close()
	assert.NoError(t, err)

	// Closing again should not panic (file will be nil)
	// Note: May return error on Windows, which is acceptable
	logger.Close()
}

func TestConsoleLogger_Creation(t *testing.T) {
	logger := NewConsoleLogger()
	assert.NotNil(t, logger)
}

func TestConsoleLogger_Logging(t *testing.T) {
	// Just verify it doesn't panic
	logger := NewConsoleLogger()
	assert.NotPanics(t, func() {
		logger.Info("Info message")
		logger.Error("Error message")
		logger.Debug("Debug message")
	})
}

package infrastructure

import (
	"fmt"
	"os"
	"path/filepath"
)

// FileSystem provides file system utilities
type FileSystem struct {
	appPath string
}

// NewFileSystem creates a new filesystem utility
func NewFileSystem() (*FileSystem, error) {
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return nil, fmt.Errorf("failed to get application path: %w", err)
	}

	return &FileSystem{
		appPath: path,
	}, nil
}

// GetAppPath returns the application directory path
func (fs *FileSystem) GetAppPath() string {
	return fs.appPath
}

// GetConfigPath returns the path to the config file
func (fs *FileSystem) GetConfigPath() string {
	return filepath.Join(fs.appPath, "server.ini")
}

// GetSystemLogPath returns the path to the system log file
func (fs *FileSystem) GetSystemLogPath() string {
	return filepath.Join(fs.appPath, "system.log")
}

// GetServerLogPath returns the path to the server log file
func (fs *FileSystem) GetServerLogPath() string {
	return filepath.Join(fs.appPath, "server.log")
}

// GetDistPath returns the path to the frontend dist directory
func (fs *FileSystem) GetDistPath() string {
	return filepath.Join(fs.appPath, "dist")
}

// Exists checks if a file or directory exists
func (fs *FileSystem) Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// EnsureDir ensures a directory exists, creating it if necessary
func (fs *FileSystem) EnsureDir(path string) error {
	if !fs.Exists(path) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}
	return nil
}

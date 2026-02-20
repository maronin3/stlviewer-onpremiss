package domain

import (
	"fmt"
)

// Config represents the application configuration entity
// This is a domain entity with business rules and validation
type Config struct {
	// Application settings
	Version string
	DevMode bool

	// Server settings
	IP   string
	Port int

	// WebView settings
	Title  string
	Width  int
	Height int
	URL    string
}

// NewDefaultConfig creates a new Config with default values
func NewDefaultConfig() *Config {
	return &Config{
		Version: "1.0.0",
		DevMode: false,
		IP:      "127.0.0.1",
		Port:    3000,
		Title:   "STLViewer",
		Width:   800,
		Height:  600,
		URL:     "http://127.0.0.1:3000",
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.IP == "" {
		return fmt.Errorf("IP address cannot be empty")
	}

	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535, got %d", c.Port)
	}

	if c.Title == "" {
		return fmt.Errorf("title cannot be empty")
	}

	if c.Width <= 0 {
		return fmt.Errorf("width must be positive, got %d", c.Width)
	}

	if c.Height <= 0 {
		return fmt.Errorf("height must be positive, got %d", c.Height)
	}

	if c.URL == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	return nil
}

// GetServerAddress returns the formatted server address
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.IP, c.Port)
}

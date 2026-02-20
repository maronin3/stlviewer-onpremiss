package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_DefaultValues(t *testing.T) {
	config := NewDefaultConfig()

	assert.Equal(t, "1.0.0", config.Version)
	assert.Equal(t, false, config.DevMode)
	assert.Equal(t, "127.0.0.1", config.IP)
	assert.Equal(t, 3000, config.Port)
	assert.Equal(t, "STLViewer", config.Title)
	assert.Equal(t, 800, config.Width)
	assert.Equal(t, 600, config.Height)
	assert.Equal(t, "http://127.0.0.1:3000", config.URL)
}

func TestConfig_Validate_ValidConfig(t *testing.T) {
	config := NewDefaultConfig()
	err := config.Validate()
	assert.NoError(t, err)
}

func TestConfig_Validate_InvalidPort(t *testing.T) {
	tests := []struct {
		name string
		port int
	}{
		{"negative port", -1},
		{"zero port", 0},
		{"port too high", 65536},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := NewDefaultConfig()
			config.Port = tt.port
			err := config.Validate()
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "port")
		})
	}
}

func TestConfig_Validate_InvalidDimensions(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
	}{
		{"negative width", -1, 600},
		{"zero width", 0, 600},
		{"negative height", 800, -1},
		{"zero height", 800, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := NewDefaultConfig()
			config.Width = tt.width
			config.Height = tt.height
			err := config.Validate()
			assert.Error(t, err)
		})
	}
}

func TestConfig_Validate_EmptyIP(t *testing.T) {
	config := NewDefaultConfig()
	config.IP = ""
	err := config.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "IP")
}

func TestConfig_Validate_EmptyURL(t *testing.T) {
	config := NewDefaultConfig()
	config.URL = ""
	err := config.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "URL")
}

func TestConfig_Validate_EmptyTitle(t *testing.T) {
	config := NewDefaultConfig()
	config.Title = ""
	err := config.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "title")
}

func TestConfig_GetServerAddress(t *testing.T) {
	config := NewDefaultConfig()
	assert.Equal(t, "127.0.0.1:3000", config.GetServerAddress())

	config.IP = "0.0.0.0"
	config.Port = 8080
	assert.Equal(t, "0.0.0.0:8080", config.GetServerAddress())
}

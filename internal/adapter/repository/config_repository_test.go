package repository

import (
	"context"
	"os"
	"path/filepath"
	"stlviewer/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestINIConfigRepository_Exists_FileExists(t *testing.T) {
	// Create temp file
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.ini")
	err := os.WriteFile(filePath, []byte(""), 0644)
	require.NoError(t, err)

	repo := NewINIConfigRepository(filePath)
	ctx := context.Background()

	exists, err := repo.Exists(ctx)
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestINIConfigRepository_Exists_FileNotExists(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "nonexistent.ini")

	repo := NewINIConfigRepository(filePath)
	ctx := context.Background()

	exists, err := repo.Exists(ctx)
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestINIConfigRepository_Save_Success(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "config.ini")

	repo := NewINIConfigRepository(filePath)
	ctx := context.Background()

	config := domain.NewDefaultConfig()
	config.Port = 8080
	config.DevMode = true

	err := repo.Save(ctx, config)
	assert.NoError(t, err)

	// Verify file was created
	_, err = os.Stat(filePath)
	assert.NoError(t, err)
}

func TestINIConfigRepository_Load_Success(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "config.ini")

	// First save a config
	repo := NewINIConfigRepository(filePath)
	ctx := context.Background()

	originalConfig := domain.NewDefaultConfig()
	originalConfig.Port = 9000
	originalConfig.DevMode = true
	originalConfig.Title = "CustomTitle"

	err := repo.Save(ctx, originalConfig)
	require.NoError(t, err)

	// Now load it back
	loadedConfig, err := repo.Load(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, loadedConfig)
	assert.Equal(t, 9000, loadedConfig.Port)
	assert.Equal(t, true, loadedConfig.DevMode)
	assert.Equal(t, "CustomTitle", loadedConfig.Title)
	assert.Equal(t, "1.0.0", loadedConfig.Version)
}

func TestINIConfigRepository_Load_FileNotExists(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "nonexistent.ini")

	repo := NewINIConfigRepository(filePath)
	ctx := context.Background()

	config, err := repo.Load(ctx)
	assert.Error(t, err)
	assert.Nil(t, config)
}

func TestINIConfigRepository_SaveAndLoad_RoundTrip(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "roundtrip.ini")

	repo := NewINIConfigRepository(filePath)
	ctx := context.Background()

	// Create custom config
	original := &domain.Config{
		Version: "2.0.0",
		DevMode: true,
		IP:      "0.0.0.0",
		Port:    8888,
		Title:   "Test Viewer",
		Width:   1024,
		Height:  768,
		URL:     "http://localhost:8888",
	}

	// Save
	err := repo.Save(ctx, original)
	require.NoError(t, err)

	// Load
	loaded, err := repo.Load(ctx)
	require.NoError(t, err)

	// Compare
	assert.Equal(t, original.Version, loaded.Version)
	assert.Equal(t, original.DevMode, loaded.DevMode)
	assert.Equal(t, original.IP, loaded.IP)
	assert.Equal(t, original.Port, loaded.Port)
	assert.Equal(t, original.Title, loaded.Title)
	assert.Equal(t, original.Width, loaded.Width)
	assert.Equal(t, original.Height, loaded.Height)
	assert.Equal(t, original.URL, loaded.URL)
}

package server

import (
	"context"
	"path/filepath"
	"stlviewer/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFiberServer_Creation(t *testing.T) {
	tmpDir := t.TempDir()
	distPath := filepath.Join(tmpDir, "dist")
	logPath := filepath.Join(tmpDir, "server.log")

	server := NewFiberServer(distPath, logPath)
	assert.NotNil(t, server)
}

func TestFiberServer_Start_Success(t *testing.T) {
	tmpDir := t.TempDir()
	distPath := tmpDir
	logPath := filepath.Join(tmpDir, "server.log")

	server := NewFiberServer(distPath, logPath)
	ctx := context.Background()

	config := domain.NewDefaultConfig()
	config.Port = 13000 // Use high port to avoid conflicts

	err := server.Start(ctx, config)
	require.NoError(t, err)

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Cleanup
	err = server.Stop(ctx)
	assert.NoError(t, err)
}

func TestFiberServer_Start_DevMode(t *testing.T) {
	tmpDir := t.TempDir()
	distPath := tmpDir
	logPath := filepath.Join(tmpDir, "server_devmode.log")

	server := NewFiberServer(distPath, logPath)
	ctx := context.Background()

	config := domain.NewDefaultConfig()
	config.Port = 13001
	config.DevMode = true

	err := server.Start(ctx, config)
	require.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	err = server.Stop(ctx)
	assert.NoError(t, err)
}

func TestFiberServer_Stop_WithoutStart(t *testing.T) {
	tmpDir := t.TempDir()
	distPath := tmpDir
	logPath := filepath.Join(tmpDir, "server.log")

	server := NewFiberServer(distPath, logPath)
	ctx := context.Background()

	// Should not error when stopping without starting
	err := server.Stop(ctx)
	assert.NoError(t, err)
}

func TestFiberServer_StartStop_Multiple(t *testing.T) {
	tmpDir := t.TempDir()
	distPath := tmpDir
	logPath := filepath.Join(tmpDir, "server.log")

	server := NewFiberServer(distPath, logPath)
	ctx := context.Background()

	config := domain.NewDefaultConfig()
	config.Port = 13002

	// Start
	err := server.Start(ctx, config)
	require.NoError(t, err)
	time.Sleep(100 * time.Millisecond)

	// Stop
	err = server.Stop(ctx)
	require.NoError(t, err)

	// Start again with new instance
	server2 := NewFiberServer(distPath, logPath)
	err = server2.Start(ctx, config)
	require.NoError(t, err)
	time.Sleep(100 * time.Millisecond)

	err = server2.Stop(ctx)
	assert.NoError(t, err)
}

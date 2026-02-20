package infrastructure

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileSystem_Creation(t *testing.T) {
	fs, err := NewFileSystem()
	require.NoError(t, err)
	assert.NotNil(t, fs)
	assert.NotEmpty(t, fs.GetAppPath())
}

func TestFileSystem_GetConfigPath(t *testing.T) {
	fs, err := NewFileSystem()
	require.NoError(t, err)

	configPath := fs.GetConfigPath()
	assert.NotEmpty(t, configPath)
	assert.Contains(t, configPath, "server.ini")
}

func TestFileSystem_GetSystemLogPath(t *testing.T) {
	fs, err := NewFileSystem()
	require.NoError(t, err)

	logPath := fs.GetSystemLogPath()
	assert.NotEmpty(t, logPath)
	assert.Contains(t, logPath, "system.log")
}

func TestFileSystem_GetServerLogPath(t *testing.T) {
	fs, err := NewFileSystem()
	require.NoError(t, err)

	logPath := fs.GetServerLogPath()
	assert.NotEmpty(t, logPath)
	assert.Contains(t, logPath, "server.log")
}

func TestFileSystem_GetDistPath(t *testing.T) {
	fs, err := NewFileSystem()
	require.NoError(t, err)

	distPath := fs.GetDistPath()
	assert.NotEmpty(t, distPath)
	assert.Contains(t, distPath, "dist")
}

func TestFileSystem_Exists_FileExists(t *testing.T) {
	fs, err := NewFileSystem()
	require.NoError(t, err)

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test"), 0644)
	require.NoError(t, err)

	assert.True(t, fs.Exists(testFile))
}

func TestFileSystem_Exists_FileNotExists(t *testing.T) {
	fs, err := NewFileSystem()
	require.NoError(t, err)

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "nonexistent.txt")

	assert.False(t, fs.Exists(testFile))
}

func TestFileSystem_EnsureDir_CreateNew(t *testing.T) {
	fs, err := NewFileSystem()
	require.NoError(t, err)

	tmpDir := t.TempDir()
	newDir := filepath.Join(tmpDir, "newdir")

	err = fs.EnsureDir(newDir)
	require.NoError(t, err)

	// Verify directory was created
	info, err := os.Stat(newDir)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestFileSystem_EnsureDir_AlreadyExists(t *testing.T) {
	fs, err := NewFileSystem()
	require.NoError(t, err)

	tmpDir := t.TempDir()

	// Should not error when directory already exists
	err = fs.EnsureDir(tmpDir)
	assert.NoError(t, err)
}

func TestFileSystem_EnsureDir_Nested(t *testing.T) {
	fs, err := NewFileSystem()
	require.NoError(t, err)

	tmpDir := t.TempDir()
	nestedDir := filepath.Join(tmpDir, "level1", "level2", "level3")

	err = fs.EnsureDir(nestedDir)
	require.NoError(t, err)

	// Verify nested directory was created
	info, err := os.Stat(nestedDir)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}

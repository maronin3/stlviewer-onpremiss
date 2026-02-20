package webview

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebViewAdapter_Creation(t *testing.T) {
	adapter := NewWebViewAdapter()
	assert.NotNil(t, adapter)
}

func TestWebViewAdapter_Initialize_NilConfig(t *testing.T) {
	adapter := NewWebViewAdapter()
	err := adapter.Initialize(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "config cannot be nil")
}

// Note: We cannot fully test webview initialization without a display
// These tests verify the adapter structure and error handling
func TestWebViewAdapter_Run_NotInitialized(t *testing.T) {
	adapter := NewWebViewAdapter()
	err := adapter.Run()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "webview not initialized")
}

func TestWebViewAdapter_Destroy_Safe(t *testing.T) {
	adapter := NewWebViewAdapter()
	// Should not panic when destroying uninitialized webview
	assert.NotPanics(t, func() {
		adapter.Destroy()
	})
}

func TestWebViewAdapter_Destroy_Multiple(t *testing.T) {
	adapter := NewWebViewAdapter()
	// Multiple destroys should be safe
	assert.NotPanics(t, func() {
		adapter.Destroy()
		adapter.Destroy()
	})
}

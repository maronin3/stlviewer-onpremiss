package webview

import (
	"fmt"
	"stlviewer/internal/domain"

	"github.com/webview/webview"
)

// WebViewAdapter implements WebViewPort using webview library
type WebViewAdapter struct {
	w webview.WebView
}

// NewWebViewAdapter creates a new webview adapter
func NewWebViewAdapter() *WebViewAdapter {
	return &WebViewAdapter{}
}

// Initialize initializes the webview with configuration
func (a *WebViewAdapter) Initialize(config *domain.Config) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	// Create webview
	a.w = webview.New(config.DevMode)
	if a.w == nil {
		return fmt.Errorf("failed to create webview")
	}

	// Configure webview
	a.w.SetTitle(config.Title)
	a.w.SetSize(config.Width, config.Height, webview.HintNone)

	// Bind functions
	a.w.Bind("quit", func() {
		a.w.Terminate()
	})

	// Navigate to URL
	a.w.Navigate(config.URL)

	return nil
}

// Run runs the webview main loop (blocking)
func (a *WebViewAdapter) Run() error {
	if a.w == nil {
		return fmt.Errorf("webview not initialized")
	}

	a.w.Run()
	return nil
}

// Destroy destroys the webview and frees resources
func (a *WebViewAdapter) Destroy() {
	if a.w != nil {
		a.w.Destroy()
		a.w = nil
	}
}

package usecase

import (
	"context"
	"stlviewer/internal/domain"
)

// WebViewPort defines the interface for webview implementations
type WebViewPort interface {
	Initialize(config *domain.Config) error
	Run() error
	Destroy()
}

// WebViewUseCase handles webview lifecycle management
type WebViewUseCase struct {
	webview WebViewPort
	logger  Logger
}

// NewWebViewUseCase creates a new webview use case
func NewWebViewUseCase(webview WebViewPort, logger Logger) *WebViewUseCase {
	return &WebViewUseCase{
		webview: webview,
		logger:  logger,
	}
}

// Run initializes and runs the webview
func (uc *WebViewUseCase) Run(ctx context.Context, config *domain.Config) error {
	uc.logger.Info("Initializing webview: %s (%dx%d)", config.Title, config.Width, config.Height)

	if err := uc.webview.Initialize(config); err != nil {
		uc.logger.Error("Failed to initialize webview: %v", err)
		return err
	}

	uc.logger.Info("Starting webview")
	if err := uc.webview.Run(); err != nil {
		uc.logger.Error("Webview error: %v", err)
		return err
	}

	uc.logger.Info("Webview closed")
	return nil
}

// Cleanup destroys the webview resources
func (uc *WebViewUseCase) Cleanup() {
	uc.logger.Info("Cleaning up webview resources")
	uc.webview.Destroy()
}

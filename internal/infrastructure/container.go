package infrastructure

import (
	"fmt"
	"stlviewer/internal/adapter/logger"
	"stlviewer/internal/adapter/repository"
	"stlviewer/internal/adapter/server"
	"stlviewer/internal/adapter/webview"
	"stlviewer/internal/usecase"
)

// Container holds all application dependencies
type Container struct {
	// Infrastructure
	FileSystem *FileSystem
	Logger     *logger.FileLogger

	// Repositories
	ConfigRepo *repository.INIConfigRepository

	// Adapters
	Server  *server.FiberServer
	WebView *webview.WebViewAdapter

	// Use Cases
	ConfigUseCase  *usecase.ConfigUseCase
	ServerUseCase  *usecase.ServerUseCase
	WebViewUseCase *usecase.WebViewUseCase
}

// NewContainer creates and wires all dependencies
func NewContainer() (*Container, error) {
	// Initialize filesystem
	fs, err := NewFileSystem()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize filesystem: %w", err)
	}

	// Initialize logger
	systemLogger, err := logger.NewFileLogger(fs.GetSystemLogPath())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Initialize repositories
	configRepo := repository.NewINIConfigRepository(fs.GetConfigPath())

	// Initialize adapters
	fiberServer := server.NewFiberServer(fs.GetDistPath(), fs.GetServerLogPath())
	webviewAdapter := webview.NewWebViewAdapter()

	// Initialize use cases
	configUseCase := usecase.NewConfigUseCase(configRepo, systemLogger)
	serverUseCase := usecase.NewServerUseCase(fiberServer, systemLogger)
	webviewUseCase := usecase.NewWebViewUseCase(webviewAdapter, systemLogger)

	return &Container{
		FileSystem:     fs,
		Logger:         systemLogger,
		ConfigRepo:     configRepo,
		Server:         fiberServer,
		WebView:        webviewAdapter,
		ConfigUseCase:  configUseCase,
		ServerUseCase:  serverUseCase,
		WebViewUseCase: webviewUseCase,
	}, nil
}

// Close cleans up resources
func (c *Container) Close() error {
	if c.Logger != nil {
		if err := c.Logger.Close(); err != nil {
			return fmt.Errorf("failed to close logger: %w", err)
		}
	}
	return nil
}

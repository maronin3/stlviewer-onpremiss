package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"stlviewer/internal/infrastructure"
	"syscall"
)

func main() {
	// Create dependency injection container
	container, err := infrastructure.NewContainer()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	defer container.Close()

	container.Logger.Info("=== STL Viewer Starting ===")
	container.Logger.Info("Application path: %s", container.FileSystem.GetAppPath())

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load or create configuration
	config, err := container.ConfigUseCase.LoadOrCreate(ctx)
	if err != nil {
		container.Logger.Error("Failed to load configuration: %v", err)
		log.Fatalf("Failed to load configuration: %v", err)
	}

	container.Logger.Info("Configuration loaded - Version: %s, DevMode: %t", config.Version, config.DevMode)
	container.Logger.Info("Server will start on %s", config.GetServerAddress())

	// Start web server in goroutine
	go func() {
		if err := container.ServerUseCase.Start(ctx, config); err != nil {
			container.Logger.Error("Server error: %v", err)
		}
	}()

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		container.Logger.Info("Shutdown signal received")
		cancel()

		// Stop server
		if err := container.ServerUseCase.Stop(ctx); err != nil {
			container.Logger.Error("Error stopping server: %v", err)
		}

		// Cleanup webview
		container.WebViewUseCase.Cleanup()

		container.Logger.Info("=== STL Viewer Stopped ===")
		os.Exit(0)
	}()

	// Run webview (blocking)
	container.Logger.Info("Starting webview")
	if err := container.WebViewUseCase.Run(ctx, config); err != nil {
		container.Logger.Error("WebView error: %v", err)
		fmt.Printf("WebView error: %v\n", err)
	}

	// Cleanup after webview closes
	container.Logger.Info("WebView closed, shutting down")
	cancel()

	if err := container.ServerUseCase.Stop(ctx); err != nil {
		container.Logger.Error("Error stopping server: %v", err)
	}

	container.Logger.Info("=== STL Viewer Stopped ===")
}

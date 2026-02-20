package usecase

import (
	"context"
	"stlviewer/internal/domain"
)

// ServerPort defines the interface for web server implementations
type ServerPort interface {
	Start(ctx context.Context, config *domain.Config) error
	Stop(ctx context.Context) error
}

// ServerUseCase handles server lifecycle management
type ServerUseCase struct {
	server ServerPort
	logger Logger
}

// NewServerUseCase creates a new server use case
func NewServerUseCase(server ServerPort, logger Logger) *ServerUseCase {
	return &ServerUseCase{
		server: server,
		logger: logger,
	}
}

// Start starts the web server with the given configuration
func (uc *ServerUseCase) Start(ctx context.Context, config *domain.Config) error {
	uc.logger.Info("Starting web server on %s", config.GetServerAddress())

	if err := uc.server.Start(ctx, config); err != nil {
		uc.logger.Error("Failed to start server: %v", err)
		return err
	}

	uc.logger.Info("Web server started successfully")
	return nil
}

// Stop gracefully stops the web server
func (uc *ServerUseCase) Stop(ctx context.Context) error {
	uc.logger.Info("Stopping web server")

	if err := uc.server.Stop(ctx); err != nil {
		uc.logger.Error("Failed to stop server: %v", err)
		return err
	}

	uc.logger.Info("Web server stopped successfully")
	return nil
}

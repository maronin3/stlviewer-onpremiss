package usecase

import (
	"context"
	"stlviewer/internal/domain"
)

// ConfigRepository defines the port for configuration persistence
type ConfigRepository interface {
	Load(ctx context.Context) (*domain.Config, error)
	Save(ctx context.Context, config *domain.Config) error
	Exists(ctx context.Context) (bool, error)
}

// Logger defines the port for logging
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
}

// ConfigUseCase handles configuration management business logic
type ConfigUseCase struct {
	repo   ConfigRepository
	logger Logger
}

// NewConfigUseCase creates a new configuration use case
func NewConfigUseCase(repo ConfigRepository, logger Logger) *ConfigUseCase {
	return &ConfigUseCase{
		repo:   repo,
		logger: logger,
	}
}

// LoadOrCreate loads existing configuration or creates default one
func (uc *ConfigUseCase) LoadOrCreate(ctx context.Context) (*domain.Config, error) {
	exists, err := uc.repo.Exists(ctx)
	if err != nil {
		uc.logger.Error("Failed to check config existence: %v", err)
		return nil, err
	}

	if exists {
		uc.logger.Info("Loading existing configuration")
		config, err := uc.repo.Load(ctx)
		if err != nil {
			uc.logger.Error("Failed to load config: %v", err)
			return nil, err
		}

		if err := config.Validate(); err != nil {
			uc.logger.Error("Invalid configuration: %v", err)
			return nil, err
		}

		uc.logger.Info("Configuration loaded successfully")
		return config, nil
	}

	uc.logger.Info("Creating default configuration")
	config := domain.NewDefaultConfig()

	if err := uc.repo.Save(ctx, config); err != nil {
		uc.logger.Error("Failed to save default config: %v", err)
		return nil, err
	}

	uc.logger.Info("Default configuration created and saved")
	return config, nil
}

// Save saves the configuration
func (uc *ConfigUseCase) Save(ctx context.Context, config *domain.Config) error {
	if err := config.Validate(); err != nil {
		uc.logger.Error("Invalid configuration: %v", err)
		return err
	}

	if err := uc.repo.Save(ctx, config); err != nil {
		uc.logger.Error("Failed to save config: %v", err)
		return err
	}

	uc.logger.Info("Configuration saved successfully")
	return nil
}

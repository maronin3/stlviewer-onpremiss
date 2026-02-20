package repository

import (
	"context"
	"fmt"
	"os"
	"stlviewer/internal/domain"

	"gopkg.in/ini.v1"
)

// INIConfigRepository implements ConfigRepository using INI files
type INIConfigRepository struct {
	filePath string
}

// NewINIConfigRepository creates a new INI-based config repository
func NewINIConfigRepository(filePath string) *INIConfigRepository {
	return &INIConfigRepository{
		filePath: filePath,
	}
}

// Exists checks if the configuration file exists
func (r *INIConfigRepository) Exists(ctx context.Context) (bool, error) {
	_, err := os.Stat(r.filePath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, fmt.Errorf("failed to check file existence: %w", err)
}

// Load loads configuration from INI file
func (r *INIConfigRepository) Load(ctx context.Context) (*domain.Config, error) {
	cfg, err := ini.Load(r.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load INI file: %w", err)
	}

	config := &domain.Config{
		Version: cfg.Section("Application").Key("Version").String(),
		DevMode: cfg.Section("Application").Key("DevMode").MustBool(false),
		IP:      cfg.Section("Server").Key("IP").String(),
		Port:    cfg.Section("Server").Key("Port").MustInt(3000),
		Title:   cfg.Section("WebView").Key("Title").String(),
		Width:   cfg.Section("WebView").Key("Width").MustInt(800),
		Height:  cfg.Section("WebView").Key("Height").MustInt(600),
		URL:     cfg.Section("WebView").Key("URL").String(),
	}

	return config, nil
}

// Save saves configuration to INI file
func (r *INIConfigRepository) Save(ctx context.Context, config *domain.Config) error {
	cfg := ini.Empty()

	cfg.Section("Application").Key("Version").SetValue(config.Version)
	cfg.Section("Application").Key("DevMode").SetValue(fmt.Sprintf("%t", config.DevMode))
	cfg.Section("Server").Key("IP").SetValue(config.IP)
	cfg.Section("Server").Key("Port").SetValue(fmt.Sprintf("%d", config.Port))
	cfg.Section("WebView").Key("Title").SetValue(config.Title)
	cfg.Section("WebView").Key("Width").SetValue(fmt.Sprintf("%d", config.Width))
	cfg.Section("WebView").Key("Height").SetValue(fmt.Sprintf("%d", config.Height))
	cfg.Section("WebView").Key("URL").SetValue(config.URL)

	if err := cfg.SaveTo(r.filePath); err != nil {
		return fmt.Errorf("failed to save INI file: %w", err)
	}

	return nil
}

package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"stlviewer/internal/domain"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// FiberServer implements ServerPort using Fiber framework
type FiberServer struct {
	app           *fiber.App
	distPath      string
	serverLogPath string
	logFile       *os.File
}

// NewFiberServer creates a new Fiber-based server
func NewFiberServer(distPath, serverLogPath string) *FiberServer {
	return &FiberServer{
		distPath:      distPath,
		serverLogPath: serverLogPath,
	}
}

// Start starts the Fiber server
func (s *FiberServer) Start(ctx context.Context, config *domain.Config) error {
	s.app = fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	// Configure middleware based on dev mode
	if config.DevMode {
		// Logger middleware
		if s.serverLogPath != "" {
			file, err := os.OpenFile(s.serverLogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				return fmt.Errorf("failed to open server log file: %w", err)
			}
			s.logFile = file

			s.app.Use(logger.New(logger.Config{
				Format:       "[${time}] ${status} - ${latency} ${method} ${path}\n",
				TimeFormat:   "15:04:05",
				TimeZone:     "Local",
				TimeInterval: 500 * time.Millisecond,
				Output:       file,
			}))
		}

		// Monitor dashboard
		s.app.Get("/dashboard", monitor.New())

		// pprof for profiling
		s.app.Use(pprof.New())
	}

	// Recover middleware
	s.app.Use(recover.New())

	// Static file serving
	s.app.Use(filesystem.New(filesystem.Config{
		Root:         http.Dir(s.distPath),
		Browse:       true,
		Index:        "index.html",
		NotFoundFile: "404.html",
		MaxAge:       3600,
	}))

	// Start server in goroutine
	go func() {
		if err := s.app.Listen(config.GetServerAddress()); err != nil {
			// Server error will be handled by context cancellation
			return
		}
	}()

	return nil
}

// Stop gracefully stops the Fiber server
func (s *FiberServer) Stop(ctx context.Context) error {
	if s.app == nil {
		return nil
	}

	if err := s.app.Shutdown(); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	// Close log file if it was opened
	if s.logFile != nil {
		s.logFile.Close()
		s.logFile = nil
	}

	return nil
}

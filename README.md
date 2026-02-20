# STL Viewer - On-Premise

A desktop application for viewing STL files, built with Clean Architecture principles.

## Architecture

This application follows **Clean Architecture** (Hexagonal Architecture / Ports & Adapters), ensuring:
- **Testability**: All components are easily testable in isolation
- **Maintainability**: Clear separation of concerns
- **Flexibility**: Easy to swap implementations
- **Independence**: Business logic independent of frameworks

### Architecture Layers

```
┌─────────────────────────────────────────┐
│         Domain Layer (Entities)         │
│    - Config entity with validation      │
│    - Pure business logic, no deps       │
└─────────────────────────────────────────┘
                    ▲
                    │
┌─────────────────────────────────────────┐
│      Use Cases Layer (App Logic)        │
│    - ConfigUseCase, ServerUseCase       │
│    - WebViewUseCase                     │
│    - Defines ports (interfaces)         │
└─────────────────────────────────────────┘
                    ▲
                    │
┌─────────────────────────────────────────┐
│    Interface Adapters (Implementations) │
│    - INI Config Repository              │
│    - File Logger                        │
│    - Fiber Server Adapter               │
│    - WebView Adapter                    │
└─────────────────────────────────────────┘
                    ▲
                    │
┌─────────────────────────────────────────┐
│      Infrastructure (DI & Utils)        │
│    - Dependency Injection Container     │
│    - File System Utilities              │
└─────────────────────────────────────────┘
```

## Directory Structure

```
stlviewer-onpremiss/
├── internal/
│   ├── domain/                     # Domain Layer
│   │   ├── config.go              # Config entity with validation
│   │   └── config_test.go         # 100% coverage
│   │
│   ├── usecase/                    # Use Cases Layer
│   │   ├── config_usecase.go      # Config management
│   │   ├── config_usecase_test.go
│   │   ├── server_usecase.go      # Server lifecycle
│   │   ├── server_usecase_test.go
│   │   ├── webview_usecase.go     # WebView lifecycle
│   │   └── webview_usecase_test.go # 96.5% coverage
│   │
│   ├── adapter/                    # Interface Adapters
│   │   ├── repository/
│   │   │   ├── config_repository.go      # INI file adapter
│   │   │   └── config_repository_test.go # 91.7% coverage
│   │   ├── logger/
│   │   │   ├── logger_adapter.go         # File/Console logger
│   │   │   └── logger_adapter_test.go    # 94.1% coverage
│   │   ├── server/
│   │   │   ├── fiber_server.go           # Fiber framework adapter
│   │   │   └── fiber_server_test.go      # 88.0% coverage
│   │   └── webview/
│   │       ├── webview_adapter.go        # WebView library adapter
│   │       └── webview_adapter_test.go
│   │
│   └── infrastructure/             # Infrastructure
│       ├── container.go            # DI container
│       ├── filesystem.go           # Path utilities
│       └── filesystem_test.go
│
├── frontend/                       # Frontend assets
├── lib/                           # Native libraries
├── main.go                        # Application entry point
├── go.mod
└── README.md
```

## Features

- **Configuration Management**: INI-based configuration with defaults
- **Web Server**: Fiber-based HTTP server for serving static files
- **WebView**: Native desktop window for UI
- **Logging**: File-based logging with dev/prod modes
- **Graceful Shutdown**: Proper cleanup of resources

## Development

### Prerequisites

- Go 1.21 or higher
- WebView2 runtime (Windows) / WebKit (Linux) / WebKit (macOS)

### Building

```bash
# Install dependencies
go mod download

# Run tests
go test ./internal/... -v

# Test coverage
go test ./internal/domain/... ./internal/usecase/... ./internal/adapter/repository/... ./internal/adapter/logger/... ./internal/adapter/server/... -cover

# Build
go build -o stlviewer.exe

# Run
./stlviewer.exe
```

### Testing

The application follows **Test-Driven Development (TDD)**:

```bash
# Run all tests
go test ./internal/... -v

# Run specific layer tests
go test ./internal/domain/... -v
go test ./internal/usecase/... -v
go test ./internal/adapter/... -v

# Coverage report
go test ./internal/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

**Test Coverage:**
- Domain Layer: 100%
- Use Cases: 96.5%
- Repository Adapter: 91.7%
- Logger Adapter: 94.1%
- Server Adapter: 88.0%

## Configuration

Configuration is stored in `server.ini`:

```ini
[Application]
Version = 1.0.0
DevMode = false

[Server]
IP = 127.0.0.1
Port = 3000

[WebView]
Title = STLViewer
Width = 800
Height = 600
URL = http://127.0.0.1:3000
```

If the file doesn't exist, it will be created with default values on first run.

## Architecture Benefits

### 1. Testability
- All business logic is testable without external dependencies
- Mocks are used for testing use cases
- Integration tests use temporary files/ports

### 2. Maintainability
- Clear separation of concerns
- Each layer has a single responsibility
- Easy to understand and modify

### 3. Flexibility
- Easy to swap INI config for JSON/YAML
- Easy to replace Fiber with another framework
- Easy to add new features

### 4. Independence
- Domain logic doesn't depend on frameworks
- Use cases don't know about implementation details
- Infrastructure is isolated

## License

[Your License Here]

package usecase

import (
	"context"
	"errors"
	"stlviewer/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockWebViewPort is a mock implementation of WebViewPort
type MockWebViewPort struct {
	mock.Mock
}

func (m *MockWebViewPort) Initialize(config *domain.Config) error {
	args := m.Called(config)
	return args.Error(0)
}

func (m *MockWebViewPort) Run() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockWebViewPort) Destroy() {
	m.Called()
}

func TestWebViewUseCase_Run_Success(t *testing.T) {
	mockWebView := new(MockWebViewPort)
	mockLogger := new(MockLogger)
	ctx := context.Background()
	config := domain.NewDefaultConfig()

	mockWebView.On("Initialize", config).Return(nil)
	mockWebView.On("Run").Return(nil)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()

	uc := NewWebViewUseCase(mockWebView, mockLogger)
	err := uc.Run(ctx, config)

	assert.NoError(t, err)
	mockWebView.AssertExpectations(t)
}

func TestWebViewUseCase_Run_InitializeError(t *testing.T) {
	mockWebView := new(MockWebViewPort)
	mockLogger := new(MockLogger)
	ctx := context.Background()
	config := domain.NewDefaultConfig()

	expectedErr := errors.New("initialization failed")
	mockWebView.On("Initialize", config).Return(expectedErr)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	uc := NewWebViewUseCase(mockWebView, mockLogger)
	err := uc.Run(ctx, config)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockWebView.AssertNotCalled(t, "Run")
}

func TestWebViewUseCase_Run_RunError(t *testing.T) {
	mockWebView := new(MockWebViewPort)
	mockLogger := new(MockLogger)
	ctx := context.Background()
	config := domain.NewDefaultConfig()

	expectedErr := errors.New("webview crashed")
	mockWebView.On("Initialize", config).Return(nil)
	mockWebView.On("Run").Return(expectedErr)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	uc := NewWebViewUseCase(mockWebView, mockLogger)
	err := uc.Run(ctx, config)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestWebViewUseCase_Cleanup(t *testing.T) {
	mockWebView := new(MockWebViewPort)
	mockLogger := new(MockLogger)

	mockWebView.On("Destroy").Return()
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()

	uc := NewWebViewUseCase(mockWebView, mockLogger)
	uc.Cleanup()

	mockWebView.AssertExpectations(t)
}

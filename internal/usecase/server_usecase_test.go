package usecase

import (
	"context"
	"errors"
	"stlviewer/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockServerPort is a mock implementation of ServerPort
type MockServerPort struct {
	mock.Mock
}

func (m *MockServerPort) Start(ctx context.Context, config *domain.Config) error {
	args := m.Called(ctx, config)
	return args.Error(0)
}

func (m *MockServerPort) Stop(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestServerUseCase_Start_Success(t *testing.T) {
	mockServer := new(MockServerPort)
	mockLogger := new(MockLogger)
	ctx := context.Background()
	config := domain.NewDefaultConfig()

	mockServer.On("Start", ctx, config).Return(nil)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()

	uc := NewServerUseCase(mockServer, mockLogger)
	err := uc.Start(ctx, config)

	assert.NoError(t, err)
	mockServer.AssertExpectations(t)
}

func TestServerUseCase_Start_Error(t *testing.T) {
	mockServer := new(MockServerPort)
	mockLogger := new(MockLogger)
	ctx := context.Background()
	config := domain.NewDefaultConfig()

	expectedErr := errors.New("port already in use")
	mockServer.On("Start", ctx, config).Return(expectedErr)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	uc := NewServerUseCase(mockServer, mockLogger)
	err := uc.Start(ctx, config)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestServerUseCase_Stop_Success(t *testing.T) {
	mockServer := new(MockServerPort)
	mockLogger := new(MockLogger)
	ctx := context.Background()

	mockServer.On("Stop", ctx).Return(nil)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()

	uc := NewServerUseCase(mockServer, mockLogger)
	err := uc.Stop(ctx)

	assert.NoError(t, err)
	mockServer.AssertExpectations(t)
}

func TestServerUseCase_Stop_Error(t *testing.T) {
	mockServer := new(MockServerPort)
	mockLogger := new(MockLogger)
	ctx := context.Background()

	expectedErr := errors.New("shutdown timeout")
	mockServer.On("Stop", ctx).Return(expectedErr)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	uc := NewServerUseCase(mockServer, mockLogger)
	err := uc.Stop(ctx)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

package usecase

import (
	"context"
	"errors"
	"stlviewer/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockConfigRepository is a mock implementation of ConfigRepository
type MockConfigRepository struct {
	mock.Mock
}

func (m *MockConfigRepository) Load(ctx context.Context) (*domain.Config, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Config), args.Error(1)
}

func (m *MockConfigRepository) Save(ctx context.Context, config *domain.Config) error {
	args := m.Called(ctx, config)
	return args.Error(0)
}

func (m *MockConfigRepository) Exists(ctx context.Context) (bool, error) {
	args := m.Called(ctx)
	return args.Bool(0), args.Error(1)
}

// MockLogger is a mock implementation of Logger
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func (m *MockLogger) Error(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func (m *MockLogger) Debug(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func TestConfigUseCase_LoadOrCreate_ExistingConfig(t *testing.T) {
	mockRepo := new(MockConfigRepository)
	mockLogger := new(MockLogger)
	ctx := context.Background()

	expectedConfig := domain.NewDefaultConfig()
	expectedConfig.Port = 8080

	mockRepo.On("Exists", ctx).Return(true, nil)
	mockRepo.On("Load", ctx).Return(expectedConfig, nil)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()

	uc := NewConfigUseCase(mockRepo, mockLogger)
	config, err := uc.LoadOrCreate(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, 8080, config.Port)
	mockRepo.AssertExpectations(t)
}

func TestConfigUseCase_LoadOrCreate_CreateDefault(t *testing.T) {
	mockRepo := new(MockConfigRepository)
	mockLogger := new(MockLogger)
	ctx := context.Background()

	mockRepo.On("Exists", ctx).Return(false, nil)
	mockRepo.On("Save", ctx, mock.AnythingOfType("*domain.Config")).Return(nil)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()

	uc := NewConfigUseCase(mockRepo, mockLogger)
	config, err := uc.LoadOrCreate(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, 3000, config.Port) // Default port
	mockRepo.AssertExpectations(t)
}

func TestConfigUseCase_LoadOrCreate_ExistsError(t *testing.T) {
	mockRepo := new(MockConfigRepository)
	mockLogger := new(MockLogger)
	ctx := context.Background()

	expectedErr := errors.New("filesystem error")
	mockRepo.On("Exists", ctx).Return(false, expectedErr)
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	uc := NewConfigUseCase(mockRepo, mockLogger)
	config, err := uc.LoadOrCreate(ctx)

	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Equal(t, expectedErr, err)
}

func TestConfigUseCase_LoadOrCreate_LoadError(t *testing.T) {
	mockRepo := new(MockConfigRepository)
	mockLogger := new(MockLogger)
	ctx := context.Background()

	expectedErr := errors.New("load error")
	mockRepo.On("Exists", ctx).Return(true, nil)
	mockRepo.On("Load", ctx).Return(nil, expectedErr)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	uc := NewConfigUseCase(mockRepo, mockLogger)
	config, err := uc.LoadOrCreate(ctx)

	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Equal(t, expectedErr, err)
}

func TestConfigUseCase_LoadOrCreate_InvalidConfig(t *testing.T) {
	mockRepo := new(MockConfigRepository)
	mockLogger := new(MockLogger)
	ctx := context.Background()

	invalidConfig := domain.NewDefaultConfig()
	invalidConfig.Port = -1 // Invalid port

	mockRepo.On("Exists", ctx).Return(true, nil)
	mockRepo.On("Load", ctx).Return(invalidConfig, nil)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	uc := NewConfigUseCase(mockRepo, mockLogger)
	config, err := uc.LoadOrCreate(ctx)

	assert.Error(t, err)
	assert.Nil(t, config)
}

func TestConfigUseCase_Save_Success(t *testing.T) {
	mockRepo := new(MockConfigRepository)
	mockLogger := new(MockLogger)
	ctx := context.Background()

	config := domain.NewDefaultConfig()
	mockRepo.On("Save", ctx, config).Return(nil)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()

	uc := NewConfigUseCase(mockRepo, mockLogger)
	err := uc.Save(ctx, config)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestConfigUseCase_Save_InvalidConfig(t *testing.T) {
	mockRepo := new(MockConfigRepository)
	mockLogger := new(MockLogger)
	ctx := context.Background()

	config := domain.NewDefaultConfig()
	config.Port = -1 // Invalid
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	uc := NewConfigUseCase(mockRepo, mockLogger)
	err := uc.Save(ctx, config)

	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "Save")
}

func TestConfigUseCase_Save_RepositoryError(t *testing.T) {
	mockRepo := new(MockConfigRepository)
	mockLogger := new(MockLogger)
	ctx := context.Background()

	config := domain.NewDefaultConfig()
	expectedErr := errors.New("save error")
	mockRepo.On("Save", ctx, config).Return(expectedErr)
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	uc := NewConfigUseCase(mockRepo, mockLogger)
	err := uc.Save(ctx, config)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

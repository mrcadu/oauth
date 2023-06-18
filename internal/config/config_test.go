package config

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockEnv struct {
	mock.Mock
}

func (m *MockEnv) LoadEnv(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

func (m *MockEnv) GetEnv(key string) string {
	args := m.Called(key)
	return args.String(0)
}

func TestGetPropertyWithValidEnvFileAndValidPropertyShouldReturnProperty(t *testing.T) {
	mockEnv := MockEnv{}
	mockEnv.On("LoadEnv", ".env").Return(nil)
	mockEnv.On("GetEnv", "key").Return("value")
	SysEnv = &mockEnv

	property := GetProperty("key")

	assert.Equal(t, property, "value")
}

func TestGetPropertyWithValidEnvFileAndInvalidPropertyShouldReturnProperty(t *testing.T) {
	mockEnv := MockEnv{}
	mockEnv.On("LoadEnv", ".env").Return(nil)
	mockEnv.On("GetEnv", "key").Return("")
	SysEnv = &mockEnv

	property := GetProperty("key")

	assert.Equal(t, property, "")
}

func TestGetPropertyWithInvalidEnvFileAndValidPropertyShouldReturnProperty(t *testing.T) {
	mockEnv := MockEnv{}
	errorMessage := "invalid env file"
	mockEnv.On("LoadEnv", ".env").Return(errors.New(errorMessage))
	mockEnv.On("GetEnv", "key").Return("")
	SysEnv = &mockEnv

	assert.Equal(t, GetProperty("key"), "")
}

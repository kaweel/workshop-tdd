package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ✅ Mock Repository
type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) GetUser(id int) string {
	args := m.Called(id)
	return args.String(0)
}

func TestGetUser(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockRepo.On("GetUser", 1).Return("John Doe")

	result := mockRepo.GetUser(1)

	assert.Equal(t, "John Doe", result)

	mockRepo.AssertExpectations(t)
}

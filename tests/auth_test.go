package sendlix_test

import (
	"context"
	"testing"

	sendlix "github.com/sendlix/go-sdk"
	"github.com/stretchr/testify/assert"
)

func TestNewAuth(t *testing.T) {
	tests := []struct {
		name        string
		apiKey      string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid API key",
			apiKey:      "secret123.456",
			expectError: false,
		},
		{
			name:        "Invalid format - no dot",
			apiKey:      "secret123456",
			expectError: true,
			errorMsg:    "invalid API key format",
		},
		{
			name:        "Invalid format - empty secret",
			apiKey:      ".456",
			expectError: true,
		},
		{
			name:        "Invalid format - empty keyID",
			apiKey:      "secret123.",
			expectError: true,
			errorMsg:    "invalid key ID",
		},
		{
			name:        "Invalid format - non-numeric keyID",
			apiKey:      "secret123.abc",
			expectError: true,
			errorMsg:    "invalid key ID",
		},
		{
			name:        "Invalid format - multiple dots",
			apiKey:      "secret.123.456",
			expectError: true,
			errorMsg:    "invalid API key format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth, err := sendlix.NewAuth(tt.apiKey)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, auth)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, auth)
			}
		})
	}
}

// MockAuth implements IAuth interface for testing
type MockAuth struct {
	Token string
	Error error
}

func (m *MockAuth) GetAuthHeader(ctx context.Context) (string, string, error) {
	if m.Error != nil {
		return "", "", m.Error
	}
	return "authorization", "Bearer " + m.Token, nil
}

func TestMockAuth(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		mockAuth := &MockAuth{Token: "test-token"}

		key, value, err := mockAuth.GetAuthHeader(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, "authorization", key)
		assert.Equal(t, "Bearer test-token", value)
	})

	t.Run("Error case", func(t *testing.T) {
		mockAuth := &MockAuth{Error: assert.AnError}

		key, value, err := mockAuth.GetAuthHeader(context.Background())

		assert.Error(t, err)
		assert.Empty(t, key)
		assert.Empty(t, value)
	})
}

func TestAuthInterface(t *testing.T) {
	// Test that MockAuth implements IAuth interface
	var auth sendlix.IAuth = &MockAuth{Token: "test"}

	key, value, err := auth.GetAuthHeader(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, "authorization", key)
	assert.Equal(t, "Bearer test", value)
}

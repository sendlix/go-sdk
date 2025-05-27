package sendlix_test

import (
	"testing"

	sendlix "github.com/sendlix/go-sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultClientConfig(t *testing.T) {
	config := sendlix.DefaultClientConfig()

	require.NotNil(t, config)
	assert.Equal(t, "api.sendlix.com:443", config.ServerAddress)
	assert.Equal(t, "sendlix-go-sdk/1.0.0", config.UserAgent)
	assert.False(t, config.Insecure)
}

func TestNewBaseClient(t *testing.T) {
	t.Run("With default config", func(t *testing.T) {
		mockAuth := &MockAuth{Token: "test-token"}

		client, err := sendlix.NewBaseClient(mockAuth, nil)

		// gRPC uses lazy connections, so client creation should succeed
		// The connection is only established when an RPC is made
		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.NotNil(t, client.GetConnection())

		// Clean up
		client.Close()
	})

	t.Run("With custom config", func(t *testing.T) {
		mockAuth := &MockAuth{Token: "test-token"}
		config := &sendlix.ClientConfig{
			ServerAddress: "localhost:8080",
			UserAgent:     "test-client/1.0.0",
			Insecure:      true,
		}

		client, err := sendlix.NewBaseClient(mockAuth, config)

		// gRPC uses lazy connections, so client creation should succeed
		// The connection is only established when an RPC is made
		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.NotNil(t, client.GetConnection())

		// Clean up
		client.Close()
	})

	t.Run("With nil auth", func(t *testing.T) {
		client, err := sendlix.NewBaseClient(nil, nil)

		// Should still try to connect and fail, but auth can be nil
		assert.Error(t, err)
		assert.Nil(t, client)
	})
}

func TestClientConfig(t *testing.T) {
	t.Run("Custom config values", func(t *testing.T) {
		config := &sendlix.ClientConfig{
			ServerAddress: "custom.example.com:9090",
			UserAgent:     "custom-agent/2.0.0",
			Insecure:      true,
		}

		assert.Equal(t, "custom.example.com:9090", config.ServerAddress)
		assert.Equal(t, "custom-agent/2.0.0", config.UserAgent)
		assert.True(t, config.Insecure)
	})

	t.Run("Zero values", func(t *testing.T) {
		config := &sendlix.ClientConfig{}

		assert.Empty(t, config.ServerAddress)
		assert.Empty(t, config.UserAgent)
		assert.False(t, config.Insecure)
	})
}

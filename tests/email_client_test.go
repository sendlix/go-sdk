package sendlix_test

import (
	"testing"

	sendlix "github.com/sendlix/go-sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmailAddress(t *testing.T) {
	t.Run("String representation with name", func(t *testing.T) {
		addr := sendlix.EmailAddress{
			Email: "test@example.com",
			Name:  "Test User",
		}

		expected := "Test User <test@example.com>"
		assert.Equal(t, expected, addr.String())
	})

	t.Run("String representation without name", func(t *testing.T) {
		addr := sendlix.EmailAddress{
			Email: "test@example.com",
			Name:  "",
		}

		expected := "test@example.com"
		assert.Equal(t, expected, addr.String())
	})

	t.Run("String representation with empty name", func(t *testing.T) {
		addr := sendlix.EmailAddress{
			Email: "test@example.com",
		}

		expected := "test@example.com"
		assert.Equal(t, expected, addr.String())
	})
}

func TestNewEmailAddress(t *testing.T) {
	t.Run("From string", func(t *testing.T) {
		email := "test@example.com"

		addr, err := sendlix.NewEmailAddress(email)

		require.NoError(t, err)
		require.NotNil(t, addr)
		assert.Equal(t, email, addr.Email)
		assert.Empty(t, addr.Name)
	})

	t.Run("From EmailAddress struct", func(t *testing.T) {
		original := sendlix.EmailAddress{
			Email: "test@example.com",
			Name:  "Test User",
		}

		addr, err := sendlix.NewEmailAddress(original)

		require.NoError(t, err)
		require.NotNil(t, addr)
		assert.Equal(t, original.Email, addr.Email)
		assert.Equal(t, original.Name, addr.Name)
	})

	t.Run("From invalid type", func(t *testing.T) {
		invalidInput := 123

		addr, err := sendlix.NewEmailAddress(invalidInput)

		assert.Error(t, err)
		assert.Nil(t, addr)
		assert.Contains(t, err.Error(), "invalid email address type")
	})

	t.Run("From nil", func(t *testing.T) {
		addr, err := sendlix.NewEmailAddress(nil)

		assert.Error(t, err)
		assert.Nil(t, addr)
		assert.Contains(t, err.Error(), "invalid email address type")
	})
}

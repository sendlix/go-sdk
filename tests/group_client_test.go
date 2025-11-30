package sendlix_test

import (
	"testing"

	sendlix "github.com/sendlix/go-sdk"
	"github.com/stretchr/testify/assert"
)

func TestGroupEntry(t *testing.T) {
	t.Run("Create GroupEntry struct", func(t *testing.T) {
		entry := sendlix.GroupEntry{
			Email: "test@example.com",
			Name:  "Test User",
			Substitutions: map[string]string{
				"first_name": "Test",
				"discount":   "10%",
			},
		}

		assert.Equal(t, "test@example.com", entry.Email)
		assert.Equal(t, "Test User", entry.Name)
		assert.Equal(t, "Test", entry.Substitutions["first_name"])
		assert.Equal(t, "10%", entry.Substitutions["discount"])
	})

	t.Run("GroupEntry with empty name", func(t *testing.T) {
		entry := sendlix.GroupEntry{
			Email: "test@example.com",
		}

		assert.Equal(t, "test@example.com", entry.Email)
		assert.Empty(t, entry.Name)
		assert.Nil(t, entry.Substitutions)
	})

	t.Run("GroupEntry with substitutions only", func(t *testing.T) {
		entry := sendlix.GroupEntry{
			Email: "test@example.com",
			Substitutions: map[string]string{
				"welcome": "true",
			},
		}

		assert.Equal(t, "test@example.com", entry.Email)
		assert.Empty(t, entry.Name)
		assert.Equal(t, "true", entry.Substitutions["welcome"])
	})
}

func TestFailureHandler(t *testing.T) {
	t.Run("FailureHandlerSkip value", func(t *testing.T) {
		assert.Equal(t, sendlix.FailureHandler(0), sendlix.FailureHandlerSkip)
	})

	t.Run("FailureHandlerAbort value", func(t *testing.T) {
		assert.Equal(t, sendlix.FailureHandler(1), sendlix.FailureHandlerAbort)
	})
}

func TestInsertOptions(t *testing.T) {
	t.Run("Create InsertOptions with Skip", func(t *testing.T) {
		options := sendlix.InsertOptions{
			OnFailure: sendlix.FailureHandlerSkip,
		}

		assert.Equal(t, sendlix.FailureHandlerSkip, options.OnFailure)
	})

	t.Run("Create InsertOptions with Abort", func(t *testing.T) {
		options := sendlix.InsertOptions{
			OnFailure: sendlix.FailureHandlerAbort,
		}

		assert.Equal(t, sendlix.FailureHandlerAbort, options.OnFailure)
	})
}

func TestUpdateResponse(t *testing.T) {
	t.Run("Success response", func(t *testing.T) {
		response := sendlix.UpdateResponse{
			Success:      true,
			Message:      "Operation completed successfully",
			AffectedRows: 5,
		}

		assert.True(t, response.Success)
		assert.Equal(t, "Operation completed successfully", response.Message)
		assert.Equal(t, int64(5), response.AffectedRows)
	})

	t.Run("Error response", func(t *testing.T) {
		response := sendlix.UpdateResponse{
			Success:      false,
			Message:      "Operation failed",
			AffectedRows: 0,
		}

		assert.False(t, response.Success)
		assert.Equal(t, "Operation failed", response.Message)
		assert.Equal(t, int64(0), response.AffectedRows)
	})
}

package sendlix_test

import (
	"testing"

	sendlix "github.com/sendlix/go-sdk"
	"github.com/stretchr/testify/assert"
)

func TestEmailData(t *testing.T) {
	t.Run("Create EmailData struct", func(t *testing.T) {
		emailData := sendlix.EmailData{
			Email: "test@example.com",
			Name:  "Test User",
		}

		assert.Equal(t, "test@example.com", emailData.Email)
		assert.Equal(t, "Test User", emailData.Name)
	})

	t.Run("EmailData with empty name", func(t *testing.T) {
		emailData := sendlix.EmailData{
			Email: "test@example.com",
		}

		assert.Equal(t, "test@example.com", emailData.Email)
		assert.Empty(t, emailData.Name)
	})
}

func TestInsertEmailToGroupResponse(t *testing.T) {
	t.Run("Success response", func(t *testing.T) {
		response := sendlix.InsertEmailToGroupResponse{
			Success:      true,
			Message:      "Emails inserted successfully",
			AffectedRows: 5,
		}

		assert.True(t, response.Success)
		assert.Equal(t, "Emails inserted successfully", response.Message)
		assert.Equal(t, int64(5), response.AffectedRows)
	})

	t.Run("Error response", func(t *testing.T) {
		response := sendlix.InsertEmailToGroupResponse{
			Success:      false,
			Message:      "Failed to insert emails",
			AffectedRows: 0,
		}

		assert.False(t, response.Success)
		assert.Equal(t, "Failed to insert emails", response.Message)
		assert.Equal(t, int64(0), response.AffectedRows)
	})
}

func TestRemoveEmailFromGroupResponse(t *testing.T) {
	t.Run("Success response", func(t *testing.T) {
		response := sendlix.RemoveEmailFromGroupResponse{
			Success:      true,
			Message:      "Emails removed successfully",
			AffectedRows: 3,
		}

		assert.True(t, response.Success)
		assert.Equal(t, "Emails removed successfully", response.Message)
		assert.Equal(t, int64(3), response.AffectedRows)
	})

	t.Run("Error response", func(t *testing.T) {
		response := sendlix.RemoveEmailFromGroupResponse{
			Success:      false,
			Message:      "Failed to remove emails",
			AffectedRows: 0,
		}

		assert.False(t, response.Success)
		assert.Equal(t, "Failed to remove emails", response.Message)
		assert.Equal(t, int64(0), response.AffectedRows)
	})
}

func TestCheckEmailInGroupResponse(t *testing.T) {
	t.Run("Email exists", func(t *testing.T) {
		response := sendlix.CheckEmailInGroupResponse{
			Exists: true,
		}

		assert.True(t, response.Exists)
	})

	t.Run("Email does not exist", func(t *testing.T) {
		response := sendlix.CheckEmailInGroupResponse{
			Exists: false,
		}

		assert.False(t, response.Exists)
	})
}

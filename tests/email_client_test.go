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

func TestMimeType(t *testing.T) {
	t.Run("MimeTypePNG value", func(t *testing.T) {
		assert.Equal(t, sendlix.MimeType(0), sendlix.MimeTypePNG)
	})

	t.Run("MimeTypeJPEG value", func(t *testing.T) {
		assert.Equal(t, sendlix.MimeType(1), sendlix.MimeTypeJPEG)
	})

	t.Run("MimeTypeGIF value", func(t *testing.T) {
		assert.Equal(t, sendlix.MimeType(2), sendlix.MimeTypeGIF)
	})
}

func TestImage(t *testing.T) {
	t.Run("Create Image struct", func(t *testing.T) {
		imageData := []byte{0x89, 0x50, 0x4E, 0x47} // PNG header bytes
		img := sendlix.Image{
			Placeholder: "{{logo}}",
			Data:        imageData,
			Type:        sendlix.MimeTypePNG,
		}

		assert.Equal(t, "{{logo}}", img.Placeholder)
		assert.Equal(t, imageData, img.Data)
		assert.Equal(t, sendlix.MimeTypePNG, img.Type)
	})

	t.Run("Create JPEG Image", func(t *testing.T) {
		imageData := []byte{0xFF, 0xD8, 0xFF} // JPEG header bytes
		img := sendlix.Image{
			Placeholder: "{{photo}}",
			Data:        imageData,
			Type:        sendlix.MimeTypeJPEG,
		}

		assert.Equal(t, "{{photo}}", img.Placeholder)
		assert.Equal(t, sendlix.MimeTypeJPEG, img.Type)
	})
}

func TestMailOptions(t *testing.T) {
	t.Run("Create MailOptions with all fields", func(t *testing.T) {
		replyTo := sendlix.EmailAddress{Email: "reply@example.com"}
		imageData := []byte{0x89, 0x50, 0x4E, 0x47}

		options := sendlix.MailOptions{
			From:     sendlix.EmailAddress{Email: "sender@example.com", Name: "Sender"},
			To:       []sendlix.EmailAddress{{Email: "recipient@example.com"}},
			CC:       []sendlix.EmailAddress{{Email: "cc@example.com"}},
			BCC:      []sendlix.EmailAddress{{Email: "bcc@example.com"}},
			Subject:  "Test Subject",
			ReplyTo:  &replyTo,
			Html:     "<h1>Hello</h1>",
			Text:     "Hello",
			Tracking: true,
			Images: []sendlix.Image{
				{Placeholder: "{{logo}}", Data: imageData, Type: sendlix.MimeTypePNG},
			},
		}

		assert.Equal(t, "sender@example.com", options.From.Email)
		assert.Equal(t, "Sender", options.From.Name)
		assert.Len(t, options.To, 1)
		assert.Len(t, options.CC, 1)
		assert.Len(t, options.BCC, 1)
		assert.Equal(t, "Test Subject", options.Subject)
		assert.Equal(t, "reply@example.com", options.ReplyTo.Email)
		assert.Equal(t, "<h1>Hello</h1>", options.Html)
		assert.Equal(t, "Hello", options.Text)
		assert.True(t, options.Tracking)
		assert.Len(t, options.Images, 1)
		assert.Equal(t, "{{logo}}", options.Images[0].Placeholder)
	})

	t.Run("Create minimal MailOptions", func(t *testing.T) {
		options := sendlix.MailOptions{
			From:    sendlix.EmailAddress{Email: "sender@example.com"},
			To:      []sendlix.EmailAddress{{Email: "recipient@example.com"}},
			Subject: "Test",
			Html:    "<p>Test</p>",
		}

		assert.Equal(t, "sender@example.com", options.From.Email)
		assert.Len(t, options.To, 1)
		assert.Equal(t, "Test", options.Subject)
		assert.Equal(t, "<p>Test</p>", options.Html)
		assert.Empty(t, options.Text)
		assert.False(t, options.Tracking)
		assert.Nil(t, options.Images)
	})
}

func TestAttachment(t *testing.T) {
	t.Run("Create Attachment struct", func(t *testing.T) {
		attachment := sendlix.Attachment{
			ContentURL:  "https://example.com/file.pdf",
			Filename:    "document.pdf",
			ContentType: "application/pdf",
		}

		assert.Equal(t, "https://example.com/file.pdf", attachment.ContentURL)
		assert.Equal(t, "document.pdf", attachment.Filename)
		assert.Equal(t, "application/pdf", attachment.ContentType)
	})
}

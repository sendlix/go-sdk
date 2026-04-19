// Package sendlix provides a comprehensive Go SDK for the Sendlix email service API.
//
// The Sendlix Go SDK enables developers to integrate email sending capabilities
// into their Go applications with support for individual emails, group emails,
// and advanced features like scheduling, attachments, and email tracking.
//
// # Documentation Scope
//
// This documentation covers only the public API of the Sendlix Go SDK.
// Generated protocol buffer files and internal implementation details
// are excluded from this documentation to provide a clean, focused
// developer experience.
//
// # Quick Start
//
// To get started with the Sendlix SDK, you'll need an API key from your Sendlix account:
//
//	import "github.com/sendlix/go-sdk"
//
//	// Create an email client with your API key
//	client, err := sendlix.NewEmailClient("your-secret.your-key-id", nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer client.Close()
//
//	// Send a simple email
//	messageIDs, err := client.SendEmail(context.Background(), sendlix.MailOptions{
//		From:    sendlix.EmailAddress{Email: "sender@example.com", Name: "Sender Name"},
//		To:      []sendlix.EmailAddress{{Email: "recipient@example.com", Name: "Recipient"}},
//		Subject: "Hello from Sendlix!",
//		Html:    "<h1>Hello World!</h1><p>This is a test email.</p>",
//		Text:    "Hello World!\n\nThis is a test email.",
//	}, nil)
//
// For more examples and advanced usage, see the individual type documentation
// and the official Sendlix API documentation.
//
// # Internal Packages
//
// Note: Protocol buffer files are located in the internal/proto/ directory
// and follow Go's internal package conventions. These files are automatically
// excluded from public documentation and should not be used directly by
// external applications. All necessary functionality is exposed through the
// public API documented here.
package sendlix

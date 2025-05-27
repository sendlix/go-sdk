# Sendlix Go SDK

The official Go SDK for the Sendlix email service API. This SDK provides a comprehensive interface for sending emails, managing email groups, and handling authentication with the Sendlix platform.

## Features

- **Email Sending**: Send individual emails with full control over recipients, content, and formatting
- **Group Management**: Manage email groups for efficient bulk email operations
- **Multiple Content Types**: Support for HTML, plain text, and EML format emails
- **Advanced Features**: Email scheduling, attachments, tracking, and categorization
- **Authentication**: Automatic JWT token management with API key authentication
- **Error Handling**: Comprehensive error reporting and quota information
- **Context Support**: Full support for Go contexts including timeouts and cancellation

## Installation

```bash
go get github.com/sendlix/go-sdk
```

## Quick Start

```go
package main

import (
    "context"
    "log"

    "github.com/sendlix/go-sdk/pkg"
)

func main() {
    // Create authentication with your API key
    auth, err := sendlix.NewAuth("your-secret.your-key-id")
    if err != nil {
        log.Fatal(err)
    }

    // Create an email client
    client, err := sendlix.NewEmailClient(auth, nil)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Send a simple email
    response, err := client.SendEmail(context.Background(), sendlix.MailOptions{
        From:    sendlix.EmailAddress{Email: "sender@example.com", Name: "Sender Name"},
        To:      []sendlix.EmailAddress{{Email: "recipient@example.com", Name: "Recipient"}},
        Subject: "Hello from Sendlix!",
        Content: sendlix.MailContent{
            HTML: "<h1>Hello World!</h1><p>This is a test email.</p>",
            Text: "Hello World!\n\nThis is a test email.",
        },
    }, nil)

    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Email sent! Message IDs: %v", response.MessageList)
    log.Printf("Emails remaining: %d", response.EmailsLeft)
}
```

## Authentication

All API operations require authentication using an API key from your Sendlix account. The API key format is `secret.keyID`:

```go
auth, err := sendlix.NewAuth("your-secret.123456")
if err != nil {
    log.Fatal(err)
}
```

The SDK automatically handles JWT token exchange and caching, so you don't need to manage tokens manually.

## Sending Emails

### Individual Emails

Send emails to specific recipients with full control over all parameters:

```go
response, err := client.SendEmail(ctx, sendlix.MailOptions{
    From:    sendlix.EmailAddress{Email: "from@example.com", Name: "Sender Name"},
    To:      []sendlix.EmailAddress{{Email: "to@example.com", Name: "Recipient"}},
    CC:      []sendlix.EmailAddress{{Email: "cc@example.com"}},
    BCC:     []sendlix.EmailAddress{{Email: "bcc@example.com"}},
    Subject: "Important Message",
    ReplyTo: &sendlix.EmailAddress{Email: "reply@example.com"},
    Content: sendlix.MailContent{
        HTML:     "<h1>HTML Content</h1><p>This is HTML content.</p>",
        Text:     "Text Content\n\nThis is plain text content.",
        Tracking: true,
    },
}, &sendlix.AdditionalOptions{
    Category: "newsletter",
    SendAt:   &futureTime,
    Attachments: []sendlix.Attachment{{
        ContentURL:  "https://example.com/document.pdf",
        Filename:    "document.pdf",
        ContentType: "application/pdf",
    }},
})
```

### Group Emails

Send emails to predefined groups for bulk operations:

```go
response, err := client.SendGroupEmail(ctx, sendlix.GroupMailData{
    GroupID: "newsletter-subscribers",
    From:    sendlix.EmailAddress{Email: "news@example.com", Name: "Newsletter"},
    Subject: "Weekly Newsletter",
    Content: sendlix.MailContent{
        HTML: "<h1>This Week's News</h1><p>Stay updated with our latest news.</p>",
        Text: "This Week's News\n\nStay updated with our latest news.",
    },
    Category: "newsletter",
})
```

### EML Format Emails

Send pre-formatted EML messages:

```go
emlContent := []byte(`From: sender@example.com
To: recipient@example.com
Subject: Test Email

This is a test email message.`)

response, err := client.SendEMLEmail(ctx, emlContent, nil)
```

## Group Management

Manage email groups for bulk operations:

```go
// Create a group client
groupClient, err := sendlix.NewGroupClient(auth, nil)
if err != nil {
    log.Fatal(err)
}
defer groupClient.Close()

// Add emails to a group
emails := []sendlix.EmailData{
    {Email: "user1@example.com", Name: "User One"},
    {Email: "user2@example.com", Name: "User Two"},
}

substitutions := map[string]string{
    "company": "Example Corp",
    "product": "Amazing Product",
}

response, err := groupClient.InsertEmailToGroup(ctx, "my-group", emails, substitutions)
if err != nil {
    log.Fatal(err)
}

// Remove an email from a group
removeResponse, err := groupClient.RemoveEmailFromGroup(ctx, "my-group", "user1@example.com")

// Check if an email exists in a group
checkResponse, err := groupClient.CheckEmailInGroup(ctx, "my-group", "user2@example.com")
if err != nil {
    log.Fatal(err)
}
if checkResponse.Exists {
    log.Println("Email is in the group")
}
```

## Configuration

Customize client behavior with configuration options:

```go
config := &sendlix.ClientConfig{
    ServerAddress: "api.sendlix.com:443",
    UserAgent:     "MyApp/1.0.0",
    Insecure:      false, // Only set to true for testing
}

client, err := sendlix.NewEmailClient(auth, config)
```

## Error Handling

The SDK provides detailed error information:

```go
response, err := client.SendEmail(ctx, options, nil)
if err != nil {
    if strings.Contains(err.Error(), "authentication") {
        log.Println("Check your API key")
    } else if strings.Contains(err.Error(), "quota") {
        log.Println("Email quota exceeded")
    } else {
        log.Printf("Email send failed: %v", err)
    }
    return
}

log.Printf("Email sent successfully!")
log.Printf("Message IDs: %v", response.MessageList)
log.Printf("Emails remaining: %d", response.EmailsLeft)
```

## Context Support

All operations support Go contexts for timeout and cancellation:

```go
// Set a timeout for the operation
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

response, err := client.SendEmail(ctx, options, nil)
```

## Best Practices

1. **Resource Management**: Always call `Close()` on clients when done to prevent resource leaks
2. **Client Reuse**: Reuse clients across multiple operations rather than creating new ones
3. **Context Usage**: Use contexts with appropriate timeouts for network operations
4. **Error Handling**: Handle errors appropriately and check quota information in responses
5. **Bulk Operations**: Use group emails for bulk operations to improve performance
6. **Validation**: Validate email addresses before sending to avoid quota waste

## License

This SDK is licensed under the Apache License 2.0. See [LICENSE](LICENSE) for details.

## Support

For support and questions:

- 📚 [Documentation](https://docs.sendlix.com)
- 📧 [Email Support](mailto:info@sendlix.com)
- 🐛 [Report Issues](https://github.com/sendlix/go-sdk/issues)

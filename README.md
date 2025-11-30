# Sendlix Go SDK

The official Go SDK for the Sendlix email service API. This SDK provides a comprehensive interface for sending emails, managing email groups, and handling authentication with the Sendlix platform.

## Features

- **Email Sending**: Send individual emails with full control over recipients, content, and formatting
- **Group Management**: Manage email groups for efficient bulk email operations
- **Multiple Content Types**: Support for HTML, plain text, and EML format emails
- **Image Embedding**: Embed images directly in emails using placeholders
- **Advanced Features**: Email scheduling, attachments, tracking, and categorization
- **Authentication**: Automatic JWT token management with API key authentication
- **Error Handling**: Comprehensive error reporting
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

    sendlix "github.com/sendlix/go-sdk"
)

func main() {
    // Create an email client with your API key
    client, err := sendlix.NewEmailClient("your-secret.your-key-id", nil)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Send a simple email
    messageIDs, err := client.SendEmail(context.Background(), sendlix.MailOptions{
        From:    sendlix.EmailAddress{Email: "sender@example.com", Name: "Sender Name"},
        To:      []sendlix.EmailAddress{{Email: "recipient@example.com", Name: "Recipient"}},
        Subject: "Hello from Sendlix!",
        Html:    "<h1>Hello World!</h1><p>This is a test email.</p>",
        Text:    "Hello World!\n\nThis is a test email.",
    }, nil)

    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Email sent! Message IDs: %v", messageIDs)
}
```

## Authentication

All API operations require authentication using an API key from your Sendlix account. The API key format is `secret.keyID`. You can pass the API key directly as a string to the client constructors:

```go
client, err := sendlix.NewEmailClient("your-secret.123456", nil)
```

Alternatively, you can use the `NewAuth` function for more control:

```go
auth, err := sendlix.NewAuth("your-secret.123456")
if err != nil {
    log.Fatal(err)
}
client, err := sendlix.NewEmailClient(auth, nil)
```

The SDK automatically handles JWT token exchange and caching, so you don't need to manage tokens manually.

## Sending Emails

### Individual Emails

Send emails to specific recipients with full control over all parameters:

```go
messageIDs, err := client.SendEmail(ctx, sendlix.MailOptions{
    From:     sendlix.EmailAddress{Email: "from@example.com", Name: "Sender Name"},
    To:       []sendlix.EmailAddress{{Email: "to@example.com", Name: "Recipient"}},
    CC:       []sendlix.EmailAddress{{Email: "cc@example.com"}},
    BCC:      []sendlix.EmailAddress{{Email: "bcc@example.com"}},
    Subject:  "Important Message",
    ReplyTo:  &sendlix.EmailAddress{Email: "reply@example.com"},
    Html:     "<h1>HTML Content</h1><p>This is HTML content.</p>",
    Text:     "Text Content\n\nThis is plain text content.",
    Tracking: true,
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

### Emails with Embedded Images

Embed images directly in your HTML emails using placeholders:

```go
// Load image data
logoData, err := os.ReadFile("logo.png")
if err != nil {
    log.Fatal(err)
}

messageIDs, err := client.SendEmail(ctx, sendlix.MailOptions{
    From:    sendlix.EmailAddress{Email: "from@example.com", Name: "Sender"},
    To:      []sendlix.EmailAddress{{Email: "to@example.com"}},
    Subject: "Welcome Email",
    Html:    "<h1>Welcome!</h1><img src=\"{{logo}}\"><p>Thanks for joining us.</p>",
    Text:    "Welcome! Thanks for joining us.",
    Images: []sendlix.Image{
        {
            Placeholder: "{{logo}}",
            Data:        logoData,
            Type:        sendlix.MimeTypePNG,
        },
    },
}, nil)
```

Supported image types:

- `sendlix.MimeTypePNG` - PNG images
- `sendlix.MimeTypeJPEG` - JPEG images
- `sendlix.MimeTypeGIF` - GIF images

### Group Emails

Send emails to predefined groups for bulk operations:

```go
err := client.SendGroupEmail(ctx, sendlix.GroupMailData{
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

messageIDs, err := client.SendEMLEmail(ctx, emlContent, nil)
```

## Group Management

Manage email groups for bulk operations:

```go
// Create a group client
groupClient, err := sendlix.NewGroupClient("your-secret.your-key-id", nil)
if err != nil {
    log.Fatal(err)
}
defer groupClient.Close()

// Add emails to a group with individual substitutions
entries := []sendlix.GroupEntry{
    {
        Email: "user1@example.com",
        Name:  "User One",
        Substitutions: map[string]string{
            "first_name": "User",
            "discount":   "10%",
        },
    },
    {
        Email: "user2@example.com",
        Name:  "User Two",
        Substitutions: map[string]string{
            "first_name": "Another",
            "discount":   "20%",
        },
    },
}

response, err := groupClient.InsertEmailsToGroup(ctx, "my-group", entries, nil)
if err != nil {
    log.Fatal(err)
}
log.Printf("Added %d emails to group", response.AffectedRows)

// Add a single email to a group
response, err = groupClient.InsertEmailToGroup(ctx, "my-group", sendlix.GroupEntry{
    Email: "newuser@example.com",
    Name:  "New User",
    Substitutions: map[string]string{
        "welcome_bonus": "15%",
    },
})

// With failure handling options
response, err = groupClient.InsertEmailsToGroup(ctx, "my-group", entries,
    &sendlix.InsertOptions{
        OnFailure: sendlix.FailureHandlerAbort, // or sendlix.FailureHandlerSkip
    })

// Remove an email from a group
removeResponse, err := groupClient.RemoveEmailFromGroup(ctx, "my-group", "user1@example.com")

// Check if an email exists in a group
exists, err := groupClient.CheckEmailInGroup(ctx, "my-group", "user2@example.com")
if err != nil {
    log.Fatal(err)
}
if exists {
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

client, err := sendlix.NewEmailClient("your-secret.your-key-id", config)
```

## Error Handling

The SDK provides detailed error information:

```go
messageIDs, err := client.SendEmail(ctx, options, nil)
if err != nil {
    if strings.Contains(err.Error(), "authentication") {
        log.Println("Check your API key")
    } else if strings.Contains(err.Error(), "from email is required") {
        log.Println("Missing sender address")
    } else {
        log.Printf("Email send failed: %v", err)
    }
    return
}

log.Printf("Email sent successfully!")
log.Printf("Message IDs: %v", messageIDs)
```

## Context Support

All operations support Go contexts for timeout and cancellation:

```go
// Set a timeout for the operation
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

messageIDs, err := client.SendEmail(ctx, options, nil)
```

## Best Practices

1. **Resource Management**: Always call `Close()` on clients when done to prevent resource leaks
2. **Client Reuse**: Reuse clients across multiple operations rather than creating new ones
3. **Context Usage**: Use contexts with appropriate timeouts for network operations
4. **Error Handling**: Handle errors appropriately for your use case
5. **Bulk Operations**: Use group emails for bulk operations to improve performance
6. **Validation**: Validate email addresses before sending to avoid errors
7. **Image Optimization**: Optimize images before embedding to reduce email size

## License

This SDK is licensed under the Apache License 2.0. See [LICENSE](LICENSE) for details.

## Support

For support and questions:

- 📚 [Documentation](https://docs.sendlix.com)
- 📧 [Email Support](mailto:info@sendlix.com)
- 🐛 [Report Issues](https://github.com/sendlix/go-sdk/issues)

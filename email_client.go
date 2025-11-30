package sendlix

import (
	"context"
	"fmt"
	"time"

	pb "github.com/sendlix/go-sdk/internal/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// EmailClient provides comprehensive email sending functionality through the Sendlix API.
// It supports various email types including individual emails, group emails, and EML format emails.
// The client handles authentication, request formatting, and response parsing automatically.
//
// EmailClient embeds BaseClient, inheriting connection management and authentication capabilities.
// All email operations require proper authentication through the configured IAuth implementation.
type EmailClient struct {
	*BaseClient
	client pb.EmailClient
}

// NewEmailClient creates a new email client with the provided authentication and configuration.
// The client establishes a gRPC connection to the Sendlix email service and is ready for immediate use.
//
// Parameters:
//   - auth: Authentication - either an IAuth implementation or an API key string (required)
//   - config: Client configuration (optional, uses defaults if nil)
//
// Returns:
//   - *EmailClient: Configured email client
//   - error: Any error encountered during client creation
//
// Example with API key string:
//
//	client, err := sendlix.NewEmailClient("secret.keyid", nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer client.Close()
//
// Example with IAuth:
//
//	auth, err := sendlix.NewAuth("secret.keyid")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	client, err := sendlix.NewEmailClient(auth, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer client.Close()
func NewEmailClient(auth interface{}, config *ClientConfig) (*EmailClient, error) {
	resolvedAuth, err := resolveAuth(auth)
	if err != nil {
		return nil, err
	}

	baseClient, err := NewBaseClient(resolvedAuth, config)
	if err != nil {
		return nil, err
	}

	return &EmailClient{
		BaseClient: baseClient,
		client:     pb.NewEmailClient(baseClient.GetConnection()),
	}, nil
}

// EmailAddress represents an email address with an optional display name.
// It provides a convenient way to specify email recipients with human-readable names.
//
// The display name is optional and when provided, creates email addresses in the format
// "Display Name <email@domain.com>". When omitted, only the email address is used.
type EmailAddress struct {
	// Email is the email address (required)
	Email string
	// Name is the optional display name for the email address
	Name string
}

// String returns a properly formatted string representation of the email address.
// If a display name is provided, it returns "Name <email@domain.com>".
// Otherwise, it returns just the email address.
//
// Returns:
//   - string: Formatted email address string
//
// Example:
//
//	addr := EmailAddress{Email: "john@example.com", Name: "John Doe"}
//	fmt.Println(addr.String()) // Output: "John Doe <john@example.com>"
//
//	addr2 := EmailAddress{Email: "jane@example.com"}
//	fmt.Println(addr2.String()) // Output: "jane@example.com"
func (e EmailAddress) String() string {
	if e.Name != "" {
		return fmt.Sprintf("%s <%s>", e.Name, e.Email)
	}
	return e.Email
}

// NewEmailAddress creates an EmailAddress from various input types.
// This function provides flexible email address creation from strings or existing EmailAddress values.
//
// Supported input types:
//   - string: Creates EmailAddress with Email field set, Name empty
//   - EmailAddress: Returns a copy of the provided EmailAddress
//   - *EmailAddress: Returns the pointer directly
//
// Parameters:
//   - addr: Email address as string, EmailAddress, or *EmailAddress
//
// Returns:
//   - *EmailAddress: Created email address
//   - error: Type conversion error for unsupported input types
//
// Example:
//
//	// From string
//	addr1, err := sendlix.NewEmailAddress("user@example.com")
//
//	// From EmailAddress struct
//	addr2, err := sendlix.NewEmailAddress(sendlix.EmailAddress{
//		Email: "user@example.com",
//		Name:  "User Name",
//	})
func NewEmailAddress(addr interface{}) (*EmailAddress, error) {
	switch v := addr.(type) {
	case string:
		return &EmailAddress{Email: v}, nil
	case EmailAddress:
		return &v, nil
	case *EmailAddress:
		return v, nil
	default:
		return nil, fmt.Errorf("invalid email address type: %T", addr)
	}
}

// MailContent represents the content and formatting options for an email message.
// It supports both HTML and plain text content, allowing for rich email formatting
// while maintaining compatibility with text-only email clients.
// Deprecated: Use Html, Text, Tracking and Images fields directly in MailOptions instead.
type MailContent struct {
	// HTML content of the email (optional)
	// Should contain valid HTML markup for rich formatting
	HTML string

	// Text content of the email (optional)
	// Plain text version for email clients that don't support HTML
	Text string

	// Tracking enables email tracking features such as open tracking
	// and click tracking when supported by the email service
	Tracking bool
}

// MimeType represents the MIME type of an embedded image.
type MimeType int

const (
	// MimeTypePNG represents PNG image format
	MimeTypePNG MimeType = iota
	// MimeTypeJPEG represents JPEG image format
	MimeTypeJPEG
	// MimeTypeGIF represents GIF image format
	MimeTypeGIF
)

// Image represents an embedded image for email content.
// Images can be embedded directly in the HTML content using placeholders.
// The placeholder in the HTML will be replaced with the actual image data.
type Image struct {
	// Placeholder is the string in HTML that will be replaced with the image
	// Example: "{{logo}}" in HTML becomes the actual image
	Placeholder string

	// Data contains the raw image bytes
	Data []byte

	// Type specifies the MIME type of the image (PNG, JPEG, or GIF)
	Type MimeType
}

// Attachment represents a file attachment for email messages.
// Attachments are referenced by URL and include metadata for proper handling.
type Attachment struct {
	// ContentURL is the URL where the attachment content can be retrieved
	ContentURL string

	// Filename is the name that will be shown for the attachment
	Filename string

	// ContentType is the MIME type of the attachment (e.g., "application/pdf")
	ContentType string
}

// MailOptions contains all the required and optional parameters for sending an email.
// This structure provides a comprehensive way to specify email details including
// recipients, content, and various email headers.
type MailOptions struct {
	// From specifies the sender's email address (required)
	From EmailAddress

	// To contains the list of primary recipients (required, at least one)
	To []EmailAddress

	// CC contains the list of carbon copy recipients (optional)
	CC []EmailAddress

	// BCC contains the list of blind carbon copy recipients (optional)
	BCC []EmailAddress

	// Subject is the email subject line (required)
	Subject string

	// ReplyTo specifies the email address for replies (optional)
	// If not set, replies will go to the From address
	ReplyTo *EmailAddress

	// Html content of the email (optional)
	// Should contain valid HTML markup for rich formatting
	Html string

	// Text content of the email (optional)
	// Plain text version for email clients that don't support HTML
	Text string

	// Tracking enables email tracking features such as open tracking
	// and click tracking when supported by the email service
	Tracking bool

	// Images contains embedded images for the email content (optional)
	// Images are embedded using placeholders in the HTML content
	Images []Image
}

// AdditionalOptions provides extended configuration options for email sending.
// These options allow for advanced features like scheduling and file attachments.
type AdditionalOptions struct {
	// Attachments is a list of files to attach to the email (optional)
	Attachments []Attachment

	// Category is used for email categorization and analytics (optional)
	Category string

	// SendAt schedules the email to be sent at a specific time (optional)
	// If nil, the email is sent immediately
	SendAt *time.Time
}

// GroupMailData represents the data structure for sending emails to predefined groups.
// This is used for bulk email operations where recipients are managed as groups.
type GroupMailData struct {
	// From specifies the sender's email address (required)
	From EmailAddress

	// GroupID identifies the recipient group (required)
	GroupID string

	// Subject is the email subject line (required)
	Subject string

	// Category is used for email categorization and analytics (optional)
	// Category is used for email categorization and analytics (optional)
	Category string

	// Content contains the email body and formatting options (required)
	Content MailContent
}

// SendEmail sends an email with the specified options and returns the result.
// This is the primary method for sending individual emails through the Sendlix API.
// It validates all required fields and handles the complete send process.
//
// The method performs comprehensive validation of email parameters including:
// - From email address presence
// - At least one recipient in To field
// - Subject line presence
// - Either HTML or text content presence
//
// Parameters:
//   - ctx: Context for the request (supports cancellation and timeouts)
//   - options: Email configuration including recipients, subject, and content
//   - additional: Optional advanced settings like attachments and scheduling
//
// Returns:
//   - []string: List of message IDs for the sent emails
//   - error: Validation or sending error
//
// Example:
//
//	response, err := client.SendEmail(ctx, sendlix.MailOptions{
//		From:     sendlix.EmailAddress{Email: "sender@example.com", Name: "Sender"},
//		To:       []sendlix.EmailAddress{{Email: "recipient@example.com"}},
//		Subject:  "Hello World",
//		Html:     "<h1>Hello World</h1>",
//		Text:     "Hello World",
//		Tracking: true,
//	}, &sendlix.AdditionalOptions{
//		Category: "newsletter",
//	})
//
// Example with embedded images:
//
//	imageData, _ := os.ReadFile("logo.png")
//	response, err := client.SendEmail(ctx, sendlix.MailOptions{
//		From:    sendlix.EmailAddress{Email: "sender@example.com"},
//		To:      []sendlix.EmailAddress{{Email: "recipient@example.com"}},
//		Subject: "Email with Logo",
//		Html:    "<h1>Welcome</h1><img src=\"{{logo}}\">",
//		Images: []sendlix.Image{
//			{Placeholder: "{{logo}}", Data: imageData, Type: sendlix.MimeTypePNG},
//		},
//	}, nil)
//
// Common errors:
//   - Missing required fields (from, to, subject, content)
//   - Invalid email addresses
//   - Authentication failures
//   - Network connectivity issues
func (c *EmailClient) SendEmail(ctx context.Context, options MailOptions, additional *AdditionalOptions) ([]string, error) {
	// Validate required fields
	if options.From.Email == "" {
		return nil, fmt.Errorf("from email is required")
	}
	if len(options.To) == 0 {
		return nil, fmt.Errorf("at least one recipient is required")
	}
	if options.Subject == "" {
		return nil, fmt.Errorf("subject is required")
	}
	if options.Html == "" && options.Text == "" {
		return nil, fmt.Errorf("either HTML or text content is required")
	}

	// Build mail content
	mailContent := &pb.MailContent{
		Html:     options.Html,
		Text:     options.Text,
		Tracking: options.Tracking,
	}

	// Add images if provided
	if len(options.Images) > 0 {
		mailContent.Images = convertImages(options.Images)
	}

	// Build request
	req := &pb.SendMailRequest{
		From:    convertEmailAddress(options.From),
		To:      convertEmailAddressList(options.To),
		Subject: options.Subject,
		Body: &pb.SendMailRequest_TextContent{
			TextContent: mailContent,
		},
	}

	// Add optional fields
	if len(options.CC) > 0 {
		req.Cc = convertEmailAddressList(options.CC)
	}
	if len(options.BCC) > 0 {
		req.Bcc = convertEmailAddressList(options.BCC)
	}
	if options.ReplyTo != nil {
		req.ReplyTo = convertEmailAddress(*options.ReplyTo)
	}

	// Add additional options
	if additional != nil {
		req.AdditionalInfos = convertAdditionalOptions(additional)
	}

	// Send request
	resp, err := c.client.SendEmail(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to send email: %v", err)
	}

	return resp.Message, nil
}

// SendEMLEmail sends an email using EML (Email Message Format) data.
// This method allows sending pre-formatted email messages that follow
// the RFC 5322 standard for email message format.
//
// EML format is useful when you have complete email messages that were
// previously generated or when integrating with other email systems
// that produce standard EML output.
//
// Parameters:
//   - ctx: Context for the request (supports cancellation and timeouts)
//   - emlData: Complete EML message data as byte array
//   - additional: Optional settings like scheduling and categorization
//
// Returns:
//   - []string: List of message IDs for the sent emails
//   - error: Parsing or sending error
//
// Example:
//
//	emlContent := []byte(`From: sender@example.com
//	To: recipient@example.com
//	Subject: Test Email
//
//	This is a test email message.`)
//
//	response, err := client.SendEMLEmail(ctx, emlContent, nil)
//
// The EML data should be a complete, valid email message including headers
// and body. Invalid EML format will result in parsing errors.
func (c *EmailClient) SendEMLEmail(ctx context.Context, emlData []byte, additional *AdditionalOptions) ([]string, error) {
	req := &pb.EmlMailRequest{
		Mail: emlData,
	}

	if additional != nil {
		req.AdditionalInfos = convertAdditionalOptions(additional)
	}

	resp, err := c.client.SendEmlEmail(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to send EML email: %v", err)
	}

	return resp.Message, nil
}

// SendGroupEmail sends an email to all members of a predefined group.
// This method is optimized for bulk email operations where recipients
// are managed as groups rather than individual addresses.
//
// Group emails are useful for newsletters, announcements, and other
// communications sent to large numbers of recipients. The group must
// be created and populated before sending emails to it.
//
// Parameters:
//   - ctx: Context for the request (supports cancellation and timeouts)
//   - data: Group email configuration including group ID and content
//
// Returns:
//   - error: Validation or sending error, nil on success
//
// Example:
//
//	err := client.SendGroupEmail(ctx, sendlix.GroupMailData{
//		GroupID: "newsletter-subscribers",
//		From:    sendlix.EmailAddress{Email: "news@example.com", Name: "Newsletter"},
//		Subject: "Weekly Newsletter",
//		Content: sendlix.MailContent{
//			HTML: "<h1>This Week's News</h1><p>...</p>",
//			Text: "This Week's News\n\n...",
//		},
//		Category: "newsletter",
//	})
//
// The group must exist and contain email addresses before calling this method.
// Empty groups will not generate an error but will result in zero emails sent.
func (c *EmailClient) SendGroupEmail(ctx context.Context, data GroupMailData) error {
	if data.GroupID == "" {
		return fmt.Errorf("group ID is required")
	}
	if data.From.Email == "" {
		return fmt.Errorf("from email is required")
	}
	if data.Subject == "" {
		return fmt.Errorf("subject is required")
	}
	if data.Content.HTML == "" && data.Content.Text == "" {
		return fmt.Errorf("either HTML or text content is required")
	}

	req := &pb.GroupMailData{
		GroupId:  data.GroupID,
		Subject:  data.Subject,
		From:     convertEmailAddress(data.From),
		Category: data.Category,
		Body: &pb.GroupMailData_TextContent{
			TextContent: &pb.MailContent{
				Html:     data.Content.HTML,
				Text:     data.Content.Text,
				Tracking: data.Content.Tracking,
			},
		},
	}

	_, err := c.client.SendGroupEmail(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to send group email: %v", err)
	}

	return nil
}

// Helper functions for converting between SDK types and protobuf types

// convertEmailAddress converts an EmailAddress to the protobuf EmailData format.
// This helper function is used internally to transform SDK types to the format
// expected by the gRPC API.
//
// Parameters:
//   - addr: EmailAddress to convert
//
// Returns:
//   - *pb.EmailData: Protobuf representation of the email address
func convertEmailAddress(addr EmailAddress) *pb.EmailData {
	return &pb.EmailData{
		Email: addr.Email,
		Name:  addr.Name,
	}
}

// convertEmailAddressList converts a slice of EmailAddress to protobuf EmailData slice.
// This helper function transforms multiple email addresses for batch operations.
//
// Parameters:
//   - addrs: Slice of EmailAddress to convert
//
// Returns:
//   - []*pb.EmailData: Slice of protobuf EmailData representations
func convertEmailAddressList(addrs []EmailAddress) []*pb.EmailData {
	result := make([]*pb.EmailData, len(addrs))
	for i, addr := range addrs {
		result[i] = convertEmailAddress(addr)
	}
	return result
}

// convertAdditionalOptions converts AdditionalOptions to protobuf AdditionalInfos format.
// This helper function handles the transformation of advanced email options including
// attachments, scheduling, and categorization settings.
//
// Parameters:
//   - opts: AdditionalOptions to convert
//
// Returns:
//   - *pb.AdditionalInfos: Protobuf representation of additional options
func convertAdditionalOptions(opts *AdditionalOptions) *pb.AdditionalInfos {
	info := &pb.AdditionalInfos{
		Category: opts.Category,
	}

	if len(opts.Attachments) > 0 {
		info.Attachments = make([]*pb.AttachmentData, len(opts.Attachments))
		for i, att := range opts.Attachments {
			info.Attachments[i] = &pb.AttachmentData{
				ContentUrl: att.ContentURL,
				Type:       att.ContentType,
				Filename:   att.Filename,
			}
		}
	}

	if opts.SendAt != nil {
		info.SendAt = timestamppb.New(*opts.SendAt)
	}

	return info
}

// convertImages converts a slice of Image to protobuf Images slice.
// This helper function transforms embedded images for email content.
//
// Parameters:
//   - images: Slice of Image to convert
//
// Returns:
//   - []*pb.Images: Slice of protobuf Images representations
func convertImages(images []Image) []*pb.Images {
	result := make([]*pb.Images, len(images))
	for i, img := range images {
		result[i] = &pb.Images{
			Placeholder: img.Placeholder,
			Image:       img.Data,
			Type:        pb.MimeType(img.Type),
		}
	}
	return result
}

// resolveAuth converts an auth parameter to an IAuth implementation.
// It accepts either an IAuth implementation directly or an API key string.
//
// Parameters:
//   - auth: Either an IAuth implementation or an API key string
//
// Returns:
//   - IAuth: Resolved authentication implementation
//   - error: Error if the auth type is invalid or API key parsing fails
func resolveAuth(auth interface{}) (IAuth, error) {
	switch v := auth.(type) {
	case IAuth:
		return v, nil
	case string:
		return NewAuth(v)
	default:
		return nil, fmt.Errorf("invalid auth type: %T, expected IAuth or string", auth)
	}
}

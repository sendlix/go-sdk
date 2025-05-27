package sendlix

import (
	"context"
	"fmt"

	pb "github.com/sendlix/go-sdk/internal/proto"
)

// GroupClient provides comprehensive group management functionality for the Sendlix API.
// It enables creating, managing, and manipulating email groups that can be used
// for bulk email operations. Groups are collections of email addresses that can
// be managed as a single unit for efficient mass communication.
//
// GroupClient embeds BaseClient, inheriting connection management and authentication capabilities.
// All group operations require proper authentication through the configured IAuth implementation.
type GroupClient struct {
	*BaseClient
	client pb.GroupClient
}

// NewGroupClient creates a new group management client with the provided authentication and configuration.
// The client establishes a gRPC connection to the Sendlix group service and is ready for immediate use.
//
// Parameters:
//   - auth: Authentication implementation (required)
//   - config: Client configuration (optional, uses defaults if nil)
//
// Returns:
//   - *GroupClient: Configured group client
//   - error: Any error encountered during client creation
//
// Example:
//
//	auth, err := sendlix.NewAuth("secret.keyid")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	client, err := sendlix.NewGroupClient(auth, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer client.Close()
func NewGroupClient(auth IAuth, config *ClientConfig) (*GroupClient, error) {
	baseClient, err := NewBaseClient(auth, config)
	if err != nil {
		return nil, err
	}

	return &GroupClient{
		BaseClient: baseClient,
		client:     pb.NewGroupClient(baseClient.GetConnection()),
	}, nil
}

// EmailData represents email address information for group operations.
// This structure is used when adding or managing email addresses within groups,
// allowing for both the email address and an optional display name.
type EmailData struct {
	// Email is the email address (required)
	Email string
	// Name is the optional display name for the email address
	Name string
}

// InsertEmailToGroupResponse represents the result of adding emails to a group.
// It provides detailed information about the operation's success and impact.
type InsertEmailToGroupResponse struct {
	// Success indicates whether the operation completed successfully
	Success bool
	// Message provides additional details about the operation result
	Message string
	// AffectedRows indicates how many email addresses were successfully added
	AffectedRows int64
}

// RemoveEmailFromGroupResponse represents the result of removing an email from a group.
// It provides detailed information about the removal operation's success and impact.
type RemoveEmailFromGroupResponse struct {
	// Success indicates whether the operation completed successfully
	Success bool
	// Message provides additional details about the operation result
	Message string
	// AffectedRows indicates how many email addresses were successfully removed
	AffectedRows int64
}

// CheckEmailInGroupResponse represents the result of checking email group membership.
// It provides a simple boolean result indicating whether the email exists in the group.
type CheckEmailInGroupResponse struct {
	// Exists indicates whether the email address is present in the group
	Exists bool
}

// InsertEmailToGroup inserts one or multiple emails into a specified group with optional substitutions.
// This method allows bulk addition of email addresses to groups, making it efficient for
// managing large subscriber lists. Each email can have associated substitution variables
// for personalized group communications.
//
// The method validates all input parameters including group ID presence and email address
// validity for each entry. All emails must have valid email addresses, while display names
// are optional.
//
// Parameters:
//   - ctx: Context for the request (supports cancellation and timeouts)
//   - groupID: Identifier of the target group (required)
//   - emails: Slice of email data to add to the group (required, at least one)
//   - substitutions: Optional key-value pairs for email personalization
//
// Returns:
//   - *InsertEmailToGroupResponse: Operation result with success status and affected rows
//   - error: Validation or operation error
//
// Example:
//
//	emails := []sendlix.EmailData{
//		{Email: "user1@example.com", Name: "User One"},
//		{Email: "user2@example.com", Name: "User Two"},
//	}
//	substitutions := map[string]string{
//		"company": "Example Corp",
//		"product": "Amazing Product",
//	}
//
//	response, err := client.InsertEmailToGroup(ctx, "newsletter-group", emails, substitutions)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Added %d emails successfully\n", response.AffectedRows)
//
// Common errors:
//   - Empty group ID
//   - Empty email list
//   - Invalid email addresses
//   - Group not found
//   - Permission denied
func (c *GroupClient) InsertEmailToGroup(ctx context.Context, groupID string, emails []EmailData, substitutions map[string]string) (*InsertEmailToGroupResponse, error) {
	if groupID == "" {
		return nil, fmt.Errorf("group ID is required")
	}
	if len(emails) == 0 {
		return nil, fmt.Errorf("at least one email is required")
	}

	// Convert emails to protobuf format
	pbEmails := make([]*pb.EmailData, len(emails))
	for i, email := range emails {
		if email.Email == "" {
			return nil, fmt.Errorf("email address is required for email at index %d", i)
		}
		pbEmails[i] = &pb.EmailData{
			Email: email.Email,
			Name:  email.Name,
		}
	}

	req := &pb.InsertEmailToGroupRequest{
		Emails:        pbEmails,
		GroupId:       groupID,
		Substitutions: substitutions,
	}

	resp, err := c.client.InsertEmailToGroup(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to insert emails to group: %v", err)
	}

	return &InsertEmailToGroupResponse{
		Success:      resp.Success,
		Message:      resp.Message,
		AffectedRows: resp.AffectedRows,
	}, nil
}

// InsertSingleEmailToGroup inserts a single email into a group with optional substitutions.
// This is a convenience method that wraps InsertEmailToGroup for single email operations,
// providing a simpler interface when only one email address needs to be added.
//
// Parameters:
//   - ctx: Context for the request (supports cancellation and timeouts)
//   - groupID: Identifier of the target group (required)
//   - email: Email data to add to the group (required)
//   - substitutions: Optional key-value pairs for email personalization
//
// Returns:
//   - *InsertEmailToGroupResponse: Operation result with success status
//   - error: Validation or operation error
//
// Example:
//
//	email := sendlix.EmailData{
//		Email: "newuser@example.com",
//		Name:  "New User",
//	}
//	substitutions := map[string]string{
//		"welcome_bonus": "20% off",
//	}
//
//	response, err := client.InsertSingleEmailToGroup(ctx, "customers", email, substitutions)
func (c *GroupClient) InsertSingleEmailToGroup(ctx context.Context, groupID string, email EmailData, substitutions map[string]string) (*InsertEmailToGroupResponse, error) {
	return c.InsertEmailToGroup(ctx, groupID, []EmailData{email}, substitutions)
}

// RemoveEmailFromGroup removes a specific email address from a group.
// This method provides targeted removal of individual email addresses from groups,
// useful for handling unsubscribes or managing group membership.
//
// The operation is idempotent - removing an email that doesn't exist in the group
// will not cause an error but will result in zero affected rows.
//
// Parameters:
//   - ctx: Context for the request (supports cancellation and timeouts)
//   - groupID: Identifier of the target group (required)
//   - email: Email address to remove from the group (required)
//
// Returns:
//   - *RemoveEmailFromGroupResponse: Operation result with success status and affected rows
//   - error: Validation or operation error
//
// Example:
//
//	response, err := client.RemoveEmailFromGroup(ctx, "newsletter-group", "user@example.com")
//	if err != nil {
//		log.Fatal(err)
//	}
//	if response.AffectedRows > 0 {
//		fmt.Println("Email successfully removed from group")
//	} else {
//		fmt.Println("Email was not found in group")
//	}
//
// Common errors:
//   - Empty group ID
//   - Empty email address
//   - Group not found
//   - Permission denied
func (c *GroupClient) RemoveEmailFromGroup(ctx context.Context, groupID string, email string) (*RemoveEmailFromGroupResponse, error) {
	if groupID == "" {
		return nil, fmt.Errorf("group ID is required")
	}
	if email == "" {
		return nil, fmt.Errorf("email address is required")
	}

	req := &pb.RemoveEmailFromGroupRequest{
		Email:   email,
		GroupId: groupID,
	}

	resp, err := c.client.RemoveEmailFromGroup(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to remove email from group: %v", err)
	}

	return &RemoveEmailFromGroupResponse{
		Success:      resp.Success,
		Message:      resp.Message,
		AffectedRows: resp.AffectedRows,
	}, nil
}

// CheckEmailInGroup checks whether a specific email address exists in a group.
// This method provides a simple way to verify group membership before performing
// other operations like sending group emails or managing subscriptions.
//
// The check is performed efficiently on the server side and returns a boolean
// result indicating membership status.
//
// Parameters:
//   - ctx: Context for the request (supports cancellation and timeouts)
//   - groupID: Identifier of the target group (required)
//   - email: Email address to check for membership (required)
//
// Returns:
//   - *CheckEmailInGroupResponse: Result containing membership status
//   - error: Validation or operation error
//
// Example:
//
//	response, err := client.CheckEmailInGroup(ctx, "newsletter-group", "user@example.com")
//	if err != nil {
//		log.Fatal(err)
//	}
//	if response.Exists {
//		fmt.Println("Email is subscribed to the group")
//	} else {
//		fmt.Println("Email is not in the group")
//	}
//
// This method is useful for:
//   - Preventing duplicate subscriptions
//   - Verifying membership before group operations
//   - Implementing subscription status checks
//   - Building user interfaces that show group membership
//
// Common errors:
//   - Empty group ID
//   - Empty email address
//   - Group not found
//   - Permission denied
func (c *GroupClient) CheckEmailInGroup(ctx context.Context, groupID string, email string) (*CheckEmailInGroupResponse, error) {
	if groupID == "" {
		return nil, fmt.Errorf("group ID is required")
	}
	if email == "" {
		return nil, fmt.Errorf("email address is required")
	}

	req := &pb.CheckEmailInGroupRequest{
		Email:   email,
		GroupId: groupID,
	}

	resp, err := c.client.CheckEmailInGroup(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to check email in group: %v", err)
	}

	return &CheckEmailInGroupResponse{
		Exists: resp.Exists,
	}, nil
}

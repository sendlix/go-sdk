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
//   - auth: Authentication - either an IAuth implementation or an API key string (required)
//   - config: Client configuration (optional, uses defaults if nil)
//
// Returns:
//   - *GroupClient: Configured group client
//   - error: Any error encountered during client creation
//
// Example with API key string:
//
//	client, err := sendlix.NewGroupClient("secret.keyid", nil)
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
//	client, err := sendlix.NewGroupClient(auth, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer client.Close()
func NewGroupClient(auth interface{}, config *ClientConfig) (*GroupClient, error) {
	resolvedAuth, err := resolveAuth(auth)
	if err != nil {
		return nil, err
	}

	baseClient, err := NewBaseClient(resolvedAuth, config)
	if err != nil {
		return nil, err
	}

	return &GroupClient{
		BaseClient: baseClient,
		client:     pb.NewGroupClient(baseClient.GetConnection()),
	}, nil
}

// FailureHandler defines how to handle failures when inserting multiple emails into a group.
type FailureHandler int

const (
	// FailureHandlerSkip skips failed entries and continues with the remaining entries
	FailureHandlerSkip FailureHandler = iota
	// FailureHandlerAbort aborts the entire operation on the first failure
	FailureHandlerAbort
)

// GroupEntry represents an email entry with optional substitution variables for a group.
// Each entry contains an email address with optional display name and personalization data.
type GroupEntry struct {
	// Email is the email address (required)
	Email string
	// Name is the optional display name for the email address
	Name string
	// Substitutions contains key-value pairs for email personalization (optional)
	Substitutions map[string]string
}

// InsertOptions provides configuration options for inserting emails into a group.
type InsertOptions struct {
	// OnFailure defines how to handle failures during bulk insert operations
	// Default is FailureHandlerSkip
	OnFailure FailureHandler
}

// UpdateResponse represents the result of a group update operation.
// It provides detailed information about the operation's success and impact.
type UpdateResponse struct {
	// Success indicates whether the operation completed successfully
	Success bool
	// Message provides additional details about the operation result
	Message string
	// AffectedRows indicates how many entries were successfully processed
	AffectedRows int64
}

// InsertEmailsToGroup inserts one or multiple email entries into a specified group.
// Each entry can have its own substitution variables for personalized group communications.
//
// This method allows bulk addition of email addresses to groups, making it efficient for
// managing large subscriber lists. The OnFailure option controls whether to skip failed
// entries or abort the entire operation.
//
// Parameters:
//   - ctx: Context for the request (supports cancellation and timeouts)
//   - groupID: Identifier of the target group (required)
//   - entries: Slice of group entries to add (required, at least one)
//   - options: Optional configuration for the insert operation
//
// Returns:
//   - *UpdateResponse: Operation result with success status and affected rows
//   - error: Validation or operation error
//
// Example:
//
//	entries := []sendlix.GroupEntry{
//		{
//			Email: "user1@example.com",
//			Name:  "User One",
//			Substitutions: map[string]string{
//				"first_name": "User",
//				"discount":   "10%",
//			},
//		},
//		{
//			Email: "user2@example.com",
//			Name:  "User Two",
//			Substitutions: map[string]string{
//				"first_name": "User",
//				"discount":   "20%",
//			},
//		},
//	}
//
//	response, err := client.InsertEmailsToGroup(ctx, "newsletter-group", entries, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Added %d emails successfully\n", response.AffectedRows)
//
// Example with failure handling:
//
//	response, err := client.InsertEmailsToGroup(ctx, "newsletter-group", entries,
//		&sendlix.InsertOptions{OnFailure: sendlix.FailureHandlerAbort})
func (c *GroupClient) InsertEmailsToGroup(ctx context.Context, groupID string, entries []GroupEntry, options *InsertOptions) (*UpdateResponse, error) {
	if groupID == "" {
		return nil, fmt.Errorf("group ID is required")
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("at least one entry is required")
	}

	// Convert entries to protobuf format
	pbEntries := make([]*pb.GroupEntry, len(entries))
	for i, entry := range entries {
		if entry.Email == "" {
			return nil, fmt.Errorf("email address is required for entry at index %d", i)
		}
		pbEntries[i] = &pb.GroupEntry{
			Email: &pb.EmailData{
				Email: entry.Email,
				Name:  entry.Name,
			},
			Substitutions: entry.Substitutions,
		}
	}

	req := &pb.InsertEmailToGroupRequest{
		Entries: pbEntries,
		GroupId: groupID,
	}

	// Set failure handler if options provided
	if options != nil {
		req.OnFailure = pb.FailureHandler(options.OnFailure)
	}

	resp, err := c.client.InsertEmailToGroup(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to insert emails to group: %v", err)
	}

	return &UpdateResponse{
		Success:      resp.Success,
		Message:      resp.Message,
		AffectedRows: resp.AffectedRows,
	}, nil
}

// InsertEmailToGroup inserts a single email into a group with optional substitutions.
// This is a convenience method for adding a single email address to a group.
//
// Parameters:
//   - ctx: Context for the request (supports cancellation and timeouts)
//   - groupID: Identifier of the target group (required)
//   - entry: Group entry to add (required)
//
// Returns:
//   - *UpdateResponse: Operation result with success status
//   - error: Validation or operation error
//
// Example:
//
//	entry := sendlix.GroupEntry{
//		Email: "newuser@example.com",
//		Name:  "New User",
//		Substitutions: map[string]string{
//			"welcome_bonus": "20% off",
//		},
//	}
//
//	response, err := client.InsertEmailToGroup(ctx, "customers", entry)
func (c *GroupClient) InsertEmailToGroup(ctx context.Context, groupID string, entry GroupEntry) (*UpdateResponse, error) {
	return c.InsertEmailsToGroup(ctx, groupID, []GroupEntry{entry}, nil)
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
//   - *UpdateResponse: Operation result with success status and affected rows
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
func (c *GroupClient) RemoveEmailFromGroup(ctx context.Context, groupID string, email string) (*UpdateResponse, error) {
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

	return &UpdateResponse{
		Success:      resp.Success,
		Message:      resp.Message,
		AffectedRows: resp.AffectedRows,
	}, nil
}

// CheckEmailInGroup checks whether a specific email address exists in a group.
// This method provides a simple way to verify group membership before performing
// other operations like sending group emails or managing subscriptions.
//
// Parameters:
//   - ctx: Context for the request (supports cancellation and timeouts)
//   - groupID: Identifier of the target group (required)
//   - email: Email address to check for membership (required)
//
// Returns:
//   - bool: true if the email exists in the group, false otherwise
//   - error: Validation or operation error
//
// Example:
//
//	exists, err := client.CheckEmailInGroup(ctx, "newsletter-group", "user@example.com")
//	if err != nil {
//		log.Fatal(err)
//	}
//	if exists {
//		fmt.Println("Email is subscribed to the group")
//	} else {
//		fmt.Println("Email is not in the group")
//	}
func (c *GroupClient) CheckEmailInGroup(ctx context.Context, groupID string, email string) (bool, error) {
	if groupID == "" {
		return false, fmt.Errorf("group ID is required")
	}
	if email == "" {
		return false, fmt.Errorf("email address is required")
	}

	req := &pb.CheckEmailInGroupRequest{
		Email:   email,
		GroupId: groupID,
	}

	resp, err := c.client.CheckEmailInGroup(ctx, req)
	if err != nil {
		return false, fmt.Errorf("failed to check email in group: %v", err)
	}

	return resp.Exists, nil
}

package sendlix

import (
	"context"
	"crypto/tls"
	"fmt"
	"strconv"
	"strings"
	"time"

	pb "github.com/sendlix/go-sdk/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// IAuth defines the authentication interface that all authentication
// implementations must satisfy. It provides a contract for generating
// authentication headers required for API requests.
type IAuth interface {
	// GetAuthHeader returns the authentication header key and value
	// that should be included in API requests.
	//
	// Parameters:
	//   - ctx: Context for the authentication request
	//
	// Returns:
	//   - string: Header key (e.g., "authorization")
	//   - string: Header value (e.g., "Bearer token")
	//   - error: Any error encountered during authentication
	GetAuthHeader(ctx context.Context) (string, string, error)
}

// Auth implements the IAuth interface for API key authentication with JWT tokens.
// It handles the exchange of API keys for JWT tokens and manages token caching
// to minimize authentication requests.
//
// The Auth struct automatically handles token refresh when tokens expire,
// providing seamless authentication for long-running applications.
type Auth struct {
	apiKey string        // The original API key in format "secret.keyID"
	keyID  int64         // Parsed key ID from the API key
	secret string        // Parsed secret from the API key
	client pb.AuthClient // gRPC client for authentication service
	token  *tokenCache   // Cached JWT token with expiration
}

// tokenCache holds a JWT token along with its expiration time
// to enable efficient token reuse and automatic refresh.
type tokenCache struct {
	token     string    // The JWT token string
	expiresAt time.Time // When the token expires
}

// NewAuth creates a new Auth instance with the provided API key.
// The API key must be in the format "secret.keyID" where secret is the
// API secret and keyID is the numeric key identifier.
//
// This constructor establishes a gRPC connection to the authentication service
// and validates the API key format. The connection is used for JWT token
// exchanges throughout the lifetime of the Auth instance.
//
// Parameters:
//   - apiKey: API key in format "secret.keyID" (e.g., "abc123.456")
//
// Returns:
//   - *Auth: Configured authentication instance
//   - error: Validation or connection error
//
// Example:
//
//	auth, err := sendlix.NewAuth("your-secret.123456")
//	if err != nil {
//		log.Fatal("Failed to create auth:", err)
//	}
//
// Common errors:
//   - Invalid API key format (missing dot separator)
//   - Empty secret portion
//   - Invalid key ID (non-numeric)
//   - Connection failure to authentication service
func NewAuth(apiKey string) (*Auth, error) {
	parts := strings.Split(apiKey, ".")

	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid API key format. Expected format: 'secret.keyID'")
	}

	secret := parts[0]

	if secret == "" {
		return nil, fmt.Errorf("invalid API key format. Secret cannot be empty")
	}

	keyID, err := strconv.ParseInt(parts[1], 10, 64)

	if err != nil {
		return nil, fmt.Errorf("invalid key ID: %v", err)
	}

	// Create gRPC connection for auth
	config := &tls.Config{}
	creds := credentials.NewTLS(config)

	conn, err := grpc.NewClient("api.sendlix.com:443",
		grpc.WithTransportCredentials(creds),
		grpc.WithUserAgent("sendlix-go-sdk/1.0.0"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %v", err)
	}

	client := pb.NewAuthClient(conn)

	return &Auth{
		apiKey: apiKey,
		keyID:  keyID,
		secret: secret,
		client: client,
	}, nil
}

// GetAuthHeader returns the authorization header for authenticated requests.
// This method implements the IAuth interface and handles JWT token retrieval
// and caching automatically.
//
// The method first checks if a valid cached token exists. If the cached token
// is still valid (not expired), it returns the cached token immediately.
// If no valid cached token exists, it requests a new JWT token from the
// authentication service and caches it for future use.
//
// Parameters:
//   - ctx: Context for the authentication request
//
// Returns:
//   - string: Header key ("authorization")
//   - string: Header value ("Bearer <token>")
//   - error: Any error encountered during token retrieval
//
// Example:
//
//	key, value, err := auth.GetAuthHeader(ctx)
//	if err != nil {
//		log.Fatal("Auth failed:", err)
//	}
//	// key = "authorization", value = "Bearer eyJ..."
//
// The returned token is automatically cached and reused until it expires,
// minimizing the number of authentication requests to the server.
func (a *Auth) GetAuthHeader(ctx context.Context) (string, string, error) {
	// Check if we have a valid cached token
	if a.token != nil && time.Now().Before(a.token.expiresAt) {
		return "authorization", "Bearer " + a.token.token, nil
	}

	// Get new token
	req := &pb.AuthRequest{
		Key: &pb.AuthRequest_ApiKey{
			ApiKey: &pb.ApiKey{
				KeyID:  a.keyID,
				Secret: a.secret,
			},
		},
	}

	resp, err := a.client.GetJwtToken(ctx, req)
	if err != nil {
		return "", "", fmt.Errorf("failed to get JWT token: %v", err)
	}

	// Cache the token
	expiresAt := resp.Expires.AsTime()
	a.token = &tokenCache{
		token:     resp.Token,
		expiresAt: expiresAt,
	}

	return "authorization", "Bearer " + resp.Token, nil
}

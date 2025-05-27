package sendlix

import (
	"context"
	"crypto/tls"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

// BaseClient provides common functionality for all API clients.
// It manages the gRPC connection, authentication, and common client configuration.
// All specific API clients (EmailClient, GroupClient, etc.) embed this type.
type BaseClient struct {
	conn   *grpc.ClientConn
	auth   IAuth
	config *ClientConfig
}

// ClientConfig holds configuration options for API clients.
// It defines connection parameters and client behavior settings.
type ClientConfig struct {
	// ServerAddress is the address of the Sendlix API server.
	// Default: "api.sendlix.com:443"
	ServerAddress string

	// UserAgent is the user agent string sent with requests.
	// Default: "sendlix-go-sdk/1.0.0"
	UserAgent string

	// Insecure determines whether to skip TLS certificate verification.
	// Only use true for testing purposes. Default: false
	Insecure bool
}

// DefaultClientConfig returns the default client configuration with
// sensible defaults for production use.
//
// Returns:
//   - ServerAddress: "api.sendlix.com:443"
//   - UserAgent: "sendlix-go-sdk/1.0.0"
//   - Insecure: false
func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		ServerAddress: "api.sendlix.com:443",
		UserAgent:     "sendlix-go-sdk/1.0.0",
		Insecure:      false,
	}
}

// NewBaseClient creates a new base client with the provided authentication and configuration.
// This function establishes a secure gRPC connection to the Sendlix API server and sets up
// automatic authentication for all requests. It is typically not called directly; instead,
// use the specific client constructors like NewEmailClient or NewGroupClient.
//
// The function performs several important setup steps:
//   - Validates that authentication is provided
//   - Applies default configuration if none is provided
//   - Establishes secure TLS connection (unless configured otherwise)
//   - Sets up automatic authentication interceptor
//
// Parameters:
//   - auth: Authentication implementation (required, cannot be nil)
//   - config: Client configuration (optional, uses defaults if nil)
//
// Returns:
//   - *BaseClient: Configured base client ready for use
//   - error: Validation or connection error
//
// Example:
//
//	auth, err := sendlix.NewAuth("secret.keyid")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	config := sendlix.DefaultClientConfig()
//	config.UserAgent = "MyApp/1.0.0"
//
//	client, err := sendlix.NewBaseClient(auth, config)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer client.Close()
//
// Common errors:
//   - nil authentication (auth parameter is required)
//   - network connectivity issues
//   - invalid server address in configuration
//   - TLS handshake failures
func NewBaseClient(auth IAuth, config *ClientConfig) (*BaseClient, error) {

	if auth == nil {
		return nil, fmt.Errorf("authentication is required")
	}

	if config == nil {
		config = DefaultClientConfig()
	}

	var creds credentials.TransportCredentials
	if config.Insecure {
		creds = credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
	} else {
		creds = credentials.NewTLS(&tls.Config{})
	}

	conn, err := grpc.NewClient(config.ServerAddress,
		grpc.WithTransportCredentials(creds),
		grpc.WithUserAgent(config.UserAgent),
		grpc.WithUnaryInterceptor(authInterceptor(auth)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %v", err)
	}

	return &BaseClient{
		conn:   conn,
		auth:   auth,
		config: config,
	}, nil
}

// Close closes the gRPC connection and releases associated resources.
// This method should be called when the client is no longer needed to prevent
// resource leaks. It's safe to call Close multiple times.
//
// Returns:
//   - error: Any error encountered while closing the connection
//
// Example:
//
//	client, err := sendlix.NewEmailClient(auth, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer client.Close() // Ensure cleanup
func (c *BaseClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GetConnection returns the underlying gRPC connection.
// This method is primarily used internally by specific API clients
// to create their respective gRPC service clients.
//
// Returns:
//   - *grpc.ClientConn: The underlying gRPC connection
func (c *BaseClient) GetConnection() *grpc.ClientConn {
	return c.conn
}

// authInterceptor creates a gRPC unary interceptor that automatically adds
// authentication headers to all outgoing requests. This interceptor retrieves
// the authentication header from the provided IAuth implementation and adds
// it to the request metadata.
//
// Parameters:
//   - auth: Authentication implementation to use for header generation
//
// Returns:
//   - grpc.UnaryClientInterceptor: Configured authentication interceptor
func authInterceptor(auth IAuth) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// Get auth header
		key, value, err := auth.GetAuthHeader(ctx)
		if err != nil {
			return fmt.Errorf("failed to get auth header: %v", err)
		}

		// Add auth header to context
		ctx = metadata.AppendToOutgoingContext(ctx, key, value)

		// Call the method
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

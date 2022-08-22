package connect

import (
	"github.com/ktr0731/grpc-web-go-client/grpcweb"
	"google.golang.org/grpc"
)

const (
	tlsError = "TLS cannot be disabled in production, only for testing suites!"
)

// Connection is an interface designed to sit between hosts and connections
// to allow use of grpcweb clients.
type Connection interface {
	// GetWebConn returns the grpcweb ClientConn for use in browsers.
	// It panics if called on a grpc client.
	GetWebConn() *grpcweb.ClientConn
	// GetGrpcConn returns the grpc ClientConn for standard use.
	// It panics if called on a grpcweb client.
	GetGrpcConn() *grpc.ClientConn
	// Connect initiates a connection with the host using connection logic
	// supplied by the underlying class.
	Connect() error
	// IsWeb returns true if the connection uses grpcweb
	IsWeb() bool

	// Close closes the underlying connection
	Close() error

	clientConnHelpers
}

// clientConnHelpers holds private helper methods exposed on the connection object
type clientConnHelpers interface {
	isAlive() bool
	disconnect()
}

// newConnection initializes a webConn and returns it wrapped as a Connection
func newConnection(isWeb bool, host *Host) Connection {
	if isWeb {
		return &webConn{
			h: host,
		}
	} else {
		return &grpcConn{
			h: host,
		}
	}
}

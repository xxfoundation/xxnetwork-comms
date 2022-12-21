package connect

import (
	"crypto/x509"
	"git.xx.network/elixxir/grpc-web-go-client/grpcweb"
	jww "github.com/spf13/jwalterweatherman"
	"google.golang.org/grpc"
	"time"
)

const (
	tlsError = "TLS cannot be disabled in production, only for testing suites!"
)

// ConnectionType is intended to act as an enum for different methods of host connection
type ConnectionType uint8

// Enumerate the extant connection methods
const (
	Grpc ConnectionType = iota
	Web
)

// Stringify connection constants
func (ct ConnectionType) String() string {
	switch ct {
	case Grpc:
		return "grpc"
	case Web:
		return "web"
	default:
		return "unknown"
	}
}

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

	IsOnline() (time.Duration, bool)

	GetRemoteCertificate() (*x509.Certificate, error)
	clientConnHelpers
}

// clientConnHelpers holds private helper methods exposed on the connection object
type clientConnHelpers interface {
	isAlive() bool
	disconnect()
}

// newConnection initializes a webConn or grpcConn and returns it wrapped as a Connection
func newConnection(t ConnectionType, host *Host) Connection {
	switch t {
	case Web:
		return &webConn{h: host}
	case Grpc:
		return &grpcConn{h: host}
	default:
		jww.ERROR.Printf("Cannot make connection of type %s", t)
		return nil
	}
}

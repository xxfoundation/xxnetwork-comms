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
	GetWebConn() *grpcweb.ClientConn
	GetGrpcConn() *grpc.ClientConn
	Connect() error
	IsWeb() bool

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

package connect

import (
	"errors"
	"fmt"
	"git.xx.network/elixxir/grpc-web-go-client/grpcweb"
	jww "github.com/spf13/jwalterweatherman"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"math"
	"net"
	"sync/atomic"
	"time"
)

// grpcConn implements the Connection interface
type grpcConn struct {
	h          *Host
	connection *grpc.ClientConn
}

// GetWebConn returns the grpcweb ClientConn object
func (gc *grpcConn) GetWebConn() *grpcweb.ClientConn {
	jww.FATAL.Panic("Cannot GetWebConn on a host that is configured for grpc connections")
	return nil
}

// GetGrpcConn returns the grpc ClientConn object
func (gc *grpcConn) GetGrpcConn() *grpc.ClientConn {
	return gc.connection
}

// Connect initializes the appropriate connection using helper functions.
func (gc *grpcConn) Connect() error {
	return gc.connectGrpcHelper()
}

// IsWeb returns true if the connection is configured for web connections
func (gc *grpcConn) IsWeb() bool {
	return false
}

// connectGrpcHelper creates a connection while not under a write lock.
// undefined behavior if the caller has not taken the write lock
func (gc *grpcConn) connectGrpcHelper() (err error) {
	// Configure TLS options
	var securityDial grpc.DialOption
	if gc.h.credentials != nil {
		// Create the gRPC client with TLS
		securityDial = grpc.WithTransportCredentials(gc.h.credentials)
	} else if TestingOnlyDisableTLS {
		// Create the gRPC client without TLS
		jww.WARN.Printf("Connecting to %v without TLS!", gc.h.GetAddress())
		securityDial = grpc.WithInsecure()
	} else {
		jww.FATAL.Panicf(tlsError)
	}

	jww.DEBUG.Printf("Attempting to establish connection to %s using"+
		" credentials: %+v", gc.h.GetAddress(), securityDial)

	// Attempt to establish a new connection
	var numRetries uint32
	//todo-remove this retry block when grpc is updated
	for numRetries = 0; numRetries < gc.h.params.MaxRetries && !gc.isAlive(); numRetries++ {
		gc.h.disconnect()

		jww.DEBUG.Printf("Connecting to %+v Attempt number %+v of %+v",
			gc.h.GetAddress(), numRetries, gc.h.params.MaxRetries)

		// If timeout is enabled, the max wait time becomes
		// ~14 seconds (with maxRetries=100)
		backoffTime := 2000 * (numRetries/16 + 1)
		if backoffTime > 15000 {
			backoffTime = 15000
		}
		ctx, cancel := newContext(time.Duration(backoffTime) * time.Millisecond)

		dialOpts := []grpc.DialOption{
			grpc.WithBlock(),
			grpc.WithKeepaliveParams(gc.h.params.KaClientOpts),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32)),
			securityDial,
		}

		windowSize := atomic.LoadInt32(gc.h.windowSize)
		if windowSize != 0 {
			dialOpts = append(dialOpts, grpc.WithInitialWindowSize(windowSize))
			dialOpts = append(dialOpts, grpc.WithInitialConnWindowSize(windowSize))
		}

		// Create the connection
		gc.connection, err = grpc.DialContext(ctx, gc.h.GetAddress(),
			dialOpts...)

		if err != nil {
			jww.DEBUG.Printf("Attempt number %+v to connect to %s failed\n",
				numRetries, gc.h.GetAddress())
		}
		cancel()
	}

	// Verify that the connection was established successfully
	if !gc.isAlive() {
		gc.h.disconnect()
		return errors.New(fmt.Sprintf(
			"Last try to connect to %s failed. Giving up",
			gc.h.GetAddress()))
	}

	// Add the successful connection to the Manager
	jww.INFO.Printf("Successfully connected to %v", gc.h.GetAddress())
	return
}

// Close calls the internal Close function on the connection
func (gc *grpcConn) Close() error {
	if gc.connection == nil {
		return nil
	}
	return gc.connection.Close()
}

// disconnect closes the grpcConn connection while not under a write lock.
// undefined behavior if the caller has not taken the write lock
func (gc *grpcConn) disconnect() {
	// it's possible to close a host which never sent so that it never made a
	// connection. In that case, we should not close a connection which does not
	// exist
	if gc.connection != nil {
		jww.INFO.Printf("Disconnected from %s at %s", gc.h.GetId(), gc.h.GetAddress())
		err := gc.connection.Close()
		if err != nil {
			jww.ERROR.Printf("Unable to close connection to %s: %+v",
				gc.h.GetAddress(), errors.New(err.Error()))
		} else {
			gc.connection = nil
		}
	}
}

// isAlive returns true if the grpcConn is non-nil and alive
// must already be under the connectionMux
func (gc *grpcConn) isAlive() bool {
	if gc.connection == nil {
		return false
	}
	state := gc.connection.GetState()
	return state == connectivity.Idle || state == connectivity.Connecting ||
		state == connectivity.Ready
}

func (gc *grpcConn) IsOnline() (time.Duration, bool) {
	addr := gc.h.GetAddress()
	start := time.Now()
	conn, err := net.DialTimeout("tcp", addr, gc.h.params.PingTimeout)
	if err != nil {
		// If we cannot connect, mark the connection as failed
		jww.DEBUG.Printf(
			"Failed to verify connectivity for address %s: %+v", addr, err)
		return 0, false
	}
	// Attempt to close the connection
	if conn != nil {
		errClose := conn.Close()
		if errClose != nil {
			jww.DEBUG.Printf(
				"Failed to close connection for address %s: %+v", addr, errClose)
		}
	}
	return time.Since(start), true
}

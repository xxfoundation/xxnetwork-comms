package connect

import (
	"errors"
	"fmt"
	"github.com/ktr0731/grpc-web-go-client/grpcweb"
	jww "github.com/spf13/jwalterweatherman"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"math"
	"sync/atomic"
	"time"
)

// webConn implements the Connection interface
type grpcConn struct {
	h        *Host
	webConn  *grpcweb.ClientConn
	grpcConn *grpc.ClientConn
}

// GetWebConn returns the grpcweb ClientConn object
func (gc *grpcConn) GetWebConn() *grpcweb.ClientConn {
	jww.FATAL.Panic("Cannot GetWebConn on a host that is configured for grpc connections")
	return nil
}

// GetGrpcConn returns the grpc ClientConn object
func (gc *grpcConn) GetGrpcConn() *grpc.ClientConn {
	return gc.grpcConn
}

// Connect initializes the appropriate connection using helper functions.
func (gc *grpcConn) Connect() error {
	return gc.connectGrpcHelper()
}

// IsWeb returns true if the webConn is configured for web connections
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
		gc.grpcConn, err = grpc.DialContext(ctx, gc.h.GetAddress(),
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

// Close calls the internal Close function on the grpcConn
func (gc *grpcConn) Close() error {
	if gc.grpcConn == nil {
		return nil
	}
	return gc.grpcConn.Close()
}

// disconnect closes the grpcConn connection while not under a write lock.
// undefined behavior if the caller has not taken the write lock
func (gc *grpcConn) disconnect() {
	// it's possible to close a host which never sent so that it never made a
	// connection. In that case, we should not close a connection which does not
	// exist
	if gc.grpcConn != nil {
		jww.INFO.Printf("Disconnected from %s at %s", gc.h.GetId(), gc.h.GetAddress())
		err := gc.grpcConn.Close()
		if err != nil {
			jww.ERROR.Printf("Unable to close connection to %s: %+v",
				gc.h.GetAddress(), errors.New(err.Error()))
		} else {
			gc.grpcConn = nil
		}
	}
}

// isAlive returns true if the grpcConn is non-nil and alive
// must already be under the connectionMux
func (gc *grpcConn) isAlive() bool {
	if gc.grpcConn == nil {
		return false
	}
	state := gc.grpcConn.GetState()
	return state == connectivity.Idle || state == connectivity.Connecting ||
		state == connectivity.Ready
}

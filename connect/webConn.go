package connect

import (
	"fmt"
	"git.xx.network/elixxir/grpc-web-go-client/grpcweb"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"google.golang.org/grpc"
	"time"
)

type WebConnParam struct {
	TlsHandshakeTimeout   time.Duration
	IdleConnTimeout       time.Duration
	ExpectContinueTimeout time.Duration
}

// webConn implements the Connection interface
type webConn struct {
	h          *Host
	connection *grpcweb.ClientConn
}

// GetWebConn returns the grpcweb ClientConn object
func (wc *webConn) GetWebConn() *grpcweb.ClientConn {
	return wc.connection
}

// GetGrpcConn returns the grpc ClientConn object
func (wc *webConn) GetGrpcConn() *grpc.ClientConn {
	jww.FATAL.Panicf("Cannot GetGrpcConn on a host that is configured for web connections")
	return nil
}

// Connect initializes the appropriate connection using helper functions.
func (wc *webConn) Connect() error {
	return wc.connectWebHelper()
}

// IsWeb returns true if the connection is configured for web connections
func (wc *webConn) IsWeb() bool {
	return true
}

// connectWebHelper initializes the grpcweb ClientConn object
// Note that until the downstream repo is fixed, this doesn't actually
// establish a connection past creating the http object.
func (wc *webConn) connectWebHelper() (err error) {
	// Configure TLS options
	var securityDial grpcweb.DialOption
	if wc.h.credentials != nil {
		securityDial = grpcweb.WithTlsCertificate(wc.h.certificate)
	} else if TestingOnlyDisableTLS {
		jww.WARN.Printf("Connecting to %v without TLS!", wc.h.GetAddress())
		securityDial = grpcweb.WithInsecure()
	} else {
		jww.FATAL.Panicf(tlsError)
	}

	jww.DEBUG.Printf("Attempting to establish connection to %s using"+
		" credentials: %+v", wc.h.GetAddress(), securityDial)

	// Attempt to establish a new connection
	var numRetries uint32
	for numRetries = 0; numRetries < wc.h.params.MaxRetries && !wc.isAlive(); numRetries++ {
		wc.h.disconnect()

		jww.DEBUG.Printf("Connecting to %+v Attempt number %+v of %+v",
			wc.h.GetAddress(), numRetries, wc.h.params.MaxRetries)

		// If timeout is enabled, the max wait time becomes
		// ~14 seconds (with maxRetries=100)
		backoffTime := 2000 * (numRetries/16 + 1)
		if backoffTime > 15000 {
			backoffTime = 15000
		}
		//ctx, cancel := newContext(time.Duration(backoffTime) * time.Millisecond)

		dialOpts := []grpcweb.DialOption{
			grpcweb.WithIdleConnTimeout(wc.h.params.WebParams.IdleConnTimeout),
			grpcweb.WithExpectContinueTimeout(wc.h.params.WebParams.ExpectContinueTimeout),
			grpcweb.WithTlsHandshakeTimeout(wc.h.params.WebParams.TlsHandshakeTimeout),
			grpcweb.WithInsecureTlsVerification(),
			grpcweb.WithDefaultCallOptions(), // grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32)),
			securityDial,
		}

		//windowSize := atomic.LoadInt32(wc.h.windowSize)
		//if windowSize != 0 {
		//	dialOpts = append(dialOpts, grpc.WithInitialWindowSize(windowSize))
		//	dialOpts = append(dialOpts, grpc.WithInitialConnWindowSize(windowSize))
		//}

		// Create the connection
		wc.connection, err = grpcweb.DialContext(wc.h.GetAddress(),
			dialOpts...)

		if err != nil {
			jww.DEBUG.Printf("Attempt number %+v to connect to %s failed\n",
				numRetries, wc.h.GetAddress())
		}
		//cancel()
	}

	// Verify that the connection was established successfully
	if !wc.isAlive() {
		wc.h.disconnect()
		return errors.New(fmt.Sprintf(
			"Last try to connect to %s failed. Giving up",
			wc.h.GetAddress()))
	}

	// Add the successful connection to the Manager
	jww.INFO.Printf("Successfully connected to %v", wc.h.GetAddress())
	return
}

// Close handles closing the http connection.
func (wc *webConn) Close() error {
	if wc.connection == nil {
		return nil
	}
	return wc.connection.Close()

}

// disconnect closes the webConn connection while not under a write lock.
// undefined behavior if the caller has not taken the write lock
func (wc *webConn) disconnect() {
	// it's possible to close a host which never sent so that it never made a
	// connection. In that case, we should not close a connection which does not
	// exist
	if wc.connection != nil {
		if err := wc.connection.Close(); err != nil {
			jww.FATAL.Panicf("Failed to disconnect web client: %+v", err)
		}
		wc.connection = nil
	}

}

// isAlive returns true if the webConn is non-nil and alive
// must already be under the connectionMux
func (wc *webConn) isAlive() bool {
	if wc.connection == nil {
		return false
	}
	return wc.connection.IsAlive()
}

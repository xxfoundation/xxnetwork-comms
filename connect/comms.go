////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Handles the basic top-level comms object used across all packages

package connect

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/pkg/errors"
	"github.com/soheilhy/cmux"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/xx_network/comms/connect/token"
	"gitlab.com/xx_network/crypto/signature/rsa"
	"gitlab.com/xx_network/primitives/id"
	"golang.org/x/crypto/cryptobyte"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"io"
	"math"
	"net"
	"net/http"
	"src.agwa.name/tlshacks"
	"strings"
	"time"
)

// MaxWindowSize 4 MB
const MaxWindowSize = math.MaxInt32
const tlsHandshakePrefixLen = 5

// TestingOnlyDisableTLS is the variable set for testing
// which allows for the disabled TLS code-path. Production
// code-path will only function with TLS enabled.
var TestingOnlyDisableTLS = false
var TestingOnlyInsecureTLSVerify = false

// KaOpts are Keepalive options for servers
// TODO: Set these via config
var KaOpts = keepalive.ServerParameters{
	// Idle for at most 60s
	MaxConnectionIdle: 60 * time.Second,
	// Reset after an hour
	MaxConnectionAge: 1 * time.Hour,
	// w/ 1m grace shutdown
	MaxConnectionAgeGrace: 1 * time.Minute,
	// Send keepAlive every Time interval
	Time: 5 * time.Second,
	// Timeout after last successful keepAlive to close connection
	Timeout: 60 * time.Second,
}

// KaEnforcement are keepalive enforcement options for servers
var KaEnforcement = keepalive.EnforcementPolicy{
	// Send keepAlive every Time interval
	MinTime: 3 * time.Second,
	// Doing KA on non-streams is OK
	PermitWithoutStream: true,
}

// MaxConcurrentStreams is the number of server-side streams to allow open
var MaxConcurrentStreams = uint32(250000)

// ProtoComms is a proto object containing a gRPC server logic.
type ProtoComms struct {
	// Inherit the Manager object
	*Manager

	// The network ID of this comms server
	networkId *id.ID

	// Private key of the local comms instance
	privateKey *rsa.PrivateKey

	// Disables the checking of authentication signatures for testing setups
	disableAuth bool

	listeningAddress string

	// SERVER-ONLY FIELDS ------------------------------------------------------

	// A map of reverse-authentication tokens
	tokens *token.Map

	// Low-level net.Listener object that listens at listeningAddress
	netListener net.Listener

	// Local network server
	grpcServer  *grpc.Server
	httpServer  *http.Server
	httpsServer *http.Server
	// GRPC credentials stored to re-initialize after restart
	grpcCreds tls.Certificate
	// Parsed grpc x509 certificate for checking incoming tls request servernames
	grpcX509 *x509.Certificate

	// CLIENT-ONLY FIELDS ------------------------------------------------------

	// Used to store the public key used for generating Client Id
	pubKeyPem []byte

	// Used to store the salt used for generating Client Id
	salt []byte

	// cmux interface for starting http listener when serving with web
	mux cmux.CMux

	// https certificate infra for replacing tls cert with no downtime
	httpsCertificate *tls.Certificate
	// Parsed https x509 certificate for checking incoming tls request servernames
	httpsX509 *x509.Certificate

	// -------------------------------------------------------------------------
}

// GetId returns a copy of the ProtoComms networkId
func (c *ProtoComms) GetId() *id.ID {
	return c.networkId.DeepCopy()
}

// GetServer returns the ProtoComms grpc.Server object
func (c *ProtoComms) GetServer() *grpc.Server {
	return c.grpcServer
}

// CreateCommClient creates a ProtoComms client-type object to be
// used in various initializers.
func CreateCommClient(id *id.ID, pubKeyPem, privKeyPem,
	salt []byte) (*ProtoComms, error) {
	// Build the ProtoComms object
	pc := &ProtoComms{
		networkId: id,
		pubKeyPem: pubKeyPem,
		salt:      salt,
		tokens:    token.NewMap(),
		Manager:   newManager(),
	}

	// Set the private key if specified
	if privKeyPem != nil {
		err := pc.setPrivateKey(privKeyPem)
		if err != nil {
			return nil, errors.Errorf("Could not set private key: %+v", err)
		}
	}
	return pc, nil
}

// StartCommServer creates a ProtoComms server-type object to be used in various initializers.
// Opens a net.Listener the local address specified by listeningAddr.
func StartCommServer(id *id.ID, listeningAddr string,
	certPEMblock, keyPEMblock []byte, preloadedHosts []*Host) (*ProtoComms, error) {

listen:
	// Listen on the given address
	lis, err := net.Listen("tcp", listeningAddr)
	if err != nil {
		if strings.Contains(err.Error(), "bind: address already in use") {
			jww.WARN.Printf("Could not listen on %s, is port in use? waiting 30s: %s", listeningAddr, err.Error())
			time.Sleep(30 * time.Second)
			goto listen
		}
		return nil, errors.New(err.Error())
	}

	// Build the comms object
	pc := &ProtoComms{
		networkId:        id,
		netListener:      lis,
		tokens:           token.NewMap(),
		Manager:          newManager(),
		listeningAddress: listeningAddr,
	}

	for _, h := range preloadedHosts {
		pc.Manager.addHost(h)
	}

	// If TLS was specified
	if certPEMblock != nil && keyPEMblock != nil {

		// Create the TLS certificate
		x509cert, err := tls.X509KeyPair(certPEMblock, keyPEMblock)
		if err != nil {
			return nil, errors.Errorf("Could not load TLS keys: %+v", err)
		}

		// Set the private key
		err = pc.setPrivateKey(keyPEMblock)
		if err != nil {
			return nil, errors.Errorf("Could not set private key: %+v", err)
		}

		// Set the public key data
		pc.pubKeyPem = certPEMblock

		// Create the gRPC server with TLS
		jww.INFO.Printf("Starting server with TLS...")
		if x509cert.Leaf == nil {
			x509cert.Leaf, err = x509.ParseCertificate(x509cert.Certificate[0])
			if err != nil {
				return nil, errors.WithMessage(err, "Could not parse x509 certificate")
			}
		}
		pc.grpcX509 = x509cert.Leaf
		creds := credentials.NewServerTLSFromCert(&x509cert)
		pc.grpcCreds = x509cert
		pc.grpcServer = grpc.NewServer(grpc.Creds(creds),
			grpc.MaxConcurrentStreams(MaxConcurrentStreams),
			grpc.MaxRecvMsgSize(math.MaxInt32),
			grpc.KeepaliveParams(KaOpts),
			grpc.KeepaliveEnforcementPolicy(KaEnforcement))
	} else if TestingOnlyDisableTLS {
		// Create the gRPC server without TLS
		jww.WARN.Printf("Starting server with TLS disabled...")
		pc.grpcServer = grpc.NewServer(
			grpc.MaxConcurrentStreams(MaxConcurrentStreams),
			grpc.MaxRecvMsgSize(math.MaxInt32),
			grpc.KeepaliveParams(KaOpts),
			grpc.KeepaliveEnforcementPolicy(KaEnforcement))
	} else {
		jww.FATAL.Panicf("TLS cannot be disabled in production, only for testing suites!")
	}

	return pc, nil
}

// Restart is a public accessor meant to allow for reuse of a host after
// Shutdown is called.  The intended use is for replacing certificates.
func (c *ProtoComms) Restart() error {
	if TestingOnlyDisableTLS {
		c.grpcServer = grpc.NewServer(
			grpc.MaxConcurrentStreams(MaxConcurrentStreams),
			grpc.MaxRecvMsgSize(math.MaxInt32),
			grpc.KeepaliveParams(KaOpts),
			grpc.KeepaliveEnforcementPolicy(KaEnforcement))
	} else {
		creds := credentials.NewServerTLSFromCert(&c.grpcCreds)
		if c.grpcCreds.Leaf == nil {
			var err error
			c.grpcCreds.Leaf, err = x509.ParseCertificate(c.grpcCreds.Certificate[0])
			if err != nil {
				return errors.WithMessage(err, "Could not parse x509 certificate")
			}
		}
		c.grpcX509 = c.grpcCreds.Leaf
		c.grpcServer = grpc.NewServer(grpc.Creds(creds),
			grpc.MaxConcurrentStreams(MaxConcurrentStreams),
			grpc.MaxRecvMsgSize(math.MaxInt32),
			grpc.KeepaliveParams(KaOpts),
			grpc.KeepaliveEnforcementPolicy(KaEnforcement))
	}

	if c.netListener != nil {
		return errors.New("ProtoComms is already listening")
	}
listen:
	// Listen on the given address
	lis, err := net.Listen("tcp", c.listeningAddress)
	if err != nil {
		if strings.Contains(err.Error(), "bind: address already in use") {
			jww.WARN.Printf("Could not listen on %s, is port in use? waiting 30s: %s", c.listeningAddress, err.Error())
			time.Sleep(30 * time.Second)
			goto listen
		}
		return errors.New(err.Error())
	}

	c.netListener = lis
	return nil
}

// Serve is a non-blocking call that begins serving content
// for GRPC. GRPC endpoints must be registered before making this call.
func (c *ProtoComms) Serve() {
	listenGRPC := func(l net.Listener) {
		// This blocks for the lifetime of the listener.
		if err := c.GetServer().Serve(l); err != nil {
			jww.FATAL.Panicf("Failed to serve GRPC: %+v", err)
		}
		jww.INFO.Printf("Shutting down GRPC server listener")
	}
	go listenGRPC(c.netListener)
}

// ServeWithWeb is a non-blocking call that begins serving content
// for grpcWeb (over HTTP) and GRPC on the same port.
// GRPC endpoints must be registered before making this call.
func (c *ProtoComms) ServeWithWeb() {
	grpcServer := c.GetServer()

	// Split netListener into two distinct listeners for GRPC and HTTP
	mux := cmux.New(c.netListener)
	grpcMatcher := c.matchGrpcTls

	if TestingOnlyDisableTLS {
		grpcMatcher = cmux.HTTP2()
	}
	httpL := mux.Match(cmux.HTTP1())
	grpcL := mux.Match(grpcMatcher)
	c.mux = mux

	listenHTTP := func(l net.Listener) {
		httpServer := grpcweb.WrapServer(grpcServer,
			grpcweb.WithOriginFunc(func(origin string) bool { return true }))
		jww.WARN.Printf("Starting HTTP server!")

		c.httpServer = &http.Server{
			Handler: httpServer,
		}

		if err := c.httpServer.Serve(l); err != nil {
			// Cannot panic here due to shared net.Listener
			jww.ERROR.Printf("Failed to serve HTTP: %+v", err)
		}
	}
	listenGRPC := func(l net.Listener) {
		// This blocks for the lifetime of the listener.
		if err := grpcServer.Serve(l); err != nil {
			// Cannot panic here due to shared net.Listener
			jww.ERROR.Printf("Failed to serve GRPC: %+v", err)
		}
		jww.INFO.Printf("Shutting down GRPC server listener")
	}
	listenPort := func() {
		if err := mux.Serve(); err != nil {
			// Cannot panic here due to shared net.Listener
			jww.ERROR.Printf("Failed to serve port: %+v", err)
		}
		jww.INFO.Printf("Shutting down port server listener")
	}
	go listenHTTP(httpL)
	go listenGRPC(grpcL)
	go listenPort()
}

// Matcher function for grpc TLS connections; parses the client hello message,
// return true if the servername contains xx.network or is recognized as a
// valid IP addressinfo or if extensions contains ALPN and if the protos
// contain 'h2'.
func (c *ProtoComms) matchGrpcTls(r io.Reader) bool {
	hello, ok := parseTlsPacket(r)
	if !ok {
		return false
	}

	if hello.Info.ServerName != nil {
		err := c.grpcX509.VerifyHostname(*hello.Info.ServerName)
		if err == nil {
			return true
		} else {
			jww.TRACE.Printf("VerifyHostname error: %+v", err)
		}
		grpcCertificateDefaultPrefix := "xx.network"
		if strings.Contains(*hello.Info.ServerName, grpcCertificateDefaultPrefix) {
			return true
		}
		if net.ParseIP(*hello.Info.ServerName) != nil {
			return true
		}
	}

	if len(hello.Info.Protocols) > 0 {
		if len(hello.Info.Protocols) == 1 && hello.Info.Protocols[0] == "h2" {
			return true
		}
	}

	return false
}

// Matcher function for grpcweb TLS connections; parses the client hello
// message, return true if the servername contains 'gwid.', if
// protocols contain http1, or as a default for other unmatched tls connections.
func (c *ProtoComms) matchWebTls(r io.Reader) bool {
	hello, ok := parseTlsPacket(r)
	if !ok {
		return false
	}

	if hello.Info.ServerName != nil {
		err := c.httpsX509.VerifyHostname(*hello.Info.ServerName)
		if err == nil {
			return true
		} else {
			jww.TRACE.Printf("VerifyHostname error: %+v", err)
		}
		snPrefix := fmt.Sprintf("%s.", base64.URLEncoding.EncodeToString(c.GetId().Marshal()))
		if strings.Contains(*hello.Info.ServerName, snPrefix) {
			return true
		}
	}

	if len(hello.Info.Protocols) > 0 {
		if len(hello.Info.Protocols) == 1 && hello.Info.Protocols[0] == "h2" {
			return false
		}
		for _, item := range hello.Info.Protocols {
			if item == "http1" {
				return true
			}
		}
	}

	return true
}

func parseTlsPacket(r io.Reader) (*tlshacks.ClientHelloInfo, bool) {
	var handshakePrefix cryptobyte.String
	handshakePrefix = make([]byte, tlsHandshakePrefixLen)
	n, err := io.ReadFull(r, handshakePrefix)
	if err != nil || n != tlsHandshakePrefixLen {
		return nil, false
	}

	jww.TRACE.Printf("Read potential prefix bytes %+v", handshakePrefix)

	var handshakeMessageType, messageVersionMinor, messageVersionMajor uint8
	var handshakeMessageLength uint16
	if !handshakePrefix.ReadUint8(&handshakeMessageType) {
		return nil, false
	}
	if !handshakePrefix.ReadUint8(&messageVersionMajor) {
		return nil, false
	}
	if !handshakePrefix.ReadUint8(&messageVersionMinor) {
		return nil, false
	}
	if !handshakePrefix.ReadUint16(&handshakeMessageLength) {
		return nil, false
	}

	jww.DEBUG.Printf("Read handshake message of type %d (%d.%d), %d bytes left to read", handshakeMessageType, messageVersionMajor, messageVersionMinor, handshakeMessageLength)

	helloMessage := make([]byte, handshakeMessageLength)
	n, err = io.ReadFull(r, helloMessage)
	if err != nil {
		return nil, false
	}

	if helloMessage[0] != 0x01 {
		return nil, false
	}

	hello := tlshacks.UnmarshalClientHello(helloMessage)
	if hello == nil {
		return nil, false
	}
	return hello, true
}

// ProvisionHttps provides a tls cert and key to the thread which serves the
// grpcweb endpoints, allowing it to serve with https.  Note that https will
// not be usable until this has been called at least once, unblocking the
// listenHTTP func in ServeWithWeb.  Future calls will be handled by the
// startUpdateCertificate thread.
func (c *ProtoComms) ServeHttps(keyPair tls.Certificate) error {
	if c.mux == nil {
		return errors.New("mux does not exist; is https enabled?")
	}

	if c.netListener == nil {
		return errors.New("ProtoComms is closed, call Restart to initialize")
	}

	httpL := c.mux.Match(c.matchWebTls)

	grpcServer := c.grpcServer
	var parsedLeafCert *x509.Certificate
	var err error
	if keyPair.Leaf == nil {
		parsedLeafCert, err = x509.ParseCertificate(keyPair.Certificate[0])
		if err != nil {
			jww.FATAL.Panicf("Failed to load TLS certificate: %+v", err)
		}
	} else {
		parsedLeafCert = keyPair.Leaf
	}

	c.httpsX509 = parsedLeafCert

	listenHTTPS := func(l net.Listener) {
		jww.INFO.Printf("Starting HTTP listener on GRPC endpoints: %+v",
			grpcweb.ListGRPCResources(grpcServer))
		httpsServer := grpcweb.WrapServer(grpcServer,
			grpcweb.WithOriginFunc(func(origin string) bool { return true }))

		// Configure TLS for this listener, using the config from
		// http.ServeTLS
		tlsConf := &tls.Config{}
		tlsConf.NextProtos = append(tlsConf.NextProtos, "h2", "http/1.1")

		// We use the GetCertificate field and a function which returns
		// c.httpsCertificate wrapped by a ReadLock to allow for future
		// changes to the certificate without downtime
		c.httpsCertificate = &keyPair
		tlsConf.Certificates = []tls.Certificate{keyPair}

		var serverName string
		serverName = parsedLeafCert.DNSNames[0]
		tlsConf.ServerName = serverName

		tlsLis := tls.NewListener(l, tlsConf)
		jww.WARN.Printf("Starting HTTPS server!")

		c.httpsServer = &http.Server{
			Handler: httpsServer,
		}

		if err := c.httpsServer.Serve(tlsLis); err != nil {
			// Cannot panic here due to shared net.Listener
			jww.WARN.Printf("HTTPS listener shutting down: %+v", err)
		}
		jww.INFO.Printf("Stopped HTTPS server listener")
	}

	go listenHTTPS(httpL)

	return nil
}

// Shutdown performs a graceful shutdown of the local server.
func (c *ProtoComms) Shutdown() {
	// Also handles closing of net.Listener
	if c.httpsServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err := c.httpsServer.Shutdown(ctx)
		if err != nil {
			jww.WARN.Printf("Failed to shutdown http server: %+v", err)
		}
		cancel()
	} else if c.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err := c.httpServer.Shutdown(ctx)
		if err != nil {
			jww.WARN.Printf("Failed to shutdown http server: %+v", err)
		}
		cancel()
	} else {
		c.grpcServer.GracefulStop()
	}

	// Close all Manager connections
	c.DisconnectAll()
	c.grpcServer = nil
	c.netListener = nil
	c.mux = nil
	jww.INFO.Printf("Comms server successfully shut down")
}

// Stringer method
func (c *ProtoComms) String() string {
	return c.listeningAddress
}

// Setter for local server's private key
func (c *ProtoComms) setPrivateKey(data []byte) error {
	key, err := rsa.LoadPrivateKeyFromPem(data)
	if err != nil {
		return errors.Errorf("Failed to form private key file from data at %s: %+v", data, err)
	}

	c.privateKey = key
	return nil
}

// GetPrivateKey is the getter for local server's private key.
func (c *ProtoComms) GetPrivateKey() *rsa.PrivateKey {
	return c.privateKey
}

// Send sets up or recovers the Host's connection,
// then runs the given transmit function.
func (c *ProtoComms) Send(host *Host, f func(conn Connection) (*any.Any,
	error)) (result *any.Any, err error) {

	jww.TRACE.Printf("Attempting to send to host: %s", host)
	fSh := func(conn Connection) (interface{}, error) {
		return f(conn)
	}

	anyFace, err := c.transmit(host, fSh)
	if err != nil {
		return nil, err
	}

	return anyFace.(*any.Any), err
}

// Stream sets up or recovers the Host's connection,
// then runs the given Stream function.
func (c *ProtoComms) Stream(host *Host, f func(conn Connection) (
	interface{}, error)) (client interface{}, err error) {

	// Ensure the connection is running
	jww.TRACE.Printf("Attempting to stream to host: %s", host)
	return c.transmit(host, f)
}

// returns true if the connection error is one of the connection errors which
// should be retried
func isConnError(err error) bool {
	return strings.Contains(err.Error(), "context deadline exceeded") ||
		strings.Contains(err.Error(), "connection refused") ||
		strings.Contains(err.Error(), "host disconnected")
}

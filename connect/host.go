////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains functionality for connecting to gateways and servers

package connect

import (
	"fmt"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/crypto/signature/rsa"
	tlsCreds "gitlab.com/elixxir/crypto/tls"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"time"
)

// Information used to describe a connection to a host
type Host struct {
	// ID used to identify the connection
	Id string

	// Address:Port being connected to
	Address string

	// PEM-format TLS Certificate
	Cert []byte

	// Indicates whether connection timeout should be disabled
	DisableTimeout bool
}

// Stores a connection and its associated information
type connection struct {
	// Address:Port being connected to
	Address string

	// GRPC connection object
	Connection *grpc.ClientConn

	// Credentials object used to establish the connection
	Creds credentials.TransportCredentials

	// RSA Public Key corresponding to the TLS Certificate
	RsaPublicKey *rsa.PublicKey
}

// Returns true if the connection is non-nil and alive
func (c *connection) isAlive() bool {
	if c.Connection == nil {
		return false
	}
	state := c.Connection.GetState()
	return state == connectivity.Idle || state == connectivity.Connecting ||
		state == connectivity.Ready
}

// Connect creates a connection
func (c *connection) connect(maxRetries int64) (err error) {

	// Configure TLS options
	var securityDial grpc.DialOption
	if c.Creds != nil {
		// Create the gRPC client with TLS
		securityDial = grpc.WithTransportCredentials(c.Creds)
	} else {
		// Create the gRPC client without TLS
		jww.WARN.Printf("Connecting to %v without TLS!", c.Address)
		securityDial = grpc.WithInsecure()
	}

	// Attempt to establish a new connection
	for numRetries := int64(0); numRetries < maxRetries && !c.isAlive(); numRetries++ {

		jww.INFO.Printf("Connecting to address %+v. Attempt number %+v of %+v",
			c.Address, numRetries, maxRetries)

		// If timeout is enabled, the max wait time becomes
		// ~14 seconds (with maxRetries=100)
		backoffTime := 2 * (numRetries/16 + 1)
		if backoffTime > 15 {
			backoffTime = 15
		}
		ctx, cancel := ConnectionContext(time.Duration(backoffTime))

		// Create the connection
		c.Connection, err = grpc.DialContext(ctx, c.Address, securityDial,
			grpc.WithBlock(), grpc.WithBackoffMaxDelay(time.Minute*5))
		if err != nil {
			jww.ERROR.Printf("Attempt number %+v to connect to %s failed: %+v\n",
				numRetries, c.Address, errors.New(err.Error()))
		}
		cancel()
	}

	// Verify that the connection was established successfully
	if !c.isAlive() {
		return errors.New(fmt.Sprintf(
			"Last try to connect to %s failed. Giving up", c.Address))
	}

	// Add the successful connection to the Manager
	jww.INFO.Printf("Successfully connected to %v", c.Address)
	return
}

// Creates TransportCredentials and RSA PublicKey objects
// using a PEM-encoded TLS Certificate
func (c *connection) setCredentials(connInfo *Host) error {

	// If no TLS Certificate specified, print a warning and do nothing
	if connInfo.Cert == nil || len(connInfo.Cert) == 0 {
		jww.WARN.Printf("No TLS Certificate specified!")
		return nil
	}

	// Obtain the DNS name included with the certificate
	dnsName := ""
	cert, err := tlsCreds.LoadCertificate(string(connInfo.Cert))
	if err != nil {
		s := fmt.Sprintf("Error forming transportCredentials: %+v", err)
		return errors.New(s)
	}
	if len(cert.DNSNames) > 0 {
		dnsName = cert.DNSNames[0]
	}

	// Create the TLS Credentials object
	c.Creds, err = tlsCreds.NewCredentialsFromPEM(string(connInfo.Cert),
		dnsName)
	if err != nil {
		s := fmt.Sprintf("Error forming transportCredentials: %+v", err)
		return errors.New(s)
	}

	// Create the RSA Public Key object
	c.RsaPublicKey, err = tlsCreds.NewPublicKeyFromPEM(connInfo.Cert)
	if err != nil {
		s := fmt.Sprintf("Error extracting PublicKey: %+v", err)
		return errors.New(s)
	}

	return nil
}

// Stringer interface for connection
func (c *connection) String() string {
	addr := c.Address
	actualConnection := c.Connection
	creds := c.Creds

	var state connectivity.State
	if actualConnection != nil {
		state = actualConnection.GetState()
	}

	serverName := "<nil>"
	protocolVersion := "<nil>"
	securityVersion := "<nil>"
	securityProtocol := "<nil>"
	if creds != nil {
		serverName = creds.Info().ServerName
		securityVersion = creds.Info().SecurityVersion
		protocolVersion = creds.Info().ProtocolVersion
		securityProtocol = creds.Info().SecurityProtocol
	}
	return fmt.Sprintf(
		"Addr: %v\tState: %v\tTLS ServerName: %v\t"+
			"TLS ProtocolVersion: %v\tTLS SecurityVersion: %v\t"+
			"TLS SecurityProtocol: %v\n",
		addr, state, serverName, protocolVersion,
		securityVersion, securityProtocol)
}

package connect

import (
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"google.golang.org/grpc"
	"strings"
)

// Sets up or recovers the Host's connection
// Then runs the given Send function
func (c *ProtoComms) transmit(host *Host, f func(conn *grpc.ClientConn) (interface{},
	error)) (result interface{}, err error) {

	if host.GetAddress() == "" {
		return nil, errors.New("Host address is blank, host might be receive only.")
	}

	numConnects, numAuths, lastEvent := 0, 0, 0

	host.sendMux.Lock()

connect:
	// Ensure the connection is running
	if !host.isAlive() {

		//do not attempt to connect again if multiple attempts have been made
		if numConnects == maxConnects {
			host.sendMux.Unlock()

			return nil, errors.WithMessage(err, "Maximum number of connects attempted")
		}

		//denote that a connection is being tried
		lastEvent = con

		//attempt to make the connection
		jww.INFO.Printf("Host %s not connected, attempting to connect...", host.id.String())
		err = host.connect()
		//if connection cannot be made, do not retry
		if err != nil {
			host.sendMux.Unlock()
			return nil, errors.WithMessage(err, "Failed to connect")
		}

		//denote the connection attempt
		numConnects++
	}

authorize:
	// Establish authentication if required
	if host.authenticationRequired() {
		//do not attempt to connect again if multiple attempts have been made
		if numAuths == maxAuths {
			host.sendMux.Unlock()

			return nil, errors.New("Maximum number of authorizations attempted")
		}

		//do not try multiple auths in a row
		if lastEvent == auth {
			host.sendMux.Unlock()

			return nil, errors.New("Cannot attempt to authorize with host multiple times in a row")
		}

		//denote that an auth is being tried
		lastEvent = auth

		jww.INFO.Printf("Attempting to establish authentication with host %s", host.id.String())
		err = c.clientHandshake(host)
		if err != nil {
			//if failure of connection, retry connection
			if isConnError(err) {
				jww.INFO.Printf("Failed to auth due to connection issue: %s", err)
				goto connect
			}
			host.sendMux.Unlock()

			//otherwise, return the error
			return nil, errors.New("Failed to authenticate")
		}

		//denote the authorization attempt
		numAuths++
	}

	//denote that a send is being tried
	lastEvent = send
	// Attempt to send to host
	host.sendMux.Unlock()

	result, err = host.transmit(f)
	// If failed to authenticate, retry negotiation by jumping to the top of the loop
	if err != nil {
		//if failure of connection, retry connection
		if isConnError(err) {
			host.sendMux.Lock()
			jww.INFO.Printf("Failed send due to connection issue: %s", err)
			goto connect
		}

		// Handle resetting authentication
		if IsAuthError(err) {
			jww.INFO.Printf("Failed send due to auth error, retrying authentication: %s", err.Error())
			host.sendMux.Lock()

			host.transmissionToken.Clear()
			goto authorize
		}

		// otherwise, return the error
		return nil, errors.WithMessage(err, "Failed to send")
	}

	return result, err
}

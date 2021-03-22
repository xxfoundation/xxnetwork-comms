///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package connect

import (
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"google.golang.org/grpc"
	"strings"
)

const MaxRetries = 3
const inCoolDownErr = "Host is in cool down. Cannot connect."

// Sets up or recovers the Host's connection
// Then runs the given Send function
func (c *ProtoComms) transmit(host *Host, f func(conn *grpc.ClientConn) (interface{},
	error)) (result interface{}, err error) {

	if host.GetAddress() == "" {
		return nil, errors.New("Host address is blank, host might be receive only.")
	}

	for numRetries := 0; numRetries < MaxRetries; numRetries++ {
		err = nil
		//reconnect if necessary
		connected, connectionCount := host.Connected()
		if !connected {
			connectionCount, err = c.connect(host, connectionCount)
			if err != nil {
				if strings.Contains(err.Error(), inCoolDownErr) {
					return nil, err
				}
				jww.WARN.Printf("Failed to connect to Host on attempt "+
					"%v/%v : %s", numRetries+1, MaxRetries, err)
				continue
			}
		}

		//transmit
		result, err = host.transmit(f)

		// if the transmission goes well or it is a domain specific error, return
		if err == nil || !(isConnError(err) || IsAuthError(err)) {
			return result, err
		}
		host.conditionalDisconnect(connectionCount)
		jww.WARN.Printf("Failed to send to Host on attempt %v/%v: %+v",
			numRetries+1, MaxRetries, err)
	}

	// Checks if the received error is a among excluded errors
	// If it is not an excluded error, update host's metrics
	if !host.IsExcludedError(err) {
		host.metrics.IncrementErrors()
	}

	return nil, err
}

func (c *ProtoComms) connect(host *Host, count uint64) (uint64, error) {
	host.sendMux.Lock()
	defer host.sendMux.Unlock()

	if host.coolOffBucket != nil {
		if host.inCoolOff {
			if host.coolOffBucket.IsEmpty() {
				host.inCoolOff = false
			} else {
				return 0, errors.New(inCoolDownErr)
			}
		}
		host.inCoolOff = !host.coolOffBucket.Add(1)
		if host.inCoolOff {
			return 0, errors.New(inCoolDownErr)
		}
	}

	//if the connection is alive return, it is possible for another transmission
	//to connect between releasing the read lock and taking the write lick
	if !host.isAlive() {
		//connect to host
		jww.INFO.Printf("Host %s not connected, attempting to connect...",
			host.id)
		err := host.connect()

		count = host.connectionCount

		//if connection cannot be made, do not retry
		if err != nil {
			host.disconnect()
			return count, errors.WithMessagef(err, "Failed to connect to Host %s",
				host.id)
		}
	}

	//check if authentication is needed
	if host.authenticationRequired() {
		jww.INFO.Printf("Attempting to establish authentication with host %s",
			host.id)
		err := c.clientHandshake(host)

		//if authentication cannot be made, do not retry
		if err != nil {
			host.disconnect()
			return count, errors.WithMessagef(err, "Failed to authenticate with "+
				"host: %s", host.id)
		}
	}

	return count, nil
}

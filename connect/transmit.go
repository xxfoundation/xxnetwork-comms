////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package connect

import (
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"strings"
)

const inCoolDownErr = "Host is in cool down. Cannot connect."
const lastTryErr = "Last try to connect to"

// transmit sets up or recovers the Host's connection
// Then runs the given Send function
// This has a bug where if a disconnect happens after "host.transmit(f)"
// and the comms was unsuccessful and is retried, this code will connect
// and then do the operation, leaving the host as connected. In a system
// like the host pool in client, this will cause untracked connections.
// Given that connections have timeouts, this is a minor issue
func (c *ProtoComms) transmit(host *Host, f func(conn Connection) (interface{},
	error)) (result interface{}, err error) {

	if host.GetAddress() == "" {
		return nil, errors.New("Host address is blank, host might be receive only.")
	}

	for numRetries := uint32(0); numRetries < host.params.MaxRetries; numRetries++ {
		err = nil
		//reconnect if necessary
		host.connectionMux.RLock()
		connected, connectionCount := host.connectedUnsafe()
		if !connected {
			host.connectionMux.RUnlock()

			// if auto-connect is not enable, return an error because
			// we cannot connect and we cannot send to a disconnected
			// host
			if host.params.DisableAutoConnect {
				return nil, errors.Errorf("Cannot send to a disconnected" +
					"host when AutoConnect is disabled")
			}
			host.connectionMux.Lock()
			connectionCount, err = c.connect(host, connectionCount)
			host.connectionMux.Unlock()
			if err != nil {
				if strings.Contains(err.Error(), inCoolDownErr) ||
					strings.Contains(err.Error(), lastTryErr) {
					return nil, err
				}
				jww.WARN.Printf("Failed to connect to Host on attempt "+
					"%v/%v : %s", numRetries+1, host.params.MaxSendRetries, err)
				continue
			}
			host.connectionMux.RLock()
		}

		// TODO: this is a temporary fix for xx-4337 (segfault in invoke)
		// Should be removed once the root cause is addressed
		if !host.IsWeb() && host.connection.GetGrpcConn() == nil {
			err = errors.New("Cannot send; connection is nil")
		} else {
			//transmit
			result, err = host.transmit(f)
		}
		host.connectionMux.RUnlock()

		// if the transmission goes well or if it is a domain specific error, return
		if err == nil || !(isConnError(err) || IsAuthError(err)) {
			return result, err
		}
		host.connectionMux.Lock()
		host.conditionalDisconnect(connectionCount)
		host.connectionMux.Unlock()
		jww.WARN.Printf("Failed to send to Host on attempt %v/%v: %+v",
			numRetries+1, host.params.MaxSendRetries, err)
	}

	return nil, err
}

func (c *ProtoComms) connect(host *Host, count uint64) (uint64, error) {
	if host.coolOffBucket != nil {
		if host.inCoolOff {
			if host.coolOffBucket.IsEmpty() {
				host.inCoolOff = false
			} else {
				return 0, errors.New(inCoolDownErr)
			}
		}
		good, _ := host.coolOffBucket.Add(1)
		host.inCoolOff = !good
		if host.inCoolOff {
			return 0, errors.New(inCoolDownErr)
		}
	}

	//if the connection is alive return, it is possible for another transmission
	//to connect between releasing the read lock and taking the write lock
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

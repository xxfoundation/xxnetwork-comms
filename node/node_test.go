////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package node

import (
	"fmt"
	"gitlab.com/elixxir/comms/connect"
	"gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/testkeys"
	"sync"
	"testing"
)

var serverAddressLock sync.Mutex
var ServerAddress = ""

var serverPortLock sync.Mutex
var serverPort = 5000

func getNextServerAddress() string {
	serverPortLock.Lock()
	defer func() {
		serverPort++
		serverPortLock.Unlock()
	}()
	return fmt.Sprintf("localhost:%d", serverPort)
}

// Tests whether the server can be connected to and run an RPC with TLS enabled
func TestTLS(t *testing.T) {
	serverAddressLock.Lock()
	defer serverAddressLock.Unlock()
	ServerAddress = getNextServerAddress()
	t.Log(ServerAddress)
	connect.ServerCertPath = testkeys.GetNodeCertPath()
	shutdown := StartServer(ServerAddress, NewImplementation(),
		testkeys.GetNodeCertPath(), testkeys.GetNodeKeyPath())
	// Reset TLS-related global variables
	defer func() {
		connect.ServerCertPath = ""
		shutdown()
	}()
	_, err := SendAskOnline(ServerAddress, &mixmessages.Ping{})
	if err != nil {
		t.Error(err)
	}
}

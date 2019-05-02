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

var serverPortLock sync.Mutex
var serverPort = 5000

// Basic implementation of Stringer for connection ID
type MockID string

func (m MockID) String() string {
	return string(m)
}

func getNextServerAddress() string {
	serverPortLock.Lock()
	defer func() {
		serverPort++
		serverPortLock.Unlock()
	}()
	return fmt.Sprintf("0.0.0.0:%d", serverPort)
}

// Tests whether the server can be connected to and run an RPC with TLS enabled
func TestTLS(t *testing.T) {
	serverAddress := getNextServerAddress()
	server := StartServer(serverAddress, NewImplementation(),
		testkeys.GetNodeCertPath(), testkeys.GetNodeKeyPath())
	serverAddress2 := getNextServerAddress()
	server2 := StartServer(serverAddress2, NewImplementation(),
		testkeys.GetNodeCertPath(), testkeys.GetNodeKeyPath())
	creds := connect.NewCredentialsFromFile(testkeys.GetNodeCertPath(),
		"*.cmix.rip")
	connectionID := MockID("server2toserver")
	// It might make more sense to call the RPC on the connection object
	// that's returned from this
	server2.ConnectToNode(connectionID,
		&connect.ConnectionInfo{
			Address:    serverAddress,
			Creds:      creds,
			Connection: nil,
		})
	// Reset TLS-related global variables
	defer server.Shutdown()
	defer server2.Shutdown()
	_, err := server2.SendAskOnline(connectionID, &mixmessages.Ping{})
	if err != nil {
		t.Error(err)
	}
}

////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package node

import (
	"fmt"
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

	keyPath := testkeys.GetNodeKeyPath()
	keyData := testkeys.LoadFromPath(keyPath)
	certPath := testkeys.GetNodeCertPath()
	certData := testkeys.LoadFromPath(certPath)

	server := StartNode(serverAddress, NewImplementation(),
		certData, keyData)
	serverAddress2 := getNextServerAddress()
	server2 := StartNode(serverAddress2, NewImplementation(),
		certData, keyData)
	connectionID := MockID("server2toserver")
	// It might make more sense to call the RPC on the connection object
	// that's returned from this
	server2.ConnectToNode(connectionID, serverAddress, certData, false)
	// Reset TLS-related global variables
	defer server.Shutdown()
	defer server2.Shutdown()
	_, err := server2.SendAskOnline(connectionID, &mixmessages.Ping{})
	if err != nil {
		t.Error(err)
	}
}

func TestBadCerts(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	Address := getNextServerAddress()

	_ = StartNode(Address, NewImplementation(),
		[]byte("bad cert"), []byte("bad key"))
}

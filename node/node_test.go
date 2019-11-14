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
	defer server.Shutdown()
	defer server2.Shutdown()
	var manager connect.Manager

	testId := "test"
	host, err := connect.NewHost(serverAddress, certData, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}
	manager.AddHost(testId, host)

	_, err = server2.SendAskOnline(host, &mixmessages.Ping{})
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

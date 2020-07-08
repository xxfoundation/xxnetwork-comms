///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package node

import (
	"fmt"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/elixxir/comms/testkeys"
	"gitlab.com/elixxir/primitives/id"
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
	testNodeID := id.NewIdFromString("test", id.Node, t)

	server := StartNode(testNodeID, serverAddress, NewImplementation(),
		certData, keyData)
	serverAddress2 := getNextServerAddress()
	server2 := StartNode(testNodeID, serverAddress2, NewImplementation(),
		certData, keyData)
	defer server.Shutdown()
	defer server2.Shutdown()
	var manager connect.Manager

	host, err := manager.AddHost(testNodeID, serverAddress, certData, false, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	_, err = server2.SendAskOnline(host)
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

	testID := id.NewIdFromString("test", id.Node, t)

	_ = StartNode(testID, Address, NewImplementation(),
		[]byte("bad cert"), []byte("bad key"))
}

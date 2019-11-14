////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package gateway

import (
	"fmt"
	"gitlab.com/elixxir/comms/connect"
	"gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/node"
	"gitlab.com/elixxir/comms/testkeys"
	"sync"
	"testing"
)

var serverPortLock sync.Mutex
var serverPort = 5500

func getNextServerAddress() string {
	serverPortLock.Lock()
	defer func() {
		serverPort++
		serverPortLock.Unlock()
	}()
	return fmt.Sprintf("0.0.0.0:%d", serverPort)
}

var gatewayPortLock sync.Mutex
var gatewayPort = 5600

func getNextGatewayAddress() string {
	gatewayPortLock.Lock()
	defer func() {
		gatewayPort++
		gatewayPortLock.Unlock()
	}()
	return fmt.Sprintf("0.0.0.0:%d", gatewayPort)
}

// Tests whether the gateway can be connected to and run an RPC with TLS enabled
func TestTLS(t *testing.T) {
	keyPath := testkeys.GetNodeKeyPath()
	keyData := testkeys.LoadFromPath(keyPath)
	certPath := testkeys.GetNodeCertPath()
	certData := testkeys.LoadFromPath(certPath)

	GatewayAddress := getNextGatewayAddress()
	gateway := StartGateway(GatewayAddress, NewImplementation(),
		certData, keyData)
	defer gateway.Shutdown()
	ServerAddress := getNextServerAddress()
	server := node.StartNode(ServerAddress, node.NewImplementation(),
		certData, keyData)
	defer server.Shutdown()
	var manager connect.Manager

	testId := "test"
	host, err := connect.NewHost(ServerAddress, certData, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}
	manager.AddHost(testId, host)

	err = gateway.PostNewBatch(host, &mixmessages.Batch{})
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

	_ = StartGateway(Address, NewImplementation(),
		[]byte("bad cert"), []byte("bad key"))
}

////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package gateway

import (
	"fmt"
	"gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/node"
	"gitlab.com/elixxir/comms/testkeys"
	"sync"
	"testing"
)

var serverPortLock sync.Mutex
var serverPort = 5500

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
	connID := MockID("gatewayToServer")
	gateway.ConnectToNode(connID,
		ServerAddress, certData)

	err := gateway.PostNewBatch(connID, &mixmessages.Batch{})
	if err != nil {
		t.Error(err)
	}
}

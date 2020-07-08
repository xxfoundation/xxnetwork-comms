///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package gateway

import (
	"fmt"
	"gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/node"
	"gitlab.com/elixxir/comms/testkeys"
	"gitlab.com/elixxir/primitives/id"
	"gitlab.com/xx_network/comms/connect"
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
	testID := id.NewIdFromString("test", id.Gateway, t)
	gateway := StartGateway(testID, GatewayAddress, NewImplementation(),
		certData, keyData)
	defer gateway.Shutdown()
	ServerAddress := getNextServerAddress()
	testNodeID := id.NewIdFromString("test", id.Node, t)
	server := node.StartNode(testNodeID, ServerAddress, node.NewImplementation(),
		certData, keyData)
	defer server.Shutdown()
	var manager connect.Manager

	testId := id.NewIdFromString("test", id.Node, t)
	host, err := manager.AddHost(testId, ServerAddress, certData, false, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

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

	testID := id.NewIdFromString("test", id.Node, t)
	_ = StartGateway(testID, Address, NewImplementation(),
		[]byte("bad cert"), []byte("bad key"))
}

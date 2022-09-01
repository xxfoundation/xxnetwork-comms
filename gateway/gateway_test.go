////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package gateway

import (
	"fmt"
	"gitlab.com/elixxir/comms/node"
	"gitlab.com/elixxir/comms/testkeys"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/gossip"
	"gitlab.com/xx_network/primitives/id"
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
		certData, keyData, gossip.DefaultManagerFlags())
	defer gateway.Shutdown()
	ServerAddress := getNextServerAddress()
	testNodeID := id.NewIdFromString("test", id.Node, t)
	server := node.StartNode(testNodeID, ServerAddress, 0, node.NewImplementation(),
		certData, keyData)
	defer server.Shutdown()
	manager := connect.NewManagerTesting(t)

	testId := id.NewIdFromString("test", id.Node, t)
	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(testId, ServerAddress, certData, params)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	_, err = gateway.GetRoundBufferInfo(host)
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
		[]byte("bad cert"), []byte("bad key"), gossip.DefaultManagerFlags())
}

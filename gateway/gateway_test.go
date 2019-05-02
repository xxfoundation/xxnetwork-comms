////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package gateway

import (
	"fmt"
	"gitlab.com/elixxir/comms/mixmessages"
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
	GatewayAddress := getNextServerAddress()
	gateway := StartGateway(GatewayAddress, NewImplementation(),
		testkeys.GetGatewayCertPath(), testkeys.GetGatewayKeyPath())
	// Reset TLS-related global variables
	defer gateway.Shutdown()
	// TODO Populate this ID by properly connecting to the gateway with a
	//  connection pool somewhere
	err := gateway.SendBatch(MockID("5"),
		[]*mixmessages.Batch{})
	if err != nil {
		t.Error(err)
	}
}

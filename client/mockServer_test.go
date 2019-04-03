////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains a dummy/mock server instance for testing purposes

package client

import (
	"fmt"
	"sync"
)

var serverPortLock sync.Mutex
var serverPort = 5700

func getNextServerAddress() string {
	serverPortLock.Lock()
	defer func() {
		serverPort++
		serverPortLock.Unlock()
	}()
	return fmt.Sprintf("localhost:%d", serverPort)
}

var gatewayPortLock sync.Mutex
var gatewayPort = 5800

func getNextGatewayAddress() string {
	gatewayPortLock.Lock()
	defer func() {
		gatewayPort++
		gatewayPortLock.Unlock()
	}()
	return fmt.Sprintf("localhost:%d", gatewayPort)
}

////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package clientregistrar

import (
	"fmt"
	"sync"
	"testing"

	"gitlab.com/elixxir/comms/client"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/testkeys"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/primitives/id"
)

var serverPortLock sync.Mutex
var serverPort = 5900

func getNextServerAddress() string {
	serverPortLock.Lock()
	defer func() {
		serverPort++
		serverPortLock.Unlock()
	}()
	return fmt.Sprintf("localhost:%d", serverPort)
}

// Tests whether the server can be connected to and run an RPC with TLS enabled
func TestTLS(t *testing.T) {
	RegAddress := getNextServerAddress()

	keyPath := testkeys.GetNodeKeyPath()
	keyData := testkeys.LoadFromPath(keyPath)
	certPath := testkeys.GetNodeCertPath()
	certData := testkeys.LoadFromPath(certPath)
	testId := id.NewIdFromString("test", id.Generic, t)

	rg := StartClientRegistrarServer(testId, RegAddress,
		NewImplementation(),
		certData, keyData)
	// Well, client shouldn't have a server type because it's not a server
	// It's a client
	// So, we need some way to add a connection to the manager for the client
	defer rg.Shutdown()
	var c client.Comms
	manager := connect.NewManagerTesting(t)

	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(testId, RegAddress, certData, params)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	_, err = c.SendRegistrationMessage(host, &pb.ClientRegistration{})
	if err != nil {
		t.Errorf("RegistrationMessage: Error received: %s", err)
	}
}

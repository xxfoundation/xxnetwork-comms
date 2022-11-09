////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package authorizer

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/node"
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

	impl := NewImplementation()
	impl.Functions.Authorize = func(auth *pb.AuthorizerAuth, ipAddr string) (err error) {
		return errors.New("This function ran")
	}

	rg := StartAuthorizerServer(testId, RegAddress,
		impl, certData, keyData)
	// Well, client shouldn't have a server type because it's not a server
	// It's a client
	// So, we need some way to add a connection to the manager for the client
	defer rg.Shutdown()
	var c node.Comms
	manager := connect.NewManagerTesting(t)

	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(testId, RegAddress, certData, params)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	_, err = c.SendAuthorizerAuth(host, &pb.AuthorizerAuth{})
	if err == nil && err.Error() != "This function ran" {
		t.Errorf("RegistrationMessage: Error received: %s", err)
	}
}

func TestBadCerts(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	RegAddress := getNextServerAddress()
	testId := id.NewIdFromString("test", id.Generic, t)

	_ = StartAuthorizerServer(testId, RegAddress, NewImplementation(),
		[]byte("bad cert"), []byte("bad key"))
}

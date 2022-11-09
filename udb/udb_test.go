////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package udb

import (
	"fmt"
	"sync"
	"testing"
)

var serverPortLock sync.Mutex
var serverPort = 5950

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
	/*
		RegAddress := getNextServerAddress()

		keyPath := testkeys.GetNodeKeyPath()
		keyData := testkeys.LoadFromPath(keyPath)
		certPath := testkeys.GetNodeCertPath()
		certData := testkeys.LoadFromPath(certPath)
		testId := id.NewIdFromString("test", id.Generic, t)

		rg := StartServer(testId, RegAddress,
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
		// Now call something... note that client call would need to be
		// implemented inside client.
		_, err = c.SendClientCall(host, &pb.UserRegistration{})
		if err != nil {
			t.Errorf("RegistrationMessage: Error received: %s", err)
		}
	*/
}

// This test doesn't start the server, and instead tests the endpoint through to
// the handler
func TestTwo(t *testing.T) {
	/*
		impl := NewImplementation()
		comms := &Comms{
			handler: impl,
		}

		IWasCalled := false
		clientCall := func(msg *pb.PermissioningPoll,
			auth *connect.Auth,
			serverAddress string) (
				*pb.PermissionPollResponse, error) {
					IWasCalled = true
					return &pb.PermissionPollResponse{}, nil
				}
		impl.Functions.ClientCall = clientCall

		// Calls ClientCall inside endpoint.go, which calls impl.ClientCall,
		// which then calls impl.Functions.ClientCall which is the dummy call
		// or the one you would replace right above.
		// note that this will technically fail because you would need a valid
		// protocomms object (not usually true for these..)
		comms.ClientCall(context.Context{}, nil)

		if IWasCalled != true {
			t.Errorf("ClientCall not called as expected")
		}
	*/
}

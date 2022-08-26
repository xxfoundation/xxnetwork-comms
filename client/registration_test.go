///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package client

import (
	"testing"

	"gitlab.com/elixxir/comms/clientregistrar"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/registration"
	"gitlab.com/elixxir/comms/testutils"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/primitives/id"
	"gitlab.com/xx_network/primitives/ndf"
)

// Smoke test SendRegistrationMessage
func TestSendRegistrationMessage(t *testing.T) {
	GatewayAddress := getNextAddress()
	testId := id.NewIdFromString("test", id.Generic, t)
	clientId := id.NewIdFromString("client", id.Generic, t)

	rg := clientregistrar.StartClientRegistrarServer(testId, GatewayAddress,
		registration.NewImplementation(), nil, nil)
	defer rg.Shutdown()

	for _, connectionType := range []connect.ConnectionType{connect.Grpc, connect.Web} {
		c, err := NewClientComms(clientId, nil, nil, nil)
		if err != nil {
			t.Errorf("Can't create client comms: %+v", err)
		}
		manager := connect.NewManagerTesting(t)

		params := connect.GetDefaultHostParams()
		params.ConnectionType = connectionType
		params.AuthEnabled = false
		host, err := manager.AddHost(testId, GatewayAddress, nil, params)
		if err != nil {
			t.Errorf("Unable to call NewHost: %+v", err)
		}

		_, err = c.SendRegistrationMessage(host, &pb.ClientRegistration{})
		if err != nil {
			t.Errorf("RegistrationMessage: Error received: %+v", err)
		}
	}
}

// Smoke test RequestNdf
func TestSendGetUpdatedNDF(t *testing.T) {
	GatewayAddress := getNextAddress()
	testId := id.NewIdFromString("test", id.Generic, t)
	clientId := id.NewIdFromString("client", id.Generic, t)

	rg := registration.StartRegistrationServer(testId, GatewayAddress, &MockRegistration{}, nil, nil, nil)
	defer rg.Shutdown()

	c, err := NewClientComms(clientId, nil, nil, nil)
	if err != nil {
		t.Errorf("Can't create client comms: %+v", err)
	}
	manager := connect.NewManagerTesting(t)

	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(testId, GatewayAddress, nil, params)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	_, err = c.RequestNdf(host, &pb.NDFHash{})
	if err != nil {
		t.Errorf("RequestNdf: Error received: %+v", err)
	}
}

// Test that Poll NDF handles all comms errors returned properly, and that it decodes and successfully returns an ndf
func TestProtoComms_PollNdf(t *testing.T) {

	// Define a client object
	clientId := id.NewIdFromString("client", id.Generic, t)
	c, err := NewClientComms(clientId, nil, nil, nil)
	if err != nil {
		t.Errorf("Can't create client comms: %+v", err)
	}

	mockPermServer := registration.StartRegistrationServer(&id.Permissioning, RegistrationAddr, RegistrationHandler, nil, nil, nil)
	defer mockPermServer.Shutdown()

	newNdf := &ndf.NetworkDefinition{}

	// Test that poll ndf fails if getHost returns an error
	GetHostErrBool = false
	RequestNdfErr = nil

	_, err = c.RetrieveNdf(newNdf)

	if err == nil {
		t.Errorf("GetHost should have failed but it didnt't: %+v", err)
		t.Fail()
	}

	// Test that pollNdf returns an error in this case
	// This enters an infinite loop is there a way to fix this test?

	// Test that pollNdf Fails if it cant decode the request msg
	RequestNdfErr = nil
	GetHostErrBool = true
	NdfToreturn.Ndf = []byte(ExampleBadNdfJSON)
	_, err = c.RetrieveNdf(newNdf)

	if err == nil {
		t.Logf("RequestNdf should have failed to parse bad ndf: %+v", err)
		t.Fail()
	}
	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	_, err = c.ProtoComms.AddHost(&id.Permissioning, RegistrationAddr, nil, params)
	if err != nil {
		t.Errorf("Failed to add permissioning as a host: %+v", err)
	}

	// Test that pollNDf Is successful with expected result
	RequestNdfErr = nil
	GetHostErrBool = true
	NdfToreturn.Ndf = []byte(testutils.ExampleJSON)
	_, err = c.RetrieveNdf(newNdf)
	// comms.mockManager.AddHost()
	if err != nil {
		t.Logf("Ndf failed to parse: %+v", err)
		t.Fail()
	}

}

// Happy path
func TestProtoComms_PollNdfRepeatedly(t *testing.T) {
	// Define a client object
	clientId := id.NewIdFromString("client", id.Generic, t)
	c, err := NewClientComms(clientId, nil, nil, nil)
	if err != nil {
		t.Errorf("Can't create client comms: %+v", err)
	}
	// Start up the mock reg server
	mockPermServer := registration.StartRegistrationServer(&id.Permissioning, RegistrationAddrErr, RegistrationError, nil, nil, nil)
	defer mockPermServer.Shutdown()

	// Add the host to the comms object
	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	_, err = c.ProtoComms.AddHost(&id.Permissioning, RegistrationAddrErr, nil, params)
	if err != nil {
		t.Errorf("Failed to add permissioning as a host: %+v", err)
	}

	newNdf := &ndf.NetworkDefinition{}

	// This should hit the loop until the number of retries is satisfied in the error handler
	_, err = c.RetrieveNdf(newNdf)
	if err != nil {
		t.Errorf("Expected error case, should not return non-error until attempt #5")
	}
}

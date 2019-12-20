////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package node

import (
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/registration"
	"gitlab.com/elixxir/comms/testkeys"
	"testing"
)

// Smoke test SendNodeRegistration
func TestSendNodeRegistration(t *testing.T) {
	RegAddress := getNextServerAddress()
	server := StartNode(getNextServerAddress(), NewImplementation(),
		nil, nil)
	reg := registration.StartRegistrationServer(RegAddress,
		registration.NewImplementation(), nil, nil)
	defer server.Shutdown()
	defer reg.Shutdown()
	var manager connect.Manager

	testId := "test"
	host, err := manager.AddHost(testId, RegAddress, nil, false, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	msgs := &pb.NodeRegistration{}
	err = server.SendNodeRegistration(host, msgs)
	if err != nil {
		t.Errorf("SendNodeTopology: Error received: %s", err)
	}
}

// Smoke test
func TestComms_RequestNdf(t *testing.T) {
	RegAddress := getNextServerAddress()
	server := StartNode(getNextServerAddress(), NewImplementation(),
		nil, nil)
	reg := registration.StartRegistrationServer(RegAddress,
		registration.NewImplementation(), nil, nil)
	defer server.Shutdown()
	defer reg.Shutdown()
	var manager connect.Manager

	testId := "test"
	host, err := manager.AddHost(testId, RegAddress, nil, false, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	msgs := &pb.NDFHash{}

	_, err = server.RequestNdf(host, msgs)
	if err != nil {
		t.Errorf("RequestNdf: Error received: %s", err)
	}
}

// Smoke test - WITH authentication
func TestComms_RequestNdfWithAuth(t *testing.T) {
	priv := testkeys.LoadFromPath(testkeys.GetNodeKeyPath())
	pub := testkeys.LoadFromPath(testkeys.GetNodeCertPath())

	RegAddress := getNextServerAddress()
	server := StartNode(getNextServerAddress(), NewImplementation(),
		pub, priv)
	reg := registration.StartRegistrationServer(RegAddress,
		registration.NewImplementation(), pub, priv)

	defer server.Shutdown()
	defer reg.Shutdown()
	var manager connect.Manager

	testId := "test"
	host, err := manager.AddHost(testId, RegAddress, pub, false, true)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	msgs := &pb.NDFHash{}

	_, err = server.RequestNdf(host, msgs)
	if err != nil {
		t.Errorf("RequestNdf: Error received: %s", err)
	}
}

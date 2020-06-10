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
	"gitlab.com/elixxir/primitives/id"
	"testing"
)

// Smoke test SendNodeRegistration
func TestSendNodeRegistration(t *testing.T) {
	RegAddress := getNextServerAddress()
	testId := id.NewIdFromString("test", id.Generic, t)
	server := StartNode(testId, getNextServerAddress(), NewImplementation(),
		nil, nil)
	reg := registration.StartRegistrationServer(testId, RegAddress,
		registration.NewImplementation(), nil, nil)
	defer server.Shutdown()
	defer reg.Shutdown()
	var manager connect.Manager

	host, err := manager.AddHost(testId, RegAddress, nil, false, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	msgs := &pb.NodeRegistration{ID: testId.Bytes()}
	err = server.SendNodeRegistration(host, msgs)
	if err != nil {
		t.Errorf("SendNodeTopology: Error received: %s", err)
	}
}

// Smoke test
func TestComms_RequestNdf(t *testing.T) {
	RegAddress := getNextServerAddress()
	testId := id.NewIdFromString("test", id.Generic, t)
	server := StartNode(testId, getNextServerAddress(), NewImplementation(),
		nil, nil)
	reg := registration.StartRegistrationServer(testId, RegAddress,
		registration.NewImplementation(), nil, nil)
	defer server.Shutdown()
	defer reg.Shutdown()
	var manager connect.Manager

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
	testId := id.NewIdFromString("test", id.Generic, t)

	RegAddress := getNextServerAddress()
	server := StartNode(testId, getNextServerAddress(), NewImplementation(),
		pub, priv)
	reg := registration.StartRegistrationServer(testId, RegAddress,
		registration.NewImplementation(), pub, priv)

	defer server.Shutdown()
	defer reg.Shutdown()

	host, err := server.AddHost(testId, RegAddress, pub, false, true)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}
	_, err = reg.AddHost(testId, RegAddress, pub, false, true)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	msgs := &pb.NDFHash{Hash: make([]byte, 0)}

	_, err = server.RequestNdf(host, msgs)
	if err != nil {
		t.Errorf("RequestNdf: Error received: %s", err)
	}
}

// Smoke test
func TestComms_SendPoll(t *testing.T) {
	RegAddress := getNextServerAddress()
	testId := id.NewIdFromString("test", id.Generic, t)
	server := StartNode(testId, getNextServerAddress(), NewImplementation(),
		nil, nil)
	reg := registration.StartRegistrationServer(testId, RegAddress,
		registration.NewImplementation(), nil, nil)
	defer server.Shutdown()
	defer reg.Shutdown()
	var manager connect.Manager

	host, err := manager.AddHost(testId, RegAddress, nil, false, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	msgs := &pb.PermissioningPoll{
		Full: &pb.NDFHash{
			Hash: make([]byte, 0),
		},
		Partial: &pb.NDFHash{
			Hash: make([]byte, 0),
		},
	}

	_, err = server.SendPoll(host, msgs)
	if err != nil {
		t.Errorf("SendPoll: Error received: %+v", err)
	}
}

func TestComms_SendRegistrationCheck(t *testing.T) {
	RegAddress := getNextServerAddress()
	testId := id.NewIdFromString("blah", id.Generic, t)
	server := StartNode(testId, getNextServerAddress(), NewImplementation(),
		nil, nil)
	reg := registration.StartRegistrationServer(testId, RegAddress,
		registration.NewImplementation(), nil, nil)
	defer server.Shutdown()
	defer reg.Shutdown()
	var manager connect.Manager

	host, err := manager.AddHost(testId, RegAddress, nil, false, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	msgs := &pb.RegisteredNodeCheck{
		ID: testId.Bytes(),
	}

	_, err = server.SendRegistrationCheck(host, msgs)
	if err != nil {
		t.Errorf("SendRegistrationCheck: Error received: %+v", err)
	}

}

///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package node

import (
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/registration"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/primitives/id"
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

	msgs := &pb.NodeRegistration{Salt: testId.Bytes()}
	err = server.SendNodeRegistration(host, msgs)
	if err != nil {
		t.Errorf("SendNodeTopology: Error received: %s", err)
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

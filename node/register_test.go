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
	host, err := manager.AddHost(testId, RegAddress, nil, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	msgs := &pb.NodeRegistration{}
	err = server.SendNodeRegistration(host, msgs)
	if err != nil {
		t.Errorf("SendNodeTopology: Error received: %s", err)
	}
}

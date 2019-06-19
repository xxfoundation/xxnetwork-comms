////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package node

import (
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/registration"
	"testing"
)

// Smoke test SendNodeRegistration
func TestSendNodeTopology(t *testing.T) {
	RegAddress := getNextServerAddress()
	server := StartNode(getNextServerAddress(), NewImplementation(),
		"", "")
	reg := registration.StartRegistrationServer(RegAddress,
		registration.NewImplementation(), "", "")
	defer server.Shutdown()
	defer reg.Shutdown()
	connID := MockID("serverToPermissioning")
	server.ConnectToRegistration(connID, RegAddress, nil)

	msgs := &pb.NodeRegistration{}
	err := server.SendNodeRegistration(connID, msgs)
	if err != nil {
		t.Errorf("SendNodeTopology: Error received: %s", err)
	}
}

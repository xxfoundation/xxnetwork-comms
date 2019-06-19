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

// Smoke test SendNodeTopology
func TestSendNodeTopology(t *testing.T) {
	ServerAddress := getNextServerAddress()
	server := StartNode(ServerAddress, NewImplementation(),
		"", "")
	reg := registration.StartRegistrationServer(getNextServerAddress(),
		registration.NewImplementation(), "", "")
	defer server.Shutdown()
	defer reg.Shutdown()
	connID := MockID("permissioningToServer")
	reg.ConnectToNode(connID, ServerAddress, nil)

	msgs := &pb.NodeRegistration{}
	err := server.SendNodeRegistration(connID, msgs)
	if err != nil {
		t.Errorf("SendNodeTopology: Error received: %s", err)
	}
}
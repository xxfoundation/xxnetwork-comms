////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package registration

import (
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/node"
	"gitlab.com/elixxir/comms/testkeys"
	"testing"
)

// Smoke test SendNodeTopology
func TestSendNodeTopology(t *testing.T) {
	ServerAddress := getNextServerAddress()
	server := node.StartNode(ServerAddress, node.NewImplementation(),
		"", "", testkeys.GetNodeCertPath())
	reg := StartRegistrationServer(getNextServerAddress(),
		NewImplementation(), testkeys.GetNodeCertPath(), testkeys.GetNodeKeyPath())
	defer server.Shutdown()
	defer reg.Shutdown()
	connID := MockID("permissioningToServer")
	reg.ConnectToNode(connID, ServerAddress, nil)

	msgs := &pb.NodeTopology{}
	err := reg.SendNodeTopology(connID, msgs)
	if err != nil {
		t.Errorf("SendNodeTopology: Error received: %s", err)
	}
}

func TestSendNodeTopologyNilKeyError(t *testing.T) {
	ServerAddress := getNextServerAddress()
	server := node.StartNode(ServerAddress, node.NewImplementation(),
		"", "", "")
	reg := StartRegistrationServer(getNextServerAddress(),
		NewImplementation(), "", "")
	defer server.Shutdown()
	defer reg.Shutdown()
	connID := MockID("permissioningToServer")
	reg.ConnectToNode(connID, ServerAddress, nil)

	msgs := &pb.NodeTopology{}
	err := reg.SendNodeTopology(connID, msgs)
	if err == nil {
		t.Errorf("SendNodeTopology: did not receive missing private key error")
	}
}

func TestSendNodeTopologyBadKeyError(t *testing.T) {
	ServerAddress := getNextServerAddress()
	server := node.StartNode(ServerAddress, node.NewImplementation(),
		"", "", "")
	reg := StartRegistrationServer(getNextServerAddress(),
		NewImplementation(), "", testkeys.GetNodeCertPath())
	defer server.Shutdown()
	defer reg.Shutdown()
	connID := MockID("permissioningToServer")
	reg.ConnectToNode(connID, ServerAddress, nil)

	msgs := &pb.NodeTopology{}
	err := reg.SendNodeTopology(connID, msgs)
	if err == nil {
		t.Errorf("SendNodeTopology: did not receive bad private key error")
	}
}

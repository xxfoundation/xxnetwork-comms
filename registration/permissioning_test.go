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
	RegAddress := getNextServerAddress()

	keyPath := testkeys.GetNodeKeyPath()
	keyData := testkeys.LoadFromPath(keyPath)
	certPath := testkeys.GetNodeCertPath()
	certData := testkeys.LoadFromPath(certPath)

	server := node.StartNode(ServerAddress, node.NewImplementation(),
		nil, nil)
	reg := StartRegistrationServer(RegAddress,
		NewImplementation(), certData, keyData)
	defer server.Shutdown()
	defer reg.Shutdown()

	connID := MockID("permissioningToServer")
	regID := MockID("Permissioning")

	err := server.ConnectToRegistration(regID, RegAddress, certData, false)
	if err != nil {
		t.Errorf("SendNodeTopology: Node could not connect to"+
			" registration: %s", err)
	}

	err = reg.ConnectToNode(connID, ServerAddress, nil, false)
	if err != nil {
		t.Errorf("SendNodeTopology: Registration could not connect to"+
			" node: %s", err)
	}

	msgs := &pb.NodeTopology{}
	err = reg.SendNodeTopology(connID, msgs)
	if err != nil {
		t.Errorf("SendNodeTopology: Error received: %s", err)
	}
}

func TestSendNodeTopologyNilKey(t *testing.T) {
	ServerAddress := getNextServerAddress()
	RegAddress := getNextServerAddress()

	server := node.StartNode(ServerAddress, node.NewImplementation(),
		nil, nil)
	reg := StartRegistrationServer(RegAddress,
		NewImplementation(), nil, nil)
	defer server.Shutdown()
	defer reg.Shutdown()

	connID := MockID("permissioningToServer")
	regID := MockID("Permissioning")

	_ = server.ConnectToRegistration(regID, RegAddress, nil, false)
	_ = reg.ConnectToNode(connID, ServerAddress, nil, false)

	msgs := &pb.NodeTopology{}
	err := reg.SendNodeTopology(connID, msgs)
	if err != nil {
		t.Errorf("Should not have tried to sign message, instead got: %+v", err)
	}
}

func TestSendNodeTopologyBadMessageError(t *testing.T) {
	ServerAddress := getNextServerAddress()
	RegAddress := getNextServerAddress()

	server := node.StartNode(ServerAddress, node.NewImplementation(),
		nil, nil)
	reg := StartRegistrationServer(RegAddress,
		NewImplementation(), nil, nil)
	defer server.Shutdown()
	defer reg.Shutdown()

	connID := MockID("permissioningToServer")
	regID := MockID("Permissioning")

	_ = server.ConnectToRegistration(regID, RegAddress, nil, false)
	_ = reg.ConnectToNode(connID, ServerAddress, nil, false)

	err := reg.SendNodeTopology(connID, nil)
	if err == nil {
		t.Errorf("SendNodeTopology: did not receive missing private key error")
	}
}

func TestSendNodeTopologyNilMessage(t *testing.T) {
	ServerAddress := getNextServerAddress()
	RegAddress := getNextServerAddress()

	server := node.StartNode(ServerAddress, node.NewImplementation(),
		nil, nil)
	reg := StartRegistrationServer(RegAddress,
		NewImplementation(), nil, nil)
	defer server.Shutdown()
	defer reg.Shutdown()

	connID := MockID("permissioningToServer")
	regID := MockID("Permissioning")

	_ = server.ConnectToRegistration(regID, RegAddress, nil, false)
	_ = reg.ConnectToNode(connID, ServerAddress, nil, false)

	//sgs := &pb.NodeTopology{}
	err := reg.SendNodeTopology(connID, nil)
	if err == nil {
		t.Errorf("Should not have tried to sign message, instead got: %+v", err)
	}
}

func TestSendNodeTopologyBadSignature(t *testing.T) {
	ServerAddress := getNextServerAddress()
	RegAddress := getNextServerAddress()

	server := node.StartNode(ServerAddress, node.NewImplementation(),
		nil, nil)
	reg := StartRegistrationServer(RegAddress,
		NewImplementation(), nil, nil)
	defer server.Shutdown()
	defer reg.Shutdown()

	connID := MockID("permissioningToServer")
	regID := MockID("Permissioning")

	_ = server.ConnectToRegistration(regID, RegAddress, nil, false)
	_ = reg.ConnectToNode(connID, ServerAddress, nil, false)

	msgs := &pb.NodeTopology{}
	err := reg.SendNodeTopology(connID, msgs)
	if err != nil {
		t.Errorf("Should not have tried to sign message, instead got: %+v", err)
	}
}

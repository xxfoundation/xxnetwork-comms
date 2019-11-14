////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package registration

import (
	"gitlab.com/elixxir/comms/connect"
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
	var manager connect.Manager

	testId := "test"
	host, err := connect.NewHost(ServerAddress, nil, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}
	manager.AddHost(testId, host)
	msgs := &pb.NodeTopology{}
	err = reg.SendNodeTopology(host, msgs)
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
	var manager connect.Manager

	testId := "test"
	host, err := connect.NewHost(ServerAddress, nil, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}
	manager.AddHost(testId, host)
	msgs := &pb.NodeTopology{}
	err = reg.SendNodeTopology(host, msgs)
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
	var manager connect.Manager

	testId := "test"
	host, err := connect.NewHost(ServerAddress, nil, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}
	manager.AddHost(testId, host)

	err = reg.SendNodeTopology(host, nil)
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
	var manager connect.Manager

	testId := "test"
	host, err := connect.NewHost(ServerAddress, nil, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}
	manager.AddHost(testId, host)
	err = reg.SendNodeTopology(host, nil)
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
	var manager connect.Manager

	testId := "test"
	host, err := connect.NewHost(ServerAddress, nil, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}
	manager.AddHost(testId, host)
	msgs := &pb.NodeTopology{}
	err = reg.SendNodeTopology(host, msgs)
	if err != nil {
		t.Errorf("Should not have tried to sign message, instead got: %+v", err)
	}
}

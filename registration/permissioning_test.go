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
	msgs := &pb.NodeTopology{}
	err := reg.SendNodeTopology(&connect.Host{
		Id:             "permissioningToServer",
		Address:        ServerAddress,
		Cert:           nil,
		DisableTimeout: false,
	}, msgs)
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
	msgs := &pb.NodeTopology{}
	err := reg.SendNodeTopology(&connect.Host{
		Id:             "permissioningToServer",
		Address:        ServerAddress,
		Cert:           nil,
		DisableTimeout: false,
	}, msgs)
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

	err := reg.SendNodeTopology(&connect.Host{
		Id:             "permissioningToServer",
		Address:        ServerAddress,
		Cert:           nil,
		DisableTimeout: false,
	}, nil)
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
	err := reg.SendNodeTopology(&connect.Host{
		Id:             "permissioningToServer",
		Address:        ServerAddress,
		Cert:           nil,
		DisableTimeout: false,
	}, nil)
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

	msgs := &pb.NodeTopology{}
	err := reg.SendNodeTopology(&connect.Host{
		Id:             "permissioningToServer",
		Address:        ServerAddress,
		Cert:           nil,
		DisableTimeout: false,
	}, msgs)
	if err != nil {
		t.Errorf("Should not have tried to sign message, instead got: %+v", err)
	}
}

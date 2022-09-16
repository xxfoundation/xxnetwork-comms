////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package gateway

import (
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/gossip"
	"gitlab.com/xx_network/primitives/id"
	"testing"
	"time"
)

// Smoke test.
func TestComms_SendPutMessage(t *testing.T) {
	gwAddress1 := getNextGatewayAddress()
	gwAddress2 := getNextGatewayAddress()
	testID1 := id.NewIdFromString("test1", id.Gateway, t)
	testID2 := id.NewIdFromString("test2", id.Gateway, t)
	gw1 := StartGateway(testID1, gwAddress1, NewImplementation(), nil, nil,
		gossip.DefaultManagerFlags())
	gw2 := StartGateway(testID2, gwAddress2, NewImplementation(), nil, nil,
		gossip.DefaultManagerFlags())
	defer gw1.Shutdown()
	defer gw2.Shutdown()
	manager := connect.NewManagerTesting(t)

	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(testID1, gwAddress2, nil, params)
	if err != nil {
		t.Fatalf("Failed to add host to manager: %+v", err)
	}

	_, err = gw1.SendPutMessageProxy(host, &pb.GatewaySlot{}, 2*time.Minute)
	if err != nil {
		t.Errorf("SendPutMessage produced an error: %+v", err)
	}
}

// Smoke test.
func TestComms_SendPutManyMessages(t *testing.T) {
	gwAddress1 := getNextGatewayAddress()
	gwAddress2 := getNextGatewayAddress()
	testID1 := id.NewIdFromString("test1", id.Gateway, t)
	testID2 := id.NewIdFromString("test2", id.Gateway, t)
	gw1 := StartGateway(testID1, gwAddress1, NewImplementation(), nil, nil,
		gossip.DefaultManagerFlags())
	gw2 := StartGateway(testID2, gwAddress2, NewImplementation(), nil, nil,
		gossip.DefaultManagerFlags())
	defer gw1.Shutdown()
	defer gw2.Shutdown()
	manager := connect.NewManagerTesting(t)

	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(testID1, gwAddress2, nil, params)
	if err != nil {
		t.Fatalf("Failed to add host to manager: %+v", err)
	}

	_, err = gw1.SendPutManyMessagesProxy(host, &pb.GatewaySlots{}, 2*time.Minute)
	if err != nil {
		t.Errorf("SendPutMessage produced an error: %+v", err)
	}
}

// Smoke test.
func TestComms_SendRequestNonce(t *testing.T) {
	gwAddress1 := getNextGatewayAddress()
	gwAddress2 := getNextGatewayAddress()
	testID1 := id.NewIdFromString("test1", id.Gateway, t)
	testID2 := id.NewIdFromString("test2", id.Gateway, t)
	gw1 := StartGateway(testID1, gwAddress1, NewImplementation(), nil, nil,
		gossip.DefaultManagerFlags())
	gw2 := StartGateway(testID2, gwAddress2, NewImplementation(), nil, nil,
		gossip.DefaultManagerFlags())
	defer gw1.Shutdown()
	defer gw2.Shutdown()
	manager := connect.NewManagerTesting(t)

	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(testID1, gwAddress2, nil, params)
	if err != nil {
		t.Fatalf("Failed to add host to manager: %+v", err)
	}

	_, err = gw1.SendRequestClientKey(host, &pb.SignedClientKeyRequest{}, 2*time.Minute)
	if err != nil {
		t.Errorf("SendRequestNonce produced an error: %+v", err)
	}
}

// Smoke test.
func TestComms_SendRequestMessages(t *testing.T) {
	gwAddress1 := getNextGatewayAddress()
	gwAddress2 := getNextGatewayAddress()
	testID1 := id.NewIdFromString("test1", id.Gateway, t)
	testID2 := id.NewIdFromString("test2", id.Gateway, t)
	gw1 := StartGateway(testID1, gwAddress1, NewImplementation(), nil, nil,
		gossip.DefaultManagerFlags())
	gw2 := StartGateway(testID2, gwAddress2, NewImplementation(), nil, nil,
		gossip.DefaultManagerFlags())
	defer gw1.Shutdown()
	defer gw2.Shutdown()
	manager := connect.NewManagerTesting(t)

	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(testID1, gwAddress2, nil, params)
	if err != nil {
		t.Fatalf("Failed to add host to manager: %+v", err)
	}

	_, err = gw1.SendRequestMessages(host, &pb.GetMessages{}, 2*time.Minute)
	if err != nil {
		t.Errorf("SendRequestMessages produced an error: %+v", err)
	}
}

///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package client

import (
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/comms/gateway"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/testkeys"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/gossip"
	"gitlab.com/xx_network/primitives/id"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	jww.SetStdoutThreshold(jww.LevelTrace)
	connect.TestingOnlyDisableTLS = true
	os.Exit(m.Run())
}

// Smoke test SendGetMessage
func TestSendPutMessage(t *testing.T) {
	gatewayAddress := getNextAddress()
	testID := id.NewIdFromString("test", id.Gateway, t)
	gw := gateway.StartGateway(testID, gatewayAddress,
		gateway.NewImplementation(), nil, nil, gossip.DefaultManagerFlags())
	defer gw.Shutdown()
	var c Comms

	for _, connectionType := range []connect.ConnectionType{connect.Grpc, connect.Web} {
		manager := connect.NewManagerTesting(t)

		params := connect.GetDefaultHostParams()
		params.ConnectionType = connectionType
		params.AuthEnabled = false
		host, err := manager.AddHost(testID, gatewayAddress, nil, params)
		if err != nil {
			t.Errorf("Unable to call NewHost: %+v", err)
		}

		_, err = c.SendPutMessage(host, &pb.GatewaySlot{}, 10*time.Second)
		if err != nil {
			t.Errorf("PutMessage: Error received: %s", err)
		}
	}
}

// Smoke test SendRequestClientKeyMessage
func TestSendRequestNonceMessage(t *testing.T) {
	gatewayAddress := getNextAddress()
	testID := id.NewIdFromString("test", id.Gateway, t)
	gw := gateway.StartGateway(testID, gatewayAddress,
		gateway.NewImplementation(), nil, nil, gossip.DefaultManagerFlags())
	defer gw.Shutdown()
	var c Comms

	for _, connectionType := range []connect.ConnectionType{connect.Grpc, connect.Web} {
		manager := connect.NewManagerTesting(t)

		params := connect.GetDefaultHostParams()
		params.ConnectionType = connectionType
		params.AuthEnabled = false
		host, err := manager.AddHost(testID, gatewayAddress, nil, params)
		if err != nil {
			t.Errorf("Unable to call NewHost: %+v", err)
		}

		_, err = c.SendRequestClientKeyMessage(host, &pb.SignedClientKeyRequest{})
		if err != nil {
			t.Errorf("SendRequestClientKeyMessage: Error received: %+v", err)
		}
	}
}

// Smoke test SendPoll
func TestComms_SendPoll(t *testing.T) {
	gatewayAddress := getNextAddress()
	testID := id.NewIdFromString("test", id.Gateway, t)
	gw := gateway.StartGateway(testID, gatewayAddress,
		mockGatewayImpl{}, nil, nil, gossip.DefaultManagerFlags())
	defer gw.Shutdown()
	var c Comms

	for _, connectionType := range []connect.ConnectionType{connect.Grpc, connect.Web} {
		manager := connect.NewManagerTesting(t)

		params := connect.GetDefaultHostParams()
		params.ConnectionType = connectionType
		params.AuthEnabled = false
		host, err := manager.AddHost(testID, gatewayAddress, nil, params)
		if err != nil {
			t.Errorf("Unable to call NewHost: %+v", err)
		}

		_, err = c.SendPoll(host,
			&pb.GatewayPoll{
				Partial: &pb.NDFHash{
					Hash: make([]byte, 0),
				},
			})
		if err != nil {
			t.Errorf("SendPoll: Error received: %+v", err)
		}
	}
}

// Smoke test RequestMessages
func TestComms_RequestMessages(t *testing.T) {
	gatewayAddress := getNextAddress()
	testID := id.NewIdFromString("test", id.Gateway, t)
	pk := testkeys.LoadFromPath(testkeys.GetGatewayKeyPath())

	gw := gateway.StartGateway(testID, gatewayAddress,
		gateway.NewImplementation(), nil, nil, gossip.DefaultManagerFlags())
	defer gw.Shutdown()

	for _, connectionType := range []connect.ConnectionType{connect.Grpc, connect.Web} {
		c, err := NewClientComms(testID, nil, pk, nil)
		if err != nil {
			t.Errorf("Could not start client: %v", err)
		}

		params := connect.GetDefaultHostParams()
		params.ConnectionType = connectionType
		params.AuthEnabled = false

		host, err := c.Manager.AddHost(testID, gatewayAddress, nil, params)
		if err != nil {
			t.Errorf("Unable to call NewHost: %+v", err)
		}

		_, err = c.RequestMessages(host,
			&pb.GetMessages{})
		if err != nil {
			t.Errorf("SendPoll: Error received: %+v", err)
		}
	}
}

// Smoke test RequestHistoricalRounds
func TestComms_RequestHistoricalRounds(t *testing.T) {
	gatewayAddress := getNextAddress()
	testID := id.NewIdFromString("test", id.Gateway, t)
	gw := gateway.StartGateway(testID, gatewayAddress,
		gateway.NewImplementation(), nil, nil, gossip.DefaultManagerFlags())
	defer gw.Shutdown()
	pk := testkeys.LoadFromPath(testkeys.GetGatewayKeyPath())

	for _, connectionType := range []connect.ConnectionType{connect.Grpc, connect.Web} {
		c, err := NewClientComms(testID, nil, pk, nil)
		if err != nil {
			t.Errorf("Could not start client: %v", err)
		}

		params := connect.GetDefaultHostParams()
		params.ConnectionType = connectionType
		params.AuthEnabled = false

		host, err := c.Manager.AddHost(testID, gatewayAddress, nil, params)
		if err != nil {
			t.Errorf("Unable to call NewHost: %+v", err)
		}

		_, err = c.RequestHistoricalRounds(host,
			&pb.HistoricalRounds{})
		if err != nil {
			t.Errorf("SendPoll: Error received: %+v", err)
		}
	}
}

type mockGatewayImpl struct{}

func (m mockGatewayImpl) PutMessageProxy(message *pb.GatewaySlot, auth *connect.Auth) (*pb.GatewaySlotResponse, error) {
	return nil, nil
}

func (m mockGatewayImpl) PutManyMessagesProxy(msgs *pb.GatewaySlots, auth *connect.Auth) (*pb.GatewaySlotResponse, error) {
	return nil, nil
}

func (m mockGatewayImpl) PutMessage(message *pb.GatewaySlot, ipAddr string) (*pb.GatewaySlotResponse, error) {
	return nil, nil
}

func (m mockGatewayImpl) PutManyMessages(msgs *pb.GatewaySlots, ipAdd string) (*pb.GatewaySlotResponse, error) {
	return nil, nil
}

func (m mockGatewayImpl) RequestClientKey(message *pb.SignedClientKeyRequest) (*pb.SignedKeyResponse, error) {
	return nil, nil
}

func (m mockGatewayImpl) Poll(msg *pb.GatewayPoll) (*pb.GatewayPollResponse, error) {
	return &pb.GatewayPollResponse{
		PartialNDF:    &NdfToreturn,
		KnownRounds:   []byte("test"),
		Filters:       nil,
		EarliestRound: 0,
	}, nil
}

func (m mockGatewayImpl) RequestHistoricalRounds(msg *pb.HistoricalRounds) (*pb.HistoricalRoundsResponse, error) {
	return nil, nil
}

func (m mockGatewayImpl) RequestMessages(msg *pb.GetMessages) (*pb.GetMessagesResponse, error) {
	return nil, nil
}

func (m mockGatewayImpl) ShareMessages(msg *pb.RoundMessages, auth *connect.Auth) error {
	return nil
}

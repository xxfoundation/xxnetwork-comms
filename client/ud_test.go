package client

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/udb"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/messages"
	"gitlab.com/xx_network/primitives/id"
	"testing"
)

// Smoke test SendGetMessage
func TestComms_SendDeleteMessage(t *testing.T) {
	udAddr := getNextAddress()
	ud := udb.StartServer(&id.UDB, udAddr, udb.NewImplementation(), nil, nil)
	_ = ud.Id
	defer ud.Shutdown()
	var c Comms
	manager := connect.NewManagerTesting(t)

	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(&id.UDB, udAddr, nil, params)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	frr, err := proto.Marshal(&pb.FactRemovalRequest{})
	if err != nil {
		t.Errorf("Unable to call proto.Marshal: %+v", err)
	}

	amm := any.Any{TypeUrl: "gitlab.com/elixxir/comms/mixmessages.FactRemovalRequest", Value: frr}

	_, err = c.SendDeleteMessage(host, &messages.AuthenticatedMessage{Message: &amm})
	if err != nil {
		t.Errorf("DeleteMessage: Error received: %s", err)
	}
}

// Smoke test SendRegisterUser
func TestComms_SendRegisterUser(t *testing.T) {
	udAddr := getNextAddress()
	ud := udb.StartServer(&id.UDB, udAddr, udb.NewImplementation(), nil, nil)
	_ = ud.Id
	defer ud.Shutdown()
	var c Comms
	manager := connect.NewManagerTesting(t)

	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(&id.UDB, udAddr, nil, params)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	frr, err := proto.Marshal(&pb.UDBUserRegistration{})
	if err != nil {
		t.Errorf("Unable to call proto.Marshal: %+v", err)
	}

	amm := any.Any{TypeUrl: "gitlab.com/elixxir/comms/mixmessages.UDBUserRegistration", Value: frr}

	_, err = c.SendRegisterUser(host, &messages.AuthenticatedMessage{Message: &amm})
	if err != nil {
		t.Errorf("DeleteMessage: Error received: %s", err)
	}
}

// Smoke test SendRegisterFact
func TestComms_SendRegisterFact(t *testing.T) {
	udAddr := getNextAddress()
	ud := udb.StartServer(&id.UDB, udAddr, udb.NewImplementation(), nil, nil)
	_ = ud.Id
	defer ud.Shutdown()
	var c Comms
	manager := connect.NewManagerTesting(t)

	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(&id.UDB, udAddr, nil, params)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	frr, err := proto.Marshal(&pb.FactRegisterRequest{})
	if err != nil {
		t.Errorf("Unable to call proto.Marshal: %+v", err)
	}

	amm := any.Any{TypeUrl: "gitlab.com/elixxir/comms/mixmessages.FactRegisterRequest", Value: frr}

	_, err = c.SendRegisterFact(host, &messages.AuthenticatedMessage{Message: &amm})
	if err != nil {
		t.Errorf("DeleteMessage: Error received: %s", err)
	}
}

// Smoke test SendRegisterFact
func TestComms_SendConfirmFact(t *testing.T) {
	udAddr := getNextAddress()
	ud := udb.StartServer(&id.UDB, udAddr, udb.NewImplementation(), nil, nil)
	_ = ud.Id
	defer ud.Shutdown()
	var c Comms
	manager := connect.NewManagerTesting(t)

	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(&id.UDB, udAddr, nil, params)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	frr, err := proto.Marshal(&pb.FactConfirmRequest{})
	if err != nil {
		t.Errorf("Unable to call proto.Marshal: %+v", err)
	}

	amm := any.Any{TypeUrl: "gitlab.com/elixxir/comms/mixmessages.FactConfirmRequest", Value: frr}

	_, err = c.SendConfirmFact(host, &messages.AuthenticatedMessage{Message: &amm})
	if err != nil {
		t.Errorf("DeleteMessage: Error received: %s", err)
	}
}

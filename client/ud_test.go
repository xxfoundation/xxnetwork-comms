package client

import (
	pb "git.xx.network/elixxir/comms/mixmessages"
	"git.xx.network/elixxir/comms/udb"
	"git.xx.network/xx_network/comms/connect"
	"git.xx.network/xx_network/primitives/id"
	"testing"
)

// Smoke test SendGetMessage
func TestComms_SendDeleteMessage(t *testing.T) {
	udAddr := getNextAddress()
	ud := udb.StartServer(&id.UDB, udAddr, udb.NewImplementation(), nil, nil)
	_ = ud.Id
	defer ud.Shutdown()
	c, err := NewClientComms(&id.DummyUser, nil, nil, nil)
	if err != nil {
		t.Error(err)
	}
	manager := connect.NewManagerTesting(t)

	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(&id.UDB, udAddr, nil, params)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	_, err = c.SendDeleteMessage(host, &pb.FactRemovalRequest{})
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
	c, err := NewClientComms(&id.DummyUser, nil, nil, nil)
	if err != nil {
		t.Error(err)
	}
	manager := connect.NewManagerTesting(t)

	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(&id.UDB, udAddr, nil, params)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	_, err = c.SendRegisterUser(host, &pb.UDBUserRegistration{})
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
	c, err := NewClientComms(&id.DummyUser, nil, nil, nil)
	if err != nil {
		t.Error(err)
	}
	manager := connect.NewManagerTesting(t)

	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(&id.UDB, udAddr, nil, params)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	_, err = c.SendRegisterFact(host, &pb.FactRegisterRequest{})
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
	c, err := NewClientComms(&id.DummyUser, nil, nil, nil)
	if err != nil {
		t.Error(err)
	}
	manager := connect.NewManagerTesting(t)

	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(&id.UDB, udAddr, nil, params)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	_, err = c.SendConfirmFact(host, &pb.FactConfirmRequest{})
	if err != nil {
		t.Errorf("DeleteMessage: Error received: %s", err)
	}
}

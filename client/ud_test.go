////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package client

import (
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/udb"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/primitives/id"
	"testing"
)

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

// Smoke test SendGetMessage
func TestComms_SendRemoveFact(t *testing.T) {
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

	_, err = c.SendRemoveFact(host, &pb.FactRemovalRequest{})
	if err != nil {
		t.Errorf("DeleteMessage: Error received: %s", err)
	}
}

// Smoke test SendGetMessage
func TestComms_SendRemoveUser(t *testing.T) {
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

	_, err = c.SendRemoveUser(host, &pb.FactRemovalRequest{})
	if err != nil {
		t.Errorf("DeleteMessage: Error received: %s", err)
	}
}

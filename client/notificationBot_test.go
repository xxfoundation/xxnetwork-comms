///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package client

import (
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/notificationBot"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/primitives/id"
	"testing"
)

// Smoke test for RegisterForNotifications
func TestRegisterForNotifications(t *testing.T) {
	testId := id.NewIdFromString("test", id.Generic, t)
	clientId := id.NewIdFromString("client", id.Generic, t)

	// Start notification bot
	nbAddress := getNextAddress()
	notificationBot := notificationBot.StartNotificationBot(testId, nbAddress,
		notificationBot.NewImplementation(), nil, nil)
	defer notificationBot.Shutdown()

	// Create client's comms object
	c, err := NewClientComms(clientId, nil, nil, nil)
	if err != nil {
		t.Errorf("Can't create client comms: %+v", err)
	}
	var manager connect.Manager

	// Add notification bot to comm's manager
	host, err := manager.AddHost(testId, nbAddress, nil, false, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	// Register client with notification bot
	_, err = c.RegisterForNotifications(host, &pb.NotificationToken{})
	if err != nil {
		t.Errorf("RegistrationMessage: Error received: %s", err)
	}

}

//Smoke test for UnregisterForNotifications
func TestUnregisterForNotifications(t *testing.T) {
	testId := id.NewIdFromString("test", id.Generic, t)
	clientId := id.NewIdFromString("client", id.Generic, t)

	// Start notification bot
	nbAddress := getNextAddress()
	notificationBot := notificationBot.StartNotificationBot(testId, nbAddress,
		notificationBot.NewImplementation(), nil, nil)
	defer notificationBot.Shutdown()

	// Create client's comms object
	c, err := NewClientComms(clientId, nil, nil, nil)
	if err != nil {
		t.Errorf("Can't create client comms: %+v", err)
	}
	var manager connect.Manager

	// Add notification bot to comm's manager
	host, err := manager.AddHost(testId, nbAddress, nil, false, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	// Unregister client with notification bot
	_, err = c.UnregisterForNotifications(host)
	if err != nil {
		t.Errorf("RegistrationMessage: Error received: %s", err)
	}

}

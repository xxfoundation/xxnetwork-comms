////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package client

import (
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/notificationBot"
	"testing"
)

// Smoke test for RegisterForNotifications
func TestRegisterForNotifications(t *testing.T) {
	testId := "test"

	// Start notification bot
	nbAddress := getNextAddress()
	notificationBot := notificationBot.StartNotificationBot(testId, nbAddress,
		notificationBot.NewImplementation(), nil, nil)
	defer notificationBot.Shutdown()

	// Create client's comms object
	c, err := NewClientComms("client", nil, nil, nil)
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
	testId := "test"

	// Start notification bot
	nbAddress := getNextAddress()
	notificationBot := notificationBot.StartNotificationBot(testId, nbAddress,
		notificationBot.NewImplementation(), nil, nil)
	defer notificationBot.Shutdown()

	// Create client's comms object
	c, err := NewClientComms("client", nil, nil, nil)
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

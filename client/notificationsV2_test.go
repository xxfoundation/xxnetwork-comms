package client

import (
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/notificationBot"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/primitives/id"
	"testing"
)

// Smoke test for UnregisterForNotifications
func TestComms_RegisterToken(t *testing.T) {
	testId := id.NewIdFromString("test", id.Generic, t)
	clientId := id.NewIdFromString("client", id.Generic, t)

	// Start notification bot
	nbAddress := getNextAddress()
	nb := notificationBot.StartNotificationBot(testId, nbAddress,
		notificationBot.NewImplementation(), nil, nil)
	defer nb.Shutdown()

	// Create client's comms object
	c, err := NewClientComms(clientId, nil, nil, nil)
	if err != nil {
		t.Errorf("Can't create client comms: %+v", err)
	}
	manager := connect.NewManagerTesting(t)

	// Add notification bot to comm's manager
	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(testId, nbAddress, nil, params)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	// Unregister client with notification bot
	_, err = c.RegisterToken(host, &pb.RegisterTokenRequest{})
	if err != nil {
		t.Errorf("RegistrationMessage: Error received: %s", err)
	}
}

// Smoke test for UnregisterForNotifications
func TestComms_RegisterTrackedID(t *testing.T) {
	testId := id.NewIdFromString("test", id.Generic, t)
	clientId := id.NewIdFromString("client", id.Generic, t)

	// Start notification bot
	nbAddress := getNextAddress()
	nb := notificationBot.StartNotificationBot(testId, nbAddress,
		notificationBot.NewImplementation(), nil, nil)
	defer nb.Shutdown()

	// Create client's comms object
	c, err := NewClientComms(clientId, nil, nil, nil)
	if err != nil {
		t.Errorf("Can't create client comms: %+v", err)
	}
	manager := connect.NewManagerTesting(t)

	// Add notification bot to comm's manager
	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(testId, nbAddress, nil, params)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	// Unregister client with notification bot
	_, err = c.RegisterTrackedID(host, &pb.TrackedIntermediaryIDRequest{})
	if err != nil {
		t.Errorf("RegistrationMessage: Error received: %s", err)
	}
}

// Smoke test for UnregisterForNotifications
func TestComms_UnregisterToken(t *testing.T) {
	testId := id.NewIdFromString("test", id.Generic, t)
	clientId := id.NewIdFromString("client", id.Generic, t)

	// Start notification bot
	nbAddress := getNextAddress()
	nb := notificationBot.StartNotificationBot(testId, nbAddress,
		notificationBot.NewImplementation(), nil, nil)
	defer nb.Shutdown()

	// Create client's comms object
	c, err := NewClientComms(clientId, nil, nil, nil)
	if err != nil {
		t.Errorf("Can't create client comms: %+v", err)
	}
	manager := connect.NewManagerTesting(t)

	// Add notification bot to comm's manager
	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(testId, nbAddress, nil, params)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	// Unregister client with notification bot
	_, err = c.UnregisterToken(host, &pb.UnregisterTokenRequest{})
	if err != nil {
		t.Errorf("RegistrationMessage: Error received: %s", err)
	}
}

// Smoke test for UnregisterForNotifications
func TestComms_UnregisterTrackedID(t *testing.T) {
	testId := id.NewIdFromString("test", id.Generic, t)
	clientId := id.NewIdFromString("client", id.Generic, t)

	// Start notification bot
	nbAddress := getNextAddress()
	nb := notificationBot.StartNotificationBot(testId, nbAddress,
		notificationBot.NewImplementation(), nil, nil)
	defer nb.Shutdown()

	// Create client's comms object
	c, err := NewClientComms(clientId, nil, nil, nil)
	if err != nil {
		t.Errorf("Can't create client comms: %+v", err)
	}
	manager := connect.NewManagerTesting(t)

	// Add notification bot to comm's manager
	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(testId, nbAddress, nil, params)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	// Unregister client with notification bot
	_, err = c.UnregisterTrackedID(host, &pb.TrackedIntermediaryIDRequest{})
	if err != nil {
		t.Errorf("RegistrationMessage: Error received: %s", err)
	}
}

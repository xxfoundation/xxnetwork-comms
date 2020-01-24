////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package notificationBot

import (
	"context"
	"gitlab.com/elixxir/comms/connect"
	"gitlab.com/elixxir/comms/gateway"
	"gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/testkeys"
	"testing"
)

// Happy path
func TestRegisterForNotifications(t *testing.T) {
	// Get keys and certs
	keyPath := testkeys.GetNodeKeyPath()
	keyData := testkeys.LoadFromPath(keyPath)
	certPath := testkeys.GetNodeCertPath()
	certData := testkeys.LoadFromPath(certPath)
	// Get ID
	testId := "test"
	// Get available port
	notificationBotAddress := getNextBotAddress()

	//Init Notification bot
	notificationBot := StartNotificationBot(testId, notificationBotAddress,
		NewImplementation(), certData, keyData)
	defer notificationBot.Shutdown()
	//Init Gateway
	gw := gateway.StartGateway(testId, getNextBotAddress(), gateway.NewImplementation(), nil, nil)
	defer gw.Shutdown()

	ctx, _ := context.WithCancel(context.Background())

	//Init host and manager
	var manager connect.Manager
	host, err := manager.AddHost(testId, notificationBotAddress,
		certData, false, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	// Create message and pack it
	msg := &mixmessages.NotificationToken{}
	authMsg, err := notificationBot.PackAuthenticatedMessage(msg, host, false)
	if err != nil {
		t.Errorf("Failed to pack authenticated message: %+v", err)
	}

	// Run comm
	_, err = notificationBot.RegisterForNotifications(ctx, authMsg)
	if err != nil {
		t.Errorf("Failed to unregister: %+v", err)
	}
}

// Happy path
func TestUnRegisterForNotifications(t *testing.T) {
	// Get keys and certs
	keyPath := testkeys.GetNodeKeyPath()
	keyData := testkeys.LoadFromPath(keyPath)
	certPath := testkeys.GetNodeCertPath()
	certData := testkeys.LoadFromPath(certPath)
	// Get Id
	testId := "test"
	// Get available port
	notificationBotAddress := getNextBotAddress()

	//Init Notification bot
	notificationBot := StartNotificationBot(testId, notificationBotAddress,
		NewImplementation(), certData, keyData)
	defer notificationBot.Shutdown()
	ctx, _ := context.WithCancel(context.Background())

	//Init host and manager
	var manager connect.Manager
	host, err := manager.AddHost(testId, notificationBotAddress,
		certData, false, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	// Create message and pack it
	msg := &mixmessages.NotificationToken{}
	authMsg, err := notificationBot.PackAuthenticatedMessage(msg, host, false)
	if err != nil {
		t.Errorf("Failed to pack authenticated message: %+v", err)
	}

	// Run comm
	_, err = notificationBot.UnregisterForNotifications(ctx, authMsg)
	if err != nil {
		t.Errorf("Failed to unregister: %+v", err)
	}
}

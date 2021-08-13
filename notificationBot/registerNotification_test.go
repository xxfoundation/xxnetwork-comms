///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package notificationBot

import (
	"git.xx.network/elixxir/comms/gateway"
	"git.xx.network/elixxir/comms/mixmessages"
	"git.xx.network/elixxir/comms/testkeys"
	"git.xx.network/elixxir/comms/testutils"
	"git.xx.network/xx_network/comms/connect"
	"git.xx.network/xx_network/comms/gossip"
	"git.xx.network/xx_network/primitives/id"
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
	testId := id.NewIdFromString("test", id.Generic, t)
	// Get available port
	notificationBotAddress := getNextAddress()

	//Init Notification bot
	notificationBot := StartNotificationBot(testId, notificationBotAddress,
		NewImplementation(), certData, keyData)
	defer notificationBot.Shutdown()
	//Init Gateway
	gw := gateway.StartGateway(testId, getNextAddress(),
		gateway.NewImplementation(), nil, nil,
		gossip.DefaultManagerFlags())
	defer gw.Shutdown()

	ctx, cancel := testutils.NewContextTesting(t)
	defer cancel()
	defer ctx.Done()

	//Init host and manager
	manager := connect.NewManagerTesting(t)
	_, err := manager.AddHost(testId, notificationBotAddress,
		certData, connect.GetDefaultHostParams())
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	// Create message and pack it
	msg := &mixmessages.NotificationRegisterRequest{}

	// Run comm
	_, err = notificationBot.RegisterForNotifications(ctx, msg)
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
	testId := id.NewIdFromString("test", id.Generic, t)
	// Get available port
	notificationBotAddress := getNextAddress()

	//Init Notification bot
	notificationBot := StartNotificationBot(testId, notificationBotAddress,
		NewImplementation(), certData, keyData)
	defer notificationBot.Shutdown()
	ctx, cancel := testutils.NewContextTesting(t)
	defer cancel()
	defer ctx.Done()

	//Init host and manager
	manager := connect.NewManagerTesting(t)
	_, err := manager.AddHost(testId, notificationBotAddress,
		certData, connect.GetDefaultHostParams())
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	// Create message and pack it
	msg := &mixmessages.NotificationUnregisterRequest{}

	// Run comm
	_, err = notificationBot.UnregisterForNotifications(ctx, msg)
	if err != nil {
		t.Errorf("Failed to unregister: %+v", err)
	}
}

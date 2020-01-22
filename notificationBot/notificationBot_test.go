////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package notificationBot

import (
	"fmt"
	"gitlab.com/elixxir/comms/connect"
	"gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/registration"
	"gitlab.com/elixxir/comms/testkeys"
	"sync"
	"testing"
)

var botPortLock sync.Mutex
var botPort = 5500

func getNextBotAddress() string {
	botPortLock.Lock()
	defer func() {
		botPort++
		botPortLock.Unlock()
	}()
	return fmt.Sprintf("0.0.0.0:%d", botPort)
}

var registrationAddressLock sync.Mutex
var registrationAddress = 1600

func getRegistrationAddress() string {
	registrationAddressLock.Lock()
	defer func() {
		registrationAddress++
		registrationAddressLock.Unlock()
	}()
	return fmt.Sprintf("0.0.0.0:%d", registrationAddress)
}

// Tests whether the notifcationBot can be connected to and run an RPC with TLS enabled
func TestTLS(t *testing.T) {
	// Pull certs & keys
	keyPath := testkeys.GetNodeKeyPath()
	keyData := testkeys.LoadFromPath(keyPath)
	certPath := testkeys.GetNodeCertPath()
	certData := testkeys.LoadFromPath(certPath)
	testId := "test"

	// Start up a registration server
	regAddress := getRegistrationAddress()
	gw := registration.StartRegistrationServer(testId, regAddress, registration.NewImplementation(),
		certData, keyData)
	defer gw.Shutdown()

	// Start up the notification bot
	notificationBotAddress := getNextBotAddress()
	notificationBot := StartNotificationBot(testId, notificationBotAddress, NewImplementation(),
		certData, keyData)
	defer notificationBot.Shutdown()
	var manager connect.Manager

	// Add the host object to the manager
	host, err := manager.AddHost(testId, regAddress, certData, false, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	// Attempt to poll NDF
	_, err = notificationBot.RequestNdf(host, &mixmessages.NDFHash{})
	if err != nil {
		t.Error(err)
	}

	notificationBot.UnregisterForNotifications()
}

// Error path: Start bot with bad certs
func TestBadCerts(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	Address := getNextBotAddress()

	// This should panic and cause the defer func above to run
	_ = StartNotificationBot("test", Address, NewImplementation(),
		[]byte("bad cert"), []byte("bad key"))
}

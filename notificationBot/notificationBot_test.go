///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package notificationBot

import (
	"fmt"
	"gitlab.com/elixxir/comms/gateway"
	"gitlab.com/elixxir/comms/testkeys"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/gossip"
	"gitlab.com/xx_network/primitives/id"
	"sync"
	"testing"
)

var botPortLock sync.Mutex
var botPort = 1500

// Helper function to prevent port collisions
func getNextAddress() string {
	botPortLock.Lock()
	defer func() {
		botPort++
		botPortLock.Unlock()
	}()
	return fmt.Sprintf("0.0.0.0:%d", botPort)
}

// Tests whether the notifcationBot can be connected to and run an RPC with TLS enabled
func TestTLS(t *testing.T) {
	// Pull certs & keys
	keyPath := testkeys.GetNodeKeyPath()
	keyData := testkeys.LoadFromPath(keyPath)
	certPath := testkeys.GetNodeCertPath()
	certData := testkeys.LoadFromPath(certPath)
	testId := id.NewIdFromString("test", id.Generic, t)

	// Start up a registration server
	regAddress := getNextAddress()
	gw := gateway.StartGateway(testId, regAddress, gateway.NewImplementation(),
		certData, keyData, gossip.DefaultManagerFlags())
	defer gw.Shutdown()

	// Start up the notification bot
	notificationBotAddress := getNextAddress()
	notificationBot := StartNotificationBot(testId, notificationBotAddress, NewImplementation(),
		certData, keyData)
	defer notificationBot.Shutdown()
	manager := connect.NewManagerTesting(t)

	// Add the host object to the manager
	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(testId, regAddress, certData, params)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	// Attempt to poll NDF
	_, err = notificationBot.RequestNotifications(host)
	if err != nil {
		t.Error(err)
	}

}

// Error path: Start bot with bad certs
func TestBadCerts(t *testing.T) {
	testID := id.NewIdFromString("test", id.Generic, t)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	Address := getNextAddress()

	// This should panic and cause the defer func above to run
	_ = StartNotificationBot(testID, Address, NewImplementation(),
		[]byte("bad cert"), []byte("bad key"))
}

func TestComms_RequestNotifications(t *testing.T) {
	GatewayAddress := getNextAddress()
	nbAddress := getNextAddress()
	testID := id.NewIdFromString("test", id.Generic, t)

	gw := gateway.StartGateway(testID, GatewayAddress, gateway.NewImplementation(), nil,
		nil, gossip.DefaultManagerFlags())
	notificationBot := StartNotificationBot(testID, nbAddress, NewImplementation(),
		nil, nil)
	defer gw.Shutdown()
	defer notificationBot.Shutdown()
	manager := connect.NewManagerTesting(t)

	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(testID, GatewayAddress, nil, params)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	_, err = notificationBot.RequestNotifications(host)
	if err != nil {
		t.Errorf("SendGetSignedCertMessage: Error received: %s", err)
	}

}

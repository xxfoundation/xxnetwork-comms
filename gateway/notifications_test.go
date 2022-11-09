////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package gateway

import (
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/notificationBot"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/gossip"
	"gitlab.com/xx_network/primitives/id"
	"testing"
	"time"
)

// Happy path.
func TestComms_SendNotificationBatch(t *testing.T) {
	jww.SetLogThreshold(jww.LevelTrace)
	jww.SetStdoutThreshold(jww.LevelTrace)

	// Set up gateway
	gwAddr := getNextGatewayAddress()
	gwID := id.NewIdFromString("TestGatewayID", id.Gateway, t)
	gateway := StartGateway(gwID, gwAddr, NewImplementation(), nil, nil, gossip.DefaultManagerFlags())
	defer gateway.Shutdown()

	// Set up notification bot
	nbAddr := getNextServerAddress()
	nbID := &id.NotificationBot
	impl := notificationBot.NewImplementation()
	receiveChan := make(chan *pb.NotificationBatch)
	impl.Functions.ReceiveNotificationBatch = func(notifBatch *pb.NotificationBatch, auth *connect.Auth) error {
		go func() { receiveChan <- notifBatch }()
		return nil
	}
	notifications := notificationBot.StartNotificationBot(nbID, nbAddr, impl, nil, nil)
	defer notifications.Shutdown()

	// Create manager and add notification bot as host
	manager := connect.NewManagerTesting(t)
	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(nbID, nbAddr, nil, params)
	if err != nil {
		t.Errorf("Failed to add host: %+v", err)
	}

	// Generate message to send
	notifBatch := &pb.NotificationBatch{
		RoundID: 42,
		Notifications: []*pb.NotificationData{
			{
				EphemeralID: 42,
				IdentityFP:  []byte("IdentityFP"),
				MessageHash: []byte("MessageHash"),
			},
		},
		// XXX_sizecache: 31,
	}

	// Send NotificationBatch to notification bot
	err = gateway.SendNotificationBatch(host, notifBatch)
	if err != nil {
		t.Errorf("SendNotificationBatch() returned an error: %+v", err)
	}

	select {
	case result := <-receiveChan:
		if notifBatch.String() != result.String() {
			t.Errorf("Failed to receive the expected NotificationData."+
				"\nexpected: %s\nreceived: %s", notifBatch, result)
		}
	case <-time.NewTimer(50 * time.Millisecond).C:
		t.Error("Timed out while waiting to receive the NotificationData.")
	}
}

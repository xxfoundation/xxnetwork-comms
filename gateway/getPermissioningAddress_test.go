////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package gateway

import (
	"gitlab.com/elixxir/comms/node"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/gossip"
	"gitlab.com/xx_network/primitives/id"
	"testing"
)

// Happy path.
func TestComms_SendGetPermissioningAddress(t *testing.T) {
	/*	jww.SetLogThreshold(jww.LevelTrace)
		jww.SetStdoutThreshold(jww.LevelTrace)*/
	GatewayAddress := getNextGatewayAddress()
	ServerAddress := getNextServerAddress()
	testID := id.NewIdFromString("test", id.Generic, t)
	gateway := StartGateway(testID, GatewayAddress, NewImplementation(), nil,
		nil, gossip.DefaultManagerFlags())
	server := node.StartNode(testID, ServerAddress, 0, node.NewImplementation(),
		nil, nil)
	defer gateway.Shutdown()
	defer server.Shutdown()
	manager := connect.NewManagerTesting(t)

	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(testID, ServerAddress, nil, params)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	_, err = gateway.SendGetPermissioningAddress(host)
	if err != nil {
		t.Errorf("SendGetPermissioningAddress returned an error: %+v", err)
	}
}

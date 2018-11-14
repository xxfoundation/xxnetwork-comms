////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package gateway

import (
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/node"
	"testing"
)

// Smoke test SendCheckMessages
func TestSendBatch(t *testing.T) {
	gwShutDown := StartGateway(GatewayAddress, NewImplementation())
	nodeShutDown := node.StartServer(ServerAddress, node.NewImplementation())
	defer gwShutDown()
	defer nodeShutDown()

	msgs := []*pb.CmixMessage{{}}
	err := SendBatch(ServerAddress, msgs)
	if err != nil {
		t.Errorf("SendBatch: Error received: %s", err)
	}
}

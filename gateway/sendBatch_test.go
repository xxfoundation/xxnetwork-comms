////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package gateway

import (
	pb "gitlab.com/privategrity/comms/mixmessages"
	"gitlab.com/privategrity/comms/node"
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

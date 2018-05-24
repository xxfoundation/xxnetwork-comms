////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package gateway

import (
	pb "gitlab.com/privategrity/comms/mixmessages"
	"testing"
)

// Smoke test SendCheckMessages
func TestSendBatch(t *testing.T) {
	msgs := []*pb.CmixMessage{{}}
	err := SendBatch(ServerAddress, msgs)
	if err != nil {
		t.Errorf("SendBatch: Error received: %s", err)
	}
}

// Smoke test SendReceiveBatch
func TestSendReceiveBatch(t *testing.T) {
	x := make([]*pb.CmixMessage, 0)
	err := SendReceiveBatch(GatewayAddress, x)
	if err != nil {
		t.Errorf("PutMessage: Error received: %s", err)
	}
}

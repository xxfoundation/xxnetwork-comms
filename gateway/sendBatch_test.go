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
	msgs := []*pb.CmixMessage{&pb.CmixMessage{}}
	err := SendBatch(SERVER_ADDRESS, msgs)
	if err != nil {
		t.Errorf("SendBatch: Error received: %s", err)
	}
}

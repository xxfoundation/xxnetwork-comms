////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package node

import (
	pb "gitlab.com/elixxir/comms/mixmessages"
	"testing"
)

// Smoke test SendPhase
func TestSendRunPhase(t *testing.T) {
	ShutDown := StartServer(ServerAddress, NewImplementation(), "", "")
	defer ShutDown()
	_, err := SendPhase(ServerAddress, &pb.CmixMessage{})
	if err != nil {
		t.Errorf("RunPhase: Error received: %s", err)
	}
}

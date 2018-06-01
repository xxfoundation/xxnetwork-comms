////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package client

import (
	pb "gitlab.com/privategrity/comms/mixmessages"
	"gitlab.com/privategrity/comms/node"
	"testing"
)

// Smoke test RequestContactList
func TestRequestContactList(t *testing.T) {
	nodeShutDown := node.StartServer(ServerAddress, node.NewImplementation())
	defer nodeShutDown()

	_, err := RequestContactList(ServerAddress, &pb.ContactPoll{})
	if err != nil {
		t.Errorf("RequestContactList: Error received: %s", err)
	}
}

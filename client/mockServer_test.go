////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// This sets up a dummy/mock server instance for testing purposes
package client

import (
	"gitlab.com/privategrity/comms/node"
	"os"
	"testing"
)

const SERVER_ADDRESS = "localhost:5556"

// Start server for testing
func TestMain(m *testing.M) {
	go node.StartServer(SERVER_ADDRESS, node.TestInterface{})
	os.Exit(m.Run())
}

////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// This sets up a dummy/mock server instance for testing purposes
package node

import (
	"os"
	"testing"
)

const ServerAddress = "localhost:5556"

// Start server for testing
func TestMain(m *testing.M) {
	go StartServer(ServerAddress, TestInterface{})
	os.Exit(m.Run())
}

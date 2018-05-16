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

const SERVER_ADDRESS = "localhost:5555"

// Start server for testing
func TestMain(m *testing.M) {
	go StartServer(SERVER_ADDRESS, TestInterface{})
	os.Exit(m.Run())
}

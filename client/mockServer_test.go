////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains a dummy/mock server instance for testing purposes

package client

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"
)

var GatewayAddress = ""
var ServerAddress = ""

// Start server for testing
func TestMain(m *testing.M) {
	rand.Seed(time.Now().Unix())
	GatewayAddress = fmt.Sprintf("localhost:%d", rand.Intn(1000)+5001)
	ServerAddress = fmt.Sprintf("localhost:%d", rand.Intn(1000)+4000)
	os.Exit(m.Run())
}

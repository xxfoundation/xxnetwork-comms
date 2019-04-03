////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package gateway

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"
)

var GatewayAddress = ""
var ServerAddress = ""

// This sets up a dummy/mock gateway instance for testing purposes
func TestMain(m *testing.M) {
	rand.Seed(time.Now().Unix())
	GatewayAddress = fmt.Sprintf("localhost:%d", rand.Intn(3000)+4000)
	ServerAddress = fmt.Sprintf("localhost:%d", rand.Intn(3000)+3000)
	// If they're the same address, keep trying until they're different
	for ServerAddress == GatewayAddress {
		ServerAddress = fmt.Sprintf("localhost:%d", rand.Intn(3000)+3000)
	}
	os.Exit(m.Run())
}

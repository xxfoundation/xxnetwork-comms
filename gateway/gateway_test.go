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
)

var GatewayAddress = ""
var ServerAddress = ""

// This sets up a dummy/mock gateway instance for testing purposes
func TestMain(m *testing.M) {
	GatewayAddress = fmt.Sprintf("localhost:%d", (rand.Intn(1000) + 5001))
	ServerAddress = fmt.Sprintf("localhost:%d", (rand.Intn(1000) + 4000))
	os.Exit(m.Run())
}

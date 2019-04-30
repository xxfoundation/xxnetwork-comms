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
	GatewayAddress = fmt.Sprintf("0.0.0.0:%d", 6001)
	ServerAddress = fmt.Sprintf("0.0.0.0:%d", 5001)
	os.Exit(m.Run())
}

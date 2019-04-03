////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package node

import (
	"fmt"
	"gitlab.com/elixxir/comms/connect"
	"gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/testkeys"
	"math/rand"
	"os"
	"testing"
	"time"
)

var ServerAddress = ""

func TestMain(m *testing.M) {
	rand.Seed(time.Now().Unix())
	ServerAddress = fmt.Sprintf("localhost:%d", rand.Intn(2000)+4000)
	os.Exit(m.Run())
}

// Tests whether the server can be connected to and run an RPC with TLS enabled
func TestTLS(t *testing.T) {
	connect.ServerCertPath = testkeys.GetNodeCertPath()
	shutdown := StartServer(ServerAddress, NewImplementation(),
		testkeys.GetNodeCertPath(), testkeys.GetNodeKeyPath())
    defer shutdown()
    _, err := SendAskOnline(ServerAddress, &mixmessages.Ping{})
	if err != nil {
		t.Error(err)
	}
}

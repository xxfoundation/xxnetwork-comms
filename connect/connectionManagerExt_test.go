package connect_test

import (
	"fmt"
	"gitlab.com/elixxir/comms/connect"
	"gitlab.com/elixxir/comms/node"
	"gitlab.com/elixxir/comms/testkeys"
	"testing"
)

// Ensures that ConnectionManager implements Stringer; otherwise, compilation
// wil fail
var _ fmt.Stringer = &connect.ConnectionManager{}

type MockID string

func (m MockID) String() string {
	return string(m)
}

// The String() method should never cause a nil dereference or panic
// even when the input isn't valid
// Putting this test in an exterior package connect_test makes it simpler to run
// test servers
func TestConnectionManager_String(t *testing.T) {
	keyPath := testkeys.GetNodeKeyPath()
	certPath := testkeys.GetNodeCertPath()

	server1 := node.StartNode(":5658", node.NewImplementation(), nil, nil)
	server2 := node.StartNode(":5659", node.NewImplementation(),
		testkeys.LoadFromPath(certPath), testkeys.LoadFromPath(keyPath))
	defer server1.Shutdown()
	defer server2.Shutdown()
	cm := &connect.ConnectionManager{}
	// A real connection will be printed correctly, though
	cm.ConnectToNode(MockID("credsNil"), ":5658", nil, false)
	t.Log(cm)
	cm.ConnectToNode(MockID("goodCreds"), ":5659",
		testkeys.LoadFromPath(certPath), false)
	t.Log(cm)
}

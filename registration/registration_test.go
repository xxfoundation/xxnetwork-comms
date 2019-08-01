package registration

import (
	"fmt"
	"gitlab.com/elixxir/comms/client"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/testkeys"
	"sync"
	"testing"
)

var serverPortLock sync.Mutex
var serverPort = 5900

func getNextServerAddress() string {
	serverPortLock.Lock()
	defer func() {
		serverPort++
		serverPortLock.Unlock()
	}()
	return fmt.Sprintf("localhost:%d", serverPort)
}

type MockID string

func (m MockID) String() string {
	return string(m)
}

// Tests whether the server can be connected to and run an RPC with TLS enabled
func TestTLS(t *testing.T) {
	RegAddress := getNextServerAddress()

	keyPath := testkeys.GetNodeKeyPath()
	keyData := testkeys.LoadFromPath(keyPath)
	certPath := testkeys.GetNodeCertPath()
	certData := testkeys.LoadFromPath(certPath)

	rg := StartRegistrationServer(RegAddress,
		NewImplementation(),
		certData, keyData)
	// Well, client shouldn't have a server type because it's not a server
	// It's a client
	// So, we need some way to add a connection to the manager for the client
	defer rg.Shutdown()
	var c client.ClientComms
	connID := MockID("clientToRegistration")
	_ = c.ConnectToRegistration(connID,
		RegAddress, certData)

	_, err := c.SendRegistrationMessage(connID, &pb.UserRegistration{})
	if err != nil {
		t.Errorf("RegistrationMessage: Error received: %s", err)
	}
}

func TestBadCerts(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	RegAddress := getNextServerAddress()

	_ = StartRegistrationServer(RegAddress, NewImplementation(),
		[]byte("bad cert"), []byte("bad key"))
}

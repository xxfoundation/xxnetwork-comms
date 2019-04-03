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
var serverPort = 5800

func getNextServerAddress() string {
	serverPortLock.Lock()
	defer func() {
		serverPort++
		serverPortLock.Unlock()
	}()
	return fmt.Sprintf("localhost:%d", serverPort)
}

// Tests whether the server can be connected to and run an RPC with TLS enabled
func TestTLS(t *testing.T) {
	RegAddress := getNextServerAddress()
	rgShutDown := StartRegistrationServer(RegAddress,
		NewImplementation(),
		testkeys.GetNodeCertPath(),
		testkeys.GetNodeKeyPath())
	defer rgShutDown()

	// Note: This line takes forever when using the gateway certs/paths both
	// for the reg server and the reg client, but succeeds when using the node
	// certs/paths for both. I don't know why it happens, but it's a bit spooky.
	_, err := client.SendRegistrationMessage(RegAddress,
		testkeys.GetNodeCertPath(),"", &pb.RegisterUserMessage{})
	if err != nil {
		t.Errorf("RegistrationMessage: Error received: %s", err)
	}
}

package message

import (
	pb "gitlab.com/privategrity/comms/mixmessages"
	"gitlab.com/privategrity/comms/mixserver"
	"os"
	"testing"
)

const SERVER_ADDRESS = "localhost:5555"

// Start server for testing
func TestMain(m *testing.M) {
	go mixserver.StartServer(SERVER_ADDRESS, TestInterface{})
	os.Exit(m.Run())
}

// Blank struct implementing ServerHandler interface for testing purposes (Passing to StartServer)
type TestInterface struct{}

func (m TestInterface) PrecompDecrypt(message *pb.PrecompDecryptMessage) {}

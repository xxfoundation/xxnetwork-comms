package message

import (
	pb "gitlab.com/privategrity/comms/mixmessages"
	"gitlab.com/privategrity/comms/mixserver"
	"os"
	"testing"
)

// Start server for testing
func TestMain(m *testing.M) {
	addr := "localhost:5555"
	go mixserver.StartServer(addr, TestInterface{})
	os.Exit(m.Run())
}

// Blank interface for testing purposes (Passing to StartServer)
type TestInterface struct{}

func (m TestInterface) PrecompDecrypt(message *pb.PrecompDecryptMessage) {}

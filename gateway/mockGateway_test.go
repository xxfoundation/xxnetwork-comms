////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package gateway

import (
	pb "gitlab.com/privategrity/comms/mixmessages"
	"gitlab.com/privategrity/comms/node"
	"os"
	"testing"
)

const GatewayAddress = "localhost:5566"
const ServerAddress = "localhost:5569"

// This sets up a dummy/mock gateway instance for testing purposes
func TestMain(m *testing.M) {
	go StartGateway(GatewayAddress, TestInterface{})
	go node.StartServer(ServerAddress, node.TestInterface{})
	os.Exit(m.Run())
}

// Blank struct implementing GatewayHandler interface for testing purposes
// (Passing to StartGateway)
type TestInterface struct{}

func (m TestInterface) GetMessage(userId uint64,
	msgId string) (*pb.CmixMessage, bool) {
	return &pb.CmixMessage{}, true
}

func (m TestInterface) CheckMessages(userId uint64, messageID string) (
	[]string, bool) {
	return make([]string, 0), true
}

func (m TestInterface) PutMessage(message *pb.CmixMessage) bool {
	return true
}

func (m TestInterface) ReceiveBatch(message *pb.OutputMessages) {
}

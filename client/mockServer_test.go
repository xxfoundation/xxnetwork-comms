////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// This sets up a dummy/mock server instance for testing purposes
package client

import (
	"gitlab.com/privategrity/comms/gateway"
	pb "gitlab.com/privategrity/comms/mixmessages"
	"gitlab.com/privategrity/comms/node"
	"os"
	"testing"
)

const ServerAddress = "localhost:5559"
const GatewayAddress = "localhost:5560"

// Start server for testing
func TestMain(m *testing.M) {
	go gateway.StartGateway(GatewayAddress, TestInterface{})
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

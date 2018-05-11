////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package gateway

import (
	pb "gitlab.com/privategrity/comms/mixmessages"
	"os"
	"testing"
)

const SERVER_ADDRESS = "localhost:5556"

// This sets up a dummy/mock gateway instance for testing purposes
func TestMain(m *testing.M) {
	go StartGateway(SERVER_ADDRESS, TestInterface{})
	os.Exit(m.Run())
}

// Blank struct implementing GatewayHandler interface for testing purposes
// (Passing to StartGateway)
type TestInterface struct{}

func (m TestInterface) GetMessage(userId uint64,
	msgId string) (*pb.CmixMessage, bool) {
	return &pb.CmixMessage{}, true
}

func (m TestInterface) CheckMessages(userId uint64) ([]string, bool) {
	return make([]string, 0), true
}

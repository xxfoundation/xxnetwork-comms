////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Dummy implementation (so you can use for tests)
package node

import (
	pb "gitlab.com/elixxir/comms/mixmessages"
)

// Blank struct implementing ServerHandler interface for testing purposes (Passing to StartServer)
type TestInterface struct{}

func (m TestInterface) NewRound(roundID string) {}

func (m TestInterface) RoundtripPing(message *pb.TimePing) {}

func (m TestInterface) GetServerMetrics(message *pb.ServerMetrics) {}

func (m TestInterface) PostPhase(message *pb.Batch) {}

func (m TestInterface) PostPrecompResult(roundID uint64,
	slots []*pb.Slot) error {
	return nil
}

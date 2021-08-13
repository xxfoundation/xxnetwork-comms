///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Dummy implementation (so you can use for tests)
package node

import (
	pb "git.xx.network/elixxir/comms/mixmessages"
)

// Blank struct implementing Handler interface for testing purposes (Passing to StartNode)
type TestInterface struct{}

func (m TestInterface) NewRound(roundID string) {}

func (m TestInterface) PostPhase(message *pb.Batch) {}

func (m TestInterface) StreamPostPhase(message pb.Node_StreamPostPhaseServer) error {
	return nil
}

func (m TestInterface) PostPrecompResult(roundID uint64,
	slots []*pb.Slot) error {
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Dummy implementation (so you can use for tests)
package node

import (
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/connect"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	jww.SetStdoutThreshold(jww.LevelTrace)
	connect.TestingOnlyDisableTLS = true
	os.Exit(m.Run())
}

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

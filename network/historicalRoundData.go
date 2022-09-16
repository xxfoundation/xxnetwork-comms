////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package network

import (
	"errors"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/primitives/id"
)

// Calls the underlying interface's function to get a specific round from history
func (i *Instance) GetHistoricalRound(id id.Round) (*pb.RoundInfo, error) {
	if i.ers != nil {
		return i.ers.Retrieve(id)
	}
	return nil, errors.New("no ExternalRoundStorage object was defined on instance creation")
}

// Calls the underlying interface's function to get specific rounds from history
func (i *Instance) GetHistoricalRounds(rounds []id.Round) ([]*pb.RoundInfo, error) {
	if i.ers != nil {
		return i.ers.RetrieveMany(rounds)
	}
	return nil, errors.New("no ExternalRoundStorage object was defined on instance creation")
}

// Calls the underlying interface's function to get a range of rounds from history
func (i *Instance) GetHistoricalRoundRange(first, last id.Round) ([]*pb.RoundInfo, error) {
	if i.ers != nil {
		return i.ers.RetrieveRange(first, last)
	}
	return nil, errors.New("no ExternalRoundStorage object was defined on instance creation")
}

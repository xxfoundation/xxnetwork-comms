///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package network

import (
	"github.com/pkg/errors"
	pb "gitlab.com/elixxir/comms/mixmessages"
	ds "gitlab.com/elixxir/comms/network/dataStructures"
	"gitlab.com/elixxir/primitives/states"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/primitives/id"
)

type FilteredUpdates struct {
	updates *ds.Updates
	comm    *connect.ProtoComms
}

func NewFilteredUpdates(comm *connect.ProtoComms) *FilteredUpdates {
	return &FilteredUpdates{
		updates: ds.NewUpdates(),
		comm:    comm,
	}
}

// Get an update ID
func (fu *FilteredUpdates) GetRoundUpdate(updateID int) (*pb.RoundInfo, error) {
	return fu.updates.GetUpdate(updateID)
}

// Get updates from a given round
func (fu *FilteredUpdates) GetRoundUpdates(id int) []*pb.RoundInfo {
	return fu.updates.GetUpdates(id)
}

// get the most recent update id
func (fu *FilteredUpdates) GetLastUpdateID() int {
	return fu.updates.GetLastUpdateID()
}

// Pluralized version of RoundUpdate
func (fu *FilteredUpdates) RoundUpdates(rounds []*pb.RoundInfo) error {
	// Process all rounds passed in
	for _, round := range rounds {
		err := fu.RoundUpdate(round)
		if err != nil {
			return err
		}
	}
	return nil
}

// Add a round to the updates filter
func (fu *FilteredUpdates) RoundUpdate(info *pb.RoundInfo) error {
	switch states.Round(info.State) {
	// Only add to filter states client cares about
	case states.COMPLETED, states.FAILED, states.QUEUED:
		perm, success := fu.comm.GetHost(&id.Permissioning)
		if !success {
			return errors.New("Could not get permissioning Public Key" +
				"for round info verification")
		}

		rnd := ds.NewRound(info, perm.GetPubKey())

		err := fu.updates.AddRound(rnd)
		if err != nil {
			return err
		}
	default:

	}

	return nil
}

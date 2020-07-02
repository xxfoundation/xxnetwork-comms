///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Stores a list of all updates in order of update id

package dataStructures

import (
	"github.com/pkg/errors"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/primitives/ring"
)

const RoundUpdatesBufLen = 1500

// Standard ring buffer, but objects come with numbering
type Updates struct {
	updates *ring.Buff
}

// Create a new Updates object
func NewUpdates() *Updates {
	// we want each updateId stored in this structure
	return &Updates{
		updates: ring.NewBuff(RoundUpdatesBufLen),
	}
}

// Add a round to the ring buffer
func (u *Updates) AddRound(info *pb.RoundInfo) error {
	return u.updates.UpsertById(int(info.ID),info)
}

// Get a given update ID from the ring buffer
func (u *Updates) GetUpdate(id int) (*pb.RoundInfo, error) {

	val, err := u.updates.GetById(id)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get update with id %d", id)
	}
	return val.(*pb.RoundInfo), nil
}

//gets all updates after a given ID
func (u *Updates) GetUpdates(id int) []*pb.RoundInfo {
	interfaceList, err := u.updates.GetNewerById(id)

	if err != nil {
		return make([]*pb.RoundInfo, 0)
	}

	infoList := make([]*pb.RoundInfo, len(interfaceList))

	for i, face := range interfaceList {
		infoList[i] = face.(*pb.RoundInfo)
	}

	return infoList
}

// Get the id of the newest update in the buffer
func (u *Updates) GetLastUpdateID() int {
	return u.updates.GetNewestId()
}

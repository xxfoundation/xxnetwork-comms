////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Stores a list of all updates in order of update id

package dataStructures

import (
	"github.com/pkg/errors"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/primitives/ring"
)

const RoundInfoBufLen = 1000
const RoundUpdatesBufLen = 10000

// Standard ring buffer, but objects come with numbering
type Updates struct {
	updates *ring.Buff
}

// Create a new Updates object
func NewUpdates() *Updates {
	// we want each updateId stored in this structure
	idFunc := func(val interface{}) int {
		if val == nil {
			return -1
		}
		return int(val.(*pb.RoundInfo).UpdateID)
	}
	return &Updates{
		updates: ring.NewBuff(RoundUpdatesBufLen, idFunc),
	}
}

// Add a round to the ring buffer
func (u *Updates) AddRound(info *pb.RoundInfo) error {

	// comparison should ensure that updates are not overwritten in the event of a duplicate
	comp := func(current interface{}, new interface{}) bool {
		if current == nil {
			return true
		}
		return false
	}

	return u.updates.UpsertById(info, comp)
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
func (u *Updates) GetUpdates(id int) ([]*pb.RoundInfo, error) {
	interfaceList, err := u.updates.GetNewerById(id)

	if err != nil {
		return nil, err
	}

	infoList := make([]*pb.RoundInfo, len(interfaceList))

	for i, face := range interfaceList {
		infoList[i] = face.(*pb.RoundInfo)
	}

	return infoList, nil
}

// Get the id of the newest update in the buffer
func (u *Updates) GetLastUpdateID() int {
	return u.updates.GetNewestId()
}

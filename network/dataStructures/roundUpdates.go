package dataStructures

import (
	"github.com/pkg/errors"
	"gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/primitives/ring"
)

const RoundInfoBufLen = 1000
const RoundUpdatesBufLen = 10000

// Standard ring buffer, but objects come with numbering
type Updates struct {
	updates *ring.Buff
}

// Add a round to the ring buffer
func (u *Updates) AddRound(info *mixmessages.RoundInfo) error {
	// comparison should ensure that updates are not overwritten in the event of a duplicate
	comp := func(current interface{}, new interface{}) bool {
		if current == nil {
			return true
		}
		return false
	}

	// we want each updateId stored in this structure
	id := func(val interface{}) int {
		if val == nil {
			return -1
		}
		return int(val.(*mixmessages.RoundInfo).UpdateID)
	}

	if u.updates == nil {
		u.updates = ring.NewBuff(RoundUpdatesBufLen, id)
	}

	return u.updates.UpsertById(info, comp)
}

// Get a given update ID from the ring buffer
func (u *Updates) GetUpdate(id int) (*mixmessages.RoundInfo, error) {
	val, err := u.updates.GetById(id)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get update with id %d", id)
	}
	return val.(*mixmessages.RoundInfo), nil
}



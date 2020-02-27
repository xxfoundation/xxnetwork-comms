package consensus

import (
	"github.com/pkg/errors"
	"gitlab.com/elixxir/comms/mixmessages"
	ds "gitlab.com/elixxir/primitives/dataStructures"
)

const RoundInfoBufLen = 1000
const RoundUpdatesBufLen = 10000

// Standard ring buffer, but objects come with numbering
type Updates struct {
	updates *ds.RingBuff
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
		u.updates = ds.NewRingBuff(RoundUpdatesBufLen, id)
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

// ID numbers can overwrite
type Data struct {
	rounds *ds.RingBuff
}

// Upsert a round into the ring bugger
func (d *Data) UpsertRound(r *mixmessages.RoundInfo) error {
	// comparison here should ensure that either the current round is nil or has a lower update id than the new round
	comp := func(current interface{}, new interface{}) bool {
		if current == nil {
			return true
		}
		if current.(*mixmessages.RoundInfo).ID < new.(*mixmessages.RoundInfo).ID {
			return true
		}
		return false
	}

	// We want data using the round ID as its primary
	id := func(val interface{}) int {
		if val == nil {
			return -1
		}
		return int(val.(*mixmessages.RoundInfo).ID)
	}

	if d.rounds == nil {
		d.rounds = ds.NewRingBuff(RoundInfoBufLen, id)
	}

	//find the round location
	//check the new state is newer then the current
	//replace the round info object
	return d.rounds.UpsertById(r, comp)
}

// Get a given round id from the ring buffer
func (d *Data) GetRound(id int) (*mixmessages.RoundInfo, error) {
	val, err := d.rounds.GetById(id)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get update with id %d", id)
	}
	return val.(*mixmessages.RoundInfo), nil
}

package consensus

import (
	"gitlab.com/elixxir/comms/mixmessages"
	ds "gitlab.com/elixxir/primitives/dataStructures"
)

const RoundInfoBufLen = 1000
const RoundUpdatesBufLen = 10000

// Standard ring buffer, but objects come with numbering
type Updates struct {
	updates *ds.RingBuff
}

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
		return int(val.(*mixmessages.RoundInfo).UpdateID)
	}

	if u.updates == nil {
		u.updates = ds.NewRingBuff(RoundUpdatesBufLen, id)
	}

	return u.updates.UpsertById(info, comp)
}

// ID numbers can overwrite
type Data struct {
	rounds *ds.RingBuff
}

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

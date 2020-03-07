////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Store a list of the most recent rounds, holding the most recent update for each

package dataStructures

import (
	"github.com/pkg/errors"
	"gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/primitives/id"
	"gitlab.com/elixxir/primitives/ring"
)

// ID numbers can overwrite
type Data struct {
	rounds *ring.Buff
}

func NewData() *Data {
	// We want data using the round ID as its primary
	idFunc := func(val interface{}) int {
		if val == nil {
			return -1
		}
		return int(val.(*mixmessages.RoundInfo).ID)
	}

	return &Data{
		rounds: ring.NewBuff(RoundInfoBufLen, idFunc),
	}
}

// Upsert a round into the ring bugger
func (d *Data) UpsertRound(r *mixmessages.RoundInfo) error {
	// comparison here should ensure that either the current round is nil or has a lower update id than the new round
	comp := func(current interface{}, new interface{}) bool {
		if current == nil {
			return true
		}
		if current.(*mixmessages.RoundInfo).UpdateID < new.(*mixmessages.RoundInfo).UpdateID {
			return true
		}
		return false
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

// Get the ID of the newest round in the buffer
func (d *Data) GetLastRoundID() id.Round {
	return id.Round(d.rounds.GetNewestId())
}

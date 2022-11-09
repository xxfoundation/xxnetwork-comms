////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Store a list of the most recent rounds, holding the most recent update for each

package dataStructures

import (
	"github.com/pkg/errors"
	"gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/primitives/id"
	"gitlab.com/xx_network/ring"
)

const RoundInfoBufLen = 1500

// ID numbers can overwrite
type Data struct {
	rounds *ring.Buff
}

// Initialize a new Data object
func NewData() *Data {
	// We want data using the round ID as its primary

	return &Data{
		rounds: ring.NewBuff(RoundInfoBufLen),
	}
}

// Upsert a round into the ring bugger
func (d *Data) UpsertRound(r *Round) error {
	//find the round location
	//check the new state is newer then the current
	//replace the round info object
	return d.rounds.UpsertById(int(r.info.ID), r)
}

// Get a given round id from the ring buffer as a roundInfo
func (d *Data) GetRound(id int) (*mixmessages.RoundInfo, error) {
	val, err := d.rounds.GetById(id)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get round by id with "+
			"%d", id)
	}

	if val == nil {
		return nil, errors.Errorf("Failed to get round by id with %d, "+
			"got nil round", id)
	}

	return val.(*Round).Get(), nil
}

// Get a given round id from the ring buffer as a round object
func (d *Data) GetWrappedRound(id int) (*Round, error) {
	val, err := d.rounds.GetById(id)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get update with id %d", id)
	}
	var rtn *Round
	if val != nil {
		rtn = val.(*Round)
	}
	return rtn, nil
}

// Get the ID of the newest round in the buffer
func (d *Data) GetLastRoundID() id.Round {
	return id.Round(d.rounds.GetNewestId())
}

// Gets the ID of the oldest roundd in the buffer
func (d *Data) GetOldestRoundID() id.Round {
	return id.Round(d.rounds.GetOldestId())
}

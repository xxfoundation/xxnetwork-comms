///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////
package dataStructures

import (
	"container/list"
	"github.com/golang-collections/collections/set"
	"github.com/pkg/errors"
	pb "git.xx.network/elixxir/comms/mixmessages"
	"git.xx.network/elixxir/primitives/states"
	"git.xx.network/xx_network/primitives/netTime"
	"sync"
	"sync/atomic"
	"time"
)

var timeOutError = errors.New("Timed out getting round furthest in the future.")

// WaitingRounds contains a list of all queued rounds ordered by which occurs
// furthest in the future with the furthest in the the back.
type WaitingRounds struct {
	rounds *list.List
	c      *sync.Cond
}

// NewWaitingRounds generates a new WaitingRounds with an empty round list.
func NewWaitingRounds() *WaitingRounds {
	wr := WaitingRounds{
		rounds: list.New(),
	}

	m := sync.Mutex{}
	wr.c = sync.NewCond(&m)

	return &wr
}

// Len returns the number of rounds in the list.
func (wr *WaitingRounds) Len() int {
	return wr.rounds.Len()
}

// Insert inserts a queued round into the list in order of its timestamp, from
// smallest to greatest. If the new round is not in a QUEUED state, then it is
// not inserted. If the new round already exists in the list but is no longer
// queued, then it is removed.
func (wr *WaitingRounds) Insert(newRound *Round) {
	// If the round is queued, then add it to the list; otherwise, remove it
	if newRound.info.GetState() == uint32(states.QUEUED) {
		wr.c.L.Lock()
		inserted := false
		// Loop through every round, starting with the furthest in the future
		for e := wr.rounds.Back(); e != nil; e = e.Prev() {
			// If the new round is larger, than add it before
			extractedRound := e.Value.(*Round)
			if getTime(newRound) > getTime(extractedRound) {
				wr.rounds.InsertAfter(newRound, e)
				inserted = true
				break
			}
		}

		// If the round's realtime is the sooner than all other rounds, then add
		// it to the beginning  of the list
		if !inserted {
			wr.rounds.PushFront(newRound)
		}
		go func() {
			wr.c.Broadcast()
			wr.c.L.Unlock()
		}()
	} else {
		wr.c.L.Lock()
		wr.remove(newRound)
		wr.c.L.Unlock()
	}
}

// getTime returns the timestamp for the round's realtime.
func getTime(round *Round) uint64 {
	return round.info.Timestamps[states.QUEUED]
}

// remove deletes the round from the list if it exists.
func (wr *WaitingRounds) remove(newRound *Round) {
	// Look for a node with a matching ID from the list
	for e := wr.rounds.Front(); e != nil; e = e.Next() {
		extractedRound := e.Value.(*Round)
		if extractedRound.info.ID == newRound.info.ID {
			wr.rounds.Remove(e)
			return
		}
	}
}

// getFurthest returns the round that will occur furthest in the future. If the
// list is empty, then nil is returned. If the round is on the exclusion list,
// then the next round is checked.
// this is assumed to be called on an operation already under the cond's lock
func (wr *WaitingRounds) getFurthest(exclude *set.Set, cutoffDelta time.Duration) *Round {
	earliestStart := netTime.Now().Add(cutoffDelta)

	// Return nil for an empty list
	if wr.Len() == 0 {
		return nil
	}

	// Return the last non-excluded round in the list
	for e := wr.rounds.Back(); e != nil; e = e.Prev() {
		r := e.Value.(*Round)
		// Cannot guarantee that the round object's pointers will be exact match
		// of value in set
		RoundStartTime := time.Unix(0, int64(r.info.Timestamps[states.QUEUED]))
		if RoundStartTime.After(earliestStart) && !isExcluded(exclude, r.info) {
			return r
		}
	}

	// If all the rounds in the list are excluded, then return nil
	return nil
}

// getClosest returns the round that will occur soonest in the future. If the
// list is empty, then nil is returned. If the round is on the exclusion list,
// then the next round is checked.
// this is assumed to be called on an operation already under the cond's lock
func (wr *WaitingRounds) getClosest(exclude *set.Set, minRoundAge time.Duration) *Round {
	earliestStart := netTime.Now().Add(minRoundAge)

	// Return nil for an empty list
	if wr.Len() == 0 {
		return nil
	}

	// Return the first non-excluded round in the list
	for e := wr.rounds.Front(); e != nil; e = e.Next() {
		r := e.Value.(*Round)
		// Cannot guarantee that the round object's pointers will be exact match
		// of value in set
		RoundStartTime := time.Unix(0, int64(r.info.Timestamps[states.QUEUED]))
		if RoundStartTime.After(earliestStart) && !isExcluded(exclude, r.info) {
			return r
		}
	}

	// If all the rounds in the list are excluded, then return nil
	return nil
}

func isExcluded(exclude *set.Set, r *pb.RoundInfo) bool {
	if exclude == nil {
		return false
	}

	return exclude.Has(r)
}

// GetSlice returns a slice of all round infos in the list that have yet to
// occur.
func (wr *WaitingRounds) GetSlice() []*pb.RoundInfo {
	wr.c.L.Lock()
	defer wr.c.L.Unlock()

	now := uint64(netTime.Now().Nanosecond())
	var roundInfos []*pb.RoundInfo
	iter := 0
	for e, i := wr.rounds.Front(), 0; e != nil; e, i = e.Next(), i+1 {
		iter++
		extractedRound := e.Value.(*Round)
		if getTime(extractedRound) > now {
			roundInfos = append(roundInfos, extractedRound.info)
		}
	}

	// Return the last round in the list, which is the furthest in the future
	return roundInfos
}

// GetUpcomingRealtime returns the round that will occur furthest in the future.
// If the list is empty, then it waits waits for a round to be added for the
// specified duration. If no round is added, then an error is returned.
func (wr *WaitingRounds) GetUpcomingRealtime(timeout time.Duration,
	exclude *set.Set, minRoundAge time.Duration) (*pb.RoundInfo, error) {

	// Start timeout timer
	timer := time.NewTimer(timeout)

	exit := uint32(0)
	defer atomic.StoreUint32(&exit, 1)

	// Start waiting for rounds to be added
	sig := make(chan *pb.RoundInfo, 1)
	go func() {
		var round *Round
		wr.c.L.Lock()
		for atomic.LoadUint32(&exit) == 0 {
			round = wr.getClosest(exclude, minRoundAge)
			if round != nil {
				break
			}
			wr.c.Wait()
		}
		wr.c.L.Unlock()
		if round != nil {
			sig <- round.Get()
		}
	}()

	// If rounds already exist in the list, then return the the correct round
	// without waiting

	// If the list is empty, then start waiting for rounds to be added.
	for {
		select {
		case <-timer.C:
			return nil, timeOutError
		case round := <-sig:
			return round, nil
		}
	}
}

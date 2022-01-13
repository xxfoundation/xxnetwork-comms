///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////
package dataStructures

import (
	"container/list"
	"github.com/pkg/errors"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/primitives/excludedRounds"
	"gitlab.com/elixxir/primitives/states"
	"gitlab.com/xx_network/primitives/netTime"
	"sync"
	"sync/atomic"
	"time"
)

var timeOutError = errors.New("Timed out getting round furthest in the future.")

// maxGetClosestTries is the maximum amount of rounds pulled by
// WaitingRounds.GetUpcomingRealtime. Exceeding this amount causes
// WaitingRounds.GetUpcomingRealtime to switch from using
// WaitingRounds.getClosest to using WaitingRounds.getFurthest
const maxGetClosestTries = 2

// WaitingRounds contains a list of all queued rounds ordered by which occurs
// furthest in the future with the furthest in the the back.
type WaitingRounds struct {
	readRounds *atomic.Value
	writeRounds *list.List
	writeMux sync.Mutex
	waitingMux sync.RWMutex
}

// NewWaitingRounds generates a new WaitingRounds with an empty round list.
func NewWaitingRounds() *WaitingRounds {
	wr := WaitingRounds{
		writeRounds: list.New(),
		readRounds: &atomic.Value{},
	}

	return &wr
}

// Len returns the number of rounds in the list.
func (wr *WaitingRounds) Len() int {
	wr.writeMux.Lock()
	defer wr.writeMux.Unlock()
	return wr.writeRounds.Len()
}

// Insert inserts a queued round into the list in order of its timestamp, from
// smallest to greatest. If the new round is not in a QUEUED state, then it is
// not inserted. If the new round already exists in the list but is no longer
// queued, then it is removed.
func (wr *WaitingRounds) Insert(newRound *Round) {
	wr.writeMux.Lock()
	defer wr.writeMux.Unlock()
	edited := false
	// If the round is queued, then add it to the list; otherwise, remove it
	if newRound.info.GetState() == uint32(states.QUEUED) {
		edited= true
		inserted := false
		// Loop through every round, starting with the furthest in the future
		for e := wr.writeRounds.Back(); e != nil; e = e.Prev() {
			// If the new round is larger, than add it before
			extractedRound := e.Value.(*Round)
			if getTime(newRound) > getTime(extractedRound) {
				wr.writeRounds.InsertAfter(newRound, e)
				inserted = true
				break
			}
		}

		// If the round's realtime is the sooner than all other rounds, then add
		// it to the beginning  of the list
		if !inserted {
			wr.writeRounds.PushFront(newRound)
		}

	} else {
		edited=true
		wr.remove(newRound)
	}
	if edited{

		go func(){
			wr.waitingMux.Lock()
			defer wr.waitingMux.Unlock()
			wr.storeReadRounds()
		}()
	}
}

func(wr *WaitingRounds)storeReadRounds(){
	roundsList := make([]*Round,0,wr.writeRounds.Len())

	for e := wr.writeRounds.Front(); e != nil; e = e.Next() {
		roundsList = append(roundsList,e.Value.(*Round))
	}
	wr.readRounds.Store(roundsList)
}

// getTime returns the timestamp for the round's realtime.
func getTime(round *Round) uint64 {
	return round.info.Timestamps[states.QUEUED]
}

// remove deletes the round from the list if it exists.
func (wr *WaitingRounds) remove(newRound *Round) {
	// Look for a node with a matching ID from the list
	for e := wr.writeRounds.Front(); e != nil; e = e.Next() {
		extractedRound := e.Value.(*Round)
		if extractedRound.info.ID == newRound.info.ID {
			wr.writeRounds.Remove(e)
			return
		}
	}
}

// getFurthest returns the round that will occur furthest in the future. If the
// list is empty, then nil is returned. If the round is on the exclusion list,
// then the next round is checked.
// this is assumed to be called on an operation already under the cond's lock
func (wr *WaitingRounds) getFurthest(exclude excludedRounds.ExcludedRounds, cutoffDelta time.Duration) *Round {
	earliestStart := netTime.Now().Add(cutoffDelta)


	roundsList := wr.readRounds.Load().([]*Round)

	// Return the last non-excluded round in the list
	for i:=len(roundsList)-1;i>=0;i-- {
		r := roundsList[i]
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
func (wr *WaitingRounds) getClosest(exclude excludedRounds.ExcludedRounds, minRoundAge time.Duration) *Round {
	earliestStart := netTime.Now().Add(minRoundAge)

	roundsList := wr.readRounds.Load().([]*Round)

	// Return the first non-excluded round in the list
	for i:=0;i<len(roundsList);i++ {
		r := roundsList[i]
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

func isExcluded(exclude excludedRounds.ExcludedRounds, r *pb.RoundInfo) bool {
	if exclude == nil {
		return false
	}

	return exclude.Has(r.GetRoundId())
}

// GetSlice returns a slice of all round infos in the list that have yet to
// occur.
func (wr *WaitingRounds) GetSlice() []*pb.RoundInfo {

	roundsList := wr.readRounds.Load().([]*Round)

	now := uint64(netTime.Now().Nanosecond())
	var roundInfos []*pb.RoundInfo
	for i:=0;i<len(roundsList);i++ {
		if getTime(roundsList[i]) > now {
			roundInfos = append(roundInfos, roundsList[i].info)
		}
	}

	// Return the last round in the list, which is the furthest in the future
	return roundInfos
}

// GetUpcomingRealtime returns the round that will occur furthest in the future.
// If the list is empty, then it waits waits for a round to be added for the
// specified duration. If no round is added, then an error is returned.
//
// The length of the excluded set indicates how many times the client has
// called GetUpcomingRealtime trying to retrieve a round to send on.
// GetUpcomingRealtime defaults to retrieving the closest non-excluded round
// from WaitingRounds. If the length of the excluded set exceeds the maximum
// attempts at pulling the closest round, GetUpcomingRealtime will retrieve
// the furthest non-excluded round from WaitingRounds.
func (wr *WaitingRounds) GetUpcomingRealtime(timeout time.Duration,
	exclude excludedRounds.ExcludedRounds, minRoundAge time.Duration) (*pb.RoundInfo, error) {

	// Start timeout timer
	timer := time.NewTimer(timeout)

	exit := uint32(0)
	defer atomic.StoreUint32(&exit, 1)

	// Start waiting for rounds to be added
	sig := make(chan *pb.RoundInfo, 1)
	go func() {
		var round *Round
		for atomic.LoadUint32(&exit) == 0 {
			if exclude.Len() < maxGetClosestTries {
				// Use getClosest when excluded set's length is small
				round = wr.getClosest(exclude, minRoundAge)
				if round != nil {
					break
				}
			} else {
				// Use getFurthest when excluded set's length exceeds maxGetClosestTries
				round = wr.getFurthest(exclude, minRoundAge)
				if round != nil {
					break
				}
			}
			wr.waitingMux.RLock()
			wr.waitingMux.RUnlock()
		}
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

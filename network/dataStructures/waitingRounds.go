////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package dataStructures

import (
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/elliotchance/orderedmap"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/primitives/excludedRounds"
	"gitlab.com/xx_network/primitives/id"
	"gitlab.com/xx_network/primitives/netTime"
)

var timeOutError = errors.New("Timed out getting round furthest in the future.")

// maxGetClosestTries is the maximum amount of rounds pulled by
// WaitingRounds.GetUpcomingRealtime. Exceeding this amount causes
// WaitingRounds.GetUpcomingRealtime to switch from using
// WaitingRounds.getClosest to using WaitingRounds.getFurthest
const maxGetClosestTries = 2

// WaitingRounds contains a list of all queued rounds ordered by which occurs
// furthest in the future with the furthest in the back.
type WaitingRounds struct {
	readRounds  *atomic.Value
	writeRounds *orderedmap.OrderedMap
	mux         sync.Mutex
	signal      chan struct{}
}

// NewWaitingRounds generates a new WaitingRounds with an empty round list.
func NewWaitingRounds() *WaitingRounds {
	wr := WaitingRounds{
		writeRounds: orderedmap.NewOrderedMap(),
		readRounds:  &atomic.Value{},
		// this is intentionally unbuffered,
		// do not change
		signal: make(chan struct{}),
	}
	roundsList := make([]*Round, 0, 0)
	wr.readRounds.Store(roundsList)

	return &wr
}

// Len returns the number of rounds in the list.
func (wr *WaitingRounds) Len() int {
	return len(wr.readRounds.Load().([]*Round))
}

// NumValidRounds returns how many rounds are, according to the local timestamp,
// ready to be sent to.
// This means they are in the "QUEUED" state and their start time is
// after the local time
func (wr *WaitingRounds) NumValidRounds(now time.Time) int {
	rounds := wr.readRounds.Load().([]*Round)

	numValid := 0

	for _, r := range rounds {
		if r.StartTime().After(now) {
			numValid++
		}
	}

	return numValid
}

// HasValidRounds returns true if there is at least one valid round
// in the queue according to its timestamp.
// This means they are in the "QUEUED" state and their start time is
// after the local time
func (wr *WaitingRounds) HasValidRounds(now time.Time) bool {
	rounds := wr.readRounds.Load().([]*Round)

	for _, r := range rounds {
		if r.StartTime().After(now) {
			return true
		}
	}

	return false
}

// Insert inserts a queued round into the list in order of its timestamp, from
// smallest to greatest. If the new round is not in a QUEUED state, then it is
// not inserted. If the new round already exists in the list but is no longer
// queued, then it is removed.
func (wr *WaitingRounds) Insert(added, removed []*Round) {
	wr.mux.Lock()
	defer wr.mux.Unlock()
	now := netTime.Now()
	// Add any round which should be added
	var addedRounds uint
	for i := range added {
		toAdd := added[i]
		if toAdd.StartTime().After(now) {
			addedRounds++
			wr.writeRounds.Set(toAdd.info.ID, toAdd)
		}
	}

	// Remove any round which should be removed
	for i := range removed {
		toRemove := removed[i]
		wr.writeRounds.Delete(toRemove.info.ID)
	}

	// If changes occurred, update the atomic
	if len(removed) > 0 || addedRounds > 0 {
		wr.storeReadRounds()
	}

	// If inserts occurred, then signal to any waiting threads
	// Only do this on inserts because only inserts will change the evaluation
	// by callers of GetUpcomingRealtime
	if addedRounds > 0 {
		go func() {
			// This will loop for as many people are waiting on the channel,
			// which is why it is in a separate function
			for {
				select {
				case wr.signal <- struct{}{}:
				default:
					// Exit when there are no listeners
					return
				}
			}
		}()
	}
}

func (wr *WaitingRounds) storeReadRounds() {
	roundsList := make([]*Round, 0, wr.writeRounds.Len())
	toDelete := make([]*Round, 0, wr.writeRounds.Len())

	now := netTime.Now()

	//filter rounds which should not be included
	for e := wr.writeRounds.Front(); e != nil; e = e.Next() {
		rnd := e.Value.(*Round)
		if now.Before(rnd.StartTime()) {
			roundsList = append(roundsList, rnd)
		} else {
			toDelete = append(toDelete, rnd)
		}
	}

	//sort the rounds list, soonest first
	sort.Slice(roundsList, func(i, j int) bool {
		return roundsList[i].StartTime().Before(roundsList[j].StartTime())
	})

	if jww.LogThreshold() == jww.LevelTrace {
		var rprint string
		for _, r := range roundsList {
			rprint += fmt.Sprintf("\n\tround: %d, startTime: %s, time to start: %s",
				r.info.ID, r.StartTime(), netTime.Until(r.StartTime()))
		}
		jww.TRACE.Printf("Rounds Order: %s", rprint)
	}

	wr.readRounds.Store(roundsList)

	if len(toDelete) > 0 {
		for i := range toDelete {
			toRemove := toDelete[i]
			wr.writeRounds.Delete(toRemove.info.ID)
		}
	}
}

// getFurthest returns the round that will occur furthest in the future. If the
// list is empty, then nil is returned. If the round is on the exclusion list,
// then the next round is checked. If it is not on the exclusion list, it is
// added.
// This is assumed to be called on an operation already under the cond's lock.
func (wr *WaitingRounds) getFurthest(exclude excludedRounds.ExcludedRounds,
	cutoffDelta time.Duration) *Round {
	earliestStart := netTime.Now().Add(cutoffDelta)

	roundsList, exists := wr.readRounds.Load().([]*Round)
	if !exists {
		return nil
	}

	// Return the last non-excluded round in the list
	for i := len(roundsList) - 1; i >= 0; i-- {
		r := roundsList[i]

		// Cannot guarantee that the round object's pointers will be exact match
		// of value in set
		if r.StartTime().After(earliestStart) {
			// If no excluded list has been passed in, do not check
			if exclude == nil {
				return r
			}

			// If the exclusion list exists, attempt and insert and return true
			// if it was a new insert, otherwise skip
			newInsertion := exclude.Insert(id.Round(r.info.ID))
			if newInsertion {
				return r
			}
		}
	}

	// If all the rounds in the list are excluded, then return nil
	return nil
}

// getClosest returns the round that will occur soonest in the future. If the
// list is empty, then nil is returned. If the round is on the exclusion list,
// then the next round is checked. If it is not on the exclusion list, it is
// added.
// This is assumed to be called on an operation already under the cond's lock.
func (wr *WaitingRounds) getClosest(exclude excludedRounds.ExcludedRounds,
	minRoundAge time.Duration) *Round {
	earliestStart := netTime.Now().Add(minRoundAge)

	roundsList, exists := wr.readRounds.Load().([]*Round)
	if !exists {
		return nil
	}

	// Return the first non-excluded round in the list
	for i := 0; i < len(roundsList); i++ {
		r := roundsList[i]

		// Cannot guarantee that the round object's pointers will be exact match
		// of value in set
		if r.StartTime().After(earliestStart) {
			// If no excluded list has been passed in, do not check
			if exclude == nil {
				return r
			}

			// If the exclusion list exists, then attempt an insert and return
			// if it was a new insert, otherwise skip
			newInsertion := exclude.Insert(id.Round(r.info.ID))
			if newInsertion {
				return r
			}
		}
	}

	// If all the rounds in the list are excluded, then return nil
	return nil
}

// GetSlice returns a slice of all round infos in the list that have yet to
// occur.
func (wr *WaitingRounds) GetSlice() []*pb.RoundInfo {
	var roundInfos []*pb.RoundInfo

	roundsList, exists := wr.readRounds.Load().([]*Round)
	if !exists {
		return roundInfos
	}

	timeNow := netTime.Now()
	for i := 0; i < len(roundsList); i++ {
		if roundsList[i].StartTime().After(timeNow) {
			roundInfos = append(roundInfos, roundsList[i].info)
		}
	}

	// Return the last round in the list, which is the furthest in the future
	return roundInfos
}

// GetUpcomingRealtime returns the round that will occur furthest in the future.
// If the list is empty, then it waits for a round to be added for the
// specified duration. If no round is added, then an error is returned.
//
// The length of the excluded set indicates how many times the client has
// called GetUpcomingRealtime trying to retrieve a round to send on.
// GetUpcomingRealtime defaults to retrieving the closest non-excluded round
// from WaitingRounds. If the length of the excluded set exceeds the maximum
// attempts at pulling the closest round, GetUpcomingRealtime will retrieve
// the furthest non-excluded round from WaitingRounds.
func (wr *WaitingRounds) GetUpcomingRealtime(timeout time.Duration,
	exclude excludedRounds.ExcludedRounds, numAttempts int, minRoundAge time.Duration) (*pb.RoundInfo, time.Duration, error) {

	// Start timeout timer
	timer := time.NewTimer(timeout)

	delay := multiply(numAttempts, minRoundAge)

	// Start seeing if an acceptable round exists
	round := wr.get(exclude, delay)
	if round != nil {
		return round, delay, nil
	}

	jww.INFO.Printf("Could not find round to send on, waiting for update")
	// If the no round exists, wait for an update to the list.
	for {
		select {
		case <-timer.C:
			return nil, 0, timeOutError
		case <-wr.signal:
			round = wr.get(exclude, 0)
			if round != nil {
				return round, 0, nil
			}
		}
	}
}

func (wr *WaitingRounds) get(exclude excludedRounds.ExcludedRounds, delay time.Duration) *pb.RoundInfo {

	round := wr.getClosest(exclude, delay)
	if round != nil {
		return round.Get()
	}

	return nil
}

type fraction struct {
	numerator   uint
	denominator uint
}

var roundAgeMultiplier = []fraction{
	fraction{1, 1},
	fraction{5, 4},
	fraction{7, 4},
	fraction{11, 4},
	fraction{19, 4},
	fraction{27, 4},
	fraction{35, 4},
}

func multiply(n int, duration time.Duration) time.Duration {
	if n >= len(roundAgeMultiplier) {
		n = len(roundAgeMultiplier) - 1
	}

	f := roundAgeMultiplier[n]
	duration = duration * time.Duration(f.numerator)
	duration = duration / time.Duration(f.denominator)
	return duration
}

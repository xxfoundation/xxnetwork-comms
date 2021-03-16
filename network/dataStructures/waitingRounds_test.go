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
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/testutils"
	"gitlab.com/elixxir/primitives/current"
	"gitlab.com/elixxir/primitives/states"
	"math/rand"
	"reflect"
	"sync"
	"testing"
	"time"
)

// Happy path of NewWaitingRounds().
func TestNewWaitingRounds(t *testing.T) {
	expectedWR := &WaitingRounds{rounds: list.New()}
	m := sync.Mutex{}
	expectedWR.c = sync.NewCond(&m)

	testWR := NewWaitingRounds()

	if !reflect.DeepEqual(expectedWR, testWR) {
		t.Errorf("NewWaitingRounds() did not return the expected object."+
			"\n\texpected: %#v\n\trecieved: %#v", expectedWR, testWR)
	}
}

// Happy path of WaitingRounds.Len().
func TestWaitingRounds_Len(t *testing.T) {
	expectedLen := list.New().Len()

	testWR := NewWaitingRounds()

	if expectedLen != testWR.Len() {
		t.Errorf("Len() did not return the expected length."+
			"\n\texpected: %d\n\trecieved: %d", expectedLen, testWR.Len())
	}
}

// Happy path of WaitingRounds.Insert().
func TestWaitingRounds_Insert(t *testing.T) {
	// Generate rounds
	expectedRounds, testRounds := createTestRoundInfos(25, t)

	// Add rounds to list
	testWR := NewWaitingRounds()
	for _, round := range testRounds {
		testWR.Insert(round)
	}

	if len(expectedRounds) != testWR.Len() {
		t.Fatalf("List does not have the expected length."+
			"\n\texpected: %v\n\trecieved: %v", len(expectedRounds), testWR.Len())
	}

	// Check if they are in the correct order
	for e, i := testWR.rounds.Front(), 0; e != nil; e, i = e.Next(), i+1 {
		if expectedRounds[i].info != e.Value.(*Round).info {
			t.Errorf("Insert() did not inser the correct element at position %d."+
				"\n\texpected: %v\n\trecieved: %v", i, expectedRounds[i], e.Value)
		}
	}
}

// Happy path of WaitingRounds.remove().
func TestWaitingRounds_remove(t *testing.T) {
	// Generate rounds
	expectedRounds, testRounds := createTestRoundInfos(25, t)

	// Add rounds to list
	testWR := NewWaitingRounds()
	for _, round := range expectedRounds {
		testWR.Insert(round)
	}

	for _, round := range testRounds {
		testWR.remove(round)
	}

	if testWR.Len() != 0 {
		t.Fatalf("remove() did not remove everything, list still has %d rounds.",
			testWR.Len())
	}
}

// Happy path of WaitingRounds.getFurthest().
func TestWaitingRounds_getFurthest(t *testing.T) {
	// Generate rounds
	expectedRounds, testRounds := createTestRoundInfos(25, t)

	// Add rounds to list
	testWR := NewWaitingRounds()
	for _, round := range testRounds {
		testWR.Insert(round)
	}

	for i := len(expectedRounds) - 1; i >= 0; i-- {
		if !reflect.DeepEqual(expectedRounds[i], testWR.getFurthest(nil)) {
			t.Errorf("getFurthest() did not return the expected round."+
				"\n\texpected: %v\n\trecieved: %v", expectedRounds[i], testWR.getFurthest(nil))
		}
		testWR.remove(expectedRounds[i])
	}

	if testWR.getFurthest(nil) != nil {
		t.Errorf("getFurthest() did not return nil on empty list.")
	}
}

// Happy path of WaitingRounds.getFurthest() with half the rounds excluded
func TestWaitingRounds_getFurthest_Exclude(t *testing.T) {
	// Generate rounds
	expectedRounds, testRounds := createTestRoundInfos(25, t)

	// Add rounds to list
	testWR := NewWaitingRounds()
	for _, round := range testRounds {
		testWR.Insert(round)
	}

	// Add rounds to exclusion list
	exclude := set.New()
	for i, round := range expectedRounds {
		if i%2 == 0 {
			exclude.Insert(round.info)
		}
	}

	for i := len(expectedRounds) - 1; i >= 0; i-- {
		if i%2 == 1 {
			received := testWR.getFurthest(exclude)
			if !reflect.DeepEqual(expectedRounds[i], received) {
				t.Errorf("getFurthest() did not return the expected round."+
					"\n\texpected: %v\n\trecieved: %v", expectedRounds[i], received)
			}
			testWR.remove(expectedRounds[i])
		}
	}
	if testWR.getFurthest(exclude) != nil {
		t.Errorf("getFurthest() did not return nil on empty list.")
	}
}

// Happy path of WaitingRounds.GetUpcomingRealtime() when the list is not empty
// and no waiting occurs.
func TestWaitingRounds_GetUpcomingRealtime_NoWait(t *testing.T) {
	// Generate rounds and add to new list
	expectedRounds, testRounds := createTestRoundInfos(25, t)
	testWR := NewWaitingRounds()

	for _, round := range testRounds {
		testutils.SignRoundInfo(round.info, t)
		testWR.Insert(round)
	}

	for i := len(expectedRounds) - 1; i >= 0; i-- {
		furthestRound, err := testWR.GetUpcomingRealtime(300*time.Millisecond, set.New())
		if err != nil {
			t.Errorf("GetUpcomingRealtime() returned an unexpected error."+
				"\n\terror: %v", err)
		}
		if expectedRounds[i].info != furthestRound {
			t.Errorf("GetUpcomingRealtime() did not return the expected round (%d)."+
				"\n\texpected: %v\n\trecieved: %v", i, expectedRounds[i], furthestRound)
		}
		testWR.remove(expectedRounds[i])
	}
}

// Tests that WaitingRounds.GetUpcomingRealtime() returns an error on an empty
// list after timeout.
func TestWaitingRounds_GetUpcomingRealtime_TimeoutError(t *testing.T) {
	testWR := NewWaitingRounds()

	furthestRound, err := testWR.GetUpcomingRealtime(300*time.Millisecond, set.New())
	if err != timeOutError {
		t.Errorf("GetUpcomingRealtime() did not time out when expected."+
			"\n\texpected: %v\n\treceived: %v", timeOutError, err)
	}
	if furthestRound != nil {
		t.Errorf("GetUpcomingRealtime() did not return nil on empty list."+
			"\n\texpected: %v\n\trecieved: %v", nil, furthestRound)
	}
}

// Happy path of WaitingRounds.GetUpcomingRealtime().
func TestWaitingRounds_GetUpcomingRealtime(t *testing.T) {
	// Generate rounds and WaitingRound
	expectedRounds, _ := createTestRoundInfos(25, t)
	testWR := NewWaitingRounds()

	for i, round := range expectedRounds {
		go func(round *Round) {
			time.Sleep(30 * time.Millisecond)
			testutils.SignRoundInfo(round.info, t)
			testWR.Insert(round)
		}(round)

		furthestRound, err := testWR.GetUpcomingRealtime(5*time.Second, set.New())

		if err != nil {
			t.Errorf("GetUpcomingRealtime() returned an unexpected error (%d)."+
				"\n\terror: %v", i, err)
		}
		if round.info != furthestRound {
			t.Errorf("GetUpcomingRealtime() did not return the expected round (%d)."+
				"\n\texpected: %v\n\trecieved: %v", i, round, furthestRound)
		}
		testWR.remove(round)
	}
}

// Happy path of WaitingRounds.GetSlice().
func TestWaitingRounds_GetSlice(t *testing.T) {
	// Generate rounds and add to new list
	expectedRounds, testRounds := createTestRoundInfos(25, t)
	testWR := NewWaitingRounds()

	ri := &pb.RoundInfo{
		ID:         rand.Uint64(),
		State:      uint32(states.QUEUED),
		Timestamps: []uint64{0, 0, 0, 0, 0},
	}

	testutils.SignRoundInfo(ri, t)

	pubKey, err := testutils.LoadPublicKeyTesting(t)
	if err != nil {
		t.Errorf("Failed to load public key: %v", err)
		t.FailNow()
	}
	rnd := NewRound(ri, pubKey)

	testRounds = append(testRounds, rnd)
	for _, round := range testRounds {
		testWR.Insert(round)
	}

	testSlice := testWR.GetSlice()

	// Convert Round slice to round info slice
	expectedRoundInfos := make([]*pb.RoundInfo, 0)
	for _, val := range expectedRounds {
		expectedRoundInfos = append(expectedRoundInfos, val.info)
	}

	if !reflect.DeepEqual(expectedRoundInfos, testSlice) {
		t.Errorf("GetSlice() returned slice with incorrect rounds."+
			"\n\texepcted: %v\n\treceived: %v",
			expectedRoundInfos, testSlice)
	}
}

// Generates two lists of round infos. The first is the expected rounds in the
// correct order after inserting the second list of random round infos.
func createTestRoundInfos(num int, t *testing.T) ([]*Round, []*Round) {
	rounds := make([]*Round, num)
	var expectedRounds []*Round
	randomRounds := make([]*Round, num)
	timeTrack := uint64(time.Now().Nanosecond()) + 100000000
	pubKey, err := testutils.LoadPublicKeyTesting(t)
	if err != nil {
		t.Errorf("Failed to load public key: %v", err)
		t.FailNow()
	}

	for i := 0; i < num; i++ {
		rounds[i] = NewRound(&pb.RoundInfo{
			ID:         rand.Uint64(),
			State:      uint32(rand.Int63n(int64(states.NUM_STATES) - 1)),
			Timestamps: make([]uint64, current.NUM_STATES),
		}, pubKey)
		timeTrack += uint64(rand.Int63n(10000))
		rounds[i].info.Timestamps[current.STANDBY] = timeTrack
		if i%2 == 1 {
			rounds[i].info.State = uint32(states.QUEUED)
		}
		if rounds[i].info.State == uint32(states.QUEUED) {
			expectedRounds = append(expectedRounds, rounds[i])
		}
	}
	perm := rand.Perm(num)
	for i, v := range perm {
		randomRounds[v] = rounds[i]
	}
	return expectedRounds, randomRounds
}

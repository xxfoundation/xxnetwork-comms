////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package network

import (
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	ds "gitlab.com/elixxir/comms/network/dataStructures"
	"gitlab.com/xx_network/primitives/id"
	"testing"
)

//region ERS Memory Map Impl
// Memory map based ExtendedRoundStorage database
type ersMemMap struct {
	rounds map[id.Round]*pb.RoundInfo
}

// Store a new round info object into the map
func (ersm *ersMemMap) Store(ri *pb.RoundInfo) error {
	rid := id.Round(ri.GetID())
	// See if the round exists, if it does we need to check that the update ID is newer than the current
	ori, err := ersm.Retrieve(rid)
	if err != nil {
		return err
	}

	if ori == nil || ori.UpdateID < ri.UpdateID {
		ersm.rounds[rid] = ri
	} else {
		jww.WARN.Printf("Passed in round update ID of %v lower than currently stored ID %v",
			ri.UpdateID, ori.UpdateID)
	}

	return nil
}

// Get a round info object from the memory map database
func (ersm *ersMemMap) Retrieve(id id.Round) (*pb.RoundInfo, error) {
	return ersm.rounds[id], nil
}

// Get multiple specific round info objects from the memory map database
func (ersm *ersMemMap) RetrieveMany(rounds []id.Round) ([]*pb.RoundInfo, error) {
	var r []*pb.RoundInfo

	for _, round := range rounds {
		ri, err := ersm.Retrieve(round)
		if err != nil {
			return nil, err
		}

		r = append(r, ri)
	}

	return r, nil
}

// Retrieve a concurrent range of round info objects from the memory map database
func (ersm *ersMemMap) RetrieveRange(first, last id.Round) ([]*pb.RoundInfo, error) {
	idrange := uint64(last - first)
	i := uint64(0)

	var r []*pb.RoundInfo

	// for some reason <= doesn't work?
	for i < idrange+1 {
		ri, err := ersm.Retrieve(id.Round(uint64(first) + i))
		if err != nil {
			return nil, err
		}

		r = append(r, ri)
		i++
	}

	return r, nil
}

//endregion

// Test we can insert a round, get it, try to update with an older ID, it doesn't update, and it does update with
// a newer ID
func TestERSStore(t *testing.T) {
	// Setup
	var ers ds.ExternalRoundStorage = &ersMemMap{rounds: make(map[id.Round]*pb.RoundInfo)}

	// Store a test round
	r := pb.RoundInfo{ID: 1, UpdateID: 5}
	err := ers.Store(&r)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Test that we get a new round with our info
	ri, err := ers.Retrieve(id.Round(r.ID))
	if err != nil {
		t.Errorf(err.Error())
	}
	if ri == nil {
		t.Fatalf("ri object is nil, Retrieve did not return a round")
	}
	if ri.ID != r.ID && ri.UpdateID != r.UpdateID {
		t.Errorf("did not return the same round ID or update ID as we inputted.")
	}

	// Update the test round
	ru1 := pb.RoundInfo{ID: 1, UpdateID: 3}
	err = ers.Store(&ru1)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Test that this updated round did not get written to the map
	riu1, err := ers.Retrieve(id.Round(r.ID))
	if err != nil {
		t.Errorf(err.Error())
	}
	if riu1 == nil {
		t.Fatalf("ri object is nil, Retrieve did not return a round")
	}
	if riu1.UpdateID == ru1.UpdateID {
		t.Errorf("stored round info was updated to have lower update ID of %v", riu1.UpdateID)
	}

	// Update the test round
	ru2 := pb.RoundInfo{ID: 1, UpdateID: 10}
	err = ers.Store(&ru2)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Test that this updated round did not get written to the map
	riu2, err := ers.Retrieve(id.Round(r.ID))
	if err != nil {
		t.Errorf(err.Error())
	}
	if riu2 == nil {
		t.Fatalf("ri object is nil, Retrieve did not return a round")
	}
	if riu2.UpdateID != ru2.UpdateID {
		t.Errorf("stored round info was not updated to have update ID of %v", riu2.UpdateID)
	}
}

// Test that Retrieve has expected behaviour if an item doesn't exist
func TestERSRetrieve(t *testing.T) {
	var ers ds.ExternalRoundStorage = &ersMemMap{rounds: make(map[id.Round]*pb.RoundInfo)}
	ri, err := ers.Retrieve(id.Round(1))
	if err != nil {
		t.Errorf(err.Error())
	}
	if ri != nil {
		t.Errorf("returned round info was not nil")
	}

	// Store a test round
	r := pb.RoundInfo{ID: 1, UpdateID: 5}
	err = ers.Store(&r)
	if err != nil {
		t.Errorf(err.Error())
	}

	nri, err := ers.Retrieve(id.Round(1))
	if err != nil {
		t.Errorf(err.Error())
	}
	if nri == nil {
		t.Fatalf("returned round info was nil")
	}
	if nri.ID != r.ID || nri.UpdateID != r.UpdateID {
		t.Errorf("Returned round or update ID did not match what we put in")
	}
}

// Test that the RetrieveMany function will get rounds that are stored, while returning nil with no error for those
// that are not
func TestERSRetrieveMany(t *testing.T) {
	// Setup
	var ers ds.ExternalRoundStorage = &ersMemMap{rounds: make(map[id.Round]*pb.RoundInfo)}

	// Store a test round
	origRound1 := pb.RoundInfo{ID: 1, UpdateID: 5}
	err := ers.Store(&origRound1)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Store another test round
	origRound2 := pb.RoundInfo{ID: 8, UpdateID: 3}
	err = ers.Store(&origRound2)
	if err != nil {
		t.Errorf(err.Error())
	}

	getRounds := []id.Round{id.Round(origRound1.ID), id.Round(origRound2.ID - 3), id.Round(origRound2.ID)}
	returnRounds, err := ers.RetrieveMany(getRounds)
	if err != nil {
		t.Errorf(err.Error())
	}

	if returnRounds[0] == nil || returnRounds[2] == nil {
		t.Fatalf("RetrieveMany did not return a round we expected to get returned")
	}
	if returnRounds[1] != nil {
		t.Errorf("Middle fake round did return a round info object")
	}
	if returnRounds[0].ID != origRound1.ID || returnRounds[0].UpdateID != origRound1.UpdateID {
		t.Errorf("First returned round and original mismatched IDs")
	}
	if returnRounds[2].ID != origRound2.ID || returnRounds[2].UpdateID != origRound2.UpdateID {
		t.Errorf("Second returned round and original mismatched IDs")
	}
}

// Test that the RetrieveRange function will get a range of rounds stored, and return nil with no error for ones that
// are not stored
func TestERSRetrieveRange(t *testing.T) {
	// Setup
	var ers ds.ExternalRoundStorage = &ersMemMap{rounds: make(map[id.Round]*pb.RoundInfo)}

	// Store a test round
	origRound1 := pb.RoundInfo{ID: 1, UpdateID: 5}
	err := ers.Store(&origRound1)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Store another test round
	origRound2 := pb.RoundInfo{ID: 3, UpdateID: 3}
	err = ers.Store(&origRound2)
	if err != nil {
		t.Errorf(err.Error())
	}

	returnRounds, err := ers.RetrieveRange(1, 3)
	if err != nil {
		t.Errorf(err.Error())
	}

	if returnRounds[0] == nil || returnRounds[2] == nil {
		t.Fatalf("RetrieveMany did not return a round we expected to get returned")
	}
	if returnRounds[1] != nil {
		t.Errorf("Middle fake round did return a round info object")
	}
	if returnRounds[0].ID != origRound1.ID || returnRounds[0].UpdateID != origRound1.UpdateID {
		t.Errorf("First returned round and original mismatched IDs")
	}
	if returnRounds[2].ID != origRound2.ID || returnRounds[2].UpdateID != origRound2.UpdateID {
		t.Errorf("Second returned round and original mismatched IDs")
	}
}

// Test that calling the GetHistoricalRound interface on Instance works
func TestInstance_GetHistoricalRound(t *testing.T) {
	i := Instance{}
	ri, err := i.GetHistoricalRound(id.Round(0))
	// This should fail since this Instance doesn't have an ERS object
	if err == nil {
		t.Errorf(err.Error())
	}

	var ers ds.ExternalRoundStorage = &ersMemMap{rounds: make(map[id.Round]*pb.RoundInfo)}
	i2 := Instance{ers: ers}
	ri, err = i2.GetHistoricalRound(id.Round(0))
	// Should return no error and blank round
	if err != nil {
		t.Errorf(err.Error())
	}
	if ri != nil {
		t.Errorf("ri contains round info")
	}
}

// Test that calling the GetHistoricalRoundRange interface on Instance works
func TestInstance_GetHistoricalRoundRange(t *testing.T) {
	i := Instance{}
	ri, err := i.GetHistoricalRoundRange(5, 10)
	// This should fail since this Instance doesn't have an ERS object
	if err == nil {
		t.Errorf(err.Error())
	}

	var ers ds.ExternalRoundStorage = &ersMemMap{rounds: make(map[id.Round]*pb.RoundInfo)}
	i2 := Instance{ers: ers}
	ri, err = i2.GetHistoricalRoundRange(5, 10)
	// Should return no error and blank round
	if err != nil {
		t.Errorf(err.Error())
	}
	for _, round := range ri {
		if round != nil {
			t.Errorf("ri contains round info")
		}
	}
}

// Test that calling the GetHistoricalRounds interface on Instance works
func TestInstance_GetHistoricalRounds(t *testing.T) {
	i := Instance{}
	getRounds := []id.Round{id.Round(1), id.Round(2), id.Round(3)}
	ri, err := i.GetHistoricalRounds(getRounds)
	// This should fail since this Instance doesn't have an ERS object
	if err == nil {
		t.Errorf(err.Error())
	}

	var ers ds.ExternalRoundStorage = &ersMemMap{rounds: make(map[id.Round]*pb.RoundInfo)}
	i2 := Instance{ers: ers}
	ri, err = i2.GetHistoricalRounds(getRounds)
	// Should return no error and blank round
	if err != nil {
		t.Errorf(err.Error())
	}
	for _, round := range ri {
		if round != nil {
			t.Errorf("ri contains round info")
		}
	}
}

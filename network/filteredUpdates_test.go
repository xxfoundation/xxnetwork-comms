///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package network

import (
	"gitlab.com/elixxir/comms/mixmessages"
	ds "gitlab.com/elixxir/comms/network/dataStructures"
	"gitlab.com/elixxir/comms/testkeys"
	"gitlab.com/elixxir/comms/testutils"
	"gitlab.com/elixxir/primitives/states"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/signature"
	"gitlab.com/xx_network/primitives/id"
	"testing"
)

func TestFilteredUpdates_RoundUpdate(t *testing.T) {
	validUpdateId := uint64(4)
	validMsg := &mixmessages.RoundInfo{
		ID:        2,
		UpdateID:  validUpdateId,
		State:     uint32(states.COMPLETED),
		BatchSize: 8,
	}
	privKey, err := testutils.LoadPrivateKeyTesting(t)
	if err != nil {
		t.Errorf("Failed to load private key: %v", err)
		t.FailNow()
	}
	err = signature.Sign(validMsg, privKey)
	testManager := connect.NewManagerTesting(t)
	pc := connect.ProtoComms{
		Manager: testManager,
	}
	testFilter := NewFilteredUpdates(&pc)
	pub := testkeys.LoadFromPath(testkeys.GetNodeCertPath())
	err = testFilter.roundUpdate(validMsg)
	if err == nil {
		t.Error("Should have failed to get perm host")
	}

	_, err = testFilter.comm.AddHost(&id.Permissioning, "0.0.0.0:4200", pub, connect.GetDefaultHostParams())
	if err != nil {
		t.Errorf("failed to add bad host: %+v", err)
	}
	err = testFilter.roundUpdate(validMsg)
	// Fixme
	/*	if err == nil {
		t.Error("Should have failed to verify")
	}*/

	retrieved, err := testFilter.GetRoundUpdate(int(validUpdateId))
	if err != nil || retrieved == nil {
		t.Errorf("Should have stored msg with state %s", states.Round(validMsg.State))
	}

	invalidUpdateId := uint64(5)
	invalidMsg := &mixmessages.RoundInfo{
		ID:        2,
		UpdateID:  invalidUpdateId,
		State:     uint32(states.PRECOMPUTING),
		BatchSize: 8,
	}

	err = testFilter.roundUpdate(invalidMsg)
	if err != nil {
		t.Errorf("Failed to update round: %v", err)
	}

	retrieved, err = testFilter.GetRoundUpdate(int(invalidUpdateId))
	if err == nil || retrieved != nil {
		t.Errorf("Should not have inserted round with state %s",
			states.Round(invalidMsg.State))
	}

}

func TestFilteredUpdates_RoundUpdates(t *testing.T) {
	validUpdateId := uint64(4)
	validMsg := &mixmessages.RoundInfo{
		ID:        2,
		UpdateID:  validUpdateId,
		State:     uint32(states.COMPLETED),
		BatchSize: 8,
	}
	privKey, err := testutils.LoadPrivateKeyTesting(t)
	if err != nil {
		t.Errorf("Failed to load private key: %v", err)
		t.FailNow()
	}
	err = signature.Sign(validMsg, privKey)
	testManager := connect.NewManagerTesting(t)
	pc := connect.ProtoComms{
		Manager: testManager,
	}
	testFilter := NewFilteredUpdates(&pc)
	pub := testkeys.LoadFromPath(testkeys.GetNodeCertPath())

	_, err = testFilter.comm.AddHost(&id.Permissioning, "0.0.0.0:4200", pub, connect.GetDefaultHostParams())
	if err != nil {
		t.Errorf("failed to add bad host: %+v", err)
	}

	invalidUpdateId := uint64(5)
	invalidMsg := &mixmessages.RoundInfo{
		ID:        2,
		UpdateID:  invalidUpdateId,
		State:     uint32(states.PRECOMPUTING),
		BatchSize: 8,
	}

	roundUpdates := []*mixmessages.RoundInfo{validMsg, invalidMsg}

	err = testFilter.RoundUpdates(roundUpdates)
	if err != nil {
		t.Error("Should have failed to get perm host")
	}

	retrieved, err := testFilter.GetRoundUpdate(int(validUpdateId))
	if err != nil || retrieved == nil {
		t.Errorf("Should have stored msg with state %s", states.Round(validMsg.State))
	}

	retrieved, err = testFilter.GetRoundUpdate(int(invalidUpdateId))
	if err == nil || retrieved != nil {
		t.Errorf("Should not have inserted round with state %s",
			states.Round(invalidMsg.State))
	}

}

func TestFilteredUpdates_GetRoundUpdate(t *testing.T) {
	i := FilteredUpdates{
		updates: ds.NewUpdates(),
	}

	ri := &mixmessages.RoundInfo{
		ID:       uint64(1),
		UpdateID: uint64(1),
		State:    uint32(states.QUEUED),
	}
	testutils.SignRoundInfo(ri, t)
	pubKey, err := testutils.LoadPublicKeyTesting(t)
	if err != nil {
		t.Errorf("Failed to load test key: %v", err)
	}
	rnd := ds.NewRound(ri, pubKey)

	_ = i.updates.AddRound(rnd)
	r, err := i.GetRoundUpdate(1)
	if err != nil || r == nil {
		t.Errorf("Failed to retrieve round update: %+v", err)
	}
}

func TestFilteredUpdates_GetRoundUpdates(t *testing.T) {
	i := FilteredUpdates{
		updates: ds.NewUpdates(),
	}
	pubKey, err := testutils.LoadPublicKeyTesting(t)
	if err != nil {
		t.Errorf("Failed to load public key: %v", err)
		t.FailNow()
	}

	roundInfoOne := &mixmessages.RoundInfo{
		ID:       uint64(1),
		UpdateID: uint64(2),
		State:    uint32(states.QUEUED),
	}
	testutils.SignRoundInfo(roundInfoOne, t)
	roundInfoTwo := &mixmessages.RoundInfo{
		ID:       uint64(2),
		UpdateID: uint64(3),
		State:    uint32(states.QUEUED),
	}
	testutils.SignRoundInfo(roundInfoTwo, t)
	roundOne := ds.NewRound(roundInfoOne, pubKey)
	roundTwo := ds.NewRound(roundInfoTwo, pubKey)

	_ = i.updates.AddRound(roundOne)
	_ = i.updates.AddRound(roundTwo)
	r := i.GetRoundUpdates(1)
	if len(r) == 0 {
		t.Errorf("Failed to retrieve round updates")
	}

	r = i.GetRoundUpdates(2)
	if len(r) == 0 {
		t.Errorf("Failed to retrieve round updates")
	}

	r = i.GetRoundUpdates(23)
	if len(r) != 0 {
		t.Errorf("Retrieved a round that was never inserted: %v", r)
	}

}

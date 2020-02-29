package network

import (
	"gitlab.com/elixxir/comms/mixmessages"
	ds "gitlab.com/elixxir/comms/network/dataStructures"
	"gitlab.com/elixxir/primitives/id"
	"testing"
)

func TestNewInstance(t *testing.T) {

}

func TestInstance_GetFullNdf(t *testing.T) {
	i := Instance{
		full: NewSecuredNdf(),
	}
	if i.GetFullNdf() == nil {
		t.Error("Failed to retrieve full ndf")
	}
}

func TestInstance_GetPartialNdf(t *testing.T) {
	i := Instance{
		partial: NewSecuredNdf(),
	}
	if i.GetPartialNdf() == nil {
		t.Error("Failed to retrieve partial ndf")
	}
}

func TestInstance_GetRound(t *testing.T) {
	i := Instance{
		roundData: &ds.Data{},
	}
	_ = i.roundData.UpsertRound(&mixmessages.RoundInfo{ID: uint64(1)})
	r, err := i.GetRound(id.Round(1))
	if err != nil || r == nil {
		t.Errorf("Failed to retrieve round: %+v", err)
	}
}

func TestInstance_GetRoundUpdate(t *testing.T) {
	i := Instance{
		roundUpdates: &ds.Updates{},
	}
	_ = i.roundUpdates.AddRound(&mixmessages.RoundInfo{ID: uint64(1), UpdateID: uint64(1)})
	r, err := i.GetRoundUpdate(1)
	if err != nil || r == nil {
		t.Errorf("Failed to retrieve round update: %+v", err)
	}
}

func TestInstance_GetRoundUpdates(t *testing.T) {
	i := Instance{
		roundUpdates: &ds.Updates{},
	}
	_ = i.roundUpdates.AddRound(&mixmessages.RoundInfo{ID: uint64(1), UpdateID: uint64(1)})
	_ = i.roundUpdates.AddRound(&mixmessages.RoundInfo{ID: uint64(1), UpdateID: uint64(2)})
	r, err := i.GetRoundUpdates(1)
	if err != nil || r == nil {
		t.Errorf("Failed to retrieve round update: %+v", err)
	}
}

func TestInstance_RoundUpdate(t *testing.T) {

}

func TestInstance_UpdateFullNdf(t *testing.T) {

}

func TestInstance_UpdatePartialNdf(t *testing.T) {

}

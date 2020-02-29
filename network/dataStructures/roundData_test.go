package dataStructures

import (
	"gitlab.com/elixxir/comms/mixmessages"
	"testing"
)

func TestData_UpsertRound(t *testing.T) {
	d := Data{}
	err := d.UpsertRound(&mixmessages.RoundInfo{
		ID:       0,
		UpdateID: 0,
	})
	if err != nil {
		t.Errorf("Failed to upsert round: %+v", err)
	}
}

func TestData_GetRound(t *testing.T) {
	d := Data{}
	_ = d.UpsertRound(&mixmessages.RoundInfo{
		ID:       0,
		UpdateID: 3,
	})
	_, err := d.GetRound(0)
	if err != nil {
		t.Errorf("Failed to get roundinfo with proper id")
	}
}

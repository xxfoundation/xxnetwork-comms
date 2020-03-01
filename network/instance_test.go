package network

import (
	"gitlab.com/elixxir/comms/connect"
	"gitlab.com/elixxir/comms/mixmessages"
	ds "gitlab.com/elixxir/comms/network/dataStructures"
	"gitlab.com/elixxir/comms/testkeys"
	"gitlab.com/elixxir/crypto/signature"
	"gitlab.com/elixxir/crypto/signature/rsa"
	"gitlab.com/elixxir/primitives/id"
	"gitlab.com/elixxir/primitives/ndf"
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

func setupComm(t *testing.T) (*Instance, *mixmessages.NDF) {
	priv := testkeys.LoadFromPath(testkeys.GetNodeKeyPath())
	privKey, err := rsa.LoadPrivateKeyFromPem(priv)
	pub := testkeys.LoadFromPath(testkeys.GetNodeCertPath())
	if err != nil {
		t.Errorf("Could not generate rsa key: %s", err)
	}

	f := &mixmessages.NDF{}

	baseNDF := ndf.NetworkDefinition{}
	f.Ndf, err = baseNDF.Marshal()

	if err != nil {
		t.Errorf("Could not generate serialized ndf: %s", err)
	}

	err = signature.Sign(f, privKey)

	pc := &connect.ProtoComms{}
	i := NewInstance(pc)

	_, err = i.comm.AddHost(id.PERMISSIONING, "0.0.0.0:4200", pub, false, true)
	if err != nil {
		t.Errorf("Failed to add permissioning host: %+v", err)
	}
	return i, f
}

func TestInstance_RoundUpdate(t *testing.T) {
	msg := &mixmessages.RoundInfo{
		ID:        2,
		UpdateID:  4,
		State:     6,
		BatchSize: 8,
	}
	priv := testkeys.LoadFromPath(testkeys.GetNodeKeyPath())
	privKey, err := rsa.LoadPrivateKeyFromPem(priv)
	err = signature.Sign(msg, privKey)

	i := NewInstance(&connect.ProtoComms{})
	pub := testkeys.LoadFromPath(testkeys.GetGatewayCertPath())
	err = i.RoundUpdate(msg)
	if err == nil {
		t.Error("Should have failed to get perm host")
	}

	_, err = i.comm.AddHost(id.PERMISSIONING, "0.0.0.0:4200", pub, false, true)
	if err != nil {
		t.Errorf("failed to add bad host: %+v", err)
	}
	err = i.RoundUpdate(msg)
	if err == nil {
		t.Error("Should have failed to verify")
	}

	i, _ = setupComm(t)

	err = i.RoundUpdate(msg)
	if err != nil {
		t.Errorf("Failed to update ndf: %+v", err)
	}
}

func TestInstance_UpdateFullNdf(t *testing.T) {
	i, f := setupComm(t)
	err := i.UpdateFullNdf(f)
	if err != nil {
		t.Errorf("Failed to update ndf: %+v", err)
	}
}

func TestInstance_UpdatePartialNdf(t *testing.T) {
	i, f := setupComm(t)
	err := i.UpdatePartialNdf(f)
	if err != nil {
		t.Errorf("Failed to update ndf: %+v", err)
	}
}

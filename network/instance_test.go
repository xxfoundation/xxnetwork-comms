////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package network

import (
	"gitlab.com/elixxir/comms/connect"
	"gitlab.com/elixxir/comms/mixmessages"
	ds "gitlab.com/elixxir/comms/network/dataStructures"
	"gitlab.com/elixxir/comms/testkeys"
	"gitlab.com/elixxir/comms/testutils"
	"gitlab.com/elixxir/crypto/signature"
	"gitlab.com/elixxir/crypto/signature/rsa"
	"gitlab.com/elixxir/primitives/id"
	"gitlab.com/elixxir/primitives/ndf"
	"testing"
)

func TestInstance_GetFullNdf(t *testing.T) {
	secured, _ := NewSecuredNdf(testutils.NDF)
	i := Instance{
		full: secured,
	}
	if i.GetFullNdf() == nil {
		t.Error("Failed to retrieve full ndf")
	}
}

func TestInstance_GetPartialNdf(t *testing.T) {
	secured, _ := NewSecuredNdf(testutils.NDF)
	i := Instance{
		partial: secured,
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
	i, err := NewInstance(pc, &baseNDF, &baseNDF)
	if err != nil {
		t.Error(nil)
	}

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

	i, err := NewInstance(&connect.ProtoComms{}, testutils.NDF, testutils.NDF)
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

func TestInstance_GetLastRoundID(t *testing.T) {
	i := Instance{
		roundData: &ds.Data{},
	}
	_ = i.roundData.UpsertRound(&mixmessages.RoundInfo{ID: uint64(1)})
	i.GetLastRoundID()
}

func TestInstance_GetLastUpdateID(t *testing.T) {
	i := Instance{
		roundUpdates: &ds.Updates{},
	}
	_ = i.roundUpdates.AddRound(&mixmessages.RoundInfo{ID: uint64(1), UpdateID: uint64(1)})
	i.GetLastUpdateID()
}

func TestInstance_UpdateGatewayConnections(t *testing.T) {
	secured, _ := NewSecuredNdf(testutils.NDF)

	i := Instance{
		full: secured,
		comm: &connect.ProtoComms{},
	}
	err := i.UpdateGatewayConnections()
	if err != nil {
		t.Errorf("Failed to update gateway connections from full: %+v", err)
	}

	i = Instance{
		partial: secured,
		comm:    &connect.ProtoComms{},
	}
	err = i.UpdateGatewayConnections()
	if err != nil {
		t.Errorf("Failed to update gateway connections from partial: %+v", err)
	}

	i = Instance{}
	err = i.UpdateGatewayConnections()
	if err == nil {
		t.Error("Should error when attempting update with no ndf")
	}
}

func TestInstance_UpdateNodeConnections(t *testing.T) {
	secured, _ := NewSecuredNdf(testutils.NDF)

	i := Instance{
		full: secured,
		comm: &connect.ProtoComms{},
	}
	err := i.UpdateNodeConnections()
	if err != nil {
		t.Errorf("Failed to update node connections from full: %+v", err)
	}

	i = Instance{
		partial: secured,
		comm:    &connect.ProtoComms{},
	}
	err = i.UpdateNodeConnections()
	if err != nil {
		t.Errorf("Failed to update node connections from partial: %+v", err)
	}

	i = Instance{}
	err = i.UpdateNodeConnections()
	if err == nil {
		t.Error("Should error when attempting update with no ndf")
	}
}

// Happy path: Tests GetPermissioningAddress with the full ndf set, the partial ndf set
// and no ndf set
func TestInstance_GetPermissioningAddress(t *testing.T) {
	// Create an instance object, setting the full ndf
	secured, _ := NewSecuredNdf(testutils.NDF)
	fullNdfInstance := Instance{
		full: secured,
	}

	// Expected address gotten from testutils.NDF
	expectedAddress := "92.42.125.61"

	// GetPermissioningAddress from the instance and compare with the expected value
	receivedAddress := fullNdfInstance.GetPermissioningAddress()
	if expectedAddress != receivedAddress {
		t.Errorf("GetPermissioningAddress did not get expected value!"+
			"\n\tExpected: %+v"+
			"\n\tReceived: %+v", expectedAddress, receivedAddress)
	}

	// Create an instance object, setting the partial ndf
	partialNdfInstance := Instance{
		partial: secured,
	}

	// GetPermissioningAddress from the instance and compare with the expected value
	receivedAddress = partialNdfInstance.GetPermissioningAddress()
	if expectedAddress != receivedAddress {
		t.Errorf("GetPermissioningAddress did not get expected value!"+
			"\n\tExpected: %+v"+
			"\n\tReceived: %+v", expectedAddress, receivedAddress)
	}

	// Create an instance object, setting no ndf
	noNdfInstance := Instance{}

	// GetPermissioningAddress, should be an empty string as no ndf's are set
	receivedAddress = noNdfInstance.GetPermissioningAddress()
	if receivedAddress != "" {
		t.Errorf("GetPermissioningAddress did not get expected value!"+
			"No ndf set, address should be an empty string. "+
			"\n\tReceived: %+v", receivedAddress)
	}

}

// Happy path: Tests GetPermissioningCert with the full ndf set, the partial ndf set
// and no ndf set
func TestInstance_GetPermissioningCert(t *testing.T) {
	// Create an instance object, setting the full ndf
	secured, _ := NewSecuredNdf(testutils.NDF)
	fullNdfInstance := Instance{
		full: secured,
	}
	// Expected cert gotten from testutils.NDF
	expectedCert := "-----BEGIN CERTIFICATE-----\nMIIDkDCCAnigAwIBAgIJAJnjosuSsP7gMA0GCSqGSIb3DQEBBQUAMHQxCzAJBgNV\nBAYTAlVTMRMwEQYDVQQIDApDYWxpZm9ybmlhMRIwEAYDVQQHDAlDbGFyZW1vbnQx\nGzAZBgNVBAoMElByaXZhdGVncml0eSBDb3JwLjEfMB0GA1UEAwwWcmVnaXN0cmF0\naW9uKi5jbWl4LnJpcDAeFw0xOTAzMDUyMTQ5NTZaFw0yOTAzMDIyMTQ5NTZaMHQx\nCzAJBgNVBAYTAlVTMRMwEQYDVQQIDApDYWxpZm9ybmlhMRIwEAYDVQQHDAlDbGFy\nZW1vbnQxGzAZBgNVBAoMElByaXZhdGVncml0eSBDb3JwLjEfMB0GA1UEAwwWcmVn\naXN0cmF0aW9uKi5jbWl4LnJpcDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoC\nggEBAOQKvqjdh35o+MECBhCwopJzPlQNmq2iPbewRNtI02bUNK3kLQUbFlYdzNGZ\nS4GYXGc5O+jdi8Slx82r1kdjz5PPCNFBARIsOP/L8r3DGeW+yeJdgBZjm1s3ylka\nmt4Ajiq/bNjysS6L/WSOp+sVumDxtBEzO/UTU1O6QRnzUphLaiWENmErGvsH0CZV\nq38Ia58k/QjCAzpUcYi4j2l1fb07xqFcQD8H6SmUM297UyQosDrp8ukdIo31Koxr\n4XDnnNNsYStC26tzHMeKuJ2Wl+3YzsSyflfM2YEcKE31sqB9DS36UkJ8J84eLsHN\nImGg3WodFAviDB67+jXDbB30NkMCAwEAAaMlMCMwIQYDVR0RBBowGIIWcmVnaXN0\ncmF0aW9uKi5jbWl4LnJpcDANBgkqhkiG9w0BAQUFAAOCAQEAF9mNzk+g+o626Rll\nt3f3/1qIyYQrYJ0BjSWCKYEFMCgZ4JibAJjAvIajhVYERtltffM+YKcdE2kTpdzJ\n0YJuUnRfuv6sVnXlVVugUUnd4IOigmjbCdM32k170CYMm0aiwGxl4FrNa8ei7AIa\nx/s1n+sqWq3HeW5LXjnoVb+s3HeCWIuLfcgrurfye8FnNhy14HFzxVYYefIKm0XL\n+DPlcGGGm/PPYt3u4a2+rP3xaihc65dTa0u5tf/XPXtPxTDPFj2JeQDFxo7QRREb\nPD89CtYnwuP937CrkvCKrL0GkW1FViXKqZY9F5uhxrvLIpzhbNrs/EbtweY35XGL\nDCCMkg==\n-----END CERTIFICATE-----"

	// GetPermissioningCert from the instance and compare with the expected value
	receivedCert := fullNdfInstance.GetPermissioningCert()
	if expectedCert != receivedCert {
		t.Errorf("GetPermissioningCert did not get expected value!"+
			"\n\tExpected: %+v"+
			"\n\tReceived: %+v", expectedCert, receivedCert)
	}

	// Create an instance object, setting the partial ndf
	partialNdfInstance := Instance{
		partial: secured,
	}

	// GetPermissioningCert from the instance and compare with the expected value
	receivedCert = partialNdfInstance.GetPermissioningCert()
	if expectedCert != receivedCert {
		t.Errorf("GetPermissioningCert did not get expected value!"+
			"\n\tExpected: %+v"+
			"\n\tReceived: %+v", expectedCert, receivedCert)
	}

	// Create an instance object, setting no ndf
	noNdfInstance := Instance{}

	// GetPermissioningCert, should be an empty string as no ndf's are set
	receivedCert = noNdfInstance.GetPermissioningCert()
	if receivedCert != "" {
		t.Errorf("GetPermissioningCert did not get expected value!"+
			"No ndf set, cert should be an empty string. "+
			"\n\tReceived: %+v", receivedCert)
	}

}

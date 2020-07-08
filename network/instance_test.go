///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package network

import (
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/elixxir/comms/mixmessages"
	ds "gitlab.com/elixxir/comms/network/dataStructures"
	"gitlab.com/elixxir/comms/testkeys"
	"gitlab.com/elixxir/comms/testutils"
	"gitlab.com/elixxir/crypto/signature"
	"gitlab.com/elixxir/crypto/signature/rsa"
	"gitlab.com/elixxir/primitives/id"
	"gitlab.com/elixxir/primitives/ndf"
	"reflect"
	"strings"
	"testing"
)

// Happy path
func TestNewInstanceTesting(t *testing.T) {
	_, err := NewInstanceTesting(&connect.ProtoComms{}, testutils.NDF, testutils.NDF, nil, nil, t)
	if err != nil {
		t.Errorf("Unable to create test instance: %+v", err)
	}
}

// Error path: pass in a non testing argument into the constructor
func TestNewInstanceTesting_Error(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	_, err := NewInstanceTesting(&connect.ProtoComms{}, testutils.NDF, testutils.NDF, nil, nil, nil)
	if err != nil {
		return
	}

	t.Errorf("Expected error case, should not be able to create instance when testing argument is nil")

}

//tests newInstance errors properly when there is no NDF
func TestNewInstance_NilNDFs(t *testing.T) {
	_, err := NewInstance(&connect.ProtoComms{}, nil, nil)
	if err == nil {
		t.Errorf("Creation of NewInstance without an ndf succeded")
	} else if !strings.Contains(err.Error(), "Cannot create a network "+
		"instance without an NDF") {
		t.Errorf("Creation of NewInstance without an ndf returned "+
			"the wrong error: %s", err.Error())
	}
}

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
		roundData: ds.NewData(),
	}
	_ = i.roundData.UpsertRound(&mixmessages.RoundInfo{ID: uint64(1)})
	r, err := i.GetRound(id.Round(1))
	if err != nil || r == nil {
		t.Errorf("Failed to retrieve round: %+v", err)
	}
}

func TestInstance_GetRoundUpdate(t *testing.T) {
	i := Instance{
		roundUpdates: ds.NewUpdates(),
	}
	_ = i.roundUpdates.AddRound(&mixmessages.RoundInfo{ID: uint64(1), UpdateID: uint64(1)})
	r, err := i.GetRoundUpdate(1)
	if err != nil || r == nil {
		t.Errorf("Failed to retrieve round update: %+v", err)
	}
}

func TestInstance_GetRoundUpdates(t *testing.T) {
	i := Instance{
		roundUpdates: ds.NewUpdates(),
	}
	_ = i.roundUpdates.AddRound(&mixmessages.RoundInfo{ID: uint64(1), UpdateID: uint64(1)})
	_ = i.roundUpdates.AddRound(&mixmessages.RoundInfo{ID: uint64(1), UpdateID: uint64(2)})
	r := i.GetRoundUpdates(1)
	if r == nil {
		t.Errorf("Failed to retrieve round updates")
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
	f.Ndf = []byte(testutils.ExampleJSON)
	baseNDF := testutils.NDF

	if err != nil {
		t.Errorf("Could not generate serialized ndf: %s", err)
	}

	err = signature.Sign(f, privKey)

	pc := &connect.ProtoComms{}
	i, err := NewInstance(pc, baseNDF, baseNDF)
	if err != nil {
		t.Error(nil)
	}

	_, err = i.comm.AddHost(&id.Permissioning, "0.0.0.0:4200", pub, false, true)
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

	_, err = i.comm.AddHost(&id.Permissioning, "0.0.0.0:4200", pub, false, true)
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

func TestInstance_UpdateFullNdf_nil(t *testing.T) {
	i, f := setupComm(t)
	i.full = nil

	err := i.UpdateFullNdf(f)
	if err == nil {
		t.Errorf("Full NDF update succeded when it shouldnt")
	} else if !strings.Contains(err.Error(),
		"Cannot update the full ndf when it is nil") {
		t.Errorf("Full NDF update when nil failed incorrectly: %s",
			err.Error())
	}
}

func TestInstance_UpdatePartialNdf(t *testing.T) {
	i, f := setupComm(t)
	err := i.UpdatePartialNdf(f)
	if err != nil {
		t.Errorf("Failed to update ndf: %+v", err)
	}
}

func TestInstance_UpdatePartialNdf_nil(t *testing.T) {
	i, f := setupComm(t)
	i.partial = nil

	err := i.UpdatePartialNdf(f)
	if err == nil {
		t.Errorf("Partial NDF update succeded when it shouldnt")
	} else if !strings.Contains(err.Error(),
		"Cannot update the partial ndf when it is nil") {
		t.Errorf("Partial NDF update when nil failed incorrectly: %s",
			err.Error())
	}
}

func TestInstance_GetLastRoundID(t *testing.T) {
	i := Instance{
		roundData: ds.NewData(),
	}
	_ = i.roundData.UpsertRound(&mixmessages.RoundInfo{ID: uint64(1)})
	i.GetLastRoundID()
}

func TestInstance_GetLastUpdateID(t *testing.T) {
	i := Instance{
		roundUpdates: ds.NewUpdates(),
	}
	_ = i.roundUpdates.AddRound(&mixmessages.RoundInfo{ID: uint64(1), UpdateID: uint64(1)})
	i.GetLastUpdateID()
}

func TestInstance_UpdateGatewayConnections(t *testing.T) {
	secured, _ := NewSecuredNdf(testutils.NDF)

	i := Instance{
		full:       secured,
		comm:       &connect.ProtoComms{},
		ipOverride: ds.NewIpOverrideList(),
	}
	err := i.UpdateGatewayConnections()
	if err != nil {
		t.Errorf("Failed to update gateway connections from full: %+v", err)
	}

	i = Instance{
		partial:    secured,
		comm:       &connect.ProtoComms{},
		ipOverride: ds.NewIpOverrideList(),
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

// Tests that UpdateGatewayConnections() returns an error when a Gateway ID
// collides with a hard coded ID.
func TestInstance_UpdateGatewayConnections_GatewayIdError(t *testing.T) {
	testDef := *testutils.NDF
	testDef.Nodes = []ndf.Node{{ID: id.TempGateway.Marshal()}}
	secured, _ := NewSecuredNdf(&testDef)

	i := Instance{
		full:       secured,
		comm:       &connect.ProtoComms{},
		ipOverride: ds.NewIpOverrideList(),
	}
	err := i.UpdateGatewayConnections()
	if err == nil {
		t.Errorf("UpdateGatewayConnections() failed to produce an error when " +
			"the Gateway ID collides with a hard coded ID.")
	}
}

func TestInstance_UpdateNodeConnections(t *testing.T) {
	secured, _ := NewSecuredNdf(testutils.NDF)

	i := Instance{
		full:       secured,
		comm:       &connect.ProtoComms{},
		ipOverride: ds.NewIpOverrideList(),
	}
	err := i.UpdateNodeConnections()
	if err != nil {
		t.Errorf("Failed to update node connections from full: %+v", err)
	}

	i = Instance{
		partial:    secured,
		comm:       &connect.ProtoComms{},
		ipOverride: ds.NewIpOverrideList(),
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

// Tests that UpdateNodeConnections() returns an error when a Node ID collides
// with a hard coded ID.
func TestInstance_UpdateNodeConnections_NodeIdError(t *testing.T) {
	testDef := *testutils.NDF
	testDef.Nodes = []ndf.Node{{ID: id.Permissioning.Marshal()}}
	secured, _ := NewSecuredNdf(&testDef)

	i := Instance{
		full:       secured,
		comm:       &connect.ProtoComms{},
		ipOverride: ds.NewIpOverrideList(),
	}
	err := i.UpdateNodeConnections()
	if err == nil {
		t.Errorf("UpdateNodeConnections() failed to produce an error when the " +
			"Node ID collides with a hard coded ID.")
	}
}

// Happy path: Tests GetPermissioningAddress with the full ndf set, the partial ndf set
// and no ndf set
func TestInstance_GetPermissioningAddress(t *testing.T) {
	// Create populated ndf (secured) and empty ndf
	secured, _ := NewSecuredNdf(testutils.NDF)

	// Create an instance object, setting full to be populated
	// and partial to be empty
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

	// Create an instance object, setting partial to be populated
	// and full to be empty
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

// Happy path
func TestInstance_GetCmixGroup(t *testing.T) {
	expectedGroup := ds.NewGroup()

	i := Instance{
		cmixGroup: expectedGroup,
	}

	receivedGroup := i.GetCmixGroup()

	if !reflect.DeepEqual(expectedGroup.Get(), receivedGroup) {
		t.Errorf("Getter didn't get expected value! "+
			"\n\tExpected: %+v"+
			"\n\tReceived: %+v", expectedGroup, receivedGroup)
	}

}

// Happy path
func TestInstance_GetE2EGroup(t *testing.T) {
	expectedGroup := ds.NewGroup()

	i := Instance{
		e2eGroup: expectedGroup,
	}

	receivedGroup := i.GetE2EGroup()

	if !reflect.DeepEqual(expectedGroup.Get(), receivedGroup) {
		t.Errorf("Getter didn't get expected value! "+
			"\n\tExpected: %+v"+
			"\n\tReceived: %+v", expectedGroup, receivedGroup)
	}
}

// Happy path: Tests GetPermissioningCert with the full ndf set, the partial ndf set
// and no ndf set
func TestInstance_GetPermissioningCert(t *testing.T) {

	// Create populated ndf (secured) and empty ndf
	secured, _ := NewSecuredNdf(testutils.NDF)
	// Create an instance object, setting full to be populated
	// and partial to be empty
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

	// Create an instance object, setting partial to be populated
	// and full to be empty
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

// Error path: nil ndf is in the instance should cause a seg fault
func TestInstance_GetPermissioningAddress_NilCase(t *testing.T) {
	// Handle expected seg fault here
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected error case, should seg fault when a nil ndf is passed through")
		}
	}()

	// Create a nil ndf
	nilNdf, _ := NewSecuredNdf(nil)

	// Create an instance object with this nil ndf
	nilNdfInstance := Instance{
		full:    nilNdf,
		partial: nilNdf,
	}

	// Attempt to call getter, should seg fault
	nilNdfInstance.GetPermissioningAddress()
}

// Error path: nil ndf is in the instance should cause a seg fault
func TestInstance_GetPermissioningCert_NilCase(t *testing.T) {
	// Handle expected seg fault here
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected error case, should seg fault when a nil ndf is passed through")
		}
	}()

	// Create a nil ndf
	nilNdf, _ := NewSecuredNdf(nil)

	// Create an instance object with this nil ndf
	nilNdfInstance := Instance{
		full:    nilNdf,
		partial: nilNdf,
	}

	// Attempt to call getter, should seg fault
	nilNdfInstance.GetPermissioningCert()
}

// GetPermissioningId should fetch the value of id.PERMISSIONING in primitives
func TestInstance_GetPermissioningId(t *testing.T) {
	// Create an instance object,
	instance := Instance{}

	receivedId := instance.GetPermissioningId()

	if receivedId != &id.Permissioning {
		t.Errorf("GetPermissioningId did not get value from primitives"+
			"\n\tExpected: %+v"+
			"\n\tReceived: %+v", id.Permissioning, receivedId)
	}
}

// Happy path
func TestInstance_UpdateGroup(t *testing.T) {
	i, f := setupComm(t)
	err := i.UpdateFullNdf(f)
	if err != nil {
		t.Errorf("Unable to initalize group: %+v", err)
	}

	// Update with same values should not cause an error
	err = i.UpdateFullNdf(f)
	if err != nil {
		t.Errorf("Unable to call update group with same values: %+v", err)
	}

}

// Error path: attempt to modify group once already initialized
func TestInstance_UpdateGroup_Error(t *testing.T) {
	i, f := setupComm(t)

	err := i.UpdateFullNdf(f)
	if err != nil {
		t.Errorf("Unable to initalize group: %+v", err)
	}

	badNdf := createBadNdf(t)

	// Update with same values should not cause an error
	err = i.UpdateFullNdf(badNdf)
	if err != nil {
		return
	}

	t.Errorf("Expected error case: Should not be able to update instance's group once initialized!")

}

// Creates a bad ndf
func createBadNdf(t *testing.T) *mixmessages.NDF {
	f := &mixmessages.NDF{}

	badGrp := ndf.Group{
		Prime:      "123",
		SmallPrime: "456",
		Generator:  "2",
	}

	baseNDF := ndf.NetworkDefinition{
		E2E:  badGrp,
		CMIX: badGrp,
	}

	var err error
	f.Ndf, err = baseNDF.Marshal()
	if err != nil {
		t.Errorf("Could not generate serialized ndf: %s", err)
	}
	priv := testkeys.LoadFromPath(testkeys.GetNodeKeyPath())
	privKey, err := rsa.LoadPrivateKeyFromPem(priv)
	if err != nil {
		t.Errorf("Could not generate rsa key: %s", err)
	}

	err = signature.Sign(f, privKey)

	return f
}

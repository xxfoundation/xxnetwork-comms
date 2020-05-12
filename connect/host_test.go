////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package connect

import (
	"bytes"
	"testing"
)

func TestHost_address(t *testing.T) {
	var mgr Manager
	testId := "test"
	testAddress := "test"
	host, err := mgr.AddHost(testId, testAddress, nil, false, false)
	if err != nil {
		t.Errorf("Unable to add host")
		return
	}

	if host.address != testAddress {
		t.Errorf("Expected addresses to match")
	}
}

func TestHost_GetCertificate(t *testing.T) {
	testCert := []byte("TEST")

	host := Host{
		address:      "",
		certificate:  testCert,
		maxRetries:   0,
		connection:   nil,
		credentials:  nil,
		rsaPublicKey: nil,
	}

	if bytes.Compare(host.certificate, testCert) != 0 {
		t.Errorf("Expected certs to match!")
	}
}

// Tests that getID returns the correct ID
func TestHost_GetId(t *testing.T) {

	id := "xXx_420No1337ScopeH4xx0r_xXx"

	host := Host{
		id: id,
	}

	if host.GetId() != id {
		t.Errorf("Correct id not returned.  Expected %s, got %s",
			id, host.id)
	}
}

// Tests that GetAddress() returns the address of the host.
func TestHost_GetAddress(t *testing.T) {
	// Test values
	testAddress := "192.167.1.1:8080"
	testHost := Host{address: testAddress}

	// Test function
	if testHost.GetAddress() != testAddress {
		t.Errorf("GetAddress() did not return the expected address."+
			"\n\texpected: %v\n\treceived: %v",
			testAddress, testHost.GetAddress())
	}
}

// Validate that dynamic host defaults to false and can be set to true
func TestHost_IsDynamicHost(t *testing.T) {

	host := Host{}

	if host.IsDynamicHost() != false {
		t.Errorf("Correct bool not returned. Expected false, got %v",
			host.dynamicHost)
	}

	host.dynamicHost = true

	if host.IsDynamicHost() != true {
		t.Errorf("Correct bool not returned. Expected true, got %v",
			host.dynamicHost)
	}
}

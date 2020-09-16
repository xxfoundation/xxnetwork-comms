///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package connect

import (
	"bytes"
	"gitlab.com/elixxir/primitives/id"
	"net"
	"testing"
)

func TestHost_address(t *testing.T) {
	mgr := newManager()
	testId := id.NewIdFromString("test", id.Node, t)
	testAddress := "test"
	host, err := mgr.AddHost(testId, testAddress, nil, GetDefaultHostParams())
	if err != nil {
		t.Errorf("Unable to add host")
		return
	}

	if host.GetAddress() != testAddress {
		t.Errorf("Expected addresses to match")
	}
}

func TestHost_GetCertificate(t *testing.T) {
	testCert := []byte("TEST")

	host := Host{
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

	testID := id.NewIdFromString("xXx_420No1337ScopeH4xx0r_xXx", id.Generic, t)

	host := Host{
		id: testID,
	}

	if host.GetId() != testID {
		t.Errorf("Correct id not returned.  Expected %s, got %s",
			testID, host.id)
	}
}

// Tests that GetAddress() returns the address of the host.
func TestHost_GetAddress(t *testing.T) {
	// Test values
	testAddress := "192.167.1.1:8080"
	testHost := Host{}
	testHost.UpdateAddress(testAddress)

	// Test function
	if testHost.GetAddress() != testAddress {
		t.Errorf("GetAddress() did not return the expected address."+
			"\n\texpected: %v\n\treceived: %v",
			testAddress, testHost.GetAddress())
	}
}

func TestHost_UpdateAddress(t *testing.T) {
	testAddress := "192.167.1.1:8080"
	testUpdatedAddress := "192.167.1.1:8080"
	testHost := Host{}
	testHost.UpdateAddress(testAddress)

	// Test function
	if testHost.GetAddress() != testAddress {
		t.Errorf("GetAddress() did not return the expected address before update."+
			"\n\texpected: %v\n\treceived: %v",
			testAddress, testHost.GetAddress())
	}

	testHost.UpdateAddress(testUpdatedAddress)

	if testHost.GetAddress() != testUpdatedAddress {
		t.Errorf("GetAddress() did not return the expected address after update."+
			"\n\texpected: %v\n\treceived: %v",
			testUpdatedAddress, testHost.GetAddress())
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

// Full test
func TestHost_IsOnline(t *testing.T) {
	addr := "0.0.0.0:10234"

	// Create the host
	host, err := NewHost(id.NewIdFromString("test", id.Gateway, t), addr, nil,
		GetDefaultHostParams())
	if err != nil {
		t.Errorf("Unable to create host: %+v", host)
		return
	}

	// Test that host is offline
	if host.IsOnline() {
		t.Errorf("Expected host to be offline!")
	}

	// Listen on the given address
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		t.Errorf("Unable to listen: %+v", err)
	}

	// Test that host is online
	if !host.IsOnline() {
		t.Errorf("Expected host to be online!")
	}

	// Close listening address
	err = lis.Close()
	if err != nil {
		t.Errorf("Unable to close listening server: %+v", err)
	}
}

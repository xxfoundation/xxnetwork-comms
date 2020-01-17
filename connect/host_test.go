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

func TestHost_GetAddress(t *testing.T) {
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

//tests that getID returns the correct ID
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

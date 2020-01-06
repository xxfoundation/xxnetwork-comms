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
	host, err := mgr.AddHost(testId, testAddress, nil, false)
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
		address:      "",
		certificate:  testCert,
		maxRetries:   0,
		connection:   nil,
		credentials:  nil,
		rsaPublicKey: nil,
	}

	if bytes.Compare(host.GetCertificate(), testCert) != 0 {
		t.Errorf("Expected certs to match!")
	}
}

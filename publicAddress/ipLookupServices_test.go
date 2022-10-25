////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package publicAddress

import (
	"net/url"
	"testing"
)

// // Tests that each lookup Service returns a valid IPv4 address.
// // Note: keep this test disabled so we don't spam the servers.
// func Test_lookupServices_ValidURLs(t *testing.T) {
// 	var myLastIp string
// 	var lastIpSet bool
// 	for _, address := range lookupServices {
// 		ip, err := getIP(address.url, connectionTimeout)
// 		if err != nil {
// 			t.Errorf("%s failed to return an IP: %+v", address, err)
// 			continue
// 		}
//
// 		if net.ParseIP(ip) == nil {
// 			t.Errorf("The IP returned by %s is invalid: %s", address, ip)
// 			continue
// 		}
//
// 		if myLastIp != ip && lastIpSet {
// 			t.Errorf("The IP returned by %s is does not match other IPs returned."+
// 				"\nexpected: %s\nreceived: %s", address, myLastIp, ip)
// 		} else {
// 			lastIpSet = true
// 			myLastIp = ip
// 		}
// 	}
// }

// Tests that there is no duplicate URLs in the lookupServices list.
func Test_lookupServices_Duplicate(t *testing.T) {
	hostMap := map[string]string{}
	for _, service := range lookupServices {
		u, err := url.Parse(service.url)
		if err != nil {
			t.Errorf("Failed to parse address: %+v", err)
		}

		if _, exists := hostMap[u.Hostname()]; exists {
			t.Errorf("Address with hostname %s already exists."+
				"\nexists:  %s\ncurrent: %s",
				u.Hostname(), hostMap[u.Hostname()], service.url)
		} else {
			hostMap[u.Hostname()] = service.url
		}
	}
}

// Tests that MakeTestLookupService creates a test server that responds with the
// expected IP address.
func TestMakeTestLookupService(t *testing.T) {
	expectedIp := "0.0.0.0"
	list, ts := MakeTestLookupService(expectedIp, t)
	defer ts.Close()

	testIp, err := getIP(list[0].url, connectionTimeout)
	if err != nil {
		t.Errorf("Failed to get IP from test server.")
	}

	if expectedIp != testIp {
		t.Errorf("Test server did not return the expected IP."+
			"\nexpected: %s\nreceived: %s", expectedIp, testIp)
	}
}

// Panic path: tests that MakeTestLookupService panics when the provided
// interface is not for testing.
func TestMakeTestLookupService_InterfaceTypePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Failed to panic when provided interface is not for testing.")
		}
	}()

	var i interface{}
	_, _ = MakeTestLookupService("", i)
}

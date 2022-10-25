////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package publicAddress

import (
	"net"
	"strconv"
	"strings"
	"testing"
)

// Happy path.
func TestGetIpOverride(t *testing.T) {
	host, port := "0.0.0.0", 11420
	expectedIP := net.JoinHostPort(host, strconv.Itoa(port))

	testIP, err := GetIpOverride(host, port)
	if err != nil {
		t.Errorf("GetIpOverride() returned an error: %+v", err)
	}

	if expectedIP != testIP {
		t.Errorf("GetIpOverride() did not return the expected IP."+
			"\nexpected: %s\nreceived: %s", expectedIP, testIP)
	}
}

// Happy path: no override IP is provided so it is looked up and combined with
// the given port.
func Test_getIpOverride_IpLookup(t *testing.T) {
	ip, port := "0.0.0.0", 11420
	expectedIP := net.JoinHostPort(ip, strconv.Itoa(port))
	testServices, ts := MakeTestLookupService(ip, t)
	defer ts.Close()

	testIP, err := getIpOverride("", port, testServices)
	if err != nil {
		t.Errorf("getIpOverride() returned an error: %+v", err)
	}

	if expectedIP != testIP {
		t.Errorf("getIpOverride() did not return the expected IP."+
			"\nexpected: %s\nreceived: %s", expectedIP, testIP)
	}
}

// Error path: IP address lookup fails.
func Test_getIpOverride_IpLookupError(t *testing.T) {
	testServices, ts := MakeTestLookupService("invalid IP", t)
	defer ts.Close()

	_, err := getIpOverride("", 11420, testServices)
	if err == nil || !strings.Contains(err.Error(), "lookup public IP address") {
		t.Errorf("getIpOverride() did not return an error for an invalid IP "+
			"response: %+v", err)
	}
}

// Happy path: an override IP is provided without a port so it is combined with
// the given port.
func TestJoinIpPort_IpOverride(t *testing.T) {
	ip, port := "1.1.1.1", 11420
	expectedIP := net.JoinHostPort(ip, strconv.Itoa(port))
	testIP, err := JoinIpPort(ip, port)
	if err != nil {
		t.Errorf("JoinIpPort() returned an error: %+v", err)
	}

	if expectedIP != testIP {
		t.Errorf("JoinIpPort() did not return the expected IP."+
			"\nexpected: %s\nreceived: %s", expectedIP, testIP)
	}
}

// Happy path: an override IP is provided with a port so it is returned as is.
func TestJoinIpPort_IpPortOverride(t *testing.T) {
	expectedIP := net.JoinHostPort("1.1.1.1", strconv.Itoa(22840))

	testIP, err := JoinIpPort(expectedIP, 11420)
	if err != nil {
		t.Errorf("JoinIpPort() returned an error: %+v", err)
	}

	if expectedIP != testIP {
		t.Errorf("JoinIpPort() did not return the expected IP."+
			"\nexpected: %s\nreceived: %s", expectedIP, testIP)
	}
}

// Error path: override IP is invalid.
func TestJoinIpPort_OverrideIpError(t *testing.T) {
	_, err := JoinIpPort("0.0.0.0::100", 11420)
	if err == nil || !strings.Contains(err.Error(), "parse public IP address") {
		t.Errorf("JoinIpPort() did not return an error for an invalid "+
			"override IP : %+v", err)
	}
}

// Happy path: supplying an empty IP results in an empty return.
func TestJoinIpPort_NoIp(t *testing.T) {
	ip, err := JoinIpPort("", 11420)
	if err != nil {
		t.Errorf("JoinIpPort() returned an error: %+v", err)
	}

	if ip != "" {
		t.Errorf("JoinIpPort() did not return the expected IP."+
			"\nexpected: %s\nreceived: %s", "", ip)
	}
}

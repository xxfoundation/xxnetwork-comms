///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package publicAddress

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

// Happy path.
func TestGetIP(t *testing.T) {
	// Create test servers
	expectedIp := "0.0.0.0"
	ts0 := newTestServer(expectedIp, t)
	ts1 := newTestServer(expectedIp, t)
	ts2 := newTestServer(expectedIp, t)
	defer func() {
		ts0.Close()
		ts1.Close()
		ts2.Close()
	}()

	lookupServices := []Service{
		{ipv4Address, ts0.URL},
		{ipv4Address, ts1.URL},
		{ipv4Address, ts2.URL},
	}
	ip, err := GetIP(lookupServices, DefaultPollTimeout)
	if err != nil {
		t.Errorf("GetIP() produced an error: %+v", err)
	}

	if expectedIp != ip {
		t.Errorf("GetIP() did not return the expected IP address."+
			"\nexpected: %s\nreceived: %s", expectedIp, ip)
	}
}

// Error path: returns an error when all servers report an IPv6 address.
func TestGetIP_IPv6Error(t *testing.T) {
	// Create test servers
	expectedIp := "0000:0000:0000:0000:0000:0000:0000:0000"
	ts0 := newTestServer(expectedIp, t)
	ts1 := newTestServer(expectedIp, t)
	ts2 := newTestServer(expectedIp, t)
	defer func() {
		ts0.Close()
		ts1.Close()
		ts2.Close()
	}()

	lookupServices := []Service{
		{ipv6Address, ts0.URL},
		{ipv6Address, ts1.URL},
		{ipv6Address, ts2.URL},
	}
	_, err := GetIP(lookupServices, DefaultPollTimeout)
	if err == nil {
		t.Errorf("GetIP() did not error when all servers return IPv6 "+
			"addresses: %+v", err)
	}
}

// Error path: functions returns a timeout error if servers take too long too
// respond.
func TestGetIP_TimeoutError(t *testing.T) {
	// Create test servers
	timeout := 50 * time.Millisecond
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * timeout)
	}))
	defer ts.Close()

	lookupServices := []Service{{ipv4Address, ts.URL}}
	_, err := GetIP(lookupServices, timeout)
	if err == nil || !strings.Contains(err.Error(), "timed out") {
		t.Errorf("GetIP() did not timeout: %+v", err)
	}
}

// Happy path.
func Test_getIpFromList(t *testing.T) {
	// Create test servers
	expectedIp := "0.0.0.0"
	ts0 := newTestServer(expectedIp, t)
	ts1 := newTestServer(expectedIp, t)
	ts2 := newTestServer(expectedIp, t)
	defer func() {
		ts0.Close()
		ts1.Close()
		ts2.Close()
	}()

	urls := []Service{
		{ipv4Address, ts0.URL},
		{ipv4Address, ts1.URL},
		{ipv4Address, ts2.URL},
	}
	ip, err := getIpFromList(urls, 50*time.Millisecond)
	if err != nil {
		t.Errorf("getIpFromList() produced an error: %+v", err)
	}

	if expectedIp != ip {
		t.Errorf("getIpFromList() did not return the expected IP address."+
			"\nexpected: %s\nreceived: %s", expectedIp, ip)
	}
}

// Error path: none of the URLs return valid IPs.
func Test_getIpFromList_NoIpError(t *testing.T) {
	// Create test servers
	expectedIp := "invalid IP"
	ts0 := newTestServer(expectedIp, t)
	ts1 := newTestServer(expectedIp, t)
	ts2 := newTestServer(expectedIp, t)
	defer func() {
		ts0.Close()
		ts1.Close()
		ts2.Close()
	}()

	urls := []Service{
		{ipv4Address, ts0.URL},
		{ipv4Address, ts1.URL},
		{ipv4Address, ts2.URL},
	}
	_, err := getIpFromList(urls, 50*time.Millisecond)
	if err == nil || !strings.Contains(err.Error(), "failed to get public IP address") {
		t.Errorf("getIpFromList() failed to return an error when all servers "+
			"should have returned invalid IP addresses: %+v", err)
	}
}

// Happy path.
func Test_getIP(t *testing.T) {
	// Create test server
	expectedIp := "0.0.0.0"
	ts := newTestServer(expectedIp, t)
	defer ts.Close()

	ip, err := getIP(ts.URL, 50*time.Millisecond)
	if err != nil {
		t.Errorf("getIP() produced an error: %+v", err)
	}

	if expectedIp != ip {
		t.Errorf("getIP() did not return the expected IP address."+
			"\nexpected: %s\nreceived: %s", expectedIp, ip)
	}
}

// Error path: Get returns an error for an invalid URL.
func Test_getIP_GetError(t *testing.T) {
	_, err := getIP("http://invalidurl", 50*time.Millisecond)
	if err == nil || !strings.Contains(err.Error(), "Get") {
		t.Errorf("getIP() did not produce an error for an invalid URL: %+v", err)
	}
}

// Error path: the response does not contain a valid IP address.
func Test_getIP_IpError(t *testing.T) {
	// Create test server
	ts := newTestServer("invalid IP", t)
	defer ts.Close()

	_, err := getIP(ts.URL, 50*time.Millisecond)
	if err == nil || !strings.Contains(err.Error(), "response could not be parsed as an IP address") {
		t.Errorf("getIP() did not produce an error for an invalid URL: %+v", err)
	}
}

// Error path: returns an error when the address returned is IPv6.
func Test_getIP_IPv6Error(t *testing.T) {
	// Create test server
	ts := newTestServer("0000:0000:0000:0000:0000:0000:0000:0000", t)
	defer ts.Close()

	_, err := getIP(ts.URL, 50*time.Millisecond)
	if err == nil || !strings.Contains(err.Error(), "IPv6") {
		t.Errorf("getIP() did not produce an error for an IPv6 address: %+v", err)
	}
}

// Happy path.
func Test_shuffleStrings(t *testing.T) {
	s := []Service{
		{ipv4Address, "A"},
		{ipv4Address, "B"},
		{ipv4Address, "C"},
		{ipv4Address, "D"},
		{ipv4Address, "E"},
		{ipv4Address, "F"},
	}
	shuffled := shuffleStrings(s)
	if reflect.DeepEqual(s, shuffled) {
		t.Errorf("shuffleStrings() failed to shuffle the list."+
			"\nlist:     %s\nshuffled: %s", s, shuffled)
	}
}

func newTestServer(response string, t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, response); err != nil {
			t.Errorf("Failed to respond: %+v", err)
		}
	}))
}

////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

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

// Tests that GetIP returns the expected IP from three different test servers.
func TestGetIP(t *testing.T) {
	// Create test servers
	expectedIp := "0.0.0.0"
	ts0 := newTestServer(expectedIp, t)
	ts1 := newTestServer(expectedIp, t)
	ts2 := newTestServer(expectedIp, t)
	defer ts0.Close()
	defer ts1.Close()
	defer ts2.Close()

	lookupServices := []Service{
		{ipv4Address, ts0.URL},
		{ipv4Address, ts1.URL},
		{ipv4Address, ts2.URL},
	}

	ip, err := GetIP(lookupServices, DefaultPollTimeout)
	if err != nil {
		t.Errorf("GetIP produced an error: %+v", err)
	}

	if expectedIp != ip {
		t.Errorf("GetIP did not return the expected IP address."+
			"\nexpected: %s\nreceived: %s", expectedIp, ip)
	}
}

// Error path: tests that GetIP returns an error when all of the test servers
// report an IPv6 address.
func TestGetIP_IPv6Error(t *testing.T) {
	// Create test servers
	expectedErr := lookupServiceErr
	expectedIp := "0000:0000:0000:0000:0000:0000:0000:0000"
	ts0 := newTestServer(expectedIp, t)
	ts1 := newTestServer(expectedIp, t)
	ts2 := newTestServer(expectedIp, t)
	defer ts0.Close()
	defer ts1.Close()
	defer ts2.Close()

	lookupServices := []Service{
		{ipv6Address, ts0.URL},
		{ipv6Address, ts1.URL},
		{ipv6Address, ts2.URL},
	}

	_, err := GetIP(lookupServices, DefaultPollTimeout)
	if err == nil || err.Error() != expectedErr {
		t.Errorf("GetIP did not return the expected error when all servers "+
			"return IPv6 addresses.\nexpected: %s\nreceived: %+v",
			expectedErr, err)
	}
}

// Error path: tests that GetIP returns a timeout error if servers take too long
// to respond.
func TestGetIP_TimeoutError(t *testing.T) {
	// Create test servers
	timeout := 50 * time.Millisecond
	expectedErr := fmt.Sprintf(findIpTimeoutErr, timeout)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * timeout)
	}))
	defer ts.Close()

	lookupServices := []Service{{ipv4Address, ts.URL}}
	_, err := GetIP(lookupServices, timeout)
	if err == nil || err.Error() != expectedErr {
		t.Errorf("GetIP did failed to timeout when the server has not "+
			"responded within the timeout period.\nexpected: %s\nreceived: %+v",
			expectedErr, err)
	}
}

// Tests that getIpFromList retrieves the expected IP from a randomly selected
// service.
func Test_getIpFromList(t *testing.T) {
	// Create test servers
	expectedIp := "0.0.0.0"
	ts0 := newTestServer(expectedIp, t)
	ts1 := newTestServer(expectedIp, t)
	ts2 := newTestServer(expectedIp, t)
	defer ts0.Close()
	defer ts1.Close()
	defer ts2.Close()

	urls := []Service{
		{ipv4Address, ts0.URL},
		{ipv4Address, ts1.URL},
		{ipv4Address, ts2.URL},
	}

	ip, err := getIpFromList(urls, 50*time.Millisecond)
	if err != nil {
		t.Errorf("getIpFromList produced an error: %+v", err)
	}

	if expectedIp != ip {
		t.Errorf("getIpFromList did not return the expected IP address."+
			"\nexpected: %s\nreceived: %s", expectedIp, ip)
	}
}

// Error path: tests that getIpFromList returns an error when none of the
// servers return valid IPs.
func Test_getIpFromList_NoIpError(t *testing.T) {
	// Create test servers
	expectedErr := lookupServiceErr
	expectedIp := "invalid IP"
	ts0 := newTestServer(expectedIp, t)
	ts1 := newTestServer(expectedIp, t)
	ts2 := newTestServer(expectedIp, t)
	defer ts0.Close()
	defer ts1.Close()
	defer ts2.Close()

	urls := []Service{
		{ipv4Address, ts0.URL},
		{ipv4Address, ts1.URL},
		{ipv4Address, ts2.URL},
	}

	_, err := getIpFromList(urls, 50*time.Millisecond)
	if err == nil || err.Error() != expectedErr {
		t.Errorf("getIpFromList failed to return an error when all servers "+
			"should have returned invalid IP addresses."+
			"\nexpected: %s\nreceived: %+v", expectedErr, err)
	}
}

// Tests that getIpMultiCheck retrieves the expected IP from a randomly selected
// service.
func Test_getIpMultiCheck(t *testing.T) {
	// Create test servers
	expectedIp := "0.0.0.0"
	ts0 := newTestServer("invalid IP", t)
	ts1 := newTestServer(expectedIp, t)
	ts2 := newTestServer(expectedIp, t)
	ts3 := newTestServer("0000:0000:0000:0000:0000:0000:0000:0000", t)
	ts4 := newTestServer(expectedIp, t)
	ts5 := newTestServer("1.1.1.1", t)
	defer ts0.Close()
	defer ts1.Close()
	defer ts2.Close()
	defer ts3.Close()
	defer ts4.Close()
	defer ts5.Close()

	urls := []Service{
		{ipv4Address, ts0.URL},
		{ipv4Address, ts1.URL},
		{ipv4Address, ts2.URL},
		{ipv6Address, ts3.URL},
		{ipv4Address, ts4.URL},
		{ipv4Address, ts5.URL},
	}

	ip, err := getIpMultiCheck(urls, 50*time.Millisecond, 3)
	if err != nil {
		t.Errorf("getIpMultiCheck produced an error: %+v", err)
	}

	if expectedIp != ip {
		t.Errorf("getIpMultiCheck did not return the expected IP address."+
			"\nexpected: %s\nreceived: %s", expectedIp, ip)
	}
}

// Tests that getIP gets the expected IP from the test server.
func Test_getIP(t *testing.T) {
	// Create test server
	expectedIp := "0.0.0.0"
	ts := newTestServer(expectedIp, t)
	defer ts.Close()

	ip, err := getIP(ts.URL, 50*time.Millisecond)
	if err != nil {
		t.Errorf("getIP produced an error: %+v", err)
	}

	if expectedIp != ip {
		t.Errorf("getIP did not return the expected IP address."+
			"\nexpected: %s\nreceived: %s", expectedIp, ip)
	}
}

// Error path: tests that getIP returns an error for an invalid URL.
func Test_getIP_GetError(t *testing.T) {
	expectedErr := strings.Split(getServiceErr, "%")[0]

	_, err := getIP("https://invalidURL", 50*time.Millisecond)
	if err == nil || !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("getIP did not produce the expected error for an invalid URL."+
			"\nexpected: %s\nreceived: %+v", expectedErr, err)
	}
}

// Error path: tests that getIP returns an error when the response does not
// contain a valid IP address.
func Test_getIP_IpError(t *testing.T) {
	response := "invalid IP"
	expectedErr := fmt.Sprintf(responseParseErr, response)

	// Create test server
	ts := newTestServer(response, t)
	defer ts.Close()

	_, err := getIP(ts.URL, 50*time.Millisecond)
	if err == nil || err.Error() != expectedErr {
		t.Errorf("getIP did not produce the expected error for an invalid URL."+
			"\nexpected: %s\nreceived: %+v", expectedErr, err)
	}
}

// Error path: tests that getIP returns an error when the address retrieved from
// the test server is an IPv6 address.
func Test_getIP_IPv6Error(t *testing.T) {
	// Create test server
	ipv6Addr := "0000:0000:0000:0000:0000:0000:0000:0000"
	expectedErr := fmt.Sprintf(receivedIPv6Err, ipv6Addr)
	ts := newTestServer(ipv6Addr, t)
	defer ts.Close()

	_, err := getIP(ts.URL, 50*time.Millisecond)
	if err == nil || err.Error() != expectedErr {
		t.Errorf("getIPdid not produce the expected error when it recieves a "+
			"IPv6 address.\nexpected: %s\nreceived: %+v", expectedErr, err)
	}
}

// Tests that shuffleStrings returns a list in a different order from the
// original.
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
		t.Errorf("shuffleStringsfailed to shuffle the list."+
			"\nlist:     %s\nshuffled: %s", s, shuffled)
	}
}

// Tests that trunc returns the expected truncated string.
func Test_trunc(t *testing.T) {
	testValues := []struct {
		str      string
		limit    int
		ellipses bool
		expected string
	}{
		{"testString", 15, false, "testString"},
		{"testString", 15, true, "testString"},
		{"testString", 10, false, "testString"},
		{"testString", 10, true, "testString"},
		{"testString", 9, false, "testStrin"},
		{"testString", 9, true, "testSt..."},
		{"testString", 7, false, "testStr"},
		{"testString", 7, true, "test..."},
		{"testStrin�", 10, false, "testStrin�"},
		{"testStrin�", 10, true, "testStr..."},
		{"testStrin�", 11, false, "testStrin�"},
	}

	for i, val := range testValues {
		truncated := trunc(val.str, val.limit, val.ellipses)
		if truncated != val.expected {
			t.Errorf("trunc did not return the expected truncated string (%d)."+
				"\nexpected: %s\nreceived: %s", i, val.expected, truncated)
		}
	}
}

// newTestServer creates a test server that returns the response.
func newTestServer(response string, t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, response); err != nil {
			t.Errorf("Failed to respond: %+v", err)
		}
	}))
}

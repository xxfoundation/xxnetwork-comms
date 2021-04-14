///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// package publicAddress contains a utility to retrieve the callers public IP
// address.
package publicAddress

import (
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/crypto/shuffle"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	DefaultPollTimeout = 30 * time.Second
	connectionTimeout  = 2 * time.Second
)

// GetIP returns the caller's public IP address. Multiple services are checked
// until one returns an IPv4 address or the timeout occurs.
func GetIP(lookupServices []Service, timeout time.Duration) (string, error) {
	ipResultChan := make(chan struct {
		ip  string
		err error
	})
	go func() {
		ip, err := getIpFromList(lookupServices, connectionTimeout)
		ipResultChan <- struct {
			ip  string
			err error
		}{ip: ip, err: err}
	}()

	timer := time.NewTimer(timeout)

	select {
	case results := <-ipResultChan:
		return results.ip, results.err
	case <-timer.C:
		return "", errors.Errorf("retreiving public IP address failed: timed "+
			"out after %s.", timeout)
	}
}

// getIpFromList returns the caller's public IP address by making a request to a
// randomly selected Service. If a Service fails to return a valid IPv4 address,
// then the next Service in the list is used.
func getIpFromList(urls []Service, timeout time.Duration) (string, error) {
	serviceList := shuffleStrings(urls)
	var ipv6 bool

	for _, service := range serviceList {
		// Skip services that return IPv6 addresses if the caller has an IPv6
		if ipv6 && service.v == ipv6Address {
			continue
		}

		ip, err := getIP(service.url, timeout)
		if err == nil {
			return ip, nil
		} else if strings.Contains(err.Error(), "IPv6") {
			ipv6 = true
		}

		jww.ERROR.Printf("Failed to get public IP address: %+v", err)
	}

	return "", errors.New("failed to get public IP address because no lookup " +
		"sources returned a valid IPv4 address")
}

// getIP requests the our public IP from the specified URL and returns it.
// Only valid IPv4 addresses are returned. An error is returned for IPv6 or
// invalid IP addresses.
func getIP(url string, timeout time.Duration) (string, error) {
	jww.INFO.Printf("Getting public IP address from %s", url)

	// Issue a GET to the URL with the specified timeout
	httpClient := http.Client{
		Timeout: timeout,
	}
	resp, err := httpClient.Get(url)
	if err != nil {
		return "", err
	}

	defer func() {
		if err = resp.Body.Close(); err != nil {
			jww.ERROR.Print(err)
		}
	}()

	// Read the response
	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Errorf("Failed to read response from %s: %+v\n", url, err)
	}

	// Ensure the response is a valid IPv4 address
	parsedIp := net.ParseIP(strings.TrimSpace(string(ip)))
	if parsedIp == nil {
		return "", errors.Errorf("response could not be parsed as an IP "+
			"address: \"%s\"", strings.ReplaceAll(string(ip[:128]), "\n", "\\n"))
	} else if strings.Contains(parsedIp.String(), ":") {
		return "", errors.Errorf("received IPv6 address instead of IPv4: %s", ip)
	}

	jww.INFO.Printf("Got public IP address: %s", parsedIp.String())

	return parsedIp.String(), nil
}

// shuffleStrings shuffles the list of strings.
func shuffleStrings(s []Service) []Service {
	// Create list of indexes to be shuffled
	indexList := make([]uint64, len(s))
	for i := range indexList {
		indexList[i] = uint64(i)
	}

	// Shuffle the index
	shuffle.Shuffle(&indexList)

	// Reorder list of strings from the new shuffled order
	shuffled := make([]Service, len(s))
	for i, j := range indexList {
		shuffled[i] = s[j]
	}

	return shuffled
}

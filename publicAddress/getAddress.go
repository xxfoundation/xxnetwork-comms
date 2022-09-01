////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Package publicAddress contains a utility to retrieve the callers public IP
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
	DefaultPollTimeout = 65 * time.Second
	connectionTimeout  = 2 * time.Second
	defaultNumChecks   = 2
)

// Error messages.
const (
	findIpTimeoutErr = "timed out looking for public IP address after %s"
	lookupServiceErr = "all public IP lookup services failed to return valid IP address"
	getServiceErr    = "failed to get service: %+v"
	responseReadErr  = "failed to read response from %s: %+v"
	responseParseErr = "response could not be parsed as valid IP address: %q"
	receivedIPv6Err  = "received IPv6 address instead of IPv4: %s"
)

// GetIP returns the caller's public IPv4 address. Multiple services are checked
// until one returns an IP address or the time out is reached.
func GetIP(lookupServices []Service, timeout time.Duration) (string, error) {
	type resultChan struct {
		ip  string
		err error
	}

	ipResultChan := make(chan resultChan)
	go func() {
		ip, err := getIpMultiCheck(lookupServices, connectionTimeout, defaultNumChecks)
		ipResultChan <- resultChan{ip, err}
	}()

	select {
	case results := <-ipResultChan:
		return results.ip, results.err
	case <-time.NewTimer(timeout).C:
		return "", errors.Errorf(findIpTimeoutErr, timeout)
	}
}

// getIpFromList returns the caller's public IP address that is provided from a
// randomly selected Service. Services are tried one by one until one return a
// valid IPv4 address.
func getIpFromList(urls []Service, timeout time.Duration) (string, error) {
	return getIpMultiCheck(urls, timeout, 1)
}

// getIpMultiCheck returns the caller's public IP address that is provided from
// multiple Services. Services are tried one by one until the given number of
//  services return the same valid IPv4 address.
func getIpMultiCheck(urls []Service, timeout time.Duration, checks int) (string, error) {
	serviceList := shuffleStrings(urls)
	var ipv6 bool
	addresses := make(map[string]int)

	for _, service := range serviceList {
		// Skip services that return IPv6 addresses if the caller has an IPv6
		if ipv6 && service.v == ipv6Address {
			continue
		}

		ip, err := getIP(service.url, timeout)
		if err == nil {
			addresses[ip]++
			if addresses[ip] >= checks {
				return ip, nil
			}
			continue
		} else if strings.Contains(err.Error(), "IPv6") {
			ipv6 = true
		}

		jww.ERROR.Printf("Failed to get public IP address from %s: %v",
			service.url, err)
	}

	return "", errors.New(lookupServiceErr)
}

// getIP requests the caller's public IP from the specified URL and returns it.
// Only valid IPv4 addresses are returned. An error is returned for IPv6 or
// invalid IP address.
func getIP(url string, timeout time.Duration) (string, error) {
	jww.INFO.Printf("Getting public IP address from %s", url)

	// Issue a GET to the URL with the specified timeout
	httpClient := http.Client{Timeout: timeout}
	resp, err := httpClient.Get(url)
	if err != nil {
		return "", errors.Errorf(getServiceErr, err)
	}

	defer func() {
		if err = resp.Body.Close(); err != nil {
			jww.ERROR.Print(err)
		}
	}()

	// Read the response
	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Errorf(responseReadErr, url, err)
	}

	// Ensure the response is a valid IPv4 address
	parsedIp := net.ParseIP(strings.TrimSpace(string(ip)))
	if parsedIp == nil {
		return "", errors.Errorf(responseParseErr, trunc(string(ip), 128, true))
	} else if strings.Contains(parsedIp.String(), ":") {
		return "", errors.Errorf(receivedIPv6Err, ip)
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

// trunc truncates a string to the given number of runes. Strings are split by
// character, not by byte. If ellipses is set, then an ellipses replaces the
// last 3 characters in the substring. This function does not properly handle
// adding ellipses at the end when the last rune is greater than one byte.
func trunc(str string, limit int, ellipses bool) string {
	var chars int
	for i := range str {
		if ellipses && chars+3 >= limit {
			if i+3 == len(str) {
				return str[:i+3]
			}
			return str[:i] + "..."
		} else if chars >= limit {
			return str[:i]
		}
		chars++
	}

	return str
}

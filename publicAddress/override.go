////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package publicAddress

import (
	"github.com/pkg/errors"
	"net"
	"strconv"
	"strings"
)

// GetIpOverride returns the public IP of the caller or the override IP, if it
// is supplied. If overrideIP is empty, then the caller's IP is looked up,
// joined with the port, and returned. If the overrideIP is supplied, then it is
// joined with the port and returned. If the overrideIP has a port, then it is
// returned as is.
func GetIpOverride(overrideIP string, port int) (string, error) {
	return getIpOverride(overrideIP, port, lookupServices)
}

// getIpOverride returns the public IP address. If an override address is
// provided, it is returned instead with the provided port unless it already
// includes a port.
func getIpOverride(overrideIP string, port int, services []Service) (string, error) {
	// If an override address was not supplied, then lookup gateway's public IP
	// and combine it with the supplied port; otherwise, use the override IP
	if overrideIP == "" {
		publicIp, err := GetIP(services, DefaultPollTimeout)
		if err != nil {
			return "", errors.Errorf("failed to lookup public IP address: %v", err)
		}

		return net.JoinHostPort(publicIp, strconv.Itoa(port)), nil
	} else {
		return JoinIpPort(overrideIP, port)
	}
}

// JoinIpPort joins the ip and port together. If the ip already has a port, it
// is returned as is.
func JoinIpPort(ip string, port int) (string, error) {
	if ip == "" {
		return ip, nil
	}

	_, _, err := net.SplitHostPort(ip)
	if err != nil {
		// If it does not have a port, then append the supplied port
		if strings.Contains(err.Error(), "missing port") {
			return net.JoinHostPort(ip, strconv.Itoa(port)), nil
		} else {
			return "", errors.Errorf("failed to parse public IP address "+
				"override \"%s\": %v", ip, err)
		}
	}

	return ip, nil
}

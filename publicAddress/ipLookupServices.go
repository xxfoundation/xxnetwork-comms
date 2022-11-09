////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package publicAddress

import (
	"fmt"
	jww "github.com/spf13/jwalterweatherman"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	ipv4Address = "ipv4"
	ipv6Address = "ipv6"
)

// Service structure contains the URL to an IP lookup Service and which IP
// version it returns.
type Service struct {
	v   string // Which IP version the lookup Service returns
	url string // URL of the Service
}

// List of URLs that respond with our public IP address.
var lookupServices = []Service{
	{ipv4Address, "https://4.echoip.de"},
	{ipv4Address, "https://api.ipify.org"},
	{ipv4Address, "https://checkip.amazonaws.com"},
	{ipv4Address, "https://ipv4.icanhazip.com"},
	{ipv4Address, "https://ifconfig.me/ip"},
	{ipv4Address, "https://ip-addr.es"},
	{ipv4Address, "https://ip4.seeip.org"},
	{ipv4Address, "https://ipaddr.site"},
	{ipv4Address, "https://ipaddress.sh"},
	{ipv4Address, "https://ipcalf.com/?format=text"},
	{ipv4Address, "https://ipecho.net/plain"},
	{ipv4Address, "https://ipinfo.io/ip"},
	{ipv4Address, "https://l2.io/ip"},
	{ipv4Address, "https://myexternalip.com/raw"},
	{ipv4Address, "https://myip.dnsomatic.com"},
	{ipv4Address, "https://sfml-dev.org/ip-provider.php"},
	{ipv4Address, "https://v4.ident.me"},
	{ipv4Address, "https://v4.ipv6-test.com/api/myip.php"},
	{ipv6Address, "https://ifconfig.co/ip"},
	{ipv6Address, "https://curlmyip.net"},
	{ipv6Address, "https://diagnostic.opendns.com/myip"},
	{ipv6Address, "https://ip.tyk.nu"},
	{ipv6Address, "https://wgetip.com"},
	{ipv6Address, "https://bot.whatismyipaddress.com/"},
	{ipv6Address, "https://ipof.in/txt"},
}

// MakeTestLookupService creates a test server and service list containing the
// IP of that server. The server will respond with the provided address. This
// function is intended for testing only so that tests do not need to reach an
// external service to function. Once the returned server is used, it should be
// closed using ts.Close().
func MakeTestLookupService(ip string, i interface{}) ([]Service, *httptest.Server) {
	switch i.(type) {
	case *testing.T, *testing.M, *testing.B, *testing.PB:
		break
	default:
		jww.FATAL.Panicf("Provided interface is not for testing: %T", i)
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, ip); err != nil {
			jww.FATAL.Panicf("Failed to write to response writer: %+v", err)
		}
	}))

	return []Service{{ipv4Address, ts.URL}, {ipv4Address, ts.URL}}, ts
}

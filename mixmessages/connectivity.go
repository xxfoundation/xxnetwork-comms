package mixmessages

import (
	"fmt"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/xx_network/comms/connect"
	"golang.org/x/net/context"
	"net"
	"strconv"
)

type ConnectivityChecker struct{}

func (cc *ConnectivityChecker) CheckConnectivity(ctx context.Context, addr *Address) (*ConnectivityResponse, error) {
	// Get gateway IP and port
	senderIp, senderPort, err := connect.GetAddressFromContext(ctx)
	if err != nil {
		return &ConnectivityResponse{}, err
	}
	return CheckConnectivity(senderIp, senderPort, addr.IP, addr.Port)
}

func CheckConnectivity(senderIp, senderPort, ip, port string) (*ConnectivityResponse, error) {
	resp := &ConnectivityResponse{}
	resp.CallerAddr = fmt.Sprintf("%s:%s", senderIp, senderPort)
	p, err := strconv.Atoi(port)
	if err != nil {
		return resp, err
	}
	if p != 0 {
		if ip != "" {
			resp.CallerAvailable = checkConn(resp.CallerAddr)

			resp.OtherAvailable = checkConn(fmt.Sprintf("%s:%s", ip, port))
		} else {
			resp.CallerAvailable = checkConn(fmt.Sprintf("%s:%s", senderIp, port))
		}
	} else {
		resp.CallerAvailable = checkConn(resp.CallerAddr)
	}
	return resp, nil
}

func checkConn(addr string) bool {
	ret := false
	conn, err := net.Dial("tcp", addr)
	if err == nil {
		ret = true
		err = conn.Close()
		if err != nil {
			jww.WARN.Printf("Error closing connection to %s: %+v", addr, err)
		}
	}
	return ret
}

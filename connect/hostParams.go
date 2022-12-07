////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package connect

import (
	"gitlab.com/xx_network/primitives/exponential"
	"google.golang.org/grpc/keepalive"
	"time"
)

// HostParams is the configuration object for Host creation
type HostParams struct {
	// Set maximum number of connection attempts
	MaxRetries uint32

	// Set maximum number of transmission attempts
	MaxSendRetries uint32

	// Toggle authorization for Host
	AuthEnabled bool

	// Toggles connection cool off
	EnableCoolOff bool

	// Number of leaky bucket sends before it stops
	NumSendsBeforeCoolOff uint32

	// Amount of time after a cool off is triggered before allowed to send again
	CoolOffTimeout time.Duration

	// Message send timeout (context deadline)
	SendTimeout time.Duration

	// Online ping timeout
	PingTimeout time.Duration

	// If set, metric handling will be enabled on this host
	EnableMetrics bool

	// If false, a connection will only be established when a comm is sent
	// else, a connection will be established immediately upon host creation
	DisableLazyConnection bool

	// If false, if transmitting to this host and there isnt a connection,
	// the host will auto connect. If true, when transmitting to a not
	// connected host, an error will be returned.
	DisableAutoConnect bool

	// List of sending errors that are deemed unimportant
	// Reception of these errors will not update the Metric state
	ExcludeMetricErrors []string

	// KeepAlive Options for Host connections
	KaClientOpts keepalive.ClientParameters

	// ProxyErrorMetricParams are the parameters used for the proxy error
	// tracker that uses exponential moving average (exponential.MovingAvg).
	ProxyErrorMetricParams exponential.MovingAvgParams

	// ConnectionType describes the method for the underlying host connection
	ConnectionType ConnectionType
	WebParams      WebConnParam
}

// GetDefaultHostParams Get default set of host params
func GetDefaultHostParams() HostParams {
	return HostParams{
		MaxRetries:            100,
		MaxSendRetries:        3,
		AuthEnabled:           true,
		EnableCoolOff:         false,
		NumSendsBeforeCoolOff: 3,
		CoolOffTimeout:        60 * time.Second,
		SendTimeout:           2 * time.Minute,
		PingTimeout:           5 * time.Second,
		EnableMetrics:         false,
		DisableLazyConnection: false,
		DisableAutoConnect:    false,
		ExcludeMetricErrors:   make([]string, 0),
		KaClientOpts: keepalive.ClientParameters{
			// Send keepAlive every Time interval
			Time: 5 * time.Second,
			// Timeout after last successful keepAlive to close connection
			Timeout: 60 * time.Second,
			// For all connections, with and without streaming
			PermitWithoutStream: true,
		},
		ProxyErrorMetricParams: exponential.DefaultMovingAvgParams(),
		ConnectionType:         GetDefaultConnectionType(),
	}
}

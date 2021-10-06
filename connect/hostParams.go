///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package connect

import (
	"google.golang.org/grpc/keepalive"
	"time"
)

// HostParams is the configuration object for Host creation
type HostParams struct {
	MaxRetries  uint32
	AuthEnabled bool

	// Toggles connection cool off
	EnableCoolOff bool

	// Number of leaky bucket sends before it stops
	NumSendsBeforeCoolOff uint32

	// Amount of time after a cool off is triggered before allowed to send again
	CoolOffTimeout time.Duration

	// Message sending timeout
	SendTimeout time.Duration

	// Online ping timeout
	PingTimeout time.Duration

	// If set, metric handling will be enabled on this host
	EnableMetrics bool

	// List of sending errors that are deemed unimportant
	// Reception of these errors will not update the Metric state
	ExcludeMetricErrors []string

	// KeepAlive Options for Host connections
	KaClientOpts keepalive.ClientParameters
}

// GetDefaultHostParams Get default set of host params
func GetDefaultHostParams() HostParams {
	return HostParams{
		MaxRetries:            100,
		AuthEnabled:           true,
		EnableCoolOff:         false,
		NumSendsBeforeCoolOff: 3,
		CoolOffTimeout:        60 * time.Second,
		SendTimeout:           2 * time.Minute,
		PingTimeout:           5 * time.Second,
		EnableMetrics:         false,
		ExcludeMetricErrors:   make([]string, 0),
		KaClientOpts: keepalive.ClientParameters{
			// Send keepAlive every Time interval
			Time: 5 * time.Second,
			// Timeout after last successful keepAlive to close connection
			Timeout: 60 * time.Second,
			// For all connections, with and without streaming
			PermitWithoutStream: true,
		},
	}
}
